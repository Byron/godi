// +build !windows

package utility

import (
	"os"
	"syscall"
)

// deviceMap maps the given paths to their device ids, effectively grouping them by device.
func DeviceMap(paths []string) map[uint64][]string {
	const defaultDevice uint64 = 0
	res := make(map[uint64][]string)

	for _, path := range paths {
		fi, ferr := os.Stat(path)
		did := defaultDevice
		if ferr == nil {
			st, ok := fi.Sys().(*syscall.Stat_t)
			if ok {
				did = uint64(st.Dev)
			}
		}
		if p, ok := res[did]; ok {
			res[did] = append(p, path)
		} else {
			p = make([]string, 1)
			p[0] = path
			res[did] = p
		}
	}

	return res
}
