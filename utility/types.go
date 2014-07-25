package utility

import (
	"fmt"
)

// A structure to keep information about what is currently going on.
// It is means to be used as shared resource, used by multiple threads, which is why
// thread-safe counters are used.
// Implementations must keep these numbers up-to-date, while async processors will digest
// and present the data in some form
type Stats struct {
	TotalFilesRead    uint32 // Amount of whole files we read so far
	TotalFilesWritten uint32 // Amount of whole files we wrote so far
	FilesBeingRead    uint32 // Amount of files currently being read
	FilesBeingWritten uint32 // Amount of files currently being written
	BytesRead         uint64 // Total of bytes read so far, counting all input streams
	BytesWritten      uint64 // Total of bytes written so far, counting all output streams
	BytesHashed       uint64 // Total of bytes hashed so far, counting all active hashers
	NumHashers        uint32 // Amount of hashers running in parallel
}

type BytesVolume uint64

// Convert ourselves into a nice and human readable representation
func (b BytesVolume) String() string {
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

	return fmt.Sprintf("%.2f%s", float64(b)/divider, unit)
}

// Prints itself as a single line full of useful information
func (s *Stats) String() string {
	out := fmt.Sprintf("%d IN(#%d with %s)", s.FilesBeingRead, s.TotalFilesRead, BytesVolume(s.BytesRead))

	if s.FilesBeingWritten > 0 {
		out += fmt.Sprintf("\t%d OUT(#%d with %s)", s.FilesBeingWritten, s.TotalFilesWritten, BytesVolume(s.BytesWritten))
	}

	out += fmt.Sprintf("\t%d HASH(%s)", s.NumHashers, BytesVolume(s.BytesHashed))
	return out
}
