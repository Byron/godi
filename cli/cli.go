package cli

import (
	"fmt"
	"os"

	"github.com/Byron/godi/api"

	"github.com/codegangsta/cli"
)

const (
	StreamsPerInputDeviceFlagName = "streams-per-input-device"
)

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
	err := api.StartEngine(cmd, LogHandler, LogHandler)
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
