package godi

import (
	"time"
)

type AggregateFinalizerState struct {
	WasCancelled        bool
	FileCount, ErrCount uint
	SizeBytes           uint64
	Elapsed             float64
}

// Aggregate is a general purpose implementation to gather fileInfo results
func Aggregate(results <-chan Result, done <-chan bool,
	resultHandler func(Result, chan<- Result) bool,
	finalizer func(chan<- Result, AggregateFinalizerState)) <-chan Result {
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

		finalizer(accumResult, s)
	}()

	return accumResult
}
