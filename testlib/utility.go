package testlib

import (
	"testing"

	"github.com/Byron/godi/api"
)

// Return a function makes tests fail if there was an error result during an operation
// If logOnly is True, errors will not make the test fail. Useful if errors are expected.
func ResultHandler(t *testing.T, logOnly bool) func(api.Result) {
	return func(res api.Result) {
		if res.Error() != nil {
			if logOnly {
				t.Log(res.Info())
			} else {
				t.Error(res.Info())
			}
		} else {
			t.Log(res.Info())
		}
	} // end resHandler
}
