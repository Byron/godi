package testlib

import (
	"testing"

	"github.com/Byron/godi/api"
)

// Return a function makes tests fail if there was an error result during an operation
func ResultHandler(t *testing.T) func(godi.Result) {
	return func(res godi.Result) {
		if res.Error() != nil {
			t.Error(res.Error())
		} else {
			t.Log(res.Info())
		}
	} // end resHandler
}
