package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/testlib"
)

const (
	plain = plainct
	delay = 50 * time.Millisecond
)

func TestRESTState(t *testing.T) {
	srv := httptest.NewServer(new(restHandler))
	defer srv.Close()
	url := srv.URL + "/"

	checkReq := func(req *http.Request, stat int, ct string, msg string) *http.Response {
		if res, err := http.DefaultClient.Do(req); err != nil {
			t.Fatal(err)
		} else if res.StatusCode != stat {
			body, _ := ioutil.ReadAll(res.Body)
			t.Fatalf("Expected status %d, got %d(%s): %s", stat, res.StatusCode, http.StatusText(res.StatusCode), string(body))
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
	checkReq(req, http.StatusBadRequest, plain, "Correct handling of unsupported methods")

	// GET
	req, _ = http.NewRequest("GET", url, nil)
	res := checkReq(req, http.StatusOK, jsonct, "Managed to get status")
	if res.Header.Get(isrwparam) != "true" {
		t.Fatalf("Unexpected RW value: '%v'", res.Header.Get(isrwparam))
	}

	// POST: Invalid state makes us fail the precondition
	req, _ = http.NewRequest("POST", url, res.Body)
	checkReq(req, http.StatusPreconditionFailed, plain, "We didn't modify anything yet, and don't own the state")

	datasetTree, _, _ := testlib.MakeDatasetOrPanic()
	defer testlib.RmTree(datasetTree)

	// Make a change
	ns := state{
		Mode:      seal.ModeSeal,
		Verbosity: api.Info.String(),
		Spid:      1,
		Sources: []string{
			datasetTree,
		},
		// INVALID FORMAT !
		Format: "FOO",
	}

	convertJson := func(s state, w io.WriteCloser) {
		go func() { s.json(w); w.Close() }()
	}

	// PUT: invalid
	r, w := io.Pipe()
	convertJson(ns, w)
	req, _ = http.NewRequest("PUT", url, r)
	checkReq(req, http.StatusPreconditionFailed, plain, "State should be unchanged")

	// DELETE without operation triggers error
	req, _ = http.NewRequest("DELETE", url, nil)
	checkReq(req, http.StatusPreconditionFailed, plain, "DELETE without operation triggers error")

	// PUT: valid
	ns.Format = codec.GobName
	r, w = io.Pipe()
	convertJson(ns, w)
	req, _ = http.NewRequest("PUT", url, r)
	res = checkReq(req, http.StatusOK, jsonct, "Should have changed the state")

	// quick comparison, ns should actually be the same. Can't compare directly though
	var s state
	if err := s.fromJson(res.Body); err != nil {
		t.Fatal(err)
	}
	if s.Format != ns.Format || s.Mode != ns.Mode || s.LastError != "" {
		t.Fatal("Unexpected format or mode, or there is a result")
	}

	// POST: Valid - empty state
	req, _ = http.NewRequest("POST", url, nil)
	res = checkReq(req, http.StatusOK, jsonct, "Should have set the machine in motion")

	// Can't change while it's going. It shouldn't change the state in that case either (something we don't check here)
	for _, m := range []string{"POST", "PUT"} {
		r, w = io.Pipe()
		convertJson(s, w)
		req, _ = http.NewRequest(m, url, r)
		checkReq(req, http.StatusPreconditionFailed, plain, fmt.Sprintf("Can't %s while operation is in progress", m))
	}

	// DELETE: abort operation - idempotent
	for i := 0; i < 2; i++ {
		req, _ = http.NewRequest("DELETE", url, nil)
		checkReq(req, http.StatusOK, plain, fmt.Sprintf("Can abort running operation attempt %d", i))
	}

	// CHECK STATUS - have to wait for it to finish (TODO: wait for websocket notification)
	s.IsRunning = true
	startedAt := time.Now()
	for s.IsRunning {
		req, _ = http.NewRequest("GET", url, nil)
		res = checkReq(req, http.StatusOK, jsonct, "Can get status after operation was cancelled")
		s.fromJson(res.Body)
		time.Sleep(delay)
	}
	if time.Now().Sub(startedAt) < delay {
		t.Fatal("Finished way to early - did it run for sure ?")
	}
	if s.LastError == "" {
		t.Fatal("Cancellation should turn out as 'Error'")
	} else {
		t.Log("Cancellation created named: %s", s.LastError)
	}
}
