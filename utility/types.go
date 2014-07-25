package utility

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	StatsClientSep = " | "
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

// MostFiles returns the greatest number of files, either the one that were read, or the ones that were written
func (s *Stats) MostFiles() uint32 {
	if s.TotalFilesRead > s.TotalFilesWritten {
		return s.TotalFilesRead
	}
	return s.TotalFilesWritten
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
// If d didn't change in a particular field, we will assume the user wants to print the speed per second so far,
// and td should have the respective value of the total program duration
// Sep is the separator to use between fields
func (s *Stats) DeltaString(d *Stats, td time.Duration, sep string) string {
	intDelta := func(cur, prev uint32) string {
		if prev == 0 {
			prev = cur
		}

		val := cur
		if prev != cur {
			val = cur - prev
		}
		return fmt.Sprintf(" #Δ%d/s", uint64(float64(val)/td.Seconds()))
	}

	bytesDelta := func(cur, prev uint64) string {
		if prev == 0 {
			prev = cur
		}
		val := cur
		if prev != cur {
			val = cur - prev
		}
		return fmt.Sprintf(" Δ%s/s", BytesVolume(float64(val)/td.Seconds()))
	}

	inOut := func(cur uint32) string {
		if cur == 0 {
			return ""
		}
		return fmt.Sprintf("%d ", cur)
	}

	if len(sep) == 0 {
		sep = "\t"
	}

	itf := atomic.LoadUint32(&s.TotalFilesRead)
	ibr := atomic.LoadUint64(&s.BytesRead)

	out := fmt.Sprintf("%s->IN #%d%s ⌰%s%s",
		inOut(atomic.LoadUint32(&s.FilesBeingRead)),
		itf,
		intDelta(itf, d.TotalFilesRead),
		BytesVolume(ibr),
		bytesDelta(ibr, d.BytesRead),
	)

	bh := atomic.LoadUint64(&s.BytesHashed)

	out += fmt.Sprintf("%s%sHASH ⌰%s%s",
		sep,
		inOut(atomic.LoadUint32(&s.NumHashers)),
		BytesVolume(bh),
		bytesDelta(bh, d.BytesHashed),
	)

	if s.TotalFilesWritten > 0 || s.FilesBeingWritten > 0 {
		otf := atomic.LoadUint32(&s.TotalFilesWritten)
		obw := atomic.LoadUint64(&s.BytesWritten)

		out += fmt.Sprintf("%s%sOUT->#%d%s ⌰%s%s",
			sep,
			inOut(atomic.LoadUint32(&s.FilesBeingWritten)),
			otf,
			intDelta(otf, d.TotalFilesWritten),
			BytesVolume(obw),
			bytesDelta(obw, d.BytesWritten),
		)
	}

	return out
}
