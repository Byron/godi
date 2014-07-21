package godi

import (
	"fmt"
	"time"
)

// Generate does all boilerplate required to be a valid generator
func Generate(
	done <-chan bool,
	generate func(chan<- FileInfo, chan<- Result)) (<-chan FileInfo, <-chan Result) {
	files := make(chan FileInfo)
	results := make(chan Result)

	go func() {
		defer close(files)
		generate(files, results)
	}()

	return files, results
}

type AggregateFinalizerState struct {
	WasCancelled        bool
	FileCount, ErrCount uint
	SizeBytes           uint64
	Elapsed             float64
}

// String generates a string with performance information
func (a *AggregateFinalizerState) String() string {
	sizeMB := float32(a.SizeBytes) / (1024.0 * 1024.0)
	return fmt.Sprintf(
		"Processed %#vMB in %vs (%#v MB/s, %d errors, cancelled=%v",
		sizeMB,
		a.Elapsed,
		float64(sizeMB)/a.Elapsed,
		a.ErrCount,
		a.WasCancelled,
	)
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
			default:
				{
					if r.Error() != nil {
						s.ErrCount += 1
						accumResult <- r
						continue
					}

					if !resultHandler(r, accumResult) {
						s.ErrCount += 1
					} else {
						s.FileCount += 1
						s.SizeBytes += uint64(r.FileInformation().Size)
					}
				} // default
			} // select
		} // range results
		s.Elapsed = time.Now().Sub(st).Seconds()

		finalizer(accumResult, &s)
	}()

	return accumResult
}
