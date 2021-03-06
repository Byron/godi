package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Byron/godi/api"
)

const (
	QPath = "path"
	QType = "type"

	TypeAll  = "all"
	TypeSeal = "sealOnly"
)

type dirHandler struct {
	stp func() State // a provider for the current state
	le  string       // last error string
}

// Simple json compatible structure which identifies an item within a directory
// It also signals if it is a file
type ItemInfo struct {
	Item  string `json:"item"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
}

func (i *ItemInfo) fromFileInfo(tree string, fi os.FileInfo) {
	i.Item = fi.Name()
	i.Path = filepath.Join(tree, i.Item)
	i.IsDir = fi.IsDir()
}

// Write the given FileInfos in a suitable format to a the given writer
func fileInfoToJson(tree string, fis []os.FileInfo, w io.Writer) error {
	infos := make([]ItemInfo, len(fis))
	for fid := range fis {
		infos[fid].fromFileInfo(tree, fis[fid])
	}
	return json.NewEncoder(w).Encode(infos)
}

// Return a list of file-info objects which have not been excluded by our filters
func (d *dirHandler) filter(fis []os.FileInfo, sealOnly bool) (out []os.FileInfo, err error) {

toNextFile:
	for _, fi := range fis {
		if sealOnly && !fi.IsDir() && !api.FilterSeals.Matches(fi.Name(), fi.Mode()) {
			continue toNextFile
		}

		for _, fname := range d.stp().Fep {
			// don't filter seals even though they are explicitly desired
			if !sealOnly && fname == api.FilterSeals.String() {
				continue
			}

			excludeFilter, err := api.ParseFileFilter(fname)
			if err != nil {
				// invalid filters  shouldn't be here in the first place.
				// Abort !
				return out, err
			}

			if excludeFilter.Matches(fi.Name(), fi.Mode()) {
				continue toNextFile
			}
		} // for each filter in current state

		out = append(out, fi)
	} // for each file-info

	return
}

func NewDirHandler(stateProvider func() State) http.Handler {
	return &dirHandler{stp: stateProvider}
}

// return a list of filesytem information entries matching the given path. We will automatically
// put it into /dir/*basename* to match everything that looks like path
// It is expected to use a leading slash to indicate the filesystem root (native for posix)
// NOTE: This function could need a special implementation on windows to allow listing all
// available drives. Otherwise you will have to know the available drives.
func filesytemItems(path string) (out []os.FileInfo, err error) {
	endsWithSep := strings.HasSuffix(path, string(filepath.Separator))
	path = filepath.Clean(path)

	dir := path
	if !endsWithSep {
		dir = filepath.Dir(path)
	}

	di, err := os.Open(dir)
	if err != nil {
		return
	}

	dirInfos, err := di.Readdir(-1)
	di.Close()
	if err != nil {
		return
	}

	glob := "*"
	if !endsWithSep {
		bs := filepath.Base(path)
		glob = ""
		for len(bs) > 0 {
			r, size := utf8.DecodeRuneInString(bs)
			glob += "*[" + string(unicode.ToUpper(r)) + string(unicode.ToLower(r)) + "]*"
			bs = bs[size:]
		}
	}

	for _, fi := range dirInfos {
		if matched, merr := filepath.Match(glob, fi.Name()); matched {
			out = append(out, fi)
		} else if merr != nil {
			return out, merr
		}
	}

	return
}

func (r *dirHandler) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	if rq.Method != "GET" {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	// VERIFY QUERY
	///////////////
	q := rq.URL.Query()
	if len(q.Get(QPath)) == 0 || len(q.Get(QType)) == 0 {
		http.Error(w, fmt.Sprintf("You have to specify the '%s' and '%s' within the query string", QPath, QType), http.StatusBadRequest)
		return
	}

	qt := q.Get(QType)
	qp := q.Get(QPath)
	sealOnly := qt == TypeSeal
	if qt != TypeSeal && qt != TypeAll {
		http.Error(w, fmt.Sprintf("Invalid request type, expected one of '%s', '%s'", TypeAll, TypeSeal), http.StatusBadRequest)
		return
	}

	// PRODUCE RESULT
	//////////////////
	if fis, err := filesytemItems(qp); err == nil {
		if fis, err = r.filter(fis, sealOnly); err == nil {
			w.Header().Set(ContentKey, JsonContent)
			if err = fileInfoToJson(filepath.Dir(qp), fis, w); err != nil {
				w.Header().Set(ContentKey, PlainContent)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, fmt.Sprintf("Problem filtering directory listing: %v", err), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, fmt.Sprintf("Failed to read path at '%s': %v", qp, err), http.StatusBadRequest)
		return
	}
}
