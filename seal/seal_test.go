package seal_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/testlib"
)

func TestSeal(t *testing.T) {
	datasetTree, dataFile, _ := testlib.MakeDatasetOrPanic()
	defer testlib.RmTree(datasetTree)
	var cmd seal.SealCommand
	var err error

	_, err = seal.NewCommand([]string{dataFile}, 1)
	if err == nil {
		t.Error("Expected it to not like files as directory")
	} else {
		t.Log(err)
	}

	cmd, err = seal.NewCommand([]string{datasetTree, filepath.Join(datasetTree, testlib.FirstSubDir, "..", testlib.FirstSubDir)}, 1)
	if err != nil {
		t.Error("Expected to not fail sanitization")
	} else if len(cmd.Trees) != 1 {
		t.Error("Trees should have been pruned, contained one should have been dropped")
	}

	cmd, err = seal.NewCommand([]string{datasetTree}, runtime.GOMAXPROCS(0))
	if err != nil {
		t.Error("Sanitize didn't like existing tree")
	}

	// Return true if we should break the loop
	resHandler := func(res godi.Result) {
		if res.Error() != nil {
			t.Error(res.Error())
		} else {
			t.Log(res.Info())
		}

	} // end resHandler

	godi.StartEngine(&cmd, uint(runtime.GOMAXPROCS(0)), resHandler, resHandler)
}
