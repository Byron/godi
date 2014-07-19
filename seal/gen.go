package seal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Byron/godi/api"
)

func (s *SealCommand) Generate(done <-chan bool) (<-chan api.FileInfo, <-chan api.Result) {
	files := make(chan api.FileInfo)
	results := make(chan api.Result)

	s.pCtrl = api.NewReadChannelController(s.nReaders)

	go func() {
		for _, tree := range s.Trees {
			if !s.traverseFilesRecursively(files, results, done, tree) {
				// interrupted usually, or there was an error
				break
			}
		}
		defer close(files)
	}()

	return files, results
}

// Traverse recursively, return false if the caller should stop traversing due to an error
func (s *SealCommand) traverseFilesRecursively(files chan<- api.FileInfo, results chan<- api.Result, done <-chan bool, tree string) bool {
	select {
	case <-done:
		return false
	default:
		{
			// read dir and, build file info, and recurse into subdirectories
			f, err := os.Open(tree)
			if err != nil {
				results <- &SealResult{nil, "", err, api.Error}
				return false
			}
			dirInfos, err := f.Readdir(-1)
			f.Close()
			if err != nil {
				results <- &SealResult{nil, "", err, api.Error}
				return false
			}

			// first generate infos
			for _, fi := range dirInfos {
				if !fi.IsDir() {
					path := filepath.Join(tree, fi.Name())
					if !fi.Mode().IsRegular() {
						results <- &SealResult{nil, fmt.Sprintf("Ignoring symbolic link: '%s'", path), nil, api.Warn}
						continue
					}
					if fi.Name()[0] == '.' {
						results <- &SealResult{nil, fmt.Sprintf("Ignoring hidden file: '%s'", path), nil, api.Warn}
						continue
					}
					files <- api.FileInfo{
						Path: path,
						Size: fi.Size(),
					}
				}
			}

			// then recurse
			for _, fi := range dirInfos {
				if fi.IsDir() {
					if !s.traverseFilesRecursively(files, results, done, filepath.Join(tree, fi.Name())) {
						return false
					}
				}
			}
		} //  default
	} // selcect

	return true
}
