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
