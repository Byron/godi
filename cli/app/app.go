package app

import (
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
		cli.IntFlag{gocli.NumReadersFlagName + ", nr", 1, "Amount of parallel read streams we can use"},
	}
	app.Version = "v0.2.0"

	return app
}
