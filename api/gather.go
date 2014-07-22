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
// Listens on done to abort early, and notifies wg when we are done to know when results can be closed.
// If wctrls is set, we will setup parallel writer which writes the bytes used for hashing
// to all controllers at the same time, which will be as slow as the slowest device
// TODO(st) wctrls must be device mapping. That way, we can parallelize writes per device.
// Right now we have a slow brute-force approach, which will make random writes to X files, but only Y at a time.
// What we want is at max Y files being written continuously at a time
func Gather(files <-chan FileInfo, results chan<- Result, wg *sync.WaitGroup, done <-chan bool,
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

	var multiWriter io.Writer

	// Build the multi-writer which will dispatch all writes to a write controller
	if isWriteMode {
		writers := make([]io.Writer, len(wctrls)+nHashes)

		// Writer with full checking enabled - it will never show anything for the hashes, but might
		// report errrs for the real writers
		// We place the hashes last, as the writers will be changed in each iteration
		writers[len(writers)-1] = sha1gen
		writers[len(writers)-2] = md5gen
		multiWriter = utility.NewParallelMultiWriter(writers)
	} else {
		multiWriter = utility.NewUncheckedParallelMultiWriter(sha1gen, md5gen)
	}

	sendResults := func(f *FileInfo, err error) {
		if !isWriteMode {
			// We have to take care of sending the read-result
			results <- makeResult(f, nil, err)
		} else {
			// check each writer for errors and return produce a result, one per  non-hash writer
			pmw := multiWriter.(*utility.ParallelMultiWriter)
			for i := 0; i < len(wctrls); i++ {
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

	// This MUST be a copy of f here, otherwise we will be in trouble thanks to the user of defer in handleHash
	// we will get f overwritten by the next iteration variable ... it's kind of special, might
	// be intersting for the mailing list.
	handleHash := func(f FileInfo) {
		rctrl, hasRCtrlForRoot := rctrls[f.Root()]
		if !hasRCtrlForRoot {
			panic(fmt.Sprintf("Couldn't find read controller for directory at '%s'", f.Root()))
		}

		// In hash-only mode, there is only one result
		var err error
		if isWriteMode {
			// in write mode, there are as many results as we have destinations
			// therefore, result handlling is needs to be done once per writer.
			pmw := multiWriter.(*utility.ParallelMultiWriter)
			for i, wctrl := range wctrls {
				// get destination path
				destPath := filepath.Join(wctrl.Tree, f.RelaPath)
				cw := wctrl.Ctrl.NewChannelWriter(&utility.LazyFileWriteCloser{Path: destPath})
				pmw.SetWriterAtIndex(i, cw)
			}
		}

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
			return
		}
		f.Sha1 = sha1gen.Sum(nil)
		f.MD5 = md5gen.Sum(nil)
		if written != f.Size {
			err = fmt.Errorf("Filesize of '%s' reported as %d, yet only %d bytes were hashed", f.Path, f.Size, written)
			sendResults(&f, err)
			return
		} else {
			// all good
			sendResults(&f, nil)
		}
	} // func() handleHash

	for f := range files {
		select {
		case <-done:
			return
		default:
			handleHash(f)
		}
	}
}
