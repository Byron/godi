package seal

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/cli"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/io"
	"github.com/Byron/godi/verify"

	gcli "github.com/codegangsta/cli"
)

const (
	Sep                    = "--"
	verifyAfterCopy        = "verify"
	streamsPerOutputDevice = "streams-per-output-device"
	formatFlag             = "format"
	sealDescription        = `
	Generate a seal for one ore more directories to allow them to be verified later.

	A seal is a file with signatures, each of them is unique for the contents of the respective file.
	If the file's contents changes, it's signature changes as well, which allows to detect changes
	in a file reliably.

	The 'verify' sub-command can test if the files on disk didn't change compared to their sealed signature.

	[arguments ...] can be files or directories, for example

	godi seal my-anniversary.mov /Volumes/backup/`

	sealedCopyDescription = `
	Seal one or more directories and copy their contents to one or more destinations.

	To help protecting your valuable data effectively, you should have at least two 
	copies around in case one of thm gets corrupted. This sub-command
	does just that, and copies the data to one or more destinations while it is being sealed.

	[arguments ...] specify the source file(s) or directories, as well as the destination(s), for example
	godi sealed-copy s/ /Volumes/a
	godi sealed-copy s1/ s2/ -- /Volumes/a /Volumes/b
	`
)

var (
	usage = fmt.Sprintf(`Please specify sealed copies like so: source/ %s destination/
	%s can be omitted if there is only one source and one destination.`, Sep, Sep)
	formatDescription = fmt.Sprintf(`The format of the produced seal file, one of %s
	%s: is a compressed binary seal format, which is temper-proof and highly efficient, 
	handling millions of files easily.
	%s: is a human-readable XML format understood by mediahashlist.org, which will 
	be inefficient for large amount of files`,
		strings.Join(codec.Names(), ", "), codec.GobName, codec.MHLName)
)

// return subcommands for our particular area of algorithms
func SubCommands() []gcli.Command {
	cmdseal := SealCommand{mode: modeSeal}
	cmdcopy := SealCommand{mode: modeCopy}

	fmt := gcli.StringFlag{
		formatFlag,
		codec.GobName,
		formatDescription,
	}

	return []gcli.Command{
		gcli.Command{
			Name:      modeSeal,
			ShortName: "",
			Usage:     sealDescription,
			Action:    func(c *gcli.Context) { cli.RunAction(&cmdseal, c) },
			Before:    func(c *gcli.Context) error { return checkSeal(&cmdseal, c) },
			Flags:     []gcli.Flag{fmt},
		},
		gcli.Command{
			Name:      modeCopy,
			ShortName: "",
			Usage:     sealedCopyDescription,
			Action:    func(c *gcli.Context) { startSealedCopy(&cmdcopy, c) },
			Before:    func(c *gcli.Context) error { return checkSealedCopy(&cmdcopy, c) },
			Flags: []gcli.Flag{
				gcli.BoolFlag{verifyAfterCopy, "Run `godi verify` on all produced seals when copy is finished"},
				gcli.IntFlag{streamsPerOutputDevice + ", spod", 1, "Amount of parallel streams per output device"},
				fmt,
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

func checkSeal(cmd *SealCommand, c *gcli.Context) error {
	cmd.format = c.String(formatFlag)
	if len(cmd.format) > 0 {
		valid := false
		for _, name := range codec.Names() {
			if name == cmd.format {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("Invalid seal format '%s', must be one of %s", cmd.format, strings.Join(codec.Names(), ", "))
		}
	}

	if err := cli.CheckCommonFlagsAndInit(cmd, c); err != nil {
		return err
	}

	return nil
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
		err := api.StartEngine(cmd, aggHandler)
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
						err = api.StartEngine(verifycmd, handler)
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

func (s *SealCommand) Init(numReaders, numWriters int, items []string, maxLogLevel api.Priority, filters []api.FileFilter) (err error) {

	if len(s.format) == 0 {
		s.format = codec.GobName
	}

	if s.mode == modeSeal {
		if len(items) == 0 {
			return errors.New("Please provide at least one source directory to work on")
		}
		items, err = api.ParseSources(items, true)
		if err != nil {
			return
		}
		s.InitBasicRunner(numReaders, items, maxLogLevel, filters)
	} else if s.mode == modeCopy {
		finishSetup := func(sources, dtrees []string) error {
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
		} // end helper

		// Parses [src, ...] -- [dst, ...]
		err = errors.New(usage)
		if len(items) < 2 {
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
				sources, e := api.ParseSources(items[:i], true)
				if e != nil {
					return e
				}

				dtrees, e := api.ParseSources(items[i+1:], false)
				if e != nil {
					return e
				}

				return finishSetup(sources, dtrees)
			}
		} // for each item

		// So there is no separator, maybe it's source and destination ?
		if len(items) == 2 {
			return finishSetup(items[:1], items[1:])
		}

		// source-destination separator not found - prints usage
		return
	} else {
		panic(fmt.Sprintf("Unsupported mode: %s", s.mode))
	}
	return
}
