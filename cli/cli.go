package cli

import (
	"github.com/codegangsta/cli"
)

const (
	NumReadersFlagName = "num-readers"
)

// Builds a Cli app from all commands we know and returns it
func NewGodiApp(cmds []cli.Command) *cli.App {
	app := cli.NewApp()
	app.Usage = "Verify data integrity and transfer data securely at highest speeds"
	app.Commands = cmds

	// general flags
	app.Flags = []cli.Flag{
		cli.IntFlag{NumReadersFlagName + ", nr", 1, "Amount of parallel read streams we can use"},
	}
	app.Version = "v0.1.0"

	return app
}
