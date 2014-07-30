package verify_test

import (
	"os"
	"testing"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/testlib"
	"github.com/Byron/godi/verify"
)

func TestVerify(t *testing.T) {
	datasetTree, file, _ := testlib.MakeDatasetOrPanic()
	defer testlib.RmTree(datasetTree)

	sealcmd, _ := seal.NewCommand([]string{datasetTree}, 1, 0)
	resHandler := testlib.ResultHandler(t, false)

	// keeps track of created indices
	var indices []string
	err := api.StartEngine(sealcmd, seal.IndexTrackingResultHandlerAdapter(&indices, resHandler))
	if err != nil {
		t.Error(err)
	}

	if len(indices) == 0 {
		t.Fatal("Didn't parse a single index")
	}

	verifycmd, _ := verify.NewCommand([]string{indices[0]}, 1)

	err = api.StartEngine(verifycmd, resHandler)
	if err != nil {
		t.Error(err)
	}

	// ALTER FILE
	fd, err := os.OpenFile(file, os.O_WRONLY, 0777)
	if err != nil {
		t.Fatal(err)
	}
	if n, err := fd.Write([]byte("a")); err != nil || n != 1 {
		t.Fatal(err)
	}
	fd.Close()

	resHandler = testlib.ResultHandler(t, true)
	verifycmd, _ = verify.NewCommand([]string{indices[0]}, 1)
	err = api.StartEngine(verifycmd, resHandler)
	if err == nil {
		t.Error("Failed to detect file with changed byte")
	}

	os.Remove(file)
	verifycmd, _ = verify.NewCommand([]string{indices[0]}, 1)
	err = api.StartEngine(verifycmd, resHandler)
	if err == nil {
		t.Error("Failed to detect a file was removed")
	}
}
