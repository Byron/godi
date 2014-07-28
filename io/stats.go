package io

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	SymbolDelta         = "Δ"
	SymbolHash          = "⌗"
	StatsClientSep      = " | "
	ElapsedData    int8 = 1 << iota
	ReadData
	WriteData
)

// A shared structure that is modified using atomic operations to keep track of what data the IO system is
// Currently processing
type Stats struct {
	// PERFORMANCE METRICS
	TotalFilesRead    uint32 // Amount of whole files we read so far
	TotalFilesWritten uint32 // Amount of whole files we wrote so far
	FilesBeingRead    uint32 // Amount of files currently being read
	FilesBeingWritten uint32 // Amount of files currently being written
	BytesRead         uint64 // Total of bytes read so far, counting all input streams
	BytesWritten      uint64 // Total of bytes written so far, counting all output streams

	// GENERAL INFORMATION
	StartedAt time.Time // The time at which we started processing

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

	d.StartedAt = s.StartedAt
}

// MostFiles returns the greatest number of files, either the one that were read, or the ones that were written
func (s *Stats) MostFiles() uint32 {
	if s.TotalFilesRead > s.TotalFilesWritten {
		return s.TotalFilesRead
	}
	return s.TotalFilesWritten
}

func (s *Stats) IntDelta(cur, prev uint32, td time.Duration, resultMode bool) string {
	if prev == cur {
		if resultMode {
			prev = 0
		} else {
			return ""
		}
	}
	return fmt.Sprintf(" %s%s%04d/s", SymbolHash, SymbolDelta, uint64(float64(cur-prev)/td.Seconds()))
}

func (s *Stats) BytesDelta(cur, prev uint64, td time.Duration, resultMode bool) string {
	if prev == cur {
		if resultMode {
			prev = 0
		} else {
			return ""
		}
	}
	return fmt.Sprintf(" %s%s/s", SymbolDelta, BytesVolume(float64(cur-prev)/td.Seconds()))
}

func (s *Stats) InOut(cur uint32) string {
	if cur == 0 {
		return "  "
	}
	return fmt.Sprintf("%d ", cur)
}

// Prints performance metrics as a single line with the given data type, including deltas of relevant metrics as compared
// to the last state d. You will also give the temporal distance which separates this stat from the previous one
// If you pass s as d, this indicates a result mode, which assumes you want the overall average throughput
// Sep is the separator to use between fields
// dataType is a combination of ElapsedData|ReadData|WriteData
// The respective field can be unset in case there is no data for it
func (s *Stats) DeltaDataString(dataType int8, d *Stats, td time.Duration, sep string) (out string) {
	resultMode := d == s

	if dataType&ElapsedData == ElapsedData {
		timeStr := fmt.Sprintf("%4.0fs", s.Elapsed().Seconds())
		if resultMode {
			timeStr = s.Elapsed().String()
		}
		out += fmt.Sprintf("%s  %s%s", SymbolWallclock, timeStr, sep)
	}

	if dataType&ReadData == ReadData && (s.TotalFilesRead > 0 || s.FilesBeingRead > 0) {
		itf := atomic.LoadUint32(&s.TotalFilesRead)
		ibr := atomic.LoadUint64(&s.BytesRead)

		out += fmt.Sprintf("%s->READ %s%04d%s %s%s%s",
			s.InOut(atomic.LoadUint32(&s.FilesBeingRead)),
			SymbolHash,
			itf,
			s.IntDelta(itf, d.TotalFilesRead, td, resultMode),
			SymbolTotal,
			BytesVolume(ibr),
			s.BytesDelta(ibr, d.BytesRead, td, resultMode),
		)
	}

	if dataType&WriteData == WriteData && (s.TotalFilesWritten > 0 || s.FilesBeingWritten > 0) {
		otf := atomic.LoadUint32(&s.TotalFilesWritten)
		obw := atomic.LoadUint64(&s.BytesWritten)

		out += fmt.Sprintf("%s%sWRITE %s%04d%s %s%s%s",
			sep,
			s.InOut(atomic.LoadUint32(&s.FilesBeingWritten)),
			SymbolHash,
			otf,
			s.IntDelta(otf, d.TotalFilesWritten, td, resultMode),
			SymbolTotal,
			BytesVolume(obw),
			s.BytesDelta(obw, d.BytesWritten, td, resultMode),
		)
	}

	return out
}

// Return amount of time elapsed since we started the operation
func (s *Stats) Elapsed() time.Duration {
	return time.Now().Sub(s.StartedAt)
}
