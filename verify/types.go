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
	api.BasicRunner
}

// Implements information about a verify operation
type VerifyResult struct {
	api.BasicResult              // will contain the actual file information from the disk file
	ifinfo          api.FileInfo // the file information we have seen in the index
}

// NewCommand returns an initialized verify command
func NewCommand(indices []string, nReaders int) (*VerifyCommand, error) {
	c := VerifyCommand{}
	return &c, c.Init(nReaders, 0, indices, api.Info)
}

func (s *VerifyCommand) Generate() (<-chan api.FileInfo, <-chan api.Result) {
	return api.Generate(func(files chan<- api.FileInfo, results chan<- api.Result) {
		for _, index := range s.Items {
			c := codec.NewByPath(index)
			if c == nil {
				panic("Should have a codec here - this was checked before")
			}

			fd, err := os.Open(index)
			if err != nil {
				results <- &VerifyResult{
					BasicResult: api.BasicResult{Err: err},
				}
				continue
			}

			// Figure out the path to use - for now we use the relative one
			// NOTE: We need to use the relative one as our read-controller device map is based on that.
			// If it was the absolute file path we use here, it could possibly point to a file far away,
			// in any case our read controller map will not yield the expected result unless we set it
			// up here, which is dangerous as it is async ! So let's not use the absolute path, ever !
			indexDir := filepath.Dir(index)
			err = c.Deserialize(fd, files, func(v *api.FileInfo) bool {
				select {
				case <-s.Done:
					return false
				default:
					{
						v.Path = filepath.Join(indexDir, v.RelaPath)
						return true
					}
				}
			})
			fd.Close()

			if err != nil {
				results <- &VerifyResult{
					BasicResult: api.BasicResult{Err: err},
				}
				continue
			}
		} // for each index
	})
}

func (s *VerifyCommand) Gather(files <-chan api.FileInfo, results chan<- api.Result, wg *sync.WaitGroup) {
	makeResult := func(f, source *api.FileInfo, err error) api.Result {
		res := VerifyResult{
			BasicResult: api.BasicResult{
				Finfo: *f,
				Prio:  api.Info,
				Err:   err,
			},
			ifinfo: *f,
		}
		return &res
	}

	api.Gather(files, results, wg, &s.Stats, makeResult, s.RootedReaders, nil)
}

func (s *VerifyCommand) Aggregate(results <-chan api.Result) <-chan api.Result {

	var signatureMismatches uint = 0
	var missingFiles uint = 0
	resultHandler := func(r api.Result, accumResult chan<- api.Result) bool {
		if r.Error() != nil {
			if os.IsNotExist(r.Error()) || os.IsPermission(r.Error()) {
				missingFiles += 1
			}
			accumResult <- r
			return false
		}

		vr := r.(*VerifyResult)

		hasError := false
		vr.Prio = api.Info
		if (len(vr.ifinfo.Sha1) > 0 && bytes.Compare(vr.ifinfo.Sha1, vr.Finfo.Sha1) != 0) ||
			(len(vr.ifinfo.MD5) > 0 && bytes.Compare(vr.ifinfo.MD5, vr.Finfo.MD5) != 0) {
			vr.Err = fmt.Errorf("HASH MISMATCH: %s", vr.Finfo.Path)
			signatureMismatches += 1
			hasError = true
			vr.Prio = api.Error
		} else {
			vr.Msg = fmt.Sprintf("OK: %s", vr.Finfo.Path)
		}
		accumResult <- vr
		return !hasError
	}

	finalizer := func(
		accumResult chan<- api.Result) {
		stats := s.Stats.DeltaString(&s.Stats, s.Stats.Elapsed(), utility.StatsClientSep)

		if signatureMismatches == 0 && missingFiles == 0 {
			accumResult <- &VerifyResult{
				BasicResult: api.BasicResult{
					Msg: fmt.Sprintf(
						"VERIFY OK: None of %d file(s) changed after sealing [%s]",
						s.Stats.MostFiles(),
						stats,
					) + s.Stats.String(),
					Prio: api.Valuable,
				},
			}
		} else {
			s.Stats.ErrCount -= signatureMismatches
			s.Stats.ErrCount -= missingFiles
			accumResult <- &VerifyResult{
				BasicResult: api.BasicResult{
					Msg: fmt.Sprintf(
						"VERIFY FAIL: %d of %d file(s) have changed on disk after sealing, %d are missing [%s]",
						signatureMismatches,
						s.Stats.MostFiles(),
						missingFiles,
						stats,
					) + s.Stats.String(),
					Prio: api.Valuable,
				},
			}
		}
	}

	return api.Aggregate(results, s.Done, resultHandler, finalizer, &s.Stats)
}
