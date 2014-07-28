package io

import (
	"fmt"
)

type BytesVolume uint64

// Convert ourselves into a nice and human readable representation
// The padding string should be something compatible to the %f format, like "6.2"
func (b BytesVolume) StringPad(pad string) string {
	var divider float64
	var unit string

	switch {
	case b < BytesVolume(1024<<10):
		divider, unit = float64(1024), "KiB"
	case b < BytesVolume(1024<<20):
		divider, unit = float64(1024<<10), "MiB"
	case b < BytesVolume(1024<<30):
		divider, unit = float64(1024<<20), "GiB"
	case b < BytesVolume(1024<<40):
		divider, unit = float64(1024<<30), "TiB"
	default:
		divider, unit = float64(1024<<40), "PiB"
	} // end switch

	return fmt.Sprintf(fmt.Sprintf("%%%sf%%s", pad), float64(b)/divider, unit)
}

func (b BytesVolume) String() string {
	return b.StringPad("6.2")
}
