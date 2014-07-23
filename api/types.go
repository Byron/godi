package godi

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Byron/godi/utility"
)

// A struct holding information about a task, including
type FileInfo struct {

	// path to file to handle
	Path string

	// Path relative to the directory it was found in
	RelaPath string

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

type Priority uint8

const (
	Progress Priority = iota
	Info
	Warn
	Error
)

type Result interface {
	// Return a string indicating the result, which can can also state an error
	// The priority show the kind of result messgae, allowing you to filter them effectively
	Info() (string, Priority)

	// Return an error instance indicating what exactly when wrong
	Error() error

	// Return the FileInformation we represent
	FileInformation() *FileInfo
}

// Implements information about any operation
// It's the minimum we need to work
type BasicResult struct {
	Finfo *FileInfo
	Msg   string
	Err   error
	Prio  Priority
}

func (s *BasicResult) Info() (string, Priority) {
	if s.Err != nil {
		return s.Err.Error(), Error
	}
	return s.Msg, s.Prio
}

func (s *BasicResult) Error() error {
	return s.Err
}

func (s *BasicResult) FileInformation() *FileInfo {
	return s.Finfo
}

// A partial implementation of a runner, which can be shared between the various commands
type BasicRunner struct {
	// Items we work on
	Items []string
	// A map of readers which maps from a root to the reader to use to read files that share the same root
	RootedReaders map[string]*utility.ReadChannelController
	// A channel to let everyone know we should finish as soon as possible - this is done by closing the channel
	Done chan bool
}

func (b *BasicRunner) NumChannels() int {
	if len(b.RootedReaders) == 0 {
		panic("NumChannels called before InitBasicRunner()")
	}
	return utility.ReadChannelDeviceMapStreams(b.RootedReaders)
}

// Initialize our Readers and items with the given information, including our cannel
func (b *BasicRunner) InitBasicRunner(numReaders int, items []string) {
	b.Items = items
	b.Done = make(chan bool)
	b.RootedReaders = utility.NewReadChannelDeviceMap(numReaders, items, b.Done)
	if len(b.RootedReaders) == 0 {
		panic("Didn't manage to build readers from input items")
	}
}

func (b *BasicRunner) CancelChannel() chan bool {
	if b.Done == nil {
		panic("CancelChannel( called before InitBasicRunner()")
	}
	return b.Done
}

type Runner interface {

	// Intialize required members to deal with controlled reading and writing. numReaders and numWriters
	// can be assumed to be valid
	// Sets the items we are supposed to be working on - must be checked by implementation, as they are
	// very generic in nature
	Init(numReaders, numWriters int, items []string) error

	// Return the amount of io-channels the runner may be using in parallel per device
	NumChannels() int

	// CancelChannel returns the channel to close when the operation should stop prematurely
	// NOTE: Only valid after Init was called, and it's an error to call it beforehand
	CancelChannel() chan bool

	// Launches a go-routine which fills the returned FileInfo channel
	// Must close FileInfo channel when done
	// Must listen for SIGTERM|SIGINT signals and abort if received
	// May report errrors or information about the progress through generateResult, which must NOT be closed when done. Return nothing
	// if there is nothing to report
	// Must listen on done and return asap
	Generate() (files <-chan FileInfo, generateResult <-chan Result)

	// Will be launched as go routine and perform whichever operation on the FileInfo received from input channel
	// Produces one result per input FileInfo and returns it in the given results channel
	// Must listen for SIGTERM|SIGINT signals
	// Use the wait group to mark when done
	// Must listen on done and return asap
	Gather(files <-chan FileInfo, results chan<- Result, wg *sync.WaitGroup)

	// Aggregate the result channel and produce whatever you have to produce from the result of the Gather steps
	// When you are done, place a single result instance into accumResult and close the channel
	// You must listen on done to know if the operation was aborted prematurely. This information should be useful
	// for your result.
	Aggregate(results <-chan Result) <-chan Result
}

// Runner Init must have been called beforehand as we don't know the values here
// The handlers receive a result of the respective stage and may perform whichever operation
// Returns the last error we received in either generator or aggregation stage
func StartEngine(runner Runner,
	generateHandler func(Result),
	aggregateHandler func(Result)) (err error) {

	nprocs := runner.NumChannels()
	results := make(chan Result, nprocs)
	done := runner.CancelChannel()

	// assure we close our done channel on signal
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		close(done)
	}()

	files, generateResult := runner.Generate()

	wg := sync.WaitGroup{}
	for i := 0; i < nprocs; i++ {
		wg.Add(1)
		go runner.Gather(files, results, &wg)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	accumResult := runner.Aggregate(results)

	mkErrPicker := func(handler func(r Result)) func(r Result) {
		return func(r Result) {
			if r.Error() != nil {
				err = r.Error()
			}
			handler(r)
		}
	}
	generateHandler = mkErrPicker(generateHandler)
	aggregateHandler = mkErrPicker(aggregateHandler)

	// Let's not hot-loop over anything, instead just process asynchronously
	rwg := sync.WaitGroup{}
	rwg.Add(1)
	go func() {
		for r := range generateResult {
			generateHandler(r)
		}
		rwg.Done()
	}()
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
