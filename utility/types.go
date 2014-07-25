package utility

import (
	"fmt"
	"sync/atomic"
	"time"
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

// CopyTo will atomically copy our fields to the destination structure. It will just read the fields atomically, and
// write it using standard means
func (s *Stats) CopyTo(d *Stats) {
	d.TotalFilesRead = atomic.LoadUint32(&s.TotalFilesRead)
	d.TotalFilesWritten = atomic.LoadUint32(&s.TotalFilesWritten)

	d.FilesBeingRead = atomic.LoadUint32(&s.FilesBeingRead)
	d.FilesBeingWritten = atomic.LoadUint32(&s.FilesBeingWritten)

	d.BytesRead = atomic.LoadUint64(&s.BytesRead)
	d.BytesWritten = atomic.LoadUint64(&s.BytesWritten)
	d.BytesHashed = atomic.LoadUint64(&s.BytesHashed)

	d.NumHashers = atomic.LoadUint32(&s.NumHashers)
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

// Prints itself as a single line full of useful information, including deltas of relevant metrics as compared
// to theh last state d. You will also give the temporal distance which separates this stat from the previous one
func (s *Stats) DeltaString(d *Stats, td time.Duration) string {
	itf := atomic.LoadUint32(&s.TotalFilesRead)
	ibr := atomic.LoadUint64(&s.BytesRead)
	itfd := itf - d.TotalFilesRead // difference of files
	ibrd := ibr - d.BytesRead      // difference of bytes

	out := fmt.Sprintf("%d->IN (#%d(%+d) ⌰%s Δ%s/%.2fs)",
		atomic.LoadUint32(&s.FilesBeingRead),
		itf,
		itfd,
		BytesVolume(ibr),
		BytesVolume(ibrd),
		td.Seconds(),
	)

	if s.FilesBeingWritten > 0 {
		otf := atomic.LoadUint32(&s.TotalFilesWritten)
		obw := atomic.LoadUint64(&s.BytesWritten)
		otfd := otf - d.TotalFilesWritten
		obwd := obw - d.BytesWritten

		out += fmt.Sprintf("\t%d OUT->(#%d(%+d) ⌰%s Δ%s/%.2fs)",
			atomic.LoadUint32(&s.FilesBeingWritten),
			otf,
			otfd,
			BytesVolume(obw),
			BytesVolume(obwd),
			td.Seconds(),
		)
	}

	bh := atomic.LoadUint64(&s.BytesHashed)
	bhd := bh - d.BytesHashed

	out += fmt.Sprintf("\t%d HASH(⌰%s Δ%s/%.2fs)",
		atomic.LoadUint32(&s.NumHashers),
		BytesVolume(bh),
		BytesVolume(bhd),
		td.Seconds(),
	)
	return out
}
