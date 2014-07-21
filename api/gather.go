package godi

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"sync"

	"github.com/Byron/godi/utility"
)

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
