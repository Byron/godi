// +build debug

package cli

import (
	"os"
	"runtime/pprof"

	"github.com/codegangsta/cli"
)

const (
	memProfileFlagName = "memory-profile"
	cpuProfileFlagName = "cpu-profile"
)

func AddAdditinalFlags(a *cli.App) {
	// nop
	a.Flags = append(a.Flags,
		cli.StringFlag{memProfileFlagName, "", "If set, the given file will contain the memory profile of the previous run"},
		cli.StringFlag{cpuProfileFlagName, "", "If set, the given file will contain the CPU profile of the previous run"},
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

	// CPU Profile
	///////////////
	if len(c.GlobalString(cpuProfileFlagName)) > 0 {
		pprof.StopCPUProfile()
	}
	return nil
}

func parseAdditionalFlags(c *cli.Context) error {
	cpuf := c.GlobalString(cpuProfileFlagName)
	if len(cpuf) > 0 {
		f, err := os.Create(cpuf)
		if err != nil {
			return err
		}
		pprof.StartCPUProfile(f)
	}
	return nil
}
