package seal_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/testlib"
	"github.com/Byron/godi/verify"
)

func TestSeal(t *testing.T) {
	datasetTree, dataFile, _ := testlib.MakeDatasetOrPanic()
	defer testlib.RmTree(datasetTree)
	var cmd *seal.SealCommand
	var err error

	_, err = seal.NewCommand([]string{dataFile}, 1, 0)
	if err == nil {
		t.Fatal("Expected it to not like files as directory")
	} else {
		t.Log(err)
	}

	cmd, err = seal.NewCommand([]string{datasetTree, filepath.Join(datasetTree, testlib.FirstSubDir, "..", testlib.FirstSubDir)}, 1, 0)
	if err != nil {
		t.Fatal("Expected to not fail sanitization")
	} else if len(cmd.Items) != 1 {
		t.Fatal("Trees should have been pruned, contained one should have been dropped")
	}

	maxProcs := runtime.GOMAXPROCS(0)
	cmd, err = seal.NewCommand([]string{datasetTree}, maxProcs, 0)
	if err != nil {
		t.Fatal("Sanitize didn't like existing tree")
	}

	// Return true if we should break the loop
	resHandler := testlib.ResultHandler(t, false)

	if err := api.StartEngine(cmd, resHandler, resHandler); err != nil {
		t.Fatal(err)
	}

	// SEALED COPY
	//////////////
	cmd, err = seal.NewCommand([]string{datasetTree}, maxProcs, 1)
	if err == nil {
		t.Fatal("Can't copy with just a source")
	} else {
		t.Log(err)
	}

	cmd, err = seal.NewCommand([]string{datasetTree, seal.Sep, filepath.Join(datasetTree, "..", filepath.Base(datasetTree))}, maxProcs, 1)
	if err == nil {
		t.Fatal("Must rule out possiblity of copying something onto itself")
	} else {
		t.Log(err)
	}

	// this one is deleted in rmTree of datasetTree - it's a subdir !
	invalidDestination := filepath.Join(datasetTree, "copy-destination")
	if err = os.Mkdir(invalidDestination, 0777); err != nil {
		t.Fatal(err)
	}
	cmd, err = seal.NewCommand([]string{datasetTree, seal.Sep, invalidDestination}, maxProcs, 1)
	if err == nil {
		t.Fatal("Cannot copy something into itself - gather would become somewhat recursive !")
	} else {
		t.Log(err)
	}

	cmd, err = seal.NewCommand([]string{datasetTree, seal.Sep, "does/not/exist"}, maxProcs, 1)
	if err == nil {
		t.Fatal("Destination must exist")
	} else {
		t.Log(err)
	}

	copyDestination1, _ := ioutil.TempDir("", "sealed-copy")
	defer testlib.RmTree(copyDestination1)
	copyDestination2, _ := ioutil.TempDir("", "sealed-copy")
	defer testlib.RmTree(copyDestination2)

	cmd, err = seal.NewCommand([]string{datasetTree, seal.Sep, copyDestination1, copyDestination2}, maxProcs, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Finally, perform the operation
	var indices []string
	if err := api.StartEngine(cmd, resHandler, seal.IndexTrackingResultHandlerAdapter(&indices, resHandler)); err != nil {
		t.Fatal(err)
	}

	if len(indices) != 2 {
		t.Fatal("Should have parsed 2 indices")
	}

	verifycmd, err := verify.NewCommand(indices, maxProcs)
	if err != nil {
		t.Fatal(err)
	}

	if err := api.StartEngine(verifycmd, resHandler, resHandler); err != nil {
		t.Fatal("Couldn't verify files that were just written")
	}
}
