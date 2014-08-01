/*
For high-level developer documentation, please see http://byron.github.io/godi/development/
*/
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Byron/godi/cli/app"
)

func main() {
	// Always use all available CPUs - the user can limit resources using GOMAXPROCS and the flags for reader- and writer-procs
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := app.NewGodiApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
