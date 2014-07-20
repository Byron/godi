package seal

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/Byron/godi/api"
)

func (s *SealCommand) Generate(done <-chan bool) (<-chan godi.FileInfo, <-chan godi.Result) {
	files := make(chan godi.FileInfo)
	results := make(chan godi.Result)

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
func (s *SealCommand) traverseFilesRecursively(files chan<- godi.FileInfo, results chan<- godi.Result, done <-chan bool, tree string) bool {
	select {
	case <-done:
		return false
	default:
		{
			// read dir and, build file info, and recurse into subdirectories
			f, err := os.Open(tree)
			if err != nil {
				results <- &godi.BasicResult{nil, "", err, godi.Error}
				return false
			}
			dirInfos, err := f.Readdir(-1)
			f.Close()
			if err != nil {
				results <- &godi.BasicResult{nil, "", err, godi.Error}
				return false
			}

			// first generate infos
			for _, fi := range dirInfos {
				if !fi.IsDir() {
					path := filepath.Join(tree, fi.Name())
					if !fi.Mode().IsRegular() {
						results <- &godi.BasicResult{nil, fmt.Sprintf("Ignoring symbolic link: '%s'", path), nil, godi.Warn}
						continue
					}
					if fr, _ := utf8.DecodeRuneInString(fi.Name()); fr == '.' {
						results <- &godi.BasicResult{nil, fmt.Sprintf("Ignoring hidden file: '%s'", path), nil, godi.Warn}
						continue
					}
					files <- godi.FileInfo{
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
