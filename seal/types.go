package seal

import (
	"github.com/Byron/godi/utility"
)

const (
	IndexBaseName = "godi"
	Name          = "seal"
)

// A type representing all arguments required to drive a Seal operation
type SealCommand struct {

	// One or more trees to seal
	// Exported just for test-cases - too lazy to make it a read-only copy through accessor
	Trees []string

	// parallel reader
	pCtrl utility.ReadChannelController
}

// NewCommand returns an initialized seal command
func NewCommand(trees []string, nReaders int) (c SealCommand, err error) {
	err = c.Init(nReaders, 0, trees)
	return
}
