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

type indexWriterResult struct {
	path string // path to the seal file
	err  error  // possible error during the seal operation
}

// Some information we store per root of files we seal
type aggregationTreeInfo struct {
	// Paths to files we have written so far - only used in sealed-copy mode
	// TODO(st): don't track these files in memory, but re-read them from the written seal file !
	// That way, we don't rely on any limited resource except for disk space
	writtenFiles []string

	// A channel to send file-infos to the attached seal serializer. Close it to finish the seal operation
	sealFInfos chan<- api.FileInfo

	// Contains the error code of the seal operation for the tree we are associated with, and the produced seal file
	// Will only yield a result one, and be closed afterwards
	sealResult <-chan indexWriterResult

	// if true, the entire tree is considered faulty, and further results won't be recorded or accepted
	hasError bool
}

// Helper to sort by longest path, descending
type byLongestPathDescending []string

func (a byLongestPathDescending) Len() int           { return len(a) }
func (a byLongestPathDescending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byLongestPathDescending) Less(i, j int) bool { return len(a[i]) > len(a[j]) }

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
func NewCommand(trees []string, nReaders, nWriters int) (*SealCommand, error) {
	c := SealCommand{}
	if nWriters == 0 {
		c.mode = modeSeal
	} else {
		c.mode = modeCopy
	}
	err := c.Init(nReaders, nWriters, trees, api.Info)
	return &c, err
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
				Prio:  api.Info,
				Err:   err,
			},
			source: s,
		}
		return &res
	}

	api.Gather(files, results, wg, s.Statistics(), makeResult, s.RootedReaders, s.rootedWriters)
}
