package seal

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Byron/godi/cli"
	"github.com/Byron/godi/utility"

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
		Action:    func(c *gcli.Context) { cli.RunAction(&cmd, c) },
		Before:    func(c *gcli.Context) error { return cli.CheckCommonArgs(&cmd, c) },
	}

	out[0] = seal
	return out
}

func (s *SealCommand) Init(numReaders, numWriters int, items []string) error {
	s.pCtrl = utility.NewReadChannelController(numReaders)
	s.Trees = items

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
		tree = path.Clean(tree)
		if !filepath.IsAbs(tree) {
			var err error
			tree, err = filepath.Abs(tree)
			if err != nil {
				return err
			}
		}
		s.Trees[i] = tree
	}

	if len(invalidTrees) > 0 {
		return errors.New("Coulnd't read at least one of the given trees to verify: " + strings.Join(invalidTrees, ", "))
	}
	if len(noTrees) > 0 {
		return errors.New("The following trees are no directory: " + strings.Join(noTrees, ", "))
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

	return nil
}
