package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/utility"

	"github.com/codegangsta/cli"
)

const (
	StatisticalLoggingInterval    = 1 * time.Second
	TimeEpsilon                   = 40 * time.Millisecond
	StreamsPerInputDeviceFlagName = "streams-per-input-device"
)

// Wraps an Aggregator handler and tracks last time the handler was called.
// If it was more than a certain amount of seconds ago, we will release a message about what's
// currently going on based on the statistical information we are passed
// NOTE: Even though it would be cleaner to just inject messages into the results channel, this way
// we wouldn't know when the last message was emitted, possibly emitting too much
func MakeStatisticalLogHandler(stats *utility.Stats, handler func(api.Result)) func(api.Result) {
	lastLoggedAt := time.Now()
	lastStat := *stats

	// An observer, printing out statistics as needed
	// We check a bit more often than the time after which to log some stats, to be more responsive
	// Lets be late at max 1/8 of a second
	ticker := time.Tick(StatisticalLoggingInterval / 4)
	go func() {
		for now := range ticker {
			td := now.Sub(lastLoggedAt) // temporal distance
			if td+TimeEpsilon < StatisticalLoggingInterval {
				continue
			}
			// Otherwise, prepare statistics
			fmt.Println(stats.DeltaString(&lastStat, td, ""))
			lastLoggedAt = now
			stats.CopyTo(&lastStat)
		}
	}()

	return func(r api.Result) {
		lastLoggedAt = time.Now()
		handler(r)
	}
}

func LogHandler(r api.Result) {
	if r.Error() != nil {
		fmt.Fprintln(os.Stderr, r.Error())
	} else {
		info, _ := r.Info()
		fmt.Fprintln(os.Stdout, info)
	}
}

// Runs a standard runner from within the cli, dealing with errors accoringly
// Both handlers may be nil to use a default one
func RunAction(cmd api.Runner, c *cli.Context) {
	// checkArgs must have initialized the seal command, so we can just run it
	handler := MakeStatisticalLogHandler(cmd.Statistics(), LogHandler)
	err := api.StartEngine(cmd, handler, handler)
	if err != nil {
		os.Exit(1)
	}
}

// As CheckCommonFlagsAndInit, but will return all parsed and verified common values, including an optional error
func CheckCommonFlags(c *cli.Context) (nr int, err error) {
	// Put parsed args in cmd and sanitize it
	nr = c.GlobalInt(StreamsPerInputDeviceFlagName)
	if nr < 1 {
		return 0, fmt.Errorf("--%v must not be smaller than 1", StreamsPerInputDeviceFlagName)
	}

	return
}

// Check common args and init a command with them.
// Further init and checking should be done in specialized function
func CheckCommonFlagsAndInit(cmd api.Runner, c *cli.Context) error {
	nr, err := CheckCommonFlags(c)
	if err != nil {
		return err
	}

	return cmd.Init(nr, nr, c.Args())
}
