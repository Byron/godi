package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/cli"
)

const (
	PROGRAMMING_ERROR = 255
	ARGUMENT_ERROR    = 1
	OTHER_ERROR       = 3
)

func main() {
	// Always use all available CPUs - the user can limit resources using GOMAXPROCS and the flags for reader- and writer-procs
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd, err := cli.ParseArgs(os.Args[1:]...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ARGUMENT_ERROR)
	}

	switch c := cmd.(type) {
	case string:
		{
			// Handle help printing
			fmt.Fprintln(os.Stdout, c)
			os.Exit(ARGUMENT_ERROR)
		}
	case cli.SubCommand:
		if err := c.SanitizeArgs(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(ARGUMENT_ERROR)
		}

		if runner, ok := cmd.(godi.Runner); !ok {
			fmt.Fprintln(os.Stderr, "Didn't get Runner interface from cli parser")
			os.Exit(PROGRAMMING_ERROR)
		} else {
			godi.StartEngine(runner, uint(runtime.GOMAXPROCS(0)))
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid command type returned - it didn't support the runner interfacea: %#v", cmd)
		os.Exit(PROGRAMMING_ERROR)
	}
}
