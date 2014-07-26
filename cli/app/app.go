package app

import (
	"fmt"

	"github.com/Byron/godi/api"
	gocli "github.com/Byron/godi/cli"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/verify"

	"github.com/codegangsta/cli"
)

// Builds a Cli app from all commands we know and returns it
func NewGodiApp() *cli.App {
	app := cli.NewApp()

	// build up all known commands
	cmds := []cli.Command{}
	cmds = append(cmds, seal.SubCommands()...)
	cmds = append(cmds, verify.SubCommands()...)

	app.Usage = "Verify data integrity and transfer data securely at highest speeds"
	app.Commands = cmds

	// general flags
	app.Flags = []cli.Flag{
		cli.IntFlag{gocli.StreamsPerInputDeviceFlagName + ", spid", 1, "Amount of parallel streams per input device"},
		cli.StringFlag{gocli.LogLevelFlagName,
			api.Progress.String(),
			fmt.Sprintf("One of %s, %s, %s, %s, or '%s' to disable all output", api.Progress, api.Info, api.Warn, api.Error, api.LogDisabled),
		},
	}
	app.Version = "v0.4.0"

	gocli.AddAdditinalFlags(app)

	return app
}
