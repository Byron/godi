package seal

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Byron/godi/cli"
	"github.com/Byron/godi/utility"

	gcli "github.com/codegangsta/cli"
)

const (
	Sep                = "---"
	usage              = "Please specify sealed copies like so: source/ => destination/"
	numWritersFlagName = "num-writers"
)

// return subcommands for our particular area of algorithms
func SubCommands() []gcli.Command {
	cmdseal := SealCommand{mode: modeSeal}
	cmdcopy := SealCommand{mode: modeCopy}

	return []gcli.Command{
		gcli.Command{
			Name:      modeSeal,
			ShortName: "",
			Usage:     "Generate a seal for one ore more directories, which allows them to be verified later on",
			Action:    func(c *gcli.Context) { cli.RunAction(&cmdseal, c) },
			Before:    func(c *gcli.Context) error { return cli.CheckCommonArgs(&cmdseal, c) },
		},
		gcli.Command{
			Name:      modeCopy,
			ShortName: "",
			Usage:     "Generate a seal for one ore more directories and copy their contents to a destination directory",
			Action:    func(c *gcli.Context) { cli.RunAction(&cmdcopy, c) },
			Before:    func(c *gcli.Context) error { return cli.CheckCommonArgs(&cmdcopy, c) },
		},
	}
}

// Parse all valid source items from the given list.
// May either be files or directories. The returned list may be shorter, as contained paths are
// skipped automatically. Paths will be normalized.
func parseSources(items []string) (res []string, err error) {
	invalidTrees := make([]string, 0, len(items))
	res = make([]string, len(items))
	noTrees := make([]string, 0, len(items))
	copy(res, items)

	for i, tree := range res {
		if stat, err := os.Stat(tree); err != nil {
			invalidTrees = append(invalidTrees, tree)
		} else if !stat.IsDir() {
			noTrees = append(noTrees, tree)
		}
		tree = path.Clean(tree)
		if !filepath.IsAbs(tree) {
			tree, err = filepath.Abs(tree)
			if err != nil {
				return nil, err
			}
		}
		res[i] = tree
	}

	if len(invalidTrees) > 0 {
		return nil, errors.New("None of the following items exists: " + strings.Join(invalidTrees, ", "))
	}
	if len(noTrees) > 0 {
		return nil, errors.New("The following items are no directory: " + strings.Join(noTrees, ", "))
	}

	// drop trees which are a sub-tree of another, or which are equal !
	if len(res) > 1 {
		validTrees := make([]string, 0, len(res))
		for i, ltree := range res {
			for x, rtree := range res {
				if i == x || strings.HasPrefix(ltree, rtree) {
					continue
				}
				validTrees = append(validTrees, ltree)
			}
		}
		if len(validTrees) == 0 {
			panic("Didn't find a single valid tree after pruning - shouldn't happen")
		}

		res = validTrees
	}

	return res, nil
}

func (s *SealCommand) NumChannels() int {
	return s.pCtrl.Streams()
}

func (s *SealCommand) Init(numReaders, numWriters int, items []string) (err error) {
	s.pCtrl = utility.NewReadChannelController(numReaders)

	if s.mode == modeSeal {
		if len(items) == 0 {
			return errors.New("Please provide at least one source directory to work on")
		}
		s.SourceTrees, err = parseSources(items)
	} else if s.mode == modeCopy {
		// Parses [src, ...] => [dst, ...]
		err = errors.New(usage)
		if len(items) < 3 {
			return
		}

		for i, item := range items {
			if item == Sep {
				if i == 0 {
					return
				}
				if i == len(items)-1 {
					return
				}
				s.SourceTrees, err = parseSources(items[:i])
				if err != nil {
					return
				}

				var dtrees []string
				dtrees, err = parseSources(items[i+1:])
				if err != nil {
					return
				}

				// Make sure we don't copy onto ourselves
				for _, stree := range s.SourceTrees {
					for _, dtree := range dtrees {
						if strings.HasPrefix(dtree, stree) {
							return fmt.Errorf("Cannot copy '%s' into it's own subdirectory or itself at '%s'", stree, dtree)
						}
					}
				}

				// build the device map with all writer destinations
				dm := utility.DeviceMap(dtrees)

				// Finally, put all actual values into our list to have a deterministic iteration order.
				// After all, we don't really care about the device from this point on
				s.pWriters = make([]utility.RootedWriteController, len(dtrees))
				c := 0
				for _, trees := range dm {
					// each device as so and so many destinations. Each destination uses the same write controller
					wctrl := utility.NewWriteChannelController(numWriters, len(trees))
					for _, tree := range trees {
						s.pWriters[c] = utility.RootedWriteController{tree, &wctrl}
						c += 1
					}
				} // for each tree set in deviceMap

				return nil
			}
		} // for each item

		// not found
		return
	} else {
		panic(fmt.Sprintf("Unsupported mode: %s", s.mode))
	}
	return
}
