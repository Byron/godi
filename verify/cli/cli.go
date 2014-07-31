/*
Package cli implements the command-line interface for the Command, for use by the cli.App
*/
package cli

import (
	"github.com/Byron/godi/cli"
	"github.com/Byron/godi/verify"

	gcli "github.com/codegangsta/cli"
)

const verifyDescription = `
	Compare stored disk-data with seal to detect changes.

	This command will read all files contained in the seal from disk and retake their signature.
	If it doesn't match the one stored in the seal file, the file on disk changed and either
	has been tempered with, or was corrupted some other way.

	Verify will clearly indicate changes in size, changes in contents, or missing files.

	[arguments ...] are one or more seal files, for example

	godi verify /Volumes/backup/godi_2014-07-30_102259.gobz path/to/godi_2012-07-10_102224.mhl
`

// return subcommands for our particular area of algorithms
func SubCommands() []gcli.Command {
	out := make([]gcli.Command, 1)
	cmd := verify.Command{}

	verify := gcli.Command{
		Name:      verify.Name,
		ShortName: "",
		Usage:     verifyDescription,
		Action:    func(c *gcli.Context) { cli.RunAction(&cmd, c) },
		Before:    func(c *gcli.Context) error { return cli.CheckCommonFlagsAndInit(&cmd, c) },
	}

	out[0] = verify
	return out
}
