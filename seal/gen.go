package seal

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/Byron/godi/api"
)

func sendErrorAtRoot(results chan<- api.Result, err error, root string) {
	results <- &SealResult{
		BasicResult: api.BasicResult{
			Err:   err,
			Prio:  api.Error,
			Finfo: api.FileInfo{Path: root},
		},
	}
}

func (s *Command) Generate() <-chan api.Result {
	generate := func(trees []string, files chan<- api.FileInfo, results chan<- api.Result) {
		for _, tree := range trees {
			// could also be a file
			if tstat, err := os.Stat(tree); err != nil {
				sendErrorAtRoot(results, fmt.Errorf("Couldn't access tree or file '%s': %v", tree, err), tree)
				continue
			} else if !tstat.IsDir() {
				// Assume it's a file and send it of like that
				files <- api.FileInfo{
					Path:     tree,
					RelaPath: filepath.Base(tree),
					Mode:     tstat.Mode(),
					Size:     tstat.Size(),
				}
				continue
			}

			cancelled, treeError := s.traverseFilesRecursively(files, results, s.Done, tree, tree)
			if cancelled {
				// interrupted usually, or there was an error
				break
			} else if treeError {
				// Just abort processing of this tree
				continue
			}
		}
	}

	return api.Generate(s.RootedReaders, s, generate)
}

// Traverse recursively, return false if the caller should stop traversing due to an error
func (s *Command) traverseFilesRecursively(files chan<- api.FileInfo, results chan<- api.Result, done <-chan bool, tree string, root string) (bool, bool) {
	select {
	case <-done:
		return true, false
	default:
	} // select

	// read dir and, build file info, and recurse into subdirectories
	f, err := os.Open(tree)
	if err != nil {
		sendErrorAtRoot(results, err, root)
		return false, true
	}

	dirInfos, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		sendErrorAtRoot(results, err, root)
		return false, true
	}

	shouldExclude := func(tree string, fi os.FileInfo, dirOnly bool) (bool, string) {
		if (fi.Mode()&os.ModeDir != os.ModeDir) == dirOnly {
			return true, ""
		}

		path := filepath.Join(tree, fi.Name())
		for _, excludeFilter := range s.Filters {
			if excludeFilter.Matches(fi.Name(), fi.Mode()) {
				atomic.AddUint32(&s.Stats.NumSkippedFiles, 1)
				results <- &SealResult{
					BasicResult: api.BasicResult{
						Msg:   fmt.Sprintf("Ignoring '%s' at '%s'", excludeFilter, path),
						Prio:  api.Info,
						Finfo: api.FileInfo{Path: root},
					},
				}
				return true, ""
			}
		}
		return false, path
	} // func shouldExclude()

	// first generate infos
	const fileOnly = false
toNextFile:
	for _, fi := range dirInfos {

		// Actually we wouldn't need atomic access here, but lets be sure the race-detector is happy with us
		// If at least one gather had errors to all destinations, there is no meaning in delivering more paths
		if atomic.LoadUint32(&s.Stats.StopTheEngines) > 0 {
			return false, true
		}

		exclude, path := shouldExclude(tree, fi, fileOnly)
		if exclude {
			continue toNextFile
		}

		files <- api.FileInfo{
			Path:     path,
			RelaPath: path[len(root)+1:],
			Mode:     fi.Mode(),
			Size:     fi.Size(),
		}
	}

	// then recurse into directories, apply a filter though
	const dirOnly = !fileOnly
toNextDir:
	for _, fi := range dirInfos {
		exclude, path := shouldExclude(tree, fi, dirOnly)
		if exclude {
			continue toNextDir
		}

		cancelled, treeError := s.traverseFilesRecursively(files, results, done, path, root)
		if cancelled || treeError {
			return cancelled, treeError
		}
	}

	return false, false
}
