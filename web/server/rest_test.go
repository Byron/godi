package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRESTState(t *testing.T) {
	s := httptest.NewServer(new(restHandler))
	defer s.Close()
	url := s.URL + "/"

	checkReq := func(req *http.Request, stat int, ct string, msg string) *http.Response {
		if res, err := http.DefaultClient.Do(req); err != nil {
			t.Fatal(err)
		} else if res.StatusCode != stat {
			t.Fatalf("Expected status %d, got %d", stat, res.StatusCode)
		} else if !strings.HasPrefix(res.Header.Get(ctkey), ct) {
			t.Fatalf("Expected content type %s, got %s", ct, res.Header.Get(ctkey))
		} else if ct == jsonct && res.ContentLength == 0 {
			t.Fatalf("Got empty json reply")
		} else {
			t.Log(msg)
			return res
		}
		panic("Shouldn't get here")
		return nil
	}

	// UNSUPPORTED METHOD
	req, _ := http.NewRequest("FOO", url, nil)
	checkReq(req, http.StatusBadRequest, "text/plain", "Correct handling of unsupported methods")

	// GET
	req, _ = http.NewRequest("GET", url, nil)
	res := checkReq(req, http.StatusOK, jsonct, "Managed to get status")

	// POST: Invalid
	req, _ = http.NewRequest("POST", url, res.Body)
	checkReq(req, http.StatusBadRequest, "text/plain", "It detected a totally invalid status")

	// POST: Valid
	// TODO

	// DELETE: abort operation

	// CHECK STATUS

}
