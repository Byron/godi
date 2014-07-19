package godi

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// A struct holding information about a task, including
type FileInfo struct {
	// path to file to handle
	Path string

	// size of file
	Size int64

	// hash of file
	Sha1 []byte

	MD5 []byte
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

type Runner interface {

	// Return maximum amount of processes we can handle.
	// It's also based on our options, and no more than MaxProcs() go routines shouldbe started for Gather
	MaxProcs() uint

	// Intialize required members to deal with controlled reading and writing. numReaders and numWriters
	// can be assumed to be valid
	// Sets the items we are supposed to be working on - must be checked by implementation, as they are
	// very generic in nature
	Init(numReaders, numWriters int, items []string) error

	// Launches a go-routine which fills the returned FileInfo channel
	// Must close FileInfo channel when done
	// Must listen for SIGTERM|SIGINT signals and abort if received
	// May report errrors or information about the progress through generateResult, which must NOT be closed when done. Return nothing
	// if there is nothing to report
	// Must listen on done and return asap
	Generate(done <-chan bool) (files <-chan FileInfo, generateResult <-chan Result)

	// Will be launched as go routine and perform whichever operation on the FileInfo received from input channel
	// Produces one result per input FileInfo and returns it in the given results channel
	// Must listen for SIGTERM|SIGINT signals
	// Use the wait group to mark when done
	// Must listen on done and return asap
	Gather(files <-chan FileInfo, results chan<- Result, wg *sync.WaitGroup, done <-chan bool)

	// Aggregate the result channel and produce whatever you have to produce from the result of the Gather steps
	// When you are done, place a single result instance into accumResult and close the channel
	// You must listen on done to know if the operation was aborted prematurely. This information should be useful
	// for your result.
	Aggregate(results <-chan Result, done <-chan bool) <-chan Result
}

// Runner Init must have been called beforehand as we don't know the values here
// The handlers receive a result of the respective stage and may perform whichever operation
// Returns the last error we received in either generator or aggregation stage
func StartEngine(runner Runner, nprocs uint,
	generateHandler func(Result),
	aggregateHandler func(Result)) (err error) {
	if nprocs > runner.MaxProcs() {
		nprocs = runner.MaxProcs()
	}
	if nprocs == 0 {
		panic("Can't use nprocs = 0 - check implementation of MaxProcs()")
	}

	results := make(chan Result, nprocs)
	done := make(chan bool)

	// assure we close our done channel on signal
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		close(done)
	}()

	files, generateResult := runner.Generate(done)

	wg := sync.WaitGroup{}
	for i := 0; uint(i) < nprocs; i++ {
		wg.Add(1)
		go runner.Gather(files, results, &wg, done)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	accumResult := runner.Aggregate(results, done)

	// Return true if we should break the loop
	resHandler := func(handler func(Result), res Result) bool {
		if res == nil {
			// channel closed, have to get out
			return true
		}
		handler(res)
		return false
	} // end resHandler

infinity:
	for {
		select {
		case r := <-generateResult:
			{
				if resHandler(generateHandler, r) {
					break infinity
				}
				err = r.Error()
			}
		case r := <-accumResult:
			{
				if resHandler(aggregateHandler, r) {
					break infinity
				}
				err = r.Error()
			}
		} // select
	} // infinty

	return
}
