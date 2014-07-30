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
	cmd := VerifyCommand{}

	verify := gcli.Command{
		Name:      Name,
		ShortName: "",
		Usage:     verifyDescription,
		Action:    func(c *gcli.Context) { cli.RunAction(&cmd, c) },
		Before:    func(c *gcli.Context) error { return cli.CheckCommonFlagsAndInit(&cmd, c) },
	}

	out[0] = verify
	return out
}

func (s *VerifyCommand) Init(numReaders, numWriters int, items []string, maxLogLevel api.Priority, filters []api.FileFilter) (err error) {
	if len(items) == 0 {
		return errors.New("Please provide at least one seal file")
	}

	validItems, err := api.ParseSources(items, true)
	if err != nil {
		return
	}

	indexDirs := make([]string, len(validItems))
	for i, index := range validItems {
		if codec := codec.NewByPath(index); codec == nil {
			return fmt.Errorf("Unknown seal file format: '%s'", index)
		}
		if _, err := os.Stat(index); err != nil {
			return fmt.Errorf("Cannot access seal file at '%s'", index)
		}
		indexDirs[i] = filepath.Dir(index)
	}

	s.InitBasicRunner(numReaders, indexDirs, maxLogLevel, filters)
	s.Items = validItems
	return nil
}
