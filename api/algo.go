package api

import (
	"fmt"
	"time"
)

// Generate does all boilerplate required to be a valid generator
func Generate(
	generate func(chan<- FileInfo, chan<- Result)) (<-chan FileInfo, <-chan Result) {
	files := make(chan FileInfo)
	results := make(chan Result)

	go func() {
		defer close(files)
		defer close(results)
		generate(files, results)
	}()

	return files, results
}

type AggregateFinalizerState struct {
	ErrCount     uint
	WasCancelled bool
	Elapsed      time.Duration
}

// String generates a string with performance information
func (a *AggregateFinalizerState) String() (out string) {
	if a.ErrCount > 0 {
		out += fmt.Sprintf("%d errors", a.ErrCount)
	}

	if a.WasCancelled {
		if len(out) == 0 {
			out = "cancelled"
		} else {
			out += ", cancelled"
		}
	}

	if len(out) != 0 {
		out = "(" + out + ")"
	}

	return
}

// Aggregate is a general purpose implementation to gather fileInfo results
func Aggregate(results <-chan Result, done <-chan bool,
	resultHandler func(Result, chan<- Result) bool,
	finalizer func(chan<- Result, *AggregateFinalizerState)) <-chan Result {
	accumResult := make(chan Result)

	go func() {
		defer close(accumResult)

		s := AggregateFinalizerState{}
		st := time.Now()

		// ACCUMULATE PATHS INFO
		/////////////////////////
		for r := range results {
			// Be sure we take note of cancellation.
			// If this happens, soon our results will be drained and we leave naturally
			select {
			case <-done:
				s.WasCancelled = true
				// fallthrough doesn't work in selects :(
			default:
				// we are just checking, but don't want to loose a result.
				// If the following code would be here, the last result
				// would not be pulled before closing
			}

			if !resultHandler(r, accumResult) {
				s.ErrCount += 1
			}
		} // range results
		s.Elapsed = time.Now().Sub(st)

		finalizer(accumResult, &s)
	}()

	return accumResult
}
