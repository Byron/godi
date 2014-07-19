package seal

import (
	"errors"
	"flag"
	"math"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/cli"

	gcli "github.com/codegangsta/cli"
)

// return subcommands for our particular area of algorithms
func SubCommands() []gcli.Command {
	out := make([]gcli.Command, 1)
	cmd := SealCommand{}

	seal := gcli.Command{
		Name:      Name,
		ShortName: "",
		Usage:     "Generate a seal for one ore more directories, which allows them to be verified later on",
		Action:    func(c *gcli.Context) { seal(&cmd, c) },
		Before:    func(c *gcli.Context) error { return checkArgs(&cmd, c) },
	}

	out[0] = seal
	return out
}

func seal(cmd *SealCommand, c *gcli.Context) {
	// checkArgs must have initialized the seal command, so we can just run it
	// TODO: Error handling
	godi.StartEngine(cmd, uint(runtime.GOMAXPROCS(0)))
}

func checkArgs(cmd *SealCommand, c *gcli.Context) error {
	// Put parsed args in cmd and sanitize it
	cmd.nReaders = c.GlobalInt(cli.NumReadersFlagName)
	cmd.SetUnparsedArgs(c.Args())
	return cmd.SanitizeArgs()
}

func (s *SealResult) Info() (string, godi.Priority) {
	if s.err != nil {
		return s.err.Error(), godi.Error
	}
	return s.msg, s.prio
}

func (s *SealResult) Error() error {
	return s.err
}

func (s *SealResult) FileInformation() *godi.FileInfo {
	return s.finfo
}

func (s *SealCommand) SetUnparsedArgs(args []string) error {
	s.Trees = args
	return nil
}

func (s *SealCommand) MaxProcs() uint {
	return uint(math.MaxUint32)
}

func (s *SealCommand) SanitizeArgs() (err error) {
	if len(s.Trees) == 0 {
		return errors.New("Please provide at least one tree to work on")
	}

	invalidTrees := make([]string, 0, len(s.Trees))
	noTrees := make([]string, 0, len(s.Trees))
	for i, tree := range s.Trees {
		if stat, err := os.Stat(tree); err != nil {
			invalidTrees = append(invalidTrees, tree)
		} else if !stat.IsDir() {
			noTrees = append(noTrees, tree)
		}
		s.Trees[i] = path.Clean(tree)
	}

	if len(invalidTrees) > 0 {
		return errors.New("Coulnd't read at least one of the given trees to verify: " + strings.Join(invalidTrees, ", "))
	}
	if len(noTrees) > 0 {
		return errors.New("The following trees are no directory: " + strings.Join(noTrees, ", "))
	}
	if s.nReaders < 1 {
		return errors.New("--num-readers must not be smaller than 1")
	}

	// drop trees which are a sub-tree of another
	if len(s.Trees) > 1 {
		validTrees := make([]string, 0, len(s.Trees))
		for i, ltree := range s.Trees {
			for x, rtree := range s.Trees {
				if i == x || strings.HasPrefix(ltree, rtree) {
					continue
				}
				validTrees = append(validTrees, ltree)
			}
		}
		if len(validTrees) == 0 {
			panic("Didn't find a single valid tree")
		}

		s.Trees = validTrees
	}

	return err
}

func (s *SealCommand) SetupParser(parser *flag.FlagSet) error {
	parser.IntVar(&s.nReaders, "num-readers", 1, "Amount of parallel read streams we can use")
	return nil
}
