package godi

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"sync"
	"time"

	"github.com/Byron/godi/utility"
)

// Generate does all boilerplate required to be a valid generator
func Generate(
	done <-chan bool,
	generate func(chan<- FileInfo, chan<- Result)) (<-chan FileInfo, <-chan Result) {
	files := make(chan FileInfo)
	results := make(chan Result)

	go func() {
		defer close(files)
		generate(files, results)
	}()

	return files, results
}

// Reads
func Gather(files <-chan FileInfo, results chan<- Result, wg *sync.WaitGroup, done <-chan bool,
	makeResult func(*FileInfo) (Result, *BasicResult),
	ctrl *utility.ReadChannelController) {
	if ctrl == nil || wg == nil {
		panic("ReadChannelController and WaitGroup must be set")
	}
	defer wg.Done()
	sha1gen := sha1.New()
	md5gen := md5.New()
	// This makes the write as slow as the slowest hash, instead of hash+hash
	allHashes := utility.UncheckedParallelMultiWriter(sha1gen, md5gen)

	// This MUST be a copy of f here, otherwise we will be in trouble thanks to the user of defer in handleHash
	// we will get f overwritten by the next iteration variable ... it's kind of special, might
	// be intersting for the mailing list.
	handleHash := func(f FileInfo) {
		sres, res := makeResult(&f)
		err := &res.Err
		defer func(res Result) { results <- res }(sres)

		// let the other end open the file and close it as well
		reader := ctrl.NewChannelReaderFromPath(f.Path)
		ctrl.Channel() <- reader

		sha1gen.Reset()
		md5gen.Reset()
		var written int64
		written, *err = reader.WriteTo(allHashes)
		if *err != nil {
			return
		}
		f.Sha1 = sha1gen.Sum(nil)
		f.MD5 = md5gen.Sum(nil)
		if written != f.Size {
			*err = fmt.Errorf("Filesize of '%s' reported as %d, yet only %d bytes were hashed", f.Path, f.Size, written)
			return
		}
	}

	for f := range files {
		select {
		case <-done:
			return
		default:
			handleHash(f)
		}
	}
}

type AggregateFinalizerState struct {
	WasCancelled        bool
	FileCount, ErrCount uint
	SizeBytes           uint64
	Elapsed             float64
}

// String generates a string with performance information
func (a *AggregateFinalizerState) String() string {
	sizeMB := float32(a.SizeBytes) / (1024.0 * 1024.0)
	return fmt.Sprintf(
		"Processed %#vMB in %vs (%#v MB/s, %d errors, cancelled=%v",
		sizeMB,
		a.Elapsed,
		float64(sizeMB)/a.Elapsed,
		a.ErrCount,
		a.WasCancelled,
	)
}

// Aggregate is a general purpose implementation to gather fileInfo results
func Aggregate(results <-chan Result, done <-chan bool,
	resultHandler func(Result, chan<- Result) bool,
	finalizer func(chan<- Result, *AggregateFinalizerState)) <-chan Result {
	accumResult := make(chan Result)

	go func() {
		defer close(accumResult)

		s := AggregateFinalizerState{}
		st := time.Now()

		// ACCUMULATE PATHS INFO
		/////////////////////////
		for r := range results {
			// Be sure we take note of cancellation.
			// If this happens, soon our results will be drained and we leave naturally
			select {
			case <-done:
				s.WasCancelled = true
			default:
				{
					if r.Error() != nil {
						s.ErrCount += 1
						accumResult <- r
						continue
					}

					if !resultHandler(r, accumResult) {
						s.ErrCount += 1
					} else {
						s.FileCount += 1
						s.SizeBytes += uint64(r.FileInformation().Size)
					}
				} // default
			} // select
		} // range results
		s.Elapsed = time.Now().Sub(st).Seconds()

		finalizer(accumResult, &s)
	}()

	return accumResult
}
