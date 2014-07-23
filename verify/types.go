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
)

const (
	Name = "verify"
)

// A type representing all arguments required to drive a Seal operation
type VerifyCommand struct {
	godi.BasicRunner
}

// Implements information about a verify operation
type VerifyResult struct {
	godi.BasicResult                // will contain the actual file information from the disk file
	ifinfo           *godi.FileInfo // the file information we have seen in the index
}

// NewCommand returns an initialized verify command
func NewCommand(indices []string, nReaders int) (c VerifyCommand, err error) {
	err = c.Init(nReaders, 0, indices)
	return
}

func (s *VerifyCommand) Generate() (<-chan godi.FileInfo, <-chan godi.Result) {
	generate := func(files chan<- godi.FileInfo, results chan<- godi.Result) {
		for _, index := range s.Items {
			c := codec.NewByPath(index)
			if c == nil {
				panic("Should have a codec here - this was checked before")
			}

			fd, err := os.Open(index)
			if err != nil {
				results <- &VerifyResult{
					BasicResult: godi.BasicResult{Err: err},
				}
				continue
			}

			fileInfos, err := c.Deserialize(fd)
			fd.Close()

			indexDir := filepath.Dir(index)
			if err == nil {
			forEachFileInfo:
				for _, fi := range fileInfos {
					// Have to be able to abort early. One must know that even though we may read super fast,
					// we will block until the gather threads have done the work.
					// Therefore it makes sense to check and abort here
					select {
					case <-s.Done:
						break forEachFileInfo
					default:
						{
							// Figure out the path to use - for now we use the relative one
							// NOTE: We need to use the relative one as our read-controller device map is based on that.
							// If it was the absolute file path we use here, it could possibly point to a file far away,
							// in any case our read controller map will not yield the expected result unless we set it
							// up here, which is dangerous as it is async ! So let's not use the absolute path, ever !
							fi.Path = filepath.Join(indexDir, fi.RelaPath)
							files <- fi
						}
					} // select
				} // for each fileInfo
			} else {
				results <- &VerifyResult{
					BasicResult: godi.BasicResult{Err: err},
				}
				continue
			}
		} // for each index
	}

	return godi.Generate(generate)
}

func (s *VerifyCommand) Gather(files <-chan godi.FileInfo, results chan<- godi.Result, wg *sync.WaitGroup) {
	makeResult := func(f, source *godi.FileInfo, err error) godi.Result {
		fcpy := *f
		res := VerifyResult{
			BasicResult: godi.BasicResult{
				Finfo: f,
				Prio:  godi.Progress,
				Err:   err,
			},
			ifinfo: &fcpy,
		}
		return &res
	}

	godi.Gather(files, results, wg, makeResult, s.RootedReaders, nil)
}

func (s *VerifyCommand) Aggregate(results <-chan godi.Result) <-chan godi.Result {

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
				BasicResult: godi.BasicResult{
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
				BasicResult: godi.BasicResult{
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

	return godi.Aggregate(results, s.Done, resultHandler, finalizer)
}
