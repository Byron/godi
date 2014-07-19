package seal

import (
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
	Trees []string

	// Amount of readers to use
	nReaders int

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

// REVIEW:
func NewCommand(trees []string, nReaders int) SealCommand {
	c := SealCommand{}
	c.Trees = trees
	c.nReaders = nReaders
	return c
}
