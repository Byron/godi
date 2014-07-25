package seal

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/Byron/godi/api"
)

func (s *SealCommand) Generate() (<-chan api.FileInfo, <-chan api.Result) {
	generate := func(files chan<- api.FileInfo, results chan<- api.Result) {
		for _, tree := range s.Items {
			if !s.traverseFilesRecursively(files, results, s.Done, tree, tree) {
				// interrupted usually, or there was an error
				break
			}
		}
	}

	return api.Generate(generate)
}

// Traverse recursively, return false if the caller should stop traversing due to an error
func (s *SealCommand) traverseFilesRecursively(files chan<- api.FileInfo, results chan<- api.Result, done <-chan bool, tree string, root string) bool {
	select {
	case <-done:
		return false
	default:
		{
			// read dir and, build file info, and recurse into subdirectories
			f, err := os.Open(tree)
			if err != nil {
				results <- &api.BasicResult{
					Err:  err,
					Prio: api.Error,
				}
				return false
			}
			dirInfos, err := f.Readdir(-1)
			f.Close()
			if err != nil {
				results <- &api.BasicResult{
					Err:  err,
					Prio: api.Error,
				}
				return false
			}

			// first generate infos
			for _, fi := range dirInfos {
				if !fi.IsDir() {
					path := filepath.Join(tree, fi.Name())
					if !fi.Mode().IsRegular() {
						results <- &api.BasicResult{
							Msg:  fmt.Sprintf("Ignoring symbolic link: '%s'", path),
							Prio: api.Warn,
						}
						continue
					}
					if fr, _ := utf8.DecodeRuneInString(fi.Name()); fr == '.' {
						results <- &api.BasicResult{
							Msg:  fmt.Sprintf("Ignoring hidden file: '%s'", path),
							Prio: api.Warn,
						}
						continue
					}
					if reIsIndexPath.Match([]byte(fi.Name())) {
						results <- &api.BasicResult{
							Msg:  fmt.Sprintf("Ignoring godi index: '%s'", path),
							Prio: api.Warn,
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
					if !s.traverseFilesRecursively(files, results, done, filepath.Join(tree, fi.Name()), root) {
						return false
					}
				}
			}
		} //  default
	} // selcect

	return true
}
