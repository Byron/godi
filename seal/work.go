package seal

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"sync"

	"github.com/Byron/godi/api"
)

func (s *SealCommand) Gather(files <-chan api.FileInfo, results chan<- api.Result, wg *sync.WaitGroup, done <-chan bool) {
	defer wg.Done()
	sha1gen := sha1.New()
	md5gen := md5.New()
	// This makes the write as slow as the slowest hash, instead of hash+hash
	allHashes := api.UncheckedParallelMultiWriter(sha1gen, md5gen)

	// This MUST be a copy of f here, otherwise we will be in trouble thanks to the user of defer in handleHash
	// we will get f overwritten by the next iteration variable ... it's kind of special, might
	// be intersting for the mailing list.
	handleHash := func(f api.FileInfo) {
		res := SealResult{&f, "", nil, api.Progress}
		err := &res.err
		defer func(res *SealResult) { results <- res }(&res)

		// let the other end open the file and close it as well
		reader := s.pCtrl.NewChannelReaderFromPath(f.Path)
		s.pCtrl.Channel() <- reader

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
