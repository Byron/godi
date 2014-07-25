package testlib

import (
	"testing"

	"github.com/Byron/godi/api"
)

// Return a function makes tests fail if there was an error result during an operation
func ResultHandler(t *testing.T) func(api.Result) bool {
	return func(res api.Result) bool {
		if res.Error() != nil {
			t.Error(res.Error())
		} else {
			t.Log(res.Info())
		}
		return true
	} // end resHandler
}
