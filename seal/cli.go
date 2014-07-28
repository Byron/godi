package seal

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/cli"
	"github.com/Byron/godi/io"
	"github.com/Byron/godi/verify"

	gcli "github.com/codegangsta/cli"
)

const (
	Sep                    = "--"
	usage                  = "Please specify sealed copies like so: source/ -- destination/"
	verifyAfterCopy        = "verify"
	streamsPerOutputDevice = "streams-per-output-device"
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
			Before:    func(c *gcli.Context) error { return cli.CheckCommonFlagsAndInit(&cmdseal, c) },
		},
		gcli.Command{
			Name:      modeCopy,
			ShortName: "",
			Usage:     "Generate a seal for one ore more directories and copy their contents to a destination directory",
			Action:    func(c *gcli.Context) { startSealedCopy(&cmdcopy, c) },
			Before:    func(c *gcli.Context) error { return checkSealedCopy(&cmdcopy, c) },
			Flags: []gcli.Flag{
				gcli.BoolFlag{verifyAfterCopy, "Run `godi verify` on all produced seals when copy is finished"},
				gcli.IntFlag{streamsPerOutputDevice + ", spod", 1, "Amount of parallel streams per output device"},
			},
		},
	}
}

// Returns a handler whichasd will track seal/index files, and call the given handler aftrewards, writing the
// into the provided slice
func IndexTrackingResultHandlerAdapter(indices *[]string, handler func(r api.Result) bool) func(r api.Result) bool {
	return func(r api.Result) (res bool) {
		res = handler(r)
		if r == nil || r.FileInformation() == nil {
			return
		}
		if r.FileInformation().Size < 0 {
			*indices = append(*indices, r.FileInformation().Path)
		}
		return
	}
}

func checkSealedCopy(cmd *SealCommand, c *gcli.Context) error {
	cmd.verify = c.Bool(verifyAfterCopy)
	// have to do init ourselves as we set amount of writers
	nr, level, filters, err := cli.CheckCommonFlags(c)
	if err != nil {
		return err
	}

	nw := c.Int(streamsPerOutputDevice)
	if nw < 1 {
		return fmt.Errorf("--%v must not be smaller than 1", streamsPerOutputDevice)
	}

	return cmd.Init(nr, nw, c.Args(), level, filters)
}

func startSealedCopy(cmd *SealCommand, c *gcli.Context) {

	// Yes, currently the post-verification is only implemented in the CLI ...
	// Testing needs to do similar things to set it up ...
	if cmd.verify {
		// Setup a aggregation result handler which tracks produced indices
		var indices []string
		cmdDone := make(chan bool)
		logHandler := cli.MakeLogHandler(cmd.LogLevel())

		handler := logHandler
		if cmd.LogLevel() != api.LogDisabled {
			handler = cli.MakeStatisticalLogHandler(&cmd.Stats, handler, cmdDone)
		}
		aggHandler := IndexTrackingResultHandlerAdapter(&indices, handler)

		// and run ourselves
		err := api.StartEngine(cmd, handler, aggHandler)
		// Make sure we don't keep logging while verification is going with its own handler
		close(cmdDone)

		if err == nil && len(indices) == 0 {
			panic("Unexpectedly I didn't see a single seal index without error")
		} else if len(indices) > 0 {
			// no matter whether we have an error, try to verify what's there
			select {
			case <-cmd.Done:
				// this does nothing, most importantly, it doesn't run verify, as we don't run it
				// after cancellation. It's arguable whether we migth want to do that anyway
				// as the index is valid !
			default:
				{
					// prepare and run a verify command
					verifycmd, err := verify.NewCommand(indices, c.GlobalInt(cli.StreamsPerInputDeviceFlagName))
					if err == nil {
						handler = logHandler
						if verifycmd.LogLevel() != api.LogDisabled {
							handler = cli.MakeStatisticalLogHandler(&verifycmd.Stats, handler, make(chan bool))
						}
						err = api.StartEngine(verifycmd, handler, handler)
					}
				}
			}
		}

		// Finally, exit with appropriate error code
		if err != nil {
			os.Exit(1)
		}
	} else {
		// copy without verify
		cli.RunAction(cmd, c)
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
				if i == x || strings.HasPrefix(ltree, rtree+string(os.PathSeparator)) {
					continue
				}
				validTrees = api.AppendUniqueString(validTrees, ltree)
			}
		}
		if len(validTrees) == 0 {
			panic("Didn't find a single valid tree after pruning - shouldn't happen")
		}

		res = validTrees
	}

	return res, nil
}

func (s *SealCommand) Init(numReaders, numWriters int, items []string, maxLogLevel api.Priority, filters []api.FileFilter) (err error) {
	if s.mode == modeSeal {
		if len(items) == 0 {
			return errors.New("Please provide at least one source directory to work on")
		}
		items, err = parseSources(items)
		if err != nil {
			return
		}
		s.InitBasicRunner(numReaders, items, maxLogLevel, filters)
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
				sources, e := parseSources(items[:i])
				if e != nil {
					return e
				}

				dtrees, e := parseSources(items[i+1:])
				if e != nil {
					return e
				}

				// Make sure we don't copy onto ourselves
				for _, stree := range sources {
					for _, dtree := range dtrees {
						if strings.HasPrefix(dtree+string(os.PathSeparator), stree) {
							return fmt.Errorf("Cannot copy '%s' into it's own subdirectory or itself at '%s'", stree, dtree)
						}
					}
				}
				s.InitBasicRunner(numReaders, sources, maxLogLevel, filters)

				// build the device map with all writer destinations
				dm := io.DeviceMap(dtrees)

				// Finally, put all actual values into our list to have a deterministic iteration order.
				// After all, we don't really care about the device from this point on
				s.rootedWriters = make([]io.RootedWriteController, len(dm))
				for did, trees := range dm {
					// each device as so and so many destinations. Each destination uses the same write controller
					s.rootedWriters[did] = io.RootedWriteController{
						Trees: trees,
						Ctrl:  io.NewWriteChannelController(numWriters, numWriters*len(trees), &s.Stats.Stats),
					}
				} // for each tree set in deviceMap

				return nil
			}
		} // for each item

		// source-destination separator not found - prints usage
		return
	} else {
		panic(fmt.Sprintf("Unsupported mode: %s", s.mode))
	}
	return
}
