package server

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/verify"
)

const (
	ctkey  = "Content-Type"
	jsonct = "application/json"
)

// A struct for json serialization and deserialization
type state struct {
	runner    api.Runner // our runner in case something is going on, nil otherwise
	Mode      string     `json:mode`
	Verbosity string     `json:verbosity`
	Spid      int        `json:spid`     // streams per input device
	Spod      int        `json:spod`     // streams per output device
	Fep       []string   `json:fep`      // file exclude patterns
	Items     []string   `json:items`    // The sources (and destinations)
	OpResult  string     `json:opresult` // result of the last operation
}

// Write ourselves to w as json
func (s *state) json(w io.Writer) error {
	return json.NewEncoder(w).Encode(s)
}

// Return true if we are doing something
func (s *state) isInProgress() bool {
	return s.runner != nil
}

// instantiate ourselves from json. Automatically verifies ourselves to assure we are a valid state
func (s *state) fromJson(r io.Reader) error {
	err := json.NewDecoder(r).Decode(s)
	if err != nil {
		return err
	}

	return s.verify()
}

// Check our own consistency and return an error if something is wrong.
// Will be nil otherwise
func (s *state) verify() error {
	if s.Mode != verify.Name && s.Mode != seal.ModeCopy && s.Mode != seal.ModeSeal {
		return fmt.Errorf("Invalid mode: '%s'", s.Mode)
	}

	if _, err := api.ParseImportance(s.Verbosity); err != nil {
		return err
	}

	if s.Spid == 0 {
		return errors.New("streams per input device must be larger than 0")
	}
	if s.Mode == seal.ModeCopy && s.Spod == 0 {
		return errors.New("streams per output device must be larger than 0")
	}

	var filterErr string
	for _, f := range s.Fep {
		if _, err := api.ParseFileFilter(f); err != nil {
			if len(filterErr) > 0 {
				filterErr += "; "
			}
			filterErr += err.Error()
		}
	}

	if len(filterErr) != 0 {
		return errors.New(filterErr)
	}

	if len(s.Items) == 0 {
		return errors.New("Didn't provide a single source or destination")
	}

	return nil
}

// An http handler to manage our simplistic rest API.
// It supports only GET, POST and DELETE
// GET returns the current state and provides information about the who last changed the state
// POST allows to change the state, and returns a token indicating the owner of the status change
// DELETE stop current operations, only valid if the user matches the one that started the operation
type restHandler struct {
	st state
	o  string //IP of original requester
	l  sync.RWMutex
}

// Returns empty string on error
func remoteToHash(remoteAddr string) (string, error) {
	if host, _, err := net.SplitHostPort(remoteAddr); err != nil {
		return "", err
	} else {
		res := md5.Sum([]byte(host))
		return hex.EncodeToString(res[:]), nil
	}
}

// Set our owner to something representing the given remoteAddress.
// If a new owner is set to a value we cannot understand, we will return an error as well.
func (r *restHandler) setOwner(remoteAddr string) error {
	if newOwner, err := remoteToHash(remoteAddr); err != nil {
		return err
	} else {
		r.o = newOwner
	}

	return nil
}

// Returns true if the given remote address is our owner.
// If the remote addr is invalid for some reason, it's always false
func (r *restHandler) isOwner(remoteAddr string) bool {
	o, err := remoteToHash(remoteAddr)
	return err == nil && o == r.o
}

// Make the given state our own.
// Return error string and status code
func (r *restHandler) applyState(ns state, remoteAddr string) (string, int) {
	r.l.Lock()
	defer r.l.Unlock()

	// If we are already owned, reject any change
	if len(r.o) != 0 && !r.isOwner(remoteAddr) {
		return "State is already owned by someone else and cannot be written", http.StatusUnauthorized
	}

	// If something is in progress, request a delete
	if r.st.isInProgress() {
		return "Cannot change state if something is currently inprogress. Abort it using the DELETE method", http.StatusBadRequest
	}

	switch ns.Mode {
	case verify.Name:
		{
			ns.runner = &verify.Command{}
		}
	case seal.ModeSeal, seal.ModeCopy:
		{
			ns.runner = &seal.Command{Mode: ns.Mode}
		}
	default:
		{
			return "Unsupported mode passed state.verify() - might need a runner factory", http.StatusInternalServerError
		}
	}

	ff := make([]api.FileFilter, len(ns.Fep))
	for fid, fep := range ns.Fep {
		// verify() would have detected invalid filters, no need to check again
		ff[fid], _ = api.ParseFileFilter(fep)
	}

	level, _ := api.ParseImportance(ns.Verbosity)
	if err := ns.runner.Init(ns.Spid, ns.Spod, ns.Items, level, ff); err != nil {
		return err.Error(), http.StatusBadRequest
	}

	// So far so good, replace our state and state and start the engine
	r.setOwner(remoteAddr)
	r.st = ns
	go r.handleOperation()

	return "", http.StatusOK
}

func (r *restHandler) aggregateHandler(res api.Result) bool {
	// TODO(st) implementation
	// Communicate state through handler ... maybe this one should be controller by someone else !
	return true
}

// An async handler for an entire godi operation
func (r *restHandler) handleOperation() {
	err := api.StartEngine(r.st.runner, r.aggregateHandler)
	r.l.Lock()
	defer r.l.Unlock()

	// Reset our state and store result
	r.st.runner = nil
	r.st.OpResult = err.Error()
	r.o = ""

	// TODO: communicate we are done through a handler
}

func (r *restHandler) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set(ctkey, jsonct)

	switch {
	case rq.Method == "GET":
		{
			r.l.RLock()
			if err := r.st.json(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			r.l.RUnlock()

			// post if remote has RW access !
			w.Header().Set("X-is-RW", fmt.Sprint(len(r.o) == 0))
		}
	case rq.Method == "POST":
		{
			var newState state
			if err := newState.fromJson(rq.Body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Otherwise, it's valid and we act on it.
			if err, status := r.applyState(newState, rq.RemoteAddr); status != http.StatusOK {
				http.Error(w, err, status)
			}
		}
	case rq.Method == "DELETE":
		{
			// TODO:
		}
	default:
		{
			http.Error(w, fmt.Sprintf("Unsupported Method %s", rq.Method), http.StatusBadRequest)
		}
	}
}
