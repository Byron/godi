package verify

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/cli"
	"github.com/Byron/godi/codec"

	gcli "github.com/codegangsta/cli"
)

// return subcommands for our particular area of algorithms
func SubCommands() []gcli.Command {
	out := make([]gcli.Command, 1)
	cmd := VerifyCommand{}

	verify := gcli.Command{
		Name:      Name,
		ShortName: "",
		Usage:     "Compare stored data with seal to detect bit errors",
		Action:    func(c *gcli.Context) { cli.RunAction(&cmd, c) },
		Before:    func(c *gcli.Context) error { return cli.CheckCommonFlagsAndInit(&cmd, c) },
	}

	out[0] = verify
	return out
}

func (s *VerifyCommand) Init(numReaders, numWriters int, items []string, maxLogLevel api.Priority, filters []api.FileFilter) error {
	if len(items) == 0 {
		return errors.New("Please provide at least one seal file")
	}

	indexDirs := make([]string, len(items))
	for i, index := range items {
		if codec := codec.NewByPath(index); codec == nil {
			return fmt.Errorf("Unknown seal file format: '%s'", index)
		}
		if _, err := os.Stat(index); err != nil {
			return fmt.Errorf("Cannot access seal file at '%s'", index)
		}
		indexDirs[i] = filepath.Dir(index)
	}

	s.InitBasicRunner(numReaders, indexDirs, maxLogLevel, filters)
	s.Items = items
	return nil
}
