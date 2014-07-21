package cli_test

import (
	"testing"

	"github.com/Byron/godi/cli/app"
)

func TestSealParsing(t *testing.T) {
	if err := app.NewGodiApp().Run([]string{"godi", "seal", "foo", "bar"}); err == nil {
		t.Error("Didn't notice that input is garbage - should have required accessible directories")
	} else {
		t.Log(err)
	}

}
