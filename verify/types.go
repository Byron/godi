// Implements verification of seal files previously written with seal command
package verify

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/io"
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

// Keeps some information on a per-tree level
type treeInfo struct {
	signatureMismatches, missingFiles, numFiles uint
	sealBroken                                  bool
}

// NewCommand returns an initialized verify command
func NewCommand(indices []string, nReaders int) (*VerifyCommand, error) {
	c := VerifyCommand{}
	return &c, c.Init(nReaders, 0, indices, api.Info, nil)
}

func (s *VerifyCommand) Generate() <-chan api.Result {
	return api.Generate(s.RootedReaders, s,
		func(trees []string, files chan<- api.FileInfo, results chan<- api.Result) {
			for _, index := range s.Items {
				// Only work in indices that are assigned to us
				found := false
				for _, tree := range trees {
					if filepath.Dir(index) == tree {
						found = true
						break
					}
				}
				if !found {
					continue
				}

				c := codec.NewByPath(index)
				if c == nil {
					panic("Should have a codec here - this was checked before")
				}

				fd, err := os.Open(index)
				if err != nil {
					results <- &VerifyResult{
						BasicResult: api.BasicResult{
							Err: &codec.DecodeError{Msg: err.Error()},
							Finfo: api.FileInfo{
								Path:     index,
								RelaPath: filepath.Base(index),
							},
						},
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
						BasicResult: api.BasicResult{
							Err: err,
							Finfo: api.FileInfo{
								Path:     index,
								RelaPath: filepath.Base(index),
							},
						},
					}
					continue
				}
			} // for each index
		})
}

func (s *VerifyCommand) Gather(rctrl *io.ReadChannelController, files <-chan api.FileInfo, results chan<- api.Result) {
	makeResult := func(f, source *api.FileInfo, err error) api.Result {
		res := VerifyResult{
			BasicResult: api.BasicResult{
				Finfo: *f,
				Prio:  api.Info,
				Err:   err,
			},
		}
		if source != nil {
			// We just copy the originally retrieved file-info
			res.ifinfo = *source
		} else {
			panic("Should have received the original fileInfo, to obtain the sealed values for hashes")
		}
		return &res
	}

	api.Gather(files, results, &s.Stats, makeResult, rctrl, nil)
}

func (s *VerifyCommand) Aggregate(results <-chan api.Result) <-chan api.Result {
	treeInfoMap := make(map[string]*treeInfo)

	resultHandler := func(r api.Result, accumResult chan<- api.Result) bool {
		vr := r.(*VerifyResult)

		ti, hasTi := treeInfoMap[vr.Finfo.Root()]
		if !hasTi {
			ti = &treeInfo{}
			treeInfoMap[vr.Finfo.Root()] = ti
		}

		if vr.Err != nil {
			if os.IsNotExist(vr.Err) || os.IsPermission(vr.Err) {
				ti.missingFiles += 1
				vr.Msg = fmt.Sprintf("MISSING %s: %s", SymbolMismatch, vr.Finfo.Path)
				accumResult <- vr
				return false
			} else if serr, isFileSizeType := vr.Err.(*api.FileSizeMismatch); isFileSizeType {
				// The file-size changed, thus the hashes will be different to. Say it accordingly.
				ti.signatureMismatches += 1
				ti.numFiles += 1
				vr.Msg = fmt.Sprintf("SIZE %s: %s sealed with size %dB, got size %dB", SymbolMismatch, serr.Path, serr.Want, serr.Got)
				accumResult <- vr
				return false
			} else if _, isSealSigMismatch := vr.Err.(*codec.SignatureMismatchError); isSealSigMismatch {
				ti.sealBroken = true
				vr.Msg = fmt.Sprintf("SEAL %s: '%s' was modified after sealing or is corrupted - don't trust the verify results", SymbolMismatch, vr.Finfo.Path)
				accumResult <- vr
				return false
			} else if _, isDecodeErr := vr.Err.(*codec.DecodeError); isDecodeErr {
				ti.sealBroken = true
				vr.Msg = fmt.Sprintf("SEAL %s", "Failed to decode seal at '%s' with error '%s' - verify results can't be trusted", SymbolFail, vr.Finfo.Path, vr.Err.Error())
				accumResult <- vr
				return false
			} else {
				// It's some other error - just push it forward
				accumResult <- vr
				return false
			}
		}

		hasError := false
		vr.Prio = api.Info
		// From here on, it must be a file with no obvious error
		ti.numFiles += 1
		if (len(vr.ifinfo.Sha1) > 0 && bytes.Compare(vr.ifinfo.Sha1, vr.Finfo.Sha1) != 0) ||
			(len(vr.ifinfo.MD5) > 0 && bytes.Compare(vr.ifinfo.MD5, vr.Finfo.MD5) != 0) {
			vr.Msg = fmt.Sprintf("HASH %s: %s flipped at least one bit", SymbolMismatch, vr.Finfo.Path)
			vr.Err = &api.FileHashMismatch{Path: vr.Finfo.Path}
			ti.signatureMismatches += 1
			hasError = true
			vr.Prio = api.Error
		} else {
			vr.Msg = fmt.Sprintf("%s: %s", SymbolOK, vr.Finfo.Path)
		}
		accumResult <- vr
		return !hasError
	}

	finalizer := func(
		accumResult chan<- api.Result) {

		count := 0
		stats := ""
		for _, ti := range treeInfoMap {
			count += 1

			s.Stats.ErrCount -= ti.signatureMismatches
			s.Stats.ErrCount -= ti.missingFiles

			// the last result we produce has the final statistics
			if count == len(treeInfoMap) {
				stats = fmt.Sprintf(" [%s]%s",
					s.Stats.DeltaString(&s.Stats, s.Stats.Elapsed(), io.StatsClientSep),
					s.Stats.String(),
				)
			}

			if ti.signatureMismatches == 0 && ti.missingFiles == 0 && !ti.sealBroken {
				accumResult <- &VerifyResult{
					BasicResult: api.BasicResult{
						Msg: fmt.Sprintf(
							"VERIFY %s: None of %d file(s) changed after sealing%s",
							SymbolSuccess,
							ti.numFiles,
							stats,
						),
						Prio: api.Valuable,
					},
				}
			} else {
				suffix := ""
				if ti.missingFiles > 0 {
					suffix = fmt.Sprintf(", with %d missing", ti.missingFiles)
				}
				accumResult <- &VerifyResult{
					BasicResult: api.BasicResult{
						Msg: fmt.Sprintf(
							"VERIFY %s: %d of %d file(s) have changed on disk after sealing%s%s",
							SymbolFail,
							ti.signatureMismatches,
							ti.numFiles,
							suffix,
							stats,
						),
						Prio: api.Valuable,
					},
				}
			}
		} // end for each treeInfo
	} // end finalizer

	return api.Aggregate(results, s.Done, resultHandler, finalizer, &s.Stats)
}
