package seal

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/Byron/godi/api"
)

func (s *SealCommand) Generate(feedback <-chan string) (<-chan api.FileInfo, <-chan api.Result) {
	generate := func(files chan<- api.FileInfo, results chan<- api.Result) {
		for _, tree := range s.Items {
			cancelled, treeError := s.traverseFilesRecursively(files, results, s.Done, feedback, tree, tree)
			if cancelled {
				// interrupted usually, or there was an error
				break
			} else if treeError {
				// Just abort processing of this tree
				continue
			}
		}
	}

	return api.Generate(generate)
}

// Traverse recursively, return false if the caller should stop traversing due to an error
func (s *SealCommand) traverseFilesRecursively(files chan<- api.FileInfo, results chan<- api.Result, done <-chan bool, feedback <-chan string, tree string, root string) (bool, bool) {
	select {
	case <-done:
		return true, false
	default:
	} // select

	// read dir and, build file info, and recurse into subdirectories
	f, err := os.Open(tree)
	if err != nil {
		results <- &api.BasicResult{
			Err:  err,
			Prio: api.Error,
		}
		return false, true
	}
	dirInfos, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		results <- &api.BasicResult{
			Err:  err,
			Prio: api.Error,
		}
		return false, true
	}

	// first generate infos
	for _, fi := range dirInfos {
		// Have to try real hard to drain the feedback channel, to assure it will never block to the writer
		select {
		case failedTree := <-feedback:
			{
				if failedTree == root {
					return false, true
				}
			}
		default:
		} // select

		if !fi.IsDir() {
			path := filepath.Join(tree, fi.Name())
			if !fi.Mode().IsRegular() {
				results <- &api.BasicResult{
					Msg:  fmt.Sprintf("Ignoring symbolic link: '%s'", path),
					Prio: api.Info,
				}
				continue
			}
			if fr, _ := utf8.DecodeRuneInString(fi.Name()); fr == '.' {
				results <- &api.BasicResult{
					Msg:  fmt.Sprintf("Ignoring hidden file: '%s'", path),
					Prio: api.Info,
				}
				continue
			}
			if reIsIndexPath.Match([]byte(fi.Name())) {
				results <- &api.BasicResult{
					Msg:  fmt.Sprintf("Ignoring godi index: '%s'", path),
					Prio: api.Info,
				}
				continue
			}

			files <- api.FileInfo{
				Path:     path,
				RelaPath: path[len(root)+1:],
				Size:     fi.Size(),
			}
		}
	}

	// then recurse
	for _, fi := range dirInfos {
		if fi.IsDir() {
			cancelled, treeError := s.traverseFilesRecursively(files, results, done, feedback, filepath.Join(tree, fi.Name()), root)
			if cancelled || treeError {
				return cancelled, treeError
			}
		}
	}

	return false, false
}
