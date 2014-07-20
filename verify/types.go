// Implements verification of seal files previously written with seal command
package verify

import (
	"os"
	"sync"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/utility"
)

const (
	Name = "verify"
)

// A type representing all arguments required to drive a Seal operation
type VerifyCommand struct {

	// The index files we are to verify
	Indices []string

	// parallel reader
	pCtrl utility.ReadChannelController
}

// Implements information about a verify operation
type VerifyResult struct {
	*godi.BasicResult
	nfinfo *godi.FileInfo // the newly gathered file information
}

// NewCommand returns an initialized verify command
func NewCommand(trees []string, nReaders int) (c VerifyCommand, err error) {
	err = c.Init(nReaders, 0, trees)
	return
}

func (s *VerifyCommand) Generate(done <-chan bool) (<-chan godi.FileInfo, <-chan godi.Result) {
	files := make(chan godi.FileInfo)
	results := make(chan godi.Result)

	go func() {
		for _, index := range s.Indices {
			c := codec.NewByPath(index)
			if c == nil {
				panic("Should have a codec here - this was checked before")
			}

			fd, err := os.Open(index)
			if err != nil {
				results <- &VerifyResult{
					BasicResult: &godi.BasicResult{Err: err},
				}
				continue
			}

			fileInfos, err := c.Deserialize(fd)
			fd.Close()
			if err == nil {
				for _, fi := range fileInfos {
					files <- fi
				}
			} else {
				results <- &VerifyResult{
					BasicResult: &godi.BasicResult{Err: err},
				}
				continue
			}
		} // for each index
	}()

	return files, results
}

// TODO: this should really be the same as in seal, but use a different result type. There should be a Gather utility function
func (s *VerifyCommand) Gather(files <-chan godi.FileInfo, results chan<- godi.Result, wg *sync.WaitGroup, done <-chan bool) {
	defer wg.Done()

	// This MUST be a copy of f here, otherwise we will be in trouble thanks to the user of defer in handleHash
	// we will get f overwritten by the next iteration variable ... it's kind of special, might
	// be intersting for the mailing list.
	handleHash := func(f godi.FileInfo) {
	}

	for f := range files {
		select {
		case <-done:
			return
		default:
			handleHash(f)
		}
	}
}

func (s *VerifyCommand) Aggregate(results <-chan godi.Result, done <-chan bool) <-chan godi.Result {

	resultHandler := func(r godi.Result, accumResult chan<- godi.Result) bool {
		return true
	}

	finalizer := func(
		accumResult chan<- godi.Result,
		st godi.AggregateFinalizerState) {

	}

	return godi.Aggregate(results, done, resultHandler, finalizer)
}
