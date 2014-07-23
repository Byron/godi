package godi

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"github.com/Byron/godi/utility"
)

// Drains FileInfos from files channel, reads them using ctrl and generates hashes.
// Creates a Result using makeResult() and sends it down the results channel.
// If wctrls is set, we will setup parallel writer which writes the bytes used for hashing
// to all controllers at the same time, which will be as slow as the slowest device
// TODO(st) wctrls must be device mapping. That way, we can parallelize writes per device.
// Right now we have a slow brute-force approach, which will make random writes to X files, but only Y at a time.
// What we want is at max Y files being written continuously at a time
func Gather(files <-chan FileInfo, results chan<- Result, wg *sync.WaitGroup,
	makeResult func(*FileInfo, *FileInfo, error) Result,
	rctrls map[string]*utility.ReadChannelController,
	wctrls []utility.RootedWriteController) {
	if len(rctrls) == 0 || wg == nil {
		panic("ReadChannelController and WaitGroup must be set")
	}
	defer wg.Done()

	sha1gen := sha1.New()
	md5gen := md5.New()
	const nHashes = 2
	isWriteMode := len(wctrls) > 0
	numParallelWriters := utility.WriteChannelDeviceMapStreams(wctrls)
	// Pre-allocated memory to keep writer allocations
	var writersPerCtrl [][]io.Writer

	var multiWriter io.Writer

	// Build the multi-writer which will dispatch all writes to a write controller
	if isWriteMode {
		// We have one controller per device, each as a number of streams
		writers := make([]io.Writer, numParallelWriters+nHashes)

		// Writer with full checking enabled - it will never show anything for the hashes, but might
		// report errrs for the real writers
		// We place the hashes last, as the writers will be changed in each iteration
		writers[len(writers)-1] = sha1gen
		writers[len(writers)-2] = md5gen
		multiWriter = utility.NewParallelMultiWriter(writers)

		// pre-allocated writers structure
		writersPerCtrl = make([][]io.Writer, len(wctrls))
		for rid, wctrl := range wctrls {
			writersPerCtrl[rid] = make([]io.Writer, len(wctrl.Trees))
		}
	} else {
		multiWriter = utility.NewUncheckedParallelMultiWriter(sha1gen, md5gen)
	}

	sendResults := func(f *FileInfo, err error) {
		if !isWriteMode {
			// We have to take care of sending the read-result
			results <- makeResult(f, nil, err)
		} else {
			// check each writer for errors and produce a result, one per non-hash writer
			pmw := multiWriter.(*utility.ParallelMultiWriter)
			for i := 0; i < numParallelWriters; i++ {
				// copy f for adjusting it's absolute path - we send it though the channel as pointer, not value
				var wfi FileInfo = *f
				w, e := pmw.WriterAtIndex(i)
				wc := w.(utility.WriteCloser)
				wc.Close()
				wfi.Path = wc.Writer().(*utility.LazyFileWriteCloser).Path

				// it doesn't matter here if there actually is an error, aggregator will handle it
				results <- makeResult(&wfi, f, e)
			} // for each write controller to write to
		} // handle write mode
	} // sendResults

	// NOTE: This loop must not return ! It must be finished !!
	for fo := range files {
		// Make a copy - we pass this on as pointer, therefore we need to assure it's not the same
		// thing after all. Range writes the same memory block all over again.
		var f FileInfo = fo
		rctrl, hasRCtrlForRoot := rctrls[f.Root()]
		if !hasRCtrlForRoot {
			panic(fmt.Sprintf("Couldn't find read controller for directory at '%s'", f.Root()))
		}

		// In hash-only mode, there is only one result
		var err error
		if isWriteMode {
			// in write mode, there are as many results as we have destinations
			// therefore, result handlling is done once per writer.
			pmw := multiWriter.(*utility.ParallelMultiWriter)
			// Current writer id, absolute to this write operation
			awid := 0
			for rid, wctrl := range wctrls {
				// per device, we can handle X streams in parallel, while having to stream to NT trees.
				// We have one lazywriter per stream S<NS-1, whereas the last stream will be a MultiWriter with
				// NT-(NS-1) sequential streams.
				// If NS is 3 and NT is 6, X0 and X1 will have one destination each, and X2 has 4 done sequentially.
				ns := wctrl.ClientStreams()
				if ns > len(wctrl.Trees) {
					panic("This can't work if there are more Channels than Trees, wctrl type has a bug")
				}

				// Amount of parallel writers to reserve
				// If the number matches the amount of trees though, there is no special handling
				reserve := 1
				if ns == len(wctrl.Trees) {
					reserve = 0
				}

				// This method assures our WriterChannels are created in order, so they are not interleaved
				// with channels currently created by another go-routine. It's important, as we
				// might deadlock otherwise.
				// The writers have to be pre-filled with the actual writers,
				writers := writersPerCtrl[rid]

				// Create lazi's for highly parallel operation
				for x := 0; x < ns-reserve; x++ {
					destPath := filepath.Join(wctrl.Trees[x], f.RelaPath)
					writers[awid] = &utility.LazyFileWriteCloser{Path: destPath}

					awid += 1
				}

				// Now let's do the remaining Trees, but we put them all into one writer
				if reserve != 0 {
					remainingTrees := len(wctrl.Trees) - (ns - reserve)
					mw := utility.MultiWriter{Writers: make([]io.Writer, remainingTrees)}
					// This loop pepares the sequential wrtiers, writing one after another
					for x := 0; x < remainingTrees; x++ {
						destPath := filepath.Join(wctrl.Trees[ns-reserve+x], f.RelaPath)
						mw.Writers[x] = &utility.LazyFileWriteCloser{Path: destPath}
					}
					//
					writers[awid] = &mw
					awid += 1
				}

				// Finally, push all the writers into our parallel pipeline
				wctrl.Ctrl.NewChannelWriters(writers)
				// writers have been replaced by channel writers
				for cid, cw := range writers {
					pmw.SetWriterAtIndex(cid, cw)
				}
			} // for each device's write controller

			if awid != numParallelWriters {
				panic("Mismatched writers")
			}
		} // handle write mode preparations

		// let the other end open the file and close it as well
		reader := rctrl.NewChannelReaderFromPath(f.Path)

		sha1gen.Reset()
		md5gen.Reset()
		var written int64
		written, err = reader.WriteTo(multiWriter)

		if err != nil {
			// This should actually never fail, the way we are implemented.
			// If it does, it's the WriteTo() implementation, and as we are decoupled from it,
			// let's make the check here anyway ...
			sendResults(&f, err)
			continue
		}
		f.Sha1 = sha1gen.Sum(nil)
		f.MD5 = md5gen.Sum(nil)
		if written != f.Size {
			err = fmt.Errorf("Filesize of '%s' reported as %d, yet only %d bytes were hashed", f.Path, f.Size, written)
			sendResults(&f, err)
			continue
		}
		// all good

		sendResults(&f, nil)
	}
}
