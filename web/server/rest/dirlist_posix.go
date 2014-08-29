// +build !windows

package rest

import (
	"os"
	"path/filepath"
	"strings"
)

// return a list of filesytem information entries matching the given path. We will automatically
// put it into /dir/*basename* to match everything that looks like path
// It is expected to use a leading slash to indicate the filesystem root (native for posix)
func filesytemItems(path string) (out []os.FileInfo, err error) {
	path = filepath.Clean(path)
	di, err := os.Open(filepath.Dir(path))
	if err != nil {
		return
	}

	dirInfos, err := di.Readdir(-1)
	di.Close()
	if err != nil {
		return
	}

	var glob string
	if strings.HasSuffix(path, "/") {
		glob = "*"
	} else {
		glob = "*" + filepath.Base(path) + "*"
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
