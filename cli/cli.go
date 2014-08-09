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

	"github.com/codegangsta/cli"
)

const (
	StreamsPerInputDeviceFlagName = "streams-per-input-device"
	LogLevelFlagName              = "verbosity"
	FileExcludePatternFlagName    = "file-exclude-patterns"
)

func MakeLogHandler(maxLogLevel api.Importance) func(r api.Result) {
	statsFilter := api.StatisticsFilter{
		LastResultShownAt:    time.Now(),
		FirstStatisticsAfter: 125 * time.Millisecond,
	}

	return func(r api.Result) {
		info, prio := r.Info()
		if !maxLogLevel.MayLog(prio) {
			return
		}

		// Prune out timed messages - we only want to see them if there was nothing else
		if !statsFilter.OK(prio) {
			return
		}

		statsFilter.LastResultShownAt = time.Now()
		if r.Error() != nil {
			fmt.Fprintln(os.Stderr, info)
		} else {
			fmt.Fprintln(os.Stdout, info)
		}
	}
}

// Runs a standard runner from within the cli, dealing with errors accoringly
func RunAction(cmd api.Runner, c *cli.Context) {
	handler := MakeLogHandler(cmd.LogLevel())
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
