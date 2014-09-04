// +build web

package app

import (
	wcli "github.com/Byron/godi/web/cli"

	"github.com/codegangsta/cli"
)

// return subcommands for our particular area of algorithms
func optionalSubCommands() []cli.Command {
	return wcli.SubCommands()
}

func setupApp(app *cli.App) {

	action := func(c *cli.Context) {
		wcli.RunWebServer(
			&wcli.ServerInfo{
				Addr:    "localhost:9078",
				MayOpen: true,
			},
		)
	}

	app.Action = action
}
