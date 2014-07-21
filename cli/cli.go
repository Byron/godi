package cli

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/Byron/godi/api"

	"github.com/codegangsta/cli"
)

const (
	NumReadersFlagName = "num-readers"
)

// Runs a standard runner from within the cli, dealing with errors accoringly
func RunAction(cmd godi.Runner, c *cli.Context) {
	// checkArgs must have initialized the seal command, so we can just run it

	logger := func(r godi.Result) {
		if r.Error() != nil {
			fmt.Fprintln(os.Stderr, r.Error())
		} else {
			info, _ := r.Info()
			fmt.Fprintln(os.Stdout, info)
		}
	}

	err := godi.StartEngine(cmd, uint(runtime.GOMAXPROCS(0)), logger, logger)
	if err != nil {
		os.Exit(1)
	}
}

// Check common args and init a command with them.
// Further init and checking should be done in specialized function
func CheckCommonArgs(cmd godi.Runner, c *cli.Context) error {
	// Put parsed args in cmd and sanitize it
	nr := c.GlobalInt(NumReadersFlagName)
	if nr < 1 {
		return errors.New("--num-readers must not be smaller than 1")
	}
	return cmd.Init(nr, 0, c.Args())
}
