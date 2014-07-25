package seal

import (
	"sync"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/utility"
)

const (
	IndexBaseName = "godi"
	Name          = "seal"

	modeSeal = Name
	modeCopy = "sealed-copy"
)

// A type representing all arguments required to drive a Seal operation
type SealCommand struct {
	api.BasicRunner

	// The type of seal operation we are supposed to perform
	mode string

	// If set, we are supposed to run in verify mode
	verify bool

	// A map of writers - there may just be one writer per device.
	// Map may be unset if we are not in write mode
	rootedWriters []utility.RootedWriteController
}

// A result which is also able to hold information about the source of a file
type SealResult struct {
	api.BasicResult
	// source of a copy operation, may be unset
	source string
}

// NewCommand returns an initialized seal command
func NewCommand(trees []string, nReaders, nWriters int) (c SealCommand, err error) {
	if nWriters == 0 {
		c.mode = modeSeal
	} else {
		c.mode = modeCopy
	}
	err = c.Init(nReaders, nWriters, trees)
	return
}

func (s *SealCommand) Gather(files <-chan api.FileInfo, results chan<- api.Result, wg *sync.WaitGroup) {
	makeResult := func(f, source *api.FileInfo, err error) api.Result {
		s := ""
		if source != nil {
			s = source.Path
		}
		res := SealResult{
			BasicResult: api.BasicResult{
				Finfo: *f,
				Prio:  api.Progress,
				Err:   err,
			},
			source: s,
		}
		return &res
	}

	api.Gather(files, results, wg, makeResult, s.RootedReaders, s.rootedWriters)
}
