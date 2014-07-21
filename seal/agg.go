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
	treePathmap := make(map[string]map[string]*godi.FileInfo)
	// Presort all paths by their root
	for _, tree := range s.SourceTrees {
		treePathmap[tree] = make(map[string]*godi.FileInfo)
	}

	resultHandler := func(r godi.Result, accumResult chan<- godi.Result) bool {
		sr := r.(*godi.BasicResult)
		// find root
		var pathmap map[string]*godi.FileInfo
		var pathTree string
		for _, tree := range s.SourceTrees {
			if strings.HasPrefix(sr.Finfo.Path, tree) {
				pathTree = tree
				pathmap = treePathmap[tree]
				break
			}
		}
		if pathmap == nil {
			panic(fmt.Sprintf("Couldn't determine root of path '%s', roots are %v", sr.Finfo.Path, s.SourceTrees))
		}
		// we store only relative paths in the map - this is all we want to serialize
		relaPath := sr.Finfo.Path[len(pathTree)+1:]

		hasError := false
		if _, pathSeen := pathmap[relaPath]; pathSeen {
			// duplicate path - shouldn't happen, but we deal with it
			sr.Err = fmt.Errorf("Path '%s' was handled multiple times - ignoring occurrence", sr.Finfo.Path)
			hasError = true
		} else {
			pathmap[relaPath] = sr.Finfo
			sr.Msg = fmt.Sprintf("DONE ...%s", relaPath)
		}

		accumResult <- sr
		return !hasError
	}

	finalizer := func(
		accumResult chan<- godi.Result,
		st *godi.AggregateFinalizerState) {
		// INDEX HANDLING
		//////////////////
		// Serialize all fileinfo structures
		// NOTE: As the parallel writer will send results only when writing finished, we can just operate serially here ...
		// For this there is also no need to optimize performance
		// However, we could use a parallel writer if we are so inclined
		if !st.WasCancelled {
			var indices []string
			var err error
			if indices, err = s.writeIndex(treePathmap); err != nil {
				// NOTE: We keep successfully written indices as each tree is consistent in itself
				st.ErrCount += 1
				accumResult <- &godi.BasicResult{Err: err, Prio: godi.Error}
			}

			// Inform about successfully written indices
			for _, index := range indices {
				accumResult <- &godi.BasicResult{
					&godi.FileInfo{Path: index, Size: -1},
					fmt.Sprintf("Wrote seal at '%s'", index), err, godi.Info,
				}
			}
		}

		accumResult <- &godi.BasicResult{
			Msg: fmt.Sprintf(
				"Sealed %d files (%v)",
				st.FileCount,
				st,
			),
			Prio: godi.Info,
		}
	}

	return godi.Aggregate(results, done, resultHandler, finalizer)
}
