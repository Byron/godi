// +build debug

package cli

import (
	"os"
	"runtime/pprof"

	"github.com/codegangsta/cli"
)

const (
	memProfileFlagName = "memory-profile"
)

func AddAdditinalFlags(a *cli.App) {
	// nop
	a.Flags = append(a.Flags, cli.StringFlag{
		memProfileFlagName, "", "If set, the given file will contain the memory profile of the previous run"},
	)
}

func CliFinishApp(c *cli.Context) error {
	mpf := c.GlobalString(memProfileFlagName)

	// MEM PROFILE
	/////////////////
	if len(mpf) > 0 {
		f, err := os.Create(mpf)
		if err != nil {
			return err
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
	return nil
}

func parseAdditionalFlags(c *cli.Context) error {
	return nil
}
