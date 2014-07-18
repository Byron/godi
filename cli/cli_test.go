package cli_test

import (
	"godi/cli"
	"godi/seal"
	"strings"
	"testing"
)

func TestParsing(t *testing.T) {
	if _, err := cli.ParseArgs("invalid_subcommand"); err == nil {
		t.Error("Shouldn't parse invalid_subcommand")
	} else {
		t.Log(err.Error())
	}

	cmds := []string{seal.Name}
	for _, cmd := range cmds {
		if res, err := cli.ParseArgs(cmd, "--help"); err != nil {
			t.Error("--help must exist in every subcommand")
		} else if str, ok := res.(string); !ok {
			t.Errorf("Didn't see string return value: %#v %v", res, err)
		} else {
			t.Log(str)
		}
	}
}

func TestSealParsing(t *testing.T) {
	sealcmd := func(args ...string) (interface{}, error) {
		nargs := make([]string, len(args)+1)
		nargs[0] = "seal"
		copy(nargs[1:], args)
		return cli.ParseArgs(nargs...)
	}

	sealcmdChecked := func(args ...string) *seal.SealCommand {
		if res, err := sealcmd(args...); err != nil {
			t.Error("Parsing shouldn't fail with no arguments")
		} else if cmd, ok := res.(*seal.SealCommand); !ok {
			t.Errorf("invalid return type: %v", res)
		} else {
			return cmd
		}
		panic("Shouldn't be here")
	}

	if res, err := sealcmd("foo", "bar"); err != nil {
		t.Errorf("seal should't fail if directory can't be read - it's part of the sanitization: %v", err)
	} else if res == nil {
		t.Error("no error, yet no args")
	} else if scmd, ok := res.(*seal.SealCommand); !ok {
		t.Errorf("Didn't get SealCommand, but %#v", res)
	} else if len(scmd.Trees) != 2 {
		t.Error("Didn't parse exactly 2 Trees")
	} else if err := scmd.SanitizeArgs(); err == nil {
		t.Error("Expected that all directories are invalid")
	} else if !strings.Contains(err.Error(), "foo, bar") {
		t.Errorf("Error string unexpected: %v", err)
	} else {
		t.Log(err)
	}

	cmd := sealcmdChecked()
	if err := cmd.SanitizeArgs(); err == nil {
		t.Error("Expected error as empty trees are disallowed")
	} else {
		t.Log(err)
	}
}
