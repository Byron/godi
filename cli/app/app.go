/*
Package app provides a codegansta/cli.App initialized with all known godi commands.
*/
package app

import (
	"fmt"

	"github.com/Byron/godi/api"
	gocli "github.com/Byron/godi/cli"
	scli "github.com/Byron/godi/seal/cli"
	vcli "github.com/Byron/godi/verify/cli"

	"github.com/codegangsta/cli"
)

var (
	excludePatternsDescription = fmt.Sprintf(`A comma separated list of patterns causing
	an input file to be excluded when matched. The filter can apply 
	to entire classes of files, or represents a glob pattern.
	%-8s: Ignore all symbolic links
	%-8s: Ignore all hidden files. Only files starting with a period are hidden
	%-8s: Ignore all godi seal files. They are matched by their default name
	%-8s: Ignore files which change a lot or are expendable,
	like '.DS_Store' on OSX. Devices like tty's match too.
	Everything else is interpreted as glob, and '*.mov' will exclude 
	all quicktime mov files.
	A filter like '%s,%s,%s,*.mov,*.dpx' would ignore all hidden 
	files, symbolic links, volatile files, as well those ending with .mov and .dpx.
	If there is nothing behind the '=' sign, no all files will be handled.
	`, api.FilterSymlinks, api.FilterHidden, api.FilterSeals, api.FilterVolatile,
		api.FilterHidden, api.FilterVolatile, api.FilterSymlinks)

	inputStreamsDescription = `Amount of parallel streams per input device.
	If you device is very fast, or if the dataset contains many small files, 
	it may inrease performance to set values of two or higher.`

	verbosityDescription = fmt.Sprintf(`"Defines the overall amount of information you see,
	may take one of '%s', '%s' and '%s'.
	%-5s: disables all output. That way, you can only use the program's 
	error code to learn about success or failure
	%-5s: shows only errors and results
	%-5s: shows all of the above, and detailed progress information,
	might be too verbose when handling many small files`,
		api.Info, api.Error, api.LogDisabled,
		api.LogDisabled, api.Error, api.Info)
)

// Builds a Cli app from all commands we know and returns it
func NewGodiApp() *cli.App {
	app := cli.NewApp()

	// build up all known commands
	cmds := []cli.Command{}
	cmds = append(cmds, scli.SubCommands()...)
	cmds = append(cmds, vcli.SubCommands()...)

	app.Usage = `Verify data integrity and transfer data securely at highest speeds.

	Read more in the web-manual at http://byron.github.io/godi
	`
	app.Commands = cmds

	// general flags
	app.Flags = []cli.Flag{
		cli.IntFlag{gocli.StreamsPerInputDeviceFlagName + ", spid", 1, inputStreamsDescription},
		cli.StringFlag{gocli.LogLevelFlagName,
			api.Error.String(),
			verbosityDescription,
		},
		cli.StringFlag{gocli.FileExcludePatternFlagName,
			api.FilterVolatile.String(),
			excludePatternsDescription,
		},
	}
	app.Version = "v1.0.1"
	app.Author = "Sebastian Thiel & Contributors"

	gocli.AddAdditinalFlags(app)

	return app
}
