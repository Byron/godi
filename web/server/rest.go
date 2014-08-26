package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/verify"
)

const (
	ctkey     = "Content-Type"
	jsonct    = "application/json"
	plainct   = "text/plain"
	isrwparam = "X-is-RW"
	// maxInactivity = 5 * time.Minute
	maxInactivity = 5 * time.Second
)

// A struct for json serialization and deserialization
type state struct {
	Mode         string   `json:"mode"`
	Verbosity    string   `json:"verbosity"`
	Spid         int      `json:"spid"`         // streams per input device
	Spod         int      `json:"spod"`         // streams per output device
	Fep          []string `json:"fep"`          // file exclude patterns
	Sources      []string `json:"sources"`      // The sources for verify and seal
	Destinations []string `json:"destinations"` // The destinations of sealed-copy
	Verify       string   `json:"verify"`       // if non-empty, verification is done after a sealed copy
	Format       string   `json:"format"`       // The serialization format of seals
	SocketURL    string   `json:"socketURL"`    // read-only URL of the web-socket people can connect to
	IsRunning    bool     `json:"status"`       // read-only, true if an operation is in progress

	LastError string `json:executionError` // error result of the last operation
}

// A struct keeping valid values for certain constants
type defaults struct {
	Modes       []string `json:"modes"`
	Verbosities []string `json:"verbosities"`
	Feps        []string `json:"feps"`
	Formats     []string `json:"formats"`
}

var valueDefaults = defaults{
	Modes:       []string{verify.Name, seal.ModeSeal, seal.ModeCopy},
	Verbosities: []string{api.Info.String(), api.Error.String()},
	Feps:        []string{api.FilterHidden.String(), api.FilterSeals.String(), api.FilterSymlinks.String(), api.FilterVolatile.String()},
	Formats:     codec.Names(),
}

// Write ourselves to w as json
func (s *state) json(w io.Writer) error {
	return json.NewEncoder(w).Encode(s)
}

// instantiate ourselves from json. Automatically verifies ourselves to assure we are a valid state
func (s *state) fromJson(r io.Reader) error {
	err := json.NewDecoder(r).Decode(s)
	if err != nil {
		return err
	}

	return s.verify(false)
}

// Check our own consistency and return an error if something is wrong.
// Will be nil otherwise.
// All values must be set if 'non-null-only' is true
func (s *state) verify(checkNullValue bool) error {
	if (len(s.Mode) > 0 || checkNullValue) && s.Mode != verify.Name && s.Mode != seal.ModeCopy && s.Mode != seal.ModeSeal {
		return fmt.Errorf("Invalid mode: '%s'", s.Mode)
	}

	if len(s.Verbosity) > 0 || checkNullValue {
		if _, err := api.ParseImportance(s.Verbosity); err != nil {
			return err
		}
	}

	if checkNullValue && s.Spid == 0 {
		return errors.New("streams per input device must be larger than 0")
	}

	if s.Mode == seal.ModeCopy {
		if checkNullValue && s.Spod == 0 {
			return errors.New("streams per output device must be larger than 0")
		}
		if checkNullValue && len(s.Destinations) == 0 {
			return errors.New("Need to provide at least one destination")
		}
	}

	if s.Mode != verify.Name {
		// defaults to something useful, so it's optional here
		if len(s.Format) > 0 && codec.NewByName(s.Format) == nil {
			return fmt.Errorf("invalid output format: '%s'", s.Format)
		}
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

	if checkNullValue && len(s.Sources) == 0 {
		return errors.New("Didn't provide a single source")
	}

	return nil
}

// Apply the given, possibly partial state, ignoring all null values
// Return an error if one of the new values is invalid. Boolean signals that we actually changed
func (s *state) apply(o state) (bool, error) {
	ns := *s
	changed := false

	// Set each value if it is not null
	if len(o.Mode) > 0 {
		changed = true
		ns.Mode = o.Mode
	}
	if len(o.Verbosity) > 0 {
		changed = true
		ns.Verbosity = o.Verbosity
	}
	if o.Spid > 0 {
		changed = true
		ns.Spid = o.Spid
	}
	if o.Spod > 0 {
		changed = true
		ns.Spod = o.Spod
	}
	if len(o.Fep) > 0 {
		changed = true
		ns.Fep = o.Fep
	}
	if len(o.Sources) > 0 {
		changed = true
		ns.Sources = o.Sources
	}
	if len(o.Destinations) > 0 {
		changed = true
		ns.Destinations = o.Destinations
	}
	if len(o.Verify) > 0 {
		changed = true
		ns.Verify = o.Verify
	}
	if len(o.Format) > 0 {
		changed = true
		ns.Format = o.Format
	}

	// Actually this shouldn't be needed here as o is already checked, but better save than sorry
	if err := ns.verify(false); err != nil {
		return changed, err
	}

	// finally replace our own values
	*s = ns
	return changed, nil
}

// An http handler to manage our simplistic rest API.
// It supports only GET, POST and DELETE
// GET returns the current state and provides information about the who last changed the state
// POST allows to change the state, and returns a token indicating the owner of the status change
// DELETE stop current operations, only valid if the user matches the one that started the operation
type restHandler struct {
	st              state
	r               api.Runner // our runner in case something is going on, nil otherwise
	o               string     //IP of original requester
	cancelRequested bool
	l               sync.RWMutex
	cb              func(bool, bool, api.Result, string) // a callback to allow others to stay informed
	lmat            time.Time                            // time at which we were modified
}

// Returns a usable REST handler.
// The callback allows to respond to state-changes, calls to it are synchronized and thus serial only.
// f(isEnd, result) - isEnd is True only when the operation is now finished, result is nil in that case
// If isEnd is false and result is nil, this indicates that our state changed
func NewRestHandler(onStateChange func(bool, bool, api.Result, string), socketURL string) http.Handler {
	if onStateChange == nil {
		panic("Callback must be set")
	}
	handler := restHandler{
		cb: onStateChange,
		st: state{
			// replicate defaults used by CLI
			Mode:      verify.Name,
			Spid:      1,
			Fep:       []string{api.FilterVolatile.String()},
			Format:    codec.GobName,
			Verbosity: api.Error.String(),
			SocketURL: socketURL,
		},
		lmat: time.Now(),
	}

	// assure that we reset the owner after inactivity
	go handler.handleInactivity()

	return &handler
}

// Returns empty string on error
func remoteToOwner(remoteAddr string) (string, error) {
	if host, _, err := net.SplitHostPort(remoteAddr); err != nil {
		return "", err
	} else {
		return host, nil
	}
}

// Make sure we reset the owner after inactivity
// Needs to be run as go-routine
func (r *restHandler) handleInactivity() {
	for tick := range time.Tick(1 * time.Second) {
		if tick.Sub(r.lmat) >= maxInactivity && len(r.o) != 0 {
			r.l.Lock()
			r.setOwner("")
			r.lmat = tick
			// signal the change
			r.cb(false, false, nil, "")
			r.l.Unlock()
		}
	} // every second
}

// Return true if we are doing something
func (r *restHandler) isInProgress() bool {
	return r.r != nil
}

// Set our owner to something representing the given remoteAddress. Using "" will reset the owner.
// If a new owner is set to a value we cannot understand, we will return an error as well.
func (r *restHandler) setOwner(remoteAddr string) error {
	if remoteAddr == "" {
		r.o = ""
		return nil
	}

	if newOwner, err := remoteToOwner(remoteAddr); err != nil {
		return err
	} else {
		r.o = newOwner
	}

	return nil
}

// Returns true if the given remote address is our owner.
// If the remote addr is invalid for some reason, it's always false
func (r *restHandler) isOwner(remoteAddr string) bool {
	o, err := remoteToOwner(remoteAddr)
	return err == nil && o == r.o
}

// Returns true if the given remoteAddr may change our state
func (r *restHandler) CanWrite(remoteAddr string) bool {
	return len(r.o) == 0 || r.isOwner(remoteAddr)
}

// Execute the state we currently have
// Return error string and status code. If there is no error, status will be StatusOK
// The new state may be partial, and contain zero values which will just be ignored
func (r *restHandler) execute(remoteAddr string) (string, int) {
	r.l.Lock()
	defer r.l.Unlock()

	// If something is in progress, request a delete
	if r.isInProgress() {
		return "Cannot change state if something is currently inprogress. Abort it using the DELETE method", http.StatusPreconditionFailed
	}

	if err := r.st.verify(true); err != nil {
		return err.Error(), http.StatusPreconditionFailed
	}

	// if we are not owned anymore, lets take ownership
	if len(r.o) == 0 {
		r.setOwner(remoteAddr)
	}

	// At this point, we expect remoteAddr to be the owner
	if !r.isOwner(remoteAddr) {
		return "You must own the state to be allowed to execute it", http.StatusUnauthorized
	}

	items := r.st.Sources
	switch r.st.Mode {
	case verify.Name:
		{
			r.r = &verify.Command{}
		}
	case seal.ModeSeal, seal.ModeCopy:
		{
			scmd := seal.Command{Mode: r.st.Mode}
			r.r = &scmd
			if r.st.Mode == seal.ModeCopy {
				items = append(items, seal.Sep)
				items = append(items, r.st.Destinations...)

				scmd.Verify = len(r.st.Verify) > 0
			}

			scmd.Format = r.st.Format
		}
	default:
		{
			return "Unsupported mode passed state.verify() - might need a runner factory", http.StatusInternalServerError
		}
	}

	ff := make([]api.FileFilter, len(r.st.Fep))
	for fid, fep := range r.st.Fep {
		// verify() would have detected invalid filters, no need to check again
		ff[fid], _ = api.ParseFileFilter(fep)
	}

	level, _ := api.ParseImportance(r.st.Verbosity)
	if err := r.r.Init(r.st.Spid, r.st.Spod, items, level, ff); err != nil {
		// This shouldn't be happening here - we verified beforehand.
		return err.Error(), http.StatusInternalServerError
	}

	// Clear previous results, we are getting a new one
	r.st.LastError = ""
	r.st.IsRunning = true
	go r.handleOperation(r.r)
	return "", http.StatusOK
}

// Apply the given (partial) state to our own.
// On error, return a string and a non-OK status code
func (r *restHandler) applyState(ns state, remoteAddr, clientID string) (string, int) {
	r.l.Lock()
	defer r.l.Unlock()
	changed := false

	if !r.CanWrite(remoteAddr) {
		return "Changing values not permitted", http.StatusUnauthorized
	}

	if r.isInProgress() {
		return "Operation already in progress", http.StatusPreconditionFailed
	}

	var err error
	if changed, err = r.st.apply(ns); err != nil {
		return err.Error(), http.StatusPreconditionFailed
	}

	// make sure we allow remoteAddr to take ownership
	if err = r.setOwner(remoteAddr); err != nil {
		return err.Error(), http.StatusBadRequest
	}

	// Inform people about the change
	if changed {
		r.lmat = time.Now()
		r.cb(false, false, nil, clientID)
	}
	return "", http.StatusOK
}

// An async handler for an entire godi operation
func (r *restHandler) handleOperation(runner api.Runner) {

	// We will filter out TimedStatistics here
	r.cb(true, false, nil, "")
	err := api.StartEngine(runner, func(res api.Result) {
		r.l.Lock()
		r.lmat = time.Now()
		r.cb(false, false, res, "")
		r.l.Unlock()
	})
	r.l.Lock()

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	{
		// Reset our state and store result. Reset owner
		r.r = nil
		r.st.LastError = errMsg
		r.st.IsRunning = false
		r.cancelRequested = false
		r.o = ""
	}

	r.l.Unlock()

	// inform that we are done - we do that after our state has changed just
	// to assure rest calls will definitely have the expected result
	r.cb(false, true, nil, "")
}

func (r *restHandler) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set(ctkey, jsonct)

	// UTILITIES
	/////////////
	doWriteState := func() {
		r.l.RLock()
		r.st.json(w)
		r.l.RUnlock()
	}

	doPut := func(writeState bool) (res bool) {
		var newState state
		if err := newState.fromJson(rq.Body); err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
		} else if msg, status := r.applyState(newState, rq.RemoteAddr, rq.Header.Get("Client-ID")); status != http.StatusOK {
			http.Error(w, msg, status)
		} else {
			// on success, return the current state
			if writeState {
				doWriteState()
			}
			res = true
		}
		return
	}

	// METHOD HANDLER
	///////////////////
	switch rq.Method {
	case "GET":
		{
			r.l.RLock()
			{
				// post if remote has RW access !
				w.Header().Set(isrwparam, fmt.Sprint(r.CanWrite(rq.RemoteAddr)))
				if err := r.st.json(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			r.l.RUnlock()
		}
	case "PUT":
		{
			doPut(true)
		}
	case "POST":
		{
			// For convenience, we allow PUT-like behaviour, in case people post all at once
			// Also saves a request
			if !doPut(false) {
				return
			}

			// Execute the current state
			if msg, status := r.execute(rq.RemoteAddr); status != http.StatusOK {
				http.Error(w, msg, status)
			} else {
				doWriteState()
			}
		}
	case "DEFAULTS":
		{
			if err := json.NewEncoder(w).Encode(&valueDefaults); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	case "DELETE":
		{
			r.l.Lock()
			{
				if !r.isInProgress() {
					http.Error(w, "No operation is currently in progress", http.StatusPreconditionFailed)
				} else if !r.cancelRequested {
					// we are idempotent if the delete operation is called multiple times, waiting for the operation
					// to actually shut down
					close(r.r.CancelChannel())
					r.cancelRequested = true
				}
				w.Header().Set(ctkey, plainct)
			}
			r.l.Unlock()
		}
	default:
		{
			http.Error(w, fmt.Sprintf("Unsupported Method %s", rq.Method), http.StatusBadRequest)
		}
	}
}
