package api

import (
	"sync"
	"time"

	"github.com/Byron/godi/io"
)

const (
	StatisticalResultInterval  = 125 * time.Millisecond
	StatisticalLoggingInterval = 1 * time.Second
	TimeEpsilon                = 40 * time.Millisecond
)

// Generate does all boilerplate required to be a valid generator
// Will produce as many generators as there are devices, each is handed a list of trees to handle
func Generate(rctrls io.RootedReadControllers,
	runner Runner,
	generate func([]string, chan<- FileInfo, chan<- Result)) <-chan Result {

	gatherToAgg := make(chan Result, runner.NumChannels())

	gatwg := sync.WaitGroup{} // wait group for gatherers

	// Spawn generators - each one has num-streams gatherers
	for _, rctrl := range rctrls {
		files := make(chan FileInfo)
		go func(trees []string, files chan<- FileInfo) {
			generate(trees, files, gatherToAgg)
			close(files)
		}(rctrl.Trees, files)

		nstreams := rctrl.Ctrl.Streams()
		for i := 0; i < nstreams; i++ {
			gatwg.Add(1)
			go func(ctrl io.ReadChannelController, files <-chan FileInfo) {
				runner.Gather(&ctrl, files, gatherToAgg)
				gatwg.Done()
			}(rctrl.Ctrl, files)
		}
	} // for each per-device controller

	go func() {
		gatwg.Wait()
		close(gatherToAgg)
	}()

	return runner.Aggregate(gatherToAgg)
}

// Aggregate is a general purpose implementation to gather fileInfo results
func Aggregate(results <-chan Result, done <-chan bool,
	resultHandler func(Result, chan<- Result) bool,
	finalizer func(chan<- Result),
	stats *Stats) <-chan Result {
	accumResult := make(chan Result)

	// For time-dependent insertion of results
	lastStat := *stats
	lastTimeResult := time.Now()

	// An observer, producing results with statistics after acertain interval
	// We set this value to be quite responsive
	ticker := time.NewTicker(StatisticalResultInterval)
	go func() {
		// We might try to send through a closed channel, especially in web-mode
		defer func() {
			recover()
		}()
		for now := range ticker.C {
			select {
			case <-done:
				ticker.Stop()
				return
			default:
			}

			// Otherwise, prepare statistics
			accumResult <- &BasicResult{
				Msg:  stats.DeltaString(&lastStat, now.Sub(lastTimeResult), io.StatsClientSep) + " " + stats.String(),
				Prio: PeriodicalStatistics,
			}
			lastTimeResult = time.Now()

			stats.CopyTo(&lastStat)
		}
	}()

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

// Utility type to determine if a Statistical result should be shown
// Assign the last time you used any result to this instance
type StatisticsFilter struct {
	LastResultShownAt    time.Time     // time at which you have used a result, whichever prio
	FirstStatisticsAfter time.Duration // time after which the first message will show
}

// Returns true if we can use the statistical information. You have to check if
// your result messsage has the right prio
func (s *StatisticsFilter) OK(prio Importance) bool {
	if prio != PeriodicalStatistics {
		s.FirstStatisticsAfter = StatisticalLoggingInterval
		return true
	}

	// Prune out timed messages - we only want to see them if there was nothing else
	td := time.Now().Sub(s.LastResultShownAt) // temporal distance
	if td+TimeEpsilon < s.FirstStatisticsAfter {
		s.FirstStatisticsAfter = StatisticalLoggingInterval
		return false
	}

	// The first time a log should be there faster to feel more responsive
	// From that point on, messages can come in more slowly
	s.FirstStatisticsAfter = StatisticalLoggingInterval
	return true
}
