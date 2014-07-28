package api

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Byron/godi/io"
)

// A structure to keep information about what is currently going on.
// It is means to be used as shared resource, used by multiple threads, which is why
// thread-safe counters are used.
// Implementations must keep these numbers up-to-date, while async processors will digest
// and present the data in some form
type Stats struct {
	io.Stats

	BytesHashed uint64 // Total of bytes hashed so far, counting all active hashers
	NumHashers  uint32 // Amount of hashers running in parallel

	//GENERATOR INFORMATION
	NumSkippedFiles uint32 // Amount of files we skipped right away
	StopTheEngines  uint32 // Amount of gather procs which had write errors on all destinations

	// AGGREGATION
	// Aggregation step is single-threaded - no atomic operation needed
	ErrCount       uint // Amount of errors that hit the aggregation step
	NumUndoneFiles uint // Amout of files removed during undo
	WasCancelled   bool // is true if the user cancelled

}

// Similar to io.Stats.CopyTo(), but with our fields
func (s *Stats) CopyTo(d *Stats) {
	s.Stats.CopyTo(&d.Stats)

	d.BytesHashed = atomic.LoadUint64(&s.BytesHashed)
	d.NumHashers = atomic.LoadUint32(&s.NumHashers)

	d.NumSkippedFiles = atomic.LoadUint32(&s.NumSkippedFiles)
	d.StopTheEngines = atomic.LoadUint32(&s.StopTheEngines)

	// Agg variables don't need to be atomic - we copy them here for completeness only
	d.ErrCount = s.ErrCount
	d.NumUndoneFiles = s.NumUndoneFiles
	d.WasCancelled = s.WasCancelled
}

// Prints performance metrics as a single line full of useful information, including deltas of relevant metrics as compared
// to the last state d. You will also give the temporal distance which separates this stat from the previous one
// If you pass s as d, this indicates a result mode, which assumes you want the overall average throughput
// Sep is the separator to use between fields
func (s *Stats) DeltaString(d *Stats, td time.Duration, sep string) string {
	resultMode := d == s

	// Embed the hashing data between read and possibly write
	out := s.DeltaDataString(io.ElapsedData|io.ReadData, &d.Stats, td, sep)

	bh := atomic.LoadUint64(&s.BytesHashed)
	out += fmt.Sprintf("%s%sHASH %s%s%s",
		sep,
		s.Stats.InOut(atomic.LoadUint32(&s.NumHashers)),
		io.SymbolHash,
		io.BytesVolume(bh),
		s.Stats.BytesDelta(bh, d.BytesHashed, td, resultMode),
	)

	out += s.DeltaDataString(io.WriteData, &d.Stats, td, sep)

	return out
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
