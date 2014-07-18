package cli

import (
	"errors"
	"flag"
	"fmt"
	"godi/seal"
	"os"
)

const usage = "godi {seal} [--help] args"

type SubCommand interface {
	// Keep the given unparsed arguments, as free argument list
	// May be discarded if unsupported, which is when an error should be provided
	SetUnparsedArgs(args []string) error

	// Check all set arguments for validity and sanity. This may involve verifying given paths are accessible
	// before trying to actually use them
	// Returns error to signal issues
	SanitizeArgs() error

	// Setup the given parser with whatever flags you require
	SetupParser(parser *flag.FlagSet) error
}

// Return a string representing detailed usage information, possibly based on the given parser
func HelpString(parser *flag.FlagSet) string {
	return usage + "\nTODO: LONG HELP"
}

// Parse the given arguments, which are assuming to NOT contain the invoked executable as first argument
// or use os.Args if there are no arguments given.
// We will parse subcommands ourselves, and then return one of the *Args types to indicate which subcomamnd
// was actually chosen.
// If there was an error parsing the arguments, it's error string will be usage information or what the problem
// was, useful for the end-user.
// If there was an error, options will be nil
// The interface return value can also be a string representing a detailed help string
// You should run SanitizeArgs yourself
func ParseArgs(args ...string) (interface{}, error) {
	if len(args) == 0 {
		args = append(args, os.Args[1:]...)
	}

	if len(args) < 1 {
		return nil, errors.New(usage)
	}

	// Parse based on subcommand
	var parser *flag.FlagSet
	var command SubCommand
	var helpFlag = false
	const helpUsage = "Prints detailed help"

	switch cmd := args[0]; cmd {
	case seal.Name:
		command = &seal.SealCommand{}
	default:
		return nil, fmt.Errorf("Invalid subcommand: %s\n%s", cmd, usage)
	}

	if command == nil {
		panic("Should have command set by now")
	}

	parser = flag.NewFlagSet(args[0], flag.ContinueOnError)
	command.SetupParser(parser)

	parser.BoolVar(&helpFlag, "help", helpFlag, helpUsage)

	if err := parser.Parse(args[1:]); err != nil {
		return nil, errors.New(usage + "\n" + err.Error())
	}

	if helpFlag {
		return HelpString(parser), nil
	}

	err := command.SetUnparsedArgs(parser.Args())
	return command, err
}
