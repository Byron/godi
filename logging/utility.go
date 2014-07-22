package logging

import (
	"fmt"
	"os"

	"github.com/Byron/godi/api"
)

func CLILogger(r godi.Result) {
	if r.Error() != nil {
		fmt.Fprintln(os.Stderr, r.Error())
	} else {
		info, _ := r.Info()
		fmt.Fprintln(os.Stdout, info)
	}
}
