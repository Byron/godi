package api

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"path/filepath"
	"sync/atomic"

	gio "github.com/Byron/godi/io"
)

// Thrown if the filesize we read didn't match with the filesize we were supposed to read
type FileSizeMismatch struct {
	Path      string
	Want, Got int64
}

// Thrown if a file hash didn't match - it's used in the verify implementation, primarily
type FileHashMismatch struct {
	Path string
}

func (f *FileSizeMismatch) Error() string {
	return fmt.Sprintf("Filesize of '%s' reported as %d, yet %d bytes were read", f.Path, f.Want, f.Got)
}

func (f *FileHashMismatch) Error() string {
	return f.Path
}

// Intercepts Write calls and updates the stats accordingly. Implements only what we need, forwrading the calls as needed
type HashStatAdapter struct {
	hash  hash.Hash
	stats *Stats
}

func (h *HashStatAdapter) Write(b []byte) (int, error) {
	n, err := h.hash.Write(b)
	atomic.AddUint64(&h.stats.BytesHashed, uint64(n))
	return n, err
}

func (h *HashStatAdapter) Reset() {
	h.hash.Reset()
}

func (h *HashStatAdapter) Sum(b []byte) []byte {
	return h.hash.Sum(b)
}

// Drains FileInfos from files channel, reads them using ctrl and generates hashes.
// Creates a Result using makeResult() and sends it down the results channel.
// If wctrls is set, we will setup parallel writer which writes the bytes used for hashing
// to all controllers at the same time, which will be as slow as the slowest device
// TODO(st) wctrls must be device mapping. That way, we can parallelize writes per device.
// Right now we have a slow brute-force approach, which will make random writes to X files, but only Y at a time.
// What we want is at max Y files being written continuously at a time
func Gather(files <-chan FileInfo, results chan<- Result, stats *Stats,
	makeResult func(*FileInfo, *FileInfo, error) Result,
	rctrl *gio.ReadChannelController,
	wctrls []gio.RootedWriteController) {
	if rctrl == nil {
		panic("ReadChannelController and WaitGroup must be set")
	}

	sha1gen := HashStatAdapter{sha1.New(), stats}
	md5gen := HashStatAdapter{md5.New(), stats}
	const nHashes = 2
	atomic.AddUint32(&stats.NumHashers, uint32(nHashes))
	isWriting := len(wctrls) > 0
	numDestinations := gio.WriteChannelDeviceMapTrees(wctrls)
	// The hgher this number, the less syscall and communication overhead we will have.
	// As we expect mostly larger files, we go for bigger buffers
	var buf [512 * 1024]byte

	var multiWriter io.Writer
	var channelWriters []gio.ChannelWriter
	var lazyWriters []gio.LazyFileWriteCloser

	// We keep an index of failed destinations, skipping them in write mode after first failure
	var isFailedDestination []bool
	numFailedDestinations := 0

	// Build the multi-writer which will dispatch all writes to a write controller
	if isWriting {
		// We have one controller per device, each as a number of streams
		writers := make([]io.Writer, numDestinations+nHashes)

		// Writer with full checking enabled - it will never show anything for the hashes, but might
		// report errrs for the real writers
		// We place the hashes last, as the writers will be changed in each iteration
		writers[len(writers)-1] = &sha1gen
		writers[len(writers)-2] = &md5gen
		multiWriter = gio.NewParallelMultiWriter(writers)

		// Keeps all Writers we are going to prepare per source file
		channelWriters = make([]gio.ChannelWriter, numDestinations)
		lazyWriters = make([]gio.LazyFileWriteCloser, numDestinations)
		isFailedDestination = make([]bool, numDestinations)

		// Init them, per controller, and set them to be used by the multi-writer right away
		// These never change, only their Destination path does
		ofs := 0
		for _, wctrl := range wctrls {
			ofse := ofs + len(wctrl.Trees)
			wctrl.Ctrl.InitChannelWriters(channelWriters[ofs:ofse])
			for x := ofs; x < ofse; x++ {
				channelWriters[x].SetWriter(&lazyWriters[x])
			}
			ofs = ofse
		}
	} else {
		multiWriter = gio.NewUncheckedParallelMultiWriter(&sha1gen, &md5gen)
	}

	// umf == unmodifiedFileInfo
	sendResults := func(f *FileInfo, umf *FileInfo, err error) {

		if !isWriting {
			// NOTE: in seal mode, we might want to communicate this ... after all, we don't get the seal done
			results <- makeResult(f, umf, err)
		} else {
			// Each parallel writer has it's own result, we just send it off
			pmw := multiWriter.(*gio.ParallelMultiWriter)
			// Make sure we keep the source intact
			forig := *f

			for i := 0; i < numDestinations; i++ {
				w, e := pmw.WriterAtIndex(i)
				if e != nil && !isFailedDestination[i] {
					isFailedDestination[i] = true
					numFailedDestinations += 1
				}

				// If the reader had an error, no write may succeed. We just don't overwrite write errors
				if e == nil && err != nil {
					e = err
				}
				// Could be a previously unset writer
				if wc, ok := w.(gio.WriteCloser); ok {
					wc.Close()
					// we may change the same instance, as it will be copied into the Result structure later on
					f.Path = wc.Writer().(*gio.LazyFileWriteCloser).Path()
					// it doesn't matter here if there actually is an error, aggregator will handle it
					results <- makeResult(f, &forig, e)
				}
			} // for each write controller to write to

			// If all of our destinations are in fail state, let the gatherer know we can't do anything
			if numFailedDestinations == numDestinations {
				atomic.AddUint32(&stats.StopTheEngines, 1)
			}

		} // handle write mode
	} // sendResults

	// NOTE: This loop must not return ! It must be finished !!
	for f := range files {
		// In hash-only mode, there is only one result
		var err error
		if isWriting {
			// in write mode, there are as many results as we have destinations
			// All that's left to be done is to fill our ChannelWriters with new lazyFileWriters
			// that write to the given paths

			// Current writer id, absolute to this write operation
			awid := 0
			for _, wctrl := range wctrls {
				// We just create one ChannelWriter per destination, and let the writers
				// deal with the parallelization and blocking
				fawid := awid // the first index
				pmw := multiWriter.(*gio.ParallelMultiWriter)

				for x := 0; x < len(wctrl.Trees); x++ {
					// reset previous errror by resetting the writer pointer
					// If the destination is known to have an error, disable its writer
					if isFailedDestination[awid] {
						pmw.SetWriterAtIndex(awid, nil)
					} else {
						pmw.SetWriterAtIndex(awid, &channelWriters[awid])
					}
					lazyWriters[awid].SetPath(filepath.Join(wctrl.Trees[awid-fawid], f.RelaPath), f.Mode)
					awid += 1
				}
			} // for each device's write controller

			if awid != numDestinations {
				panic("Mismatched writers")
			}

		} // handle write mode preparations

		// let the other end open the file and close it as well
		reader := rctrl.NewChannelReaderFromPath(f.Path, f.Mode, buf[:])

		sha1gen.Reset()
		md5gen.Reset()
		var written int64
		written, err = reader.WriteTo(multiWriter)

		if err != nil {
			// This should actually never fail, the way we are implemented.
			// If it does, it's the WriteTo() implementation, and as we are decoupled from it,
			// let's make the check here anyway ...
			sendResults(&f, &f, err)
			continue
		}

		// Always keep the hash as far as we have it - it's a value to preserve
		umf := f
		f.Sha1 = sha1gen.Sum(nil)
		f.MD5 = md5gen.Sum(nil)

		if written != f.Size {
			err = &FileSizeMismatch{f.Path, f.Size, written}
			sendResults(&f, &umf, err)
			continue
		}

		// all good

		sendResults(&f, &umf, nil)
	} // end for each file to process

	// Keep the count valid ...
	atomic.AddUint32(&stats.NumHashers, ^uint32(nHashes-1))
}
