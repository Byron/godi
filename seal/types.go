package seal

import (
	"math"

	"github.com/Byron/godi/api"
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

// Implements information about a seal operation
type SealResult struct {
	finfo *godi.FileInfo
	msg   string
	err   error
	prio  godi.Priority
}

func (s *SealResult) Info() (string, godi.Priority) {
	if s.err != nil {
		return s.err.Error(), godi.Error
	}
	return s.msg, s.prio
}

func (s *SealResult) Error() error {
	return s.err
}

func (s *SealResult) FileInformation() *godi.FileInfo {
	return s.finfo
}

func (s *SealCommand) MaxProcs() uint {
	return uint(math.MaxUint32)
}

// REVIEW:
func NewCommand(trees []string, nReaders int) (c SealCommand, err error) {
	err = c.Init(nReaders, 0, trees)
	return
}
