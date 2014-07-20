package verify_test

import (
	"runtime"
	"testing"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/testlib"
	"github.com/Byron/godi/verify"
)

func TestVerify(t *testing.T) {
	datasetTree, _, _ := testlib.MakeDatasetOrPanic()
	defer testlib.RmTree(datasetTree)

	sealcmd, _ := seal.NewCommand([]string{datasetTree}, 1)
	resHandler := testlib.ResultHandler(t)

	// keeps track of created indices
	var index string
	aggHandler := func(res godi.Result) {
		resHandler(res)
		if res == nil || res.FileInformation() == nil {
			return
		}
		if res.FileInformation().Size < 0 {
			if len(index) > 0 {
				t.Fatal("Can only keep one index right now")
			}
			index = res.FileInformation().Path
		}
	}

	var maxProcs uint = uint(runtime.GOMAXPROCS(0))
	err := godi.StartEngine(&sealcmd, maxProcs, resHandler, aggHandler)
	if err != nil {
		t.Error(err)
	}

	if len(index) == 0 {
		t.Fatal("Didn't parse a single index")
	}

	verifycmd, _ := verify.NewCommand([]string{index}, 1)

	err = godi.StartEngine(&verifycmd, maxProcs, resHandler, resHandler)
	if err != nil {
		t.Error(err)
	}
}
