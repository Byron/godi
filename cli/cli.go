/*
Package cli provides common utilities for commandline handling in conjunction with codegansta/cli
*/
package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/io"

	"github.com/codegangsta/cli"
)

const (
	StatisticalLoggingInterval    = 1 * time.Second
	TimeEpsilon                   = 40 * time.Millisecond
	StreamsPerInputDeviceFlagName = "streams-per-input-device"
	LogLevelFlagName              = "verbosity"
	FileExcludePatternFlagName    = "file-exclude-patterns"
)

// Wraps an Aggregator handler and tracks last time the handler was called.
// If it was more than a certain amount of seconds ago, we will release a message about what's
// currently going on based on the statistical information we are passed
// NOTE: Even though it would be cleaner to just inject messages into the results channel, this way
// we wouldn't know when the last message was emitted, possibly emitting too much
// Done should be called to signal that we should stop
func MakeStatisticalLogHandler(stats *api.Stats, handler func(api.Result) bool, done <-chan bool) func(api.Result) bool {
	lastLoggedAt := time.Now()
	lastStat := *stats
	minLogInterval := StatisticalLoggingInterval / 6

	// An observer, printing out statistics as needed
	// We check a bit more often than the time after which to log some stats, to be more responsive
	// Lets be late at max 1/8 of a second
	ticker := time.NewTicker(minLogInterval)
	go func() {
		for now := range ticker.C {
			select {
			case <-done:
				ticker.Stop()
			default:
			}

			td := now.Sub(lastLoggedAt) // temporal distance
			if td+TimeEpsilon < minLogInterval {
				continue
			}
			// The first time a log should be there faster to feel more responsive
			// From that point on, messages can come in more slowly
			minLogInterval = StatisticalLoggingInterval
			// Otherwise, prepare statistics
			fmt.Println(stats.DeltaString(&lastStat, td, io.StatsClientSep), stats.String())
			lastLoggedAt = now
			stats.CopyTo(&lastStat)
		}
	}()

	return func(r api.Result) bool {
		hasLogged := handler(r)
		if hasLogged {
			lastLoggedAt = time.Now()
		}
		return hasLogged
	}
}

func MakeLogHandler(maxLogLevel api.Importance) func(r api.Result) bool {
	return func(r api.Result) bool {
		info, prio := r.Info()
		if !maxLogLevel.MayLog(prio) {
			return false
		}

		if r.Error() != nil {
			fmt.Fprintln(os.Stderr, info)
		} else {
			fmt.Fprintln(os.Stdout, info)
		}
		return true
	}
}

// Runs a standard runner from within the cli, dealing with errors accoringly
// Both handlers may be nil to use a default one
func RunAction(cmd api.Runner, c *cli.Context) {
	// checkArgs must have initialized the seal command, so we can just run it
	handler := MakeLogHandler(cmd.LogLevel())
	if cmd.LogLevel() != api.LogDisabled {
		handler = MakeStatisticalLogHandler(cmd.Statistics(), handler, make(chan bool))
	}
	err := api.StartEngine(cmd, handler)
	nerr := CliFinishApp(c)
	if err != nil || nerr != nil {
		os.Exit(1)
	}
}

// As CheckCommonFlagsAndInit, but will return all parsed and verified common values, including an optional error
func CheckCommonFlags(c *cli.Context) (nr int, level api.Importance, filters []api.FileFilter, err error) {
	// Put parsed args in cmd and sanitize it
	nr = c.GlobalInt(StreamsPerInputDeviceFlagName)
	if nr < 1 {
		return 0, level, filters, fmt.Errorf("--%v must not be smaller than 1", StreamsPerInputDeviceFlagName)
	}

	level, err = api.ParseImportance(c.GlobalString(LogLevelFlagName))
	if err != nil {
		return
	}

	filterStr := c.GlobalString(FileExcludePatternFlagName)
	for _, fstr := range strings.Split(filterStr, ",") {
		if f, e := api.ParseFileFilter(fstr); e != nil {
			err = e
			return
		} else {
			filters = append(filters, f)
		}
	}

	err = parseAdditionalFlags(c)
	return
}

// Check common args and init a command with them.
// Further init and checking should be done in specialized function
func CheckCommonFlagsAndInit(cmd api.Runner, c *cli.Context) error {
	nr, level, filters, err := CheckCommonFlags(c)
	if err != nil {
		return err
	}

	return cmd.Init(nr, nr, c.Args(), level, filters)
}
