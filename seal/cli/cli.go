/*
Package cli implements the commandline interface for the Command, ready for digestion by the cli.Appf
*/
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/cli"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/verify"

	gcli "github.com/codegangsta/cli"
)

const (
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
	formatDescription = fmt.Sprintf(`The format of the produced seal file, one of %s
	%s: is a compressed binary seal format, which is temper-proof and highly efficient, 
	handling millions of files easily.
	%s: is a human-readable XML format understood by mediahashlist.org, which will 
	be inefficient for large amount of files`,
		strings.Join(codec.Names(), ", "), codec.GobName, codec.MHLName)
)

// return subcommands for our particular area of algorithms
func SubCommands() []gcli.Command {
	cmdseal := seal.Command{Mode: seal.ModeSeal}
	cmdcopy := seal.Command{Mode: seal.ModeCopy}

	fmt := gcli.StringFlag{
		Name:  formatFlag,
		Value: codec.GobName,
		Usage: formatDescription,
	}

	return []gcli.Command{
		gcli.Command{
			Name:      seal.ModeSeal,
			ShortName: "",
			Usage:     sealDescription,
			Action:    func(c *gcli.Context) { cli.RunAction(&cmdseal, c) },
			Before:    func(c *gcli.Context) error { return checkSeal(&cmdseal, c) },
			Flags:     []gcli.Flag{fmt},
		},
		gcli.Command{
			Name:      seal.ModeCopy,
			ShortName: "",
			Usage:     sealedCopyDescription,
			Action:    func(c *gcli.Context) { startSealedCopy(&cmdcopy, c) },
			Before:    func(c *gcli.Context) error { return checkSealedCopy(&cmdcopy, c) },
			Flags: []gcli.Flag{
				gcli.BoolFlag{
					Name:  verifyAfterCopy,
					Usage: "Run `godi verify` on all produced seals when copy is finished"},
				gcli.IntFlag{
					Name:  streamsPerOutputDevice + ", spod",
					Value: 1,
					Usage: "Amount of parallel streams per output device"},
				fmt,
			},
		},
	}
}

func checkSeal(cmd *seal.Command, c *gcli.Context) error {
	cmd.Format = c.String(formatFlag)
	if len(cmd.Format) > 0 {
		valid := false
		for _, name := range codec.Names() {
			if name == cmd.Format {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("Invalid seal format '%s', must be one of %s", cmd.Format, strings.Join(codec.Names(), ", "))
		}
	}

	if err := cli.CheckCommonFlagsAndInit(cmd, c); err != nil {
		return err
	}

	return nil
}

func checkSealedCopy(cmd *seal.Command, c *gcli.Context) error {
	cmd.Verify = c.Bool(verifyAfterCopy)
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

func startSealedCopy(cmd *seal.Command, c *gcli.Context) {

	// Yes, currently the post-verification is only implemented in the CLI ...
	// Testing needs to do similar things to set it up ...
	if cmd.Verify {
		// Setup a aggregation result handler which tracks produced indices
		var indices []string
		cmdDone := make(chan bool)
		logHandler := cli.MakeLogHandler(cmd.LogLevel())

		handler := logHandler
		if cmd.LogLevel() != api.LogDisabled {
			handler = cli.MakeStatisticalLogHandler(&cmd.Stats, handler, cmdDone)
		}
		aggHandler := api.IndexTrackingResultHandlerAdapter(&indices, handler)

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
