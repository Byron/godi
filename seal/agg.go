package seal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
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
func (s *SealCommand) writeIndex(commonTree string, pathMap map[string]*godi.FileInfo) (string, error) {
	// Serialize all fileinfo structures
	// NOTE: As the parallel writer will send results only when writing finished, we can just operate serially here ...
	// For this there is also no need to optimize performance
	// However, we could use a parallel writer if we are so inclined
	// For now, we do only gob
	encoder := codec.Gob{}

	// It's currently possible to have no paths as we pre-allocate these and don't care if we are in copy mode
	if len(pathMap) == 0 {
		return "", fmt.Errorf("No paths provided - cannot write seal for nothing")
	}

	// This will and should fail if the file already exists
	fd, err := os.OpenFile(s.IndexPath(commonTree, encoder.Extension()), os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return "", err
	}

	err = encoder.Serialize(pathMap, fd)
	fd.Close()
	if err != nil {
		// remove possibly half-written file
		os.Remove(fd.Name())
		return "", err
	}
	return fd.Name(), nil
}

func (s *SealCommand) Aggregate(results <-chan godi.Result) <-chan godi.Result {
	treePathmap := make(map[string]map[string]*godi.FileInfo)
	// Presort all paths by their root
	for _, tree := range s.Items {
		treePathmap[tree] = make(map[string]*godi.FileInfo)
	}

	// Fill the root-map with the write-roots, if available
	for _, rctrl := range s.rootedWriters {
		for _, tree := range rctrl.Trees {
			treePathmap[tree] = make(map[string]*godi.FileInfo)
		}
	}

	resultHandler := func(r godi.Result, accumResult chan<- godi.Result) bool {
		sr := r.(*SealResult)

		// find root
		pathmap := treePathmap[sr.Finfo.Root()]

		// We will keep track of the file even if it reported an error.
		// That way, we can later determine what to cleanup
		relaPath := sr.Finfo.RelaPath
		if r.Error() != nil {
			pathmap[relaPath] = sr.Finfo
			accumResult <- r
			return false
		}

		if pathmap == nil {
			panic(fmt.Sprintf("Didn't find map matching root '%s'", sr.Finfo.Root()))
		}
		// we store only relative paths in the map - this is all we want to serialize
		hasError := false

		if _, pathSeen := pathmap[relaPath]; pathSeen {
			// duplicate path - shouldn't happen, but we deal with it
			sr.Err = fmt.Errorf("Path '%s' was handled multiple times - ignoring occurrence", sr.Finfo.Path)
			hasError = true
		} else {
			pathmap[relaPath] = sr.Finfo
			if len(sr.source) == 0 {
				sr.Msg = fmt.Sprintf("DONE ...%s", relaPath)
			} else {
				sr.Msg = fmt.Sprintf("DONE CP %s -> %s", sr.source, sr.Finfo.Path)
			}
		}

		accumResult <- sr
		return !hasError
	} // end resultHandler()

	finalizer := func(
		accumResult chan<- godi.Result,
		st *godi.AggregateFinalizerState) {

		// Check each destination tree for errors. If there are some, and if we are in write mode,
		// remove files we have written so far. Otherwise, create an index file
		for tree, pathMap := range treePathmap {
			// Especially source path maps will be empty - the only results we see is the destination paths
			if len(pathMap) == 0 {
				continue
			}

			foundError := false
			for _, fi := range pathMap {
				// we have to mark an error if we don't have a sha or md5
				if fi.HasError() {
					foundError = true
					break
				}
			}

			if foundError {
				// Remove all previously written files in this tree if we are in write mode
				// TODO(st): use parallel per-ctrl writer to do that
				if len(s.rootedWriters) > 0 {
					for _, fi := range pathMap {
						err := os.Remove(fi.Path)

						var msg string
						prio := godi.Error
						if err == nil {
							msg = fmt.Sprintf("Removed '%s'", fi.Path)
							prio = godi.Error
						}

						accumResult <- &godi.BasicResult{
							Finfo: fi,
							Msg:   msg,
							Err:   err,
							Prio:  prio,
						}
					}
				} // handle errors in write mode
			} else if !st.WasCancelled {
				// INDEX HANDLING
				//////////////////
				// Serialize all fileinfo structures
				if index, err := s.writeIndex(tree, pathMap); err != nil {
					st.ErrCount += 1
					accumResult <- &godi.BasicResult{Err: err, Prio: godi.Error}
				} else {
					accumResult <- &godi.BasicResult{
						Finfo: &godi.FileInfo{Path: index, Size: -1},
						Msg:   fmt.Sprintf("Wrote seal at '%s'", index),
						Prio:  godi.Info,
					}
				} // handle index writing errors
			} // handle errors in tree
		} // for each item in treePathMap

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
