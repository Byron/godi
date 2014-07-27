package api

import (
	"github.com/Byron/godi/utility"
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

// Aggregate is a general purpose implementation to gather fileInfo results
func Aggregate(results <-chan Result, done <-chan bool,
	resultHandler func(Result, chan<- Result) bool,
	finalizer func(chan<- Result),
	stats *utility.Stats) <-chan Result {
	accumResult := make(chan Result)

	go func() {
		defer close(accumResult)

		// ACCUMULATE PATHS INFO
		/////////////////////////
		for r := range results {
			// Be sure we take note of cancellation.
			// If this happens, soon our results will be drained and we leave naturally
			select {
			case <-done:
				stats.WasCancelled = true
				// fallthrough doesn't work in selects :(
			default:
				// we are just checking, but don't want to loose a result.
				// If the following code would be here, the last result
				// would not be pulled before closing
			}

			if !resultHandler(r, accumResult) {
				stats.ErrCount += 1
			}
		} // range results

		finalizer(accumResult)
	}()

	return accumResult
}
