package api

import (
	"sync"

	"github.com/Byron/godi/utility"
)

// Generate does all boilerplate required to be a valid generator
// Will produce as many generators as there are devices, each is handed a list of trees to handle
func Generate(rctrls []utility.RootedReadController,
	runner Runner,
	generate func([]string, chan<- FileInfo, chan<- Result)) (<-chan Result, <-chan Result) {

	genResults := make(chan Result)
	gatherToAgg := make(chan Result, runner.NumChannels())

	genwg := sync.WaitGroup{} // wait group for generators
	gatwg := sync.WaitGroup{} // wait group for gatherers

	// Spawn generators - each one has num-streams gatherers
	for _, rctrl := range rctrls {
		files := make(chan FileInfo)
		genwg.Add(1)
		go func(trees []string, files chan<- FileInfo) {
			generate(trees, files, genResults)
			close(files)
			genwg.Done()
		}(rctrl.Trees, files)

		nstreams := rctrl.Ctrl.Streams()
		for i := 0; i < nstreams; i++ {
			gatwg.Add(1)
			go func(ctrl utility.ReadChannelController, files <-chan FileInfo) {
				runner.Gather(&ctrl, files, gatherToAgg)
				gatwg.Done()
			}(rctrl.Ctrl, files)
		}
	} // for each per-device controller

	go func() {
		// Cleans up when all are done
		genwg.Wait()
		defer close(genResults)
	}()

	go func() {
		gatwg.Wait()
		close(gatherToAgg)
	}()

	return genResults, runner.Aggregate(gatherToAgg)
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
