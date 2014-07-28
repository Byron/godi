// +build !windows

package io

import (
	"os"
	"syscall"
)

// deviceMap maps the given paths to their device ids, effectively grouping them by device.
// We use a simple array for this as actual device IDs are not relevant
func DeviceMap(paths []string) [][]string {
	const defaultDevice uint64 = 0
	m := make(map[uint64][]string)

	for _, path := range paths {
		fi, ferr := os.Stat(path)
		did := defaultDevice
		if ferr == nil {
			st, ok := fi.Sys().(*syscall.Stat_t)
			if ok {
				did = uint64(st.Dev)
			}
		}
		if p, ok := m[did]; ok {
			m[did] = append(p, path)
		} else {
			p = make([]string, 1)
			p[0] = path
			m[did] = p
		}
	}

	res := make([][]string, len(m))
	c := 0
	for _, trees := range m {
		res[c] = trees
		c += 1
	}

	return res
}
