package seal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/utility"
)

// Must be kept in sync with IndexPath generator
var reIsIndexPath = regexp.MustCompile(fmt.Sprintf(`%s_\d{4}-\d{2}-\d{2}_\d{2}\d{2}\d{2}\.%s`, IndexBaseName, codec.GobExtension))

// return a path to an index file residing at tree
func (s *SealCommand) IndexPath(tree string, extension string) string {
	n := time.Now()
	return filepath.Join(tree, fmt.Sprintf("%s_%04d-%02d-%02d_%02d%02d%02d.%s",
		IndexBaseName,
		n.Year(),
		n.Month(),
		n.Day(),
		n.Hour(),
		n.Minute(),
		n.Second(),
		extension))
}

// When called, we have seen no error in the given mapping of relativePaths to FileInfos
// Returns error in case we failed to produce an index
// It's up to the caller to remove existing files on error
func (s *SealCommand) writeIndex(commonTree string, paths []codec.SerializableFileInfo) (string, error) {
	// Serialize all fileinfo structures
	// NOTE: As the parallel writer will send results only when writing finished, we can just operate serially here ...
	// For this there is also no need to optimize performance
	// However, we could use a parallel writer if we are so inclined
	// For now, we do only gob
	encoder := codec.Gob{}

	// It's currently possible to have no paths as we pre-allocate these and don't care if we are in copy mode
	if len(paths) == 0 {
		return "", fmt.Errorf("No paths provided - cannot write seal for nothing")
	}

	// This will and should fail if the file already exists
	fd, err := os.OpenFile(s.IndexPath(commonTree, encoder.Extension()), os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return "", err
	}

	err = encoder.Serialize(paths, fd)
	fd.Close()
	if err != nil {
		// remove possibly half-written file
		os.Remove(fd.Name())
		return "", err
	}
	return fd.Name(), nil
}

func (s *SealCommand) Aggregate(results <-chan godi.Result) <-chan godi.Result {
	treePathsmap := make(map[string][]codec.SerializableFileInfo)

	resultHandler := func(r godi.Result, accumResult chan<- godi.Result) bool {
		sr := r.(*SealResult)

		// Keep the file-information, in any case
		pathInfos := treePathsmap[sr.Finfo.Root()]
		pathInfos = append(pathInfos, codec.SerializableFileInfo{sr.Finfo, r.Error()})
		treePathsmap[sr.Finfo.Root()] = pathInfos

		// We will keep track of the file even if it reported an error.
		// That way, we can later determine what to cleanup
		relaPath := sr.Finfo.RelaPath
		if r.Error() != nil {
			accumResult <- r
			return false
		}

		// we store only relative paths in the map - this is all we want to serialize
		hasError := false

		if len(sr.source) == 0 {
			sr.Msg = fmt.Sprintf("DONE ...%s", relaPath)
		} else {
			sr.Msg = fmt.Sprintf("DONE CP %s -> %s", sr.source, sr.Finfo.Path)
		}

		accumResult <- sr
		return !hasError
	} // end resultHandler()

	finalizer := func(
		accumResult chan<- godi.Result,
		st *godi.AggregateFinalizerState) {

		// ROLLBACK HANDLING
		/////////////////////
		// Check each destination tree for errors. If there are some, and if we are in write mode,
		// remove files we have written so far. Otherwise, create an index file
		// We will do one delete operation per write-device, in parallel. Each device worker will
		// operate on all trees on that device in order
		// Pre-allocate a bunch of tree strings - it's at max the total amount of destinations, which might
		// all be on one device
		// We natually don't do anything in non-copy mode as we have no writers
		treesWithError := make([]string, utility.WriteChannelDeviceMapTrees(s.rootedWriters))
		nit := 0 // num invalid trees
		var wg sync.WaitGroup
		for _, wctrl := range s.rootedWriters {
			init := nit // initial nit
			for _, tree := range wctrl.Trees {
				for _, sfi := range treePathsmap[tree] {
					if sfi.Err != nil {
						treesWithError[nit] = tree
						nit += 1
						break
					}
				} // for each file-info below tree
			} // for each destination of write controller

			// found some - shoot off go routine
			if nit != init {
				wg.Add(1)
				go func(trees []string) {
					for _, tree := range trees {
						pathInfos := treePathsmap[tree]

						// Remove all previously written files in this tree if we are in write mode
						accumResult <- &godi.BasicResult{
							Msg:  fmt.Sprintf("Rolling back changes at copy destination '%s'", tree),
							Prio: godi.Info,
						}

						// For path deletion to work correctly, we need it sorted
						sort.Sort(codec.ByLongestPathDescending(pathInfos))

						for _, sfi := range pathInfos {
							// We may only remove it if the error wasn't a 'Existed' one, or we kill a file
							// that wasn't created in this run.
							var msg string
							prio := godi.Error
							err := sfi.Err
							if err != nil && os.IsExist(err) {
								msg = fmt.Sprintf("Won't remove existing file: '%s'", sfi.Path)
								prio = godi.Info
								err = nil
							} else {
								err = os.Remove(sfi.Path)
								if err == nil {
									msg = fmt.Sprintf("Removed '%s'", sfi.Path)
									prio = godi.Info

									// try to remove the directory - will fail if non-empty.
									// only do that if we wouldn't remove the tree.
									// Also crawl upwards
									var derr error
									for dir := filepath.Dir(sfi.Path); dir != tree && derr == nil; dir = filepath.Dir(dir) {
										derr = os.Remove(dir)
									}
								}
							}

							accumResult <- &godi.BasicResult{
								Finfo: sfi.FileInfo,
								Msg:   msg,
								Err:   err,
								Prio:  prio,
							}
						} // for each path info
					} // for each tree to handle
					wg.Done()
				}(treesWithError[init:nit]) // go handle errors in write mode
			} // if we have invalid trees on that device
		} // for each write controller (per device)

		// INDEX HANDLING
		//////////////////
		// The cool thing is, that cleanup happens in parallel !
		for tree, pathInfos := range treePathsmap {
			// Especially source path maps will be empty - the only results we see is the destination paths
			if len(pathInfos) == 0 {
				continue
			}

			// See if we have any error below this root
			foundError := false
			for _, invalidTree := range treesWithError[:nit] {
				if tree == invalidTree {
					foundError = true
					break
				}
			}

			if !foundError && !st.WasCancelled {
				// Only done if there are no errors below the current tree
				// Serialize all fileinfo structures
				if index, err := s.writeIndex(tree, pathInfos); err != nil {
					st.ErrCount += 1
					accumResult <- &godi.BasicResult{Err: err, Prio: godi.Error}
				} else {
					accumResult <- &godi.BasicResult{
						Finfo: godi.FileInfo{Path: index, Size: -1},
						Msg:   fmt.Sprintf("Wrote seal at '%s'", index),
						Prio:  godi.Info,
					}
				} // handle index writing errors
			} // handle errors in tree
		} // for each item in treePathMap

		// Wait for cleanup jobs
		wg.Wait()

		accumResult <- &godi.BasicResult{
			Msg: fmt.Sprintf(
				"Sealed %d files (%v)",
				st.FileCount,
				st,
			),
			Prio: godi.Info,
		}
	} // end finalizer()

	return godi.Aggregate(results, s.Done, resultHandler, finalizer)
}
