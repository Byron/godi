// Do-nothing standin for non-developer builds
// +build !debug

package cli

import (
	"github.com/codegangsta/cli"
)

func AddAdditinalFlags(a *cli.App) {
	// nop
}

func CliFinishApp(c *cli.Context) error {
	// nop
	return nil
}

func parseAdditionalFlags(c *cli.Context) error {
	// nop
	return nil
}
