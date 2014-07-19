package seal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
)

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

// When called, we have seen no error and everything can be assumed to be in order
// Returns error in case we can't write, and all index files written so far.
// It's up to the caller to remove existing files on error
func (s *SealCommand) writeIndex(treeMap map[string]map[string]*godi.FileInfo) ([]string, error) {
	// Serialize all fileinfo structures
	// NOTE: As the parallel writer will send results only when writing finished, we can just operate serially here ...
	// For this there is also no need to optimize performance
	// However, we could use a parallel writer if we are so inclined
	// For now, we do only gob
	encoder := codec.Gob{}
	indices := make([]string, 0, len(treeMap))
	for tree, paths := range treeMap {
		// This will and should fail if the file already exists
		fd, err := os.OpenFile(s.IndexPath(tree, encoder.Extension()), os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return indices, err
		}

		err = encoder.Serialize(paths, fd)
		fd.Close()
		if err != nil {
			// remove possibly half-written file
			os.Remove(fd.Name())
			return indices, err
		}
		indices = append(indices, fd.Name())
	}

	return indices, nil
}

func (s *SealCommand) Aggregate(results <-chan godi.Result, done <-chan bool) <-chan godi.Result {
	accumResult := make(chan godi.Result)

	go func() {
		defer close(accumResult)
		treePathmap := make(map[string]map[string]*godi.FileInfo)

		// Presort all paths by their root
		for _, tree := range s.Trees {
			treePathmap[tree] = make(map[string]*godi.FileInfo)
		}

		var count uint = 0
		var errCount uint = 0
		var size uint64 = 0
		st := time.Now()
		wasCancelled := false

		// ACCUMULATE PATHS INFO
		/////////////////////////
		for r := range results {
			if r.Error() != nil {
				errCount += 1
				accumResult <- r
			}

			// Be sure we take note of cancellation.
			// If this happens, soon our results will be drained and we leave naturally
			select {
			case <-done:
				wasCancelled = true
			default:
				{
					sr := r.(*SealResult)
					// find root
					var pathmap map[string]*godi.FileInfo
					var pathTree string
					for _, tree := range s.Trees {
						if strings.HasPrefix(sr.finfo.Path, tree) {
							pathTree = tree
							pathmap = treePathmap[tree]
							break
						}
					}
					if pathmap == nil {
						panic(fmt.Sprintf("Couldn't determine root of path '%s', roots are %v", sr.finfo.Path, s.Trees))
					}
					// we store only relative paths in the map - this is all we want to serialize
					relaPath := sr.finfo.Path[len(pathTree)+1:]

					_, ok := pathmap[relaPath]
					if ok {
						// duplicate path - shouldn't happen, but we deal with it
						sr.err = fmt.Errorf("Path '%s' was handled multiple times - ignoring occurrence", sr.finfo.Path)
						errCount += 1
						accumResult <- sr
					} else {
						pathmap[relaPath] = sr.finfo
						count += 1
						size += uint64(sr.finfo.Size)
						accumResult <- &SealResult{nil, fmt.Sprintf("DONE ...%s", relaPath), nil, godi.Progress}
					}
				} // default
			} // select
		} // range results
		elapsed := time.Now().Sub(st).Seconds()
		sizeMB := float32(size) / (1024.0 * 1024.0)

		// INDEX HANDLING
		//////////////////
		// Serialize all fileinfo structures
		// NOTE: As the parallel writer will send results only when writing finished, we can just operate serially here ...
		// For this there is also no need to optimize performance
		// However, we could use a parallel writer if we are so inclined
		if !wasCancelled {
			var indices []string
			var err error
			if indices, err = s.writeIndex(treePathmap); err != nil {
				// NOTE: We keep successfully written indices as each tree is consistent in itself
				accumResult <- &SealResult{nil, "", err, godi.Error}
			}

			// Inform about successfully written indices
			for _, index := range indices {
				accumResult <- &SealResult{nil, fmt.Sprintf("Wrote seal at '%s'", index), err, godi.Info}
			}
		}

		accumResult <- &SealResult{nil, fmt.Sprintf("Sealed %d files with total size of %#vMB in %vs (%#v MB/s, %d errors, cancelled=%v)", count, sizeMB, elapsed, float64(sizeMB)/elapsed, errCount, wasCancelled), nil, godi.Info}
	}()

	return accumResult
}
