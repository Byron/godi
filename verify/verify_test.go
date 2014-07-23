package verify_test

import (
	"testing"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/testlib"
	"github.com/Byron/godi/verify"
)

func TestVerify(t *testing.T) {
	datasetTree, _, _ := testlib.MakeDatasetOrPanic()
	defer testlib.RmTree(datasetTree)

	sealcmd, _ := seal.NewCommand([]string{datasetTree}, 1, 0)
	resHandler := testlib.ResultHandler(t)

	// keeps track of created indices
	var indices []string
	err := godi.StartEngine(&sealcmd, resHandler, seal.IndexTrackingResultHandlerAdapter(&indices, resHandler))
	if err != nil {
		t.Error(err)
	}

	if len(indices) == 0 {
		t.Fatal("Didn't parse a single index")
	}

	verifycmd, _ := verify.NewCommand([]string{indices[0]}, 1)

	err = godi.StartEngine(&verifycmd, resHandler, resHandler)
	if err != nil {
		t.Error(err)
	}

	// TODO(st): Alter File

	// TODO(st): Remove File
}
