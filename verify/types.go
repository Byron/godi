// Implements verification of seal files previously written with seal command
package verify

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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
	*godi.BasicResult                // will contain the actual file information from the disk file
	ifinfo            *godi.FileInfo // the file information we have seen in the index
}

// NewCommand returns an initialized verify command
func NewCommand(trees []string, nReaders int) (c VerifyCommand, err error) {
	err = c.Init(nReaders, 0, trees)
	return
}

func (v *VerifyCommand) NumChannels() int {
	return v.pCtrl.Streams()
}

func (s *VerifyCommand) Generate(done <-chan bool) (<-chan godi.FileInfo, <-chan godi.Result) {
	generate := func(files chan<- godi.FileInfo, results chan<- godi.Result) {
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

			indexDir := filepath.Dir(index)
			if err == nil {
				for _, fi := range fileInfos {
					// Figure out the path to use - for now we use the relative one
					fi.Path = filepath.Join(indexDir, fi.RelaPath)
					files <- fi
				}
			} else {
				results <- &VerifyResult{
					BasicResult: &godi.BasicResult{Err: err},
				}
				continue
			}
		} // for each index
	}

	return godi.Generate(done, generate)
}

func (s *VerifyCommand) Gather(files <-chan godi.FileInfo, results chan<- godi.Result, wg *sync.WaitGroup, done <-chan bool) {
	makeResult := func(f *godi.FileInfo) (godi.Result, *godi.BasicResult) {
		fcpy := *f
		res := VerifyResult{
			BasicResult: &godi.BasicResult{
				Finfo: f,
				Prio:  godi.Progress,
			},
			ifinfo: &fcpy,
		}
		return &res, res.BasicResult
	}

	godi.Gather(files, results, wg, done, makeResult, &s.pCtrl, nil)
}

func (s *VerifyCommand) Aggregate(results <-chan godi.Result, done <-chan bool) <-chan godi.Result {

	var signatureMismatches uint = 0
	resultHandler := func(r godi.Result, accumResult chan<- godi.Result) bool {
		vr := r.(*VerifyResult)

		hasError := false
		if (len(vr.ifinfo.Sha1) > 0 && bytes.Compare(vr.ifinfo.Sha1, vr.Finfo.Sha1) != 0) ||
			(len(vr.ifinfo.MD5) > 0 && bytes.Compare(vr.ifinfo.MD5, vr.Finfo.MD5) != 0) {
			vr.Err = fmt.Errorf("HASH MISMATCH: %s", vr.Finfo.Path)
			signatureMismatches += 1
			hasError = true
		} else {
			vr.Msg = fmt.Sprintf("OK: %s", vr.Finfo.Path)
		}
		accumResult <- vr
		return !hasError
	}

	finalizer := func(
		accumResult chan<- godi.Result,
		st *godi.AggregateFinalizerState) {

		if signatureMismatches == 0 {
			accumResult <- &VerifyResult{
				BasicResult: &godi.BasicResult{
					Msg: fmt.Sprintf(
						"All %d files did not change after sealing (%v)",
						st.FileCount,
						st,
					),
					Prio: godi.Info,
				},
			}
		} else {
			st.ErrCount -= signatureMismatches
			accumResult <- &VerifyResult{
				BasicResult: &godi.BasicResult{
					Msg: fmt.Sprintf(
						"%d of %d files have changed on disk after sealing (%v)",
						signatureMismatches,
						st.FileCount,
						st,
					),
					Prio: godi.Info,
				},
			}
		}
	}

	return godi.Aggregate(results, done, resultHandler, finalizer)
}
