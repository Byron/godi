package utility

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

const (
	StatsClientSep = " | "
)

// Append elm if it is not yet on dest
func AppendUniqueString(dest []string, elm string) []string {
	for _, d := range dest {
		if d == elm {
			return dest
		}
	}
	return append(dest, elm)
}

// A structure to keep information about what is currently going on.
// It is means to be used as shared resource, used by multiple threads, which is why
// thread-safe counters are used.
// Implementations must keep these numbers up-to-date, while async processors will digest
// and present the data in some form
// NOTE: Even though not
type Stats struct {
	// PERFORMANCE METRICS
	TotalFilesRead    uint32 // Amount of whole files we read so far
	TotalFilesWritten uint32 // Amount of whole files we wrote so far
	FilesBeingRead    uint32 // Amount of files currently being read
	FilesBeingWritten uint32 // Amount of files currently being written
	BytesRead         uint64 // Total of bytes read so far, counting all input streams
	BytesWritten      uint64 // Total of bytes written so far, counting all output streams
	BytesHashed       uint64 // Total of bytes hashed so far, counting all active hashers
	NumHashers        uint32 // Amount of hashers running in parallel

	// GENERAL INFORMATION
	StartedAt       time.Time // The time at which we started processing
	NumSkippedFiles uint32    // Amount of files we skipped right away
	StopTheEngines  uint32    // Amount of gather procs which had write errors on all destinations

	// AGGREGATION
	// Aggregation step is single-threaded - no atomic operation needed
	ErrCount       uint // Amount of errors that hit the aggregation step
	NumUndoneFiles uint // Amout of files removed during undo
	WasCancelled   bool // is true if the user cancelled
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

	d.StartedAt = s.StartedAt
	d.NumSkippedFiles = atomic.LoadUint32(&s.NumSkippedFiles)
	d.StopTheEngines = atomic.LoadUint32(&s.StopTheEngines)

	// Agg variables don't need to be atomic - we copy them here for completeness only
	d.ErrCount = s.ErrCount
	d.NumUndoneFiles = s.NumUndoneFiles
	d.WasCancelled = s.WasCancelled
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

	return fmt.Sprintf("%6.2f%s", float64(b)/divider, unit)
}

// Prints performance metrics as a single line full of useful information, including deltas of relevant metrics as compared
// to the last state d. You will also give the temporal distance which separates this stat from the previous one
// If you pass s as d, this indicates a result mode, which assumes you want the overall average throughput
// Sep is the separator to use between fields
func (s *Stats) DeltaString(d *Stats, td time.Duration, sep string) string {
	resultMode := d == s

	intDelta := func(cur, prev uint32) string {
		if prev == cur {
			if resultMode {
				prev = 0
			} else {
				return ""
			}
		}
		return fmt.Sprintf(" #Δ%5d/s", uint64(float64(cur-prev)/td.Seconds()))
	}

	bytesDelta := func(cur, prev uint64) string {
		if prev == cur {
			if resultMode {
				prev = 0
			} else {
				return ""
			}
		}
		return fmt.Sprintf(" Δ%s/s", BytesVolume(float64(cur-prev)/td.Seconds()))
	}

	inOut := func(cur uint32) string {
		if cur == 0 {
			return "  "
		}
		return fmt.Sprintf("%d ", cur)
	}

	if len(sep) == 0 {
		sep = "\t"
	}

	itf := atomic.LoadUint32(&s.TotalFilesRead)
	ibr := atomic.LoadUint64(&s.BytesRead)

	timeStr := fmt.Sprintf("%4.0fs", s.Elapsed().Seconds())
	if resultMode {
		timeStr = s.Elapsed().String()
	}
	out := fmt.Sprintf("🕑  %s%s", timeStr, sep)

	out += fmt.Sprintf("%s->READ #%d%s ⌰%s%s",
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

		out += fmt.Sprintf("%s%sWRITE #%d%s ⌰%s%s",
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

// Return amount of time elapsed since we started the operation
func (s *Stats) Elapsed() time.Duration {
	return time.Now().Sub(s.StartedAt)
}

// String generates a string with general information
func (s *Stats) String() (out string) {
	var tokens []string

	if s.ErrCount > 0 {
		tokens = append(tokens, fmt.Sprintf("%d errors", s.ErrCount))
	}
	if s.NumSkippedFiles > 0 {
		tokens = append(tokens, fmt.Sprintf("%d skipped", s.NumSkippedFiles))
	}
	if s.NumUndoneFiles > 0 {
		tokens = append(tokens, fmt.Sprintf("%d Undone", s.NumUndoneFiles))
	}
	if s.WasCancelled {
		tokens = append(tokens, "cancelled")
	}

	if len(tokens) > 0 {
		out = strings.Join(tokens, ", ")
		out = " (" + out + ")"
	}

	return
}
