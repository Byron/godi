package verify_test

import (
	"testing"

	"github.com/Byron/godi/testlib"
)

func TestVerify(t *testing.T) {
	datasetTree, _, _ := testlib.MakeDatasetOrPanic()
	defer testlib.RmTree(datasetTree)

	t.Fail()
}
