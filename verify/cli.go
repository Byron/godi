package verify

import (
	"errors"
	"fmt"

	"github.com/Byron/godi/cli"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/utility"

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
		Before:    func(c *gcli.Context) error { return cli.CheckCommonArgs(&cmd, c) },
	}

	out[0] = verify
	return out
}

func (s *VerifyCommand) Init(numReaders, numWriters int, items []string) error {
	s.pCtrl = utility.NewReadChannelController(numReaders)
	s.Indices = items

	if len(s.Indices) == 0 {
		return errors.New("Please provide at least one seal file")
	}

	for _, index := range s.Indices {
		if codec := codec.NewByPath(index); codec == nil {
			return fmt.Errorf("Unknown seal file format: '%s'", index)
		}
	}

	return nil
}
