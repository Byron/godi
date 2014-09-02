// +build !windows

package rest

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

// return a list of filesytem information entries matching the given path. We will automatically
// put it into /dir/*basename* to match everything that looks like path
// It is expected to use a leading slash to indicate the filesystem root (native for posix)
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
