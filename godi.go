package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Byron/godi/cli"
	"github.com/Byron/godi/seal"

	gcli "github.com/codegangsta/cli"
)

const (
	PROGRAMMING_ERROR = 255
	ARGUMENT_ERROR    = 1
	OTHER_ERROR       = 3
)

func main() {

	// build up all known commands
	cmds := []gcli.Command{}
	cmds = append(cmds, seal.SubCommands()...)

	// Always use all available CPUs - the user can limit resources using GOMAXPROCS and the flags for reader- and writer-procs
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewGodiApp(cmds)
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(OTHER_ERROR)
	}
}
