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

	// One or more trees to seal
	// Exported just for test-cases - too lazy to make it a read-only copy through accessor
	SourceTrees []string

	// The type of seal operation we are supposed to perform
	mode string

	// parallel reader
	pCtrl utility.ReadChannelController

	// A map of writers - there may just be one writer per device.
	// Map may be unset if we are not in write mode
	pWriters []utility.RootedWriteController
}

// A result which is also able to hold information about the source of a file
type SealResult struct {
	godi.BasicResult
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

func (s *SealCommand) Gather(files <-chan godi.FileInfo, results chan<- godi.Result, wg *sync.WaitGroup, done <-chan bool) {
	makeResult := func(f, source *godi.FileInfo, err error) godi.Result {
		s := ""
		if source != nil {
			s = source.Path
		}
		res := SealResult{
			BasicResult: godi.BasicResult{
				Finfo: f,
				Prio:  godi.Progress,
				Err:   err,
			},
			source: s,
		}
		return &res
	}

	godi.Gather(files, results, wg, done, makeResult, &s.pCtrl, s.pWriters)
}
