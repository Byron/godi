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

	excludePatternsDescription := fmt.Sprintf(`A comma separated list of patterns that will cause
	an input file to be excluded when matched. The filter can apply to entire classes
	of files, or represents a glob pattern.
	%-8s: Ignore all symbolic links
	%-8s: Ignore all hidden files. Only files starting with a period are considered to match.
	%-8s: Ignore all godi seal files. Those are matched by their default name only.
	%-8s: Ignore files which are known to changed a lot and used for bookkeeping only, like '.DS_Store' on OSX. Devices like tty's match too.
	Everything else is interpreted as glob, and '*.mov' will exclude all quicktime mov files.
	A filter like '%s,%s,%s,*.mov,*.dpx' would ignore all hidden files, symbolic links, volatile files, 
	as well those having the .mov and .dpx extension.
	If there is nothing behind the '=' sign, no all files will be handled.
	`, api.FilterSymlinks, api.FilterHidden, api.FilterSeals, api.FilterVolatile,
		api.FilterHidden, api.FilterVolatile, api.FilterSymlinks)

	// general flags
	app.Flags = []cli.Flag{
		cli.IntFlag{gocli.StreamsPerInputDeviceFlagName + ", spid", 1, "Amount of parallel streams per input device"},
		cli.StringFlag{gocli.LogLevelFlagName,
			api.Error.String(),
			fmt.Sprintf("One of %s, %s, or '%s' to disable all output", api.Info, api.Error, api.LogDisabled),
		},
		cli.StringFlag{gocli.FileExcludePatternFlagName,
			api.FilterVolatile.String(),
			excludePatternsDescription,
		},
	}
	app.Version = "v0.9.0"

	gocli.AddAdditinalFlags(app)

	return app
}
