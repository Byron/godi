package godi

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"sync"

	"github.com/Byron/godi/utility"
)

// Drains FileInfos from files channel, reads them using ctrl and generates hashes.
// Creates a Result using makeResult() and sends it down the results channel.
// Listens on done to abort early, and notifies wg when we are done to know when results can be closed.
// If wmap is set, we will setup a parallel writer which writes the bytes used for hashing
// to all controllers at the same time, which will be as slow as the slowest device
func Gather(files <-chan FileInfo, results chan<- Result, wg *sync.WaitGroup, done <-chan bool,
	makeResult func(*FileInfo) (Result, *BasicResult),
	rctrl *utility.ReadChannelController,
	wmap map[uint64]utility.RootedWriteController) {
	if rctrl == nil || wg == nil {
		panic("ReadChannelController and WaitGroup must be set")
	}
	defer wg.Done()
	sha1gen := sha1.New()
	md5gen := md5.New()
	var deviceWriter *utility.UncheckedParallelMultiWriter

	multiWriter := utility.NewUncheckedParallelMultiWriter(sha1gen, md5gen)

	// This makes the write as slow as the slowest hash, instead of hash+hash
	if deviceWriter != nil {
		multiWriter.Writers = append(multiWriter.Writers, deviceWriter)
	}

	sendErrorResult := func(f *FileInfo, err error) {
		sres, res := makeResult(f)
		res.Err = err
		results <- sres
	}

	// This MUST be a copy of f here, otherwise we will be in trouble thanks to the user of defer in handleHash
	// we will get f overwritten by the next iteration variable ... it's kind of special, might
	// be intersting for the mailing list.
	handleHash := func(f FileInfo) {
		// In hash-only mode, there is only one result
		var err error
		if deviceWriter != nil {
			// in write mode, there are as many results as we have destinations
			// therefore, result handlling is not done by the writer itself
		}

		// let the other end open the file and close it as well
		reader := rctrl.NewChannelReaderFromPath(f.Path)
		rctrl.Channel() <- reader

		sha1gen.Reset()
		md5gen.Reset()
		var written int64
		written, err = reader.WriteTo(&multiWriter)
		if err != nil {
			sendErrorResult(&f, err)
			return
		}
		f.Sha1 = sha1gen.Sum(nil)
		f.MD5 = md5gen.Sum(nil)
		if written != f.Size {
			err = fmt.Errorf("Filesize of '%s' reported as %d, yet only %d bytes were hashed", f.Path, f.Size, written)
			sendErrorResult(&f, err)
			return
		}

		if deviceWriter == nil {
			// We have to take care of sending the read-result
			sres, _ := makeResult(&f)
			results <- sres
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
