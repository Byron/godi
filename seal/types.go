package seal

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/io"
)

const (
	Name = "seal"
	Sep  = "--"

	ModeSeal = Name
	ModeCopy = "sealed-copy"
)

var (
	usage = fmt.Sprintf(`Please specify sealed copies like so: source/ %s destination/
	%s can be omitted if there is only one source and one destination.`, Sep, Sep)
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

	// A possible result we might have gotten due to an early seal error
	lsr indexWriterResult // lastSealResult

	// if true, the entire tree is considered faulty, and further results won't be recorded or accepted
	hasError bool
}

// Helper to sort by longest path, descending
type byLongestPathDescending []string

func (a byLongestPathDescending) Len() int           { return len(a) }
func (a byLongestPathDescending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byLongestPathDescending) Less(i, j int) bool { return len(a[i]) > len(a[j]) }

// A type representing all arguments required to drive a Seal operation
type Command struct {
	api.BasicRunner

	// The type of seal operation we are supposed to perform
	Mode string

	// If set, we are supposed to run in verify mode
	Verify bool

	// The name of the seal format to use
	Format string

	// A map of writers - there may just be one writer per device.
	// Map may be unset if we are not in write mode
	rootedWriters []io.RootedWriteController
}

// A result which is also able to hold information about the source of a file
type SealResult struct {
	api.BasicResult
	// source of a copy operation, may be unset
	source string
}

// Returns true if this result was sent from a generator. The latter sends the root as Path, but doesn't set a RelaPath
func (s *SealResult) FromGenerator() bool {
	return len(s.Finfo.RelaPath) == 0
}

// NewCommand returns an initialized seal command
func NewCommand(trees []string, nReaders, nWriters int) (*Command, error) {
	c := Command{}
	if nWriters == 0 {
		c.Mode = ModeSeal
	} else {
		c.Mode = ModeCopy
	}
	err := c.Init(nReaders, nWriters, trees, api.Info, []api.FileFilter{api.FilterSeals})
	return &c, err
}

func (s *Command) Gather(rctrl *io.ReadChannelController, files <-chan api.FileInfo, results chan<- api.Result) {
	makeResult := func(f, source *api.FileInfo, err error) api.Result {
		s := ""
		if source != nil && source.Path != f.Path {
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

	api.Gather(files, results, s.Statistics(), makeResult, rctrl, s.rootedWriters)
}

func (s *Command) Init(numReaders, numWriters int, items []string, maxLogLevel api.Priority, filters []api.FileFilter) (err error) {

	if len(s.Format) == 0 {
		s.Format = codec.GobName
	}

	if s.Mode == ModeSeal {
		if len(items) == 0 {
			return errors.New("Please provide at least one source directory to work on")
		}
		items, err = api.ParseSources(items, true)
		if err != nil {
			return
		}
		s.InitBasicRunner(numReaders, items, maxLogLevel, filters)
	} else if s.Mode == ModeCopy {
		finishSetup := func(sources, dtrees []string) error {
			// Make sure we don't copy onto ourselves
			for _, stree := range sources {
				for _, dtree := range dtrees {
					if strings.HasPrefix(dtree+string(os.PathSeparator), stree) {
						return fmt.Errorf("Cannot copy '%s' into it's own subdirectory or itself at '%s'", stree, dtree)
					}
				}
			}
			s.InitBasicRunner(numReaders, sources, maxLogLevel, filters)

			// build the device map with all writer destinations
			dm := io.DeviceMap(dtrees)

			// Finally, put all actual values into our list to have a deterministic iteration order.
			// After all, we don't really care about the device from this point on
			s.rootedWriters = make([]io.RootedWriteController, len(dm))
			for did, trees := range dm {
				// each device as so and so many destinations. Each destination uses the same write controller
				s.rootedWriters[did] = io.RootedWriteController{
					Trees: trees,
					Ctrl:  io.NewWriteChannelController(numWriters, numWriters*len(trees), &s.Stats.Stats),
				}
			} // for each tree set in deviceMap
			return nil
		} // end helper

		// Parses [src, ...] -- [dst, ...]
		err = errors.New(usage)
		if len(items) < 2 {
			return
		}

		for i, item := range items {
			if item == Sep {
				if i == 0 {
					return
				}
				if i == len(items)-1 {
					return
				}
				sources, e := api.ParseSources(items[:i], true)
				if e != nil {
					return e
				}

				dtrees, e := api.ParseSources(items[i+1:], false)
				if e != nil {
					return e
				}

				return finishSetup(sources, dtrees)
			}
		} // for each item

		// So there is no separator, maybe it's source and destination ?
		if len(items) == 2 {
			return finishSetup(items[:1], items[1:])
		}

		// source-destination separator not found - prints usage
		return
	} else {
		panic(fmt.Sprintf("Unsupported mode: %s", s.Mode))
	}
	return
}
