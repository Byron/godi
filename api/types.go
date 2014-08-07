package api

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/Byron/godi/io"
)

const IndexBaseName = "godi"

const (
	filterModeSymlinks int8 = iota
	filterModeHidden
	filterModeSeals
	filterModeVolatile
	filterModeFnMatch
)

// Must be kept in sync with indexPath() generator
var reIsIndexPath = regexp.MustCompile(fmt.Sprintf(`%s_\d{4}-\d{2}-\d{2}_\d{2}\d{2}\d{2}\..*`, IndexBaseName))

// return a path to an index file residing at tree
func IndexPath(tree string, extension string) string {
	n := time.Now()
	return filepath.Join(tree, fmt.Sprintf("%s_%04d-%02d-%02d_%02d%02d%02d.%s",
		IndexBaseName,
		n.Year(),
		n.Month(),
		n.Day(),
		n.Hour(),
		n.Minute(),
		n.Second(),
		extension))
}

// A utility to encapsulate a file-filter
// These exist in special modes to filter entire classes of files, and as FNMatch compatible string
type FileFilter struct {
	fnFilter string // fnmatch compatible string, used if mode is filterFnMatch
	kind     int8   // The kind of files we apply to
}

func (f FileFilter) String() string {
	switch f.kind {
	case filterModeSymlinks:
		return "SYMLINK"
	case filterModeHidden:
		return "HIDDEN"
	case filterModeSeals:
		return "SEALS"
	case filterModeVolatile:
		return "VOLATILE"
	case filterModeFnMatch:
		return f.fnFilter
	default:
		panic("Not implemented")
	}
}

func (f *FileFilter) Matches(name string, mode os.FileMode) bool {
	switch f.kind {
	case filterModeSymlinks:
		return mode&os.ModeSymlink == os.ModeSymlink
	case filterModeHidden:
		if fr, _ := utf8.DecodeRuneInString(name); fr == '.' {
			return true
		}
	case filterModeSeals:
		return reIsIndexPath.Match([]byte(name))
	case filterModeVolatile:
		return ((mode&os.ModeSymlink != os.ModeSymlink) && !mode.IsRegular()) || name == ".DS_Store"
	case filterModeFnMatch:
		{
			// We assume the patten was already checked for correctness
			res, _ := filepath.Match(f.fnFilter, name)
			return res
		}
	default:
		panic("unknown kind")
	} // kind switch

	// select may fall through, so defeault is no match
	return false
}

// Return a new FileFilter matching the given string.
// Every string which is not a special kind of filter will be interpreted as fnmatch filter. Err is returned if
// the glob is invalid
func ParseFileFilter(name string) (FileFilter, error) {
	for _, f := range [...]FileFilter{FilterSymlinks, FilterHidden, FilterSeals, FilterVolatile} {
		if f.String() == name {
			return f, nil
		}
	}

	if _, err := filepath.Match(name, "empty"); err != nil {
		return FileFilter{}, err
	}

	return FileFilter{
		fnFilter: name,
		kind:     filterModeFnMatch,
	}, nil
}

var (
	FilterSymlinks = FileFilter{kind: filterModeSymlinks}
	FilterHidden   = FileFilter{kind: filterModeHidden}
	FilterSeals    = FileFilter{kind: filterModeSeals}
	FilterVolatile = FileFilter{kind: filterModeVolatile}
)

// A struct holding information about a task, including
type FileInfo struct {

	// path to file to handle
	Path string

	// Path relative to the directory it was found in
	RelaPath string

	// Provides information about the type of the file
	Mode os.FileMode

	// size of file
	Size int64

	// hashes of file
	Sha1 []byte
	MD5  []byte
}

// Compute the root of this file - it is the top-level directory used to specify all files to process
func (f *FileInfo) Root() string {
	return f.Path[:len(f.Path)-len(f.RelaPath)-1]
}

type Importance uint8

const (
	Progress Importance = iota
	Info
	Warn
	Error
	Valuable
	LogDisabled
)

// MayLog returns true if the given priority may be logged as seen from our log-level.
// Results may always be logged
func (p Importance) MayLog(op Importance) bool {
	if p == LogDisabled {
		return false
	}

	if op == Valuable {
		return true
	}
	return op >= p
}

func (p Importance) String() string {
	switch {
	case p == Progress:
		return "progress"
	case p == Info:
		return "info"
	case p == Warn:
		return "warn"
	case p == Error:
		return "error"
	case p == Valuable:
		return "result"
	case p == LogDisabled:
		return "off"
	default:
		panic("Unknown log level")
	}
}

// Parse a priority from the given string. error will be set if this fails
func ParseImportance(p string) (Importance, error) {
	for _, t := range [...]Importance{Progress, Info, Warn, Error, Valuable, LogDisabled} {
		if t.String() == p {
			return t, nil
		}
	}
	return LogDisabled, fmt.Errorf("Unknown verbosity level: '%s'", p)
}

type Result interface {
	// Return a string indicating the result, which can can also state an error
	// The priority show the kind of result messgae, allowing you to filter them effectively
	Info() (string, Importance)

	// Return an error instance indicating what exactly when wrong
	Error() error

	// Return the FileInformation we represent
	FileInformation() *FileInfo
}

// Implements information about any operation
// It's the minimum we need to work
type BasicResult struct {
	Finfo FileInfo
	Msg   string
	Err   error
	Prio  Importance
}

func (s *BasicResult) Info() (string, Importance) {
	if s.Err != nil {
		msg := s.Err.Error()
		if len(s.Msg) > 0 {
			msg = s.Msg
		}
		return msg, Error
	}
	return s.Msg, s.Prio
}

func (s *BasicResult) Error() error {
	return s.Err
}

func (s *BasicResult) FileInformation() *FileInfo {
	return &s.Finfo
}

// A partial implementation of a runner, which can be shared between the various commands
type BasicRunner struct {
	// Items we work on
	Items []string
	// A map of readers which maps from a root to the reader to use to read files that share the same root
	RootedReaders io.RootedReadControllers
	// A channel to let everyone know we should finish as soon as possible - this is done by closing the channel
	Done chan bool

	// our statistics instance
	Stats Stats

	// The maximum log-level. We just keep this value here because the cli makes a difference between CHECK and RUN !
	// This member shouldn't be needed as logging is not done by the runner anyway - it's all done by result handlers.
	// Only they are concerned, which is a function of the CLI entirely
	// TODO(st) Fork codegangsa/CLI and make the fix, use the fork from that point on ... .
	Level   Importance
	Filters []FileFilter
}

func (b *BasicRunner) LogLevel() Importance {
	return b.Level
}

func (b *BasicRunner) Statistics() *Stats {
	return &b.Stats
}

func (b *BasicRunner) NumChannels() int {
	if len(b.RootedReaders) == 0 {
		panic("NumChannels called before InitBasicRunner()")
	}
	return b.RootedReaders.Streams()
}

// Initialize our Readers and items with the given information, including our cannel
func (b *BasicRunner) InitBasicRunner(numReaders int, items []string, maxLogLevel Importance, filters []FileFilter) {
	b.Items = items
	b.Done = make(chan bool)
	b.RootedReaders = io.NewDeviceReadControllers(numReaders, items, &b.Stats.Stats, b.Done)
	if len(b.RootedReaders) == 0 {
		panic("Didn't manage to build readers from input items")
	}
	b.Level = maxLogLevel
	b.Filters = filters
}

func (b *BasicRunner) CancelChannel() chan bool {
	if b.Done == nil {
		panic("CancelChannel( called before InitBasicRunner()")
	}
	return b.Done
}

// An interface to help implementing types which read one ore more data streams, run an operation on them
// whose result is aggregated and provided to the caller.
type Runner interface {

	// Intialize required members to deal with controlled reading and writing. numReaders and numWriters
	// can be assumed to be valid
	// Sets the items we are supposed to be working on - must be checked by implementation, as they are
	// very generic in nature
	Init(numReaders, numWriters int, items []string, maxLogLevel Importance, filters []FileFilter) error

	// Return the amount of io-channels the runner may be using in parallel per device
	NumChannels() int

	// Return the minimum allowed level for logging
	// TODO(st): get rid of this method !
	LogLevel() Importance

	// Statistics returns the commands shared statistics structure
	Statistics() *Stats

	// CancelChannel returns the channel to close when the operation should stop prematurely
	// NOTE: Only valid after Init was called, and it's an error to call it beforehand
	CancelChannel() chan bool

	// Launches generators, gatherers and an aggregator, setting up their connections to fit.
	// Must close FileInfo channel when done
	// Must listen for SIGTERM|SIGINT signals and abort if received
	// May report errrors or information about the progress through generateResult, which must NOT be closed when done. Return nothing
	// if there is nothing to report
	// Must listen on done and return asap
	Generate() (aggregationResult <-chan Result)

	// Will be launched as go routine and perform whichever operation on the FileInfo received from input channel
	// Produces one result per input FileInfo and returns it in the given results channel
	// Must listen for SIGTERM|SIGINT signals
	// Use the wait group to mark when done, which is when the results channel need to be closed.
	// Must listen on done and return asap
	Gather(rctrl *io.ReadChannelController, files <-chan FileInfo, results chan<- Result)

	// Aggregate the result channel and produce whatever you have to produce from the result of the Gather steps
	// When you are done, place a single result instance into accumResult and close the channel
	// You must listen on done to know if the operation was aborted prematurely. This information should be useful
	// for your result.
	Aggregate(results <-chan Result) <-chan Result
}

// Runner Init must have been called beforehand as we don't know the values here
// The handlers receive a result of the respective stage and may perform whichever operation. It returns true if
// it used the result it obtained, false otherwise.
// Returns the last error we received in either generator or aggregation stage
func StartEngine(runner Runner,
	aggregateHandler func(Result) bool) (err error) {

	runner.Statistics().StartedAt = time.Now()

	done := runner.CancelChannel()

	// assure we close our done channel on signal
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		close(done)
	}()

	accumResult := runner.Generate()

	mkErrPicker := func(handler func(r Result) bool) func(r Result) bool {
		return func(r Result) bool {
			if r.Error() != nil {
				err = r.Error()
			}
			return handler(r)
		}
	}
	aggregateHandler = mkErrPicker(aggregateHandler)

	// Let's not hot-loop over anything, instead just process asynchronously
	rwg := sync.WaitGroup{}
	rwg.Add(1)
	go func() {
		for r := range accumResult {
			aggregateHandler(r)
		}
		rwg.Done()
	}()
	rwg.Wait()

	return
}
