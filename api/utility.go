package api

import (
	"io"
	"os"
	"sync"
)

// Size for allocated buffers
const bufSize = 64 * 1024
const readChannelSize = 0

// Keeps enough buffers to ensure continuous reading without interfering with users
var bufPool = sync.Pool{New: func() interface{} { return make([]byte, bufSize) }}

// The result of a read operation, similar to what Reader.Read returns
type readResult struct {
	buf []byte
	n   int
	err error
}

// Contains all information about a file or reader to be read
type ChannelReader struct {
	// An optional path, which will be opened for reading when Reader is nil
	path string

	// A Reader interface, in case Path is unset. Use this if you want to open the file or provide your
	// own custom reader
	reader io.Reader

	// The channel to transport read results
	results chan readResult
}

// Return a new channel reader
// You should set either path
func NewChannelReaderFromPath(path string) ChannelReader {
	// NOTE: size of this channel controls how much we can cache into memory before we block
	// as the consumer doesn't keep up
	return ChannelReader{path, nil, make(chan readResult, readChannelSize)}
}

func NewChannelReaderFromReader(reader io.Reader) ChannelReader {
	return ChannelReader{"", reader, make(chan readResult, readChannelSize)}
}

// Allows to use a ChannelReader as source for io.Copy operations
// This should be preferred as it will save a copy operation
// WriteTo will block until a Reader is ready to serve us bytes
// Note that the read operation is performed by N reader routines - we just receive the data
// and pass it on
func (p *ChannelReader) WriteTo(w io.Writer) (n int64, err error) {
	// We are just consuming, and assume the channel is closed when the reading is finished
	var written int
	for res := range p.results {
		// Write what's possible
		if res.n > 0 {
			written, err = w.Write(res.buf)
			bufPool.Put(res.buf)
			n += int64(written)
		}
		// I could think of plenty of assertion here ... but there is no such thing
	}

	// whatever is held in n, err, we return
	return
}

// Create a new parallel reader with nprocs go-routines and return a channel to it.
// Feed the channel with ChannelReader structures and listen on it's channel to read bytes until EOF, which
// is when the channel will be closed by the reader
func NewReadChannel(nprocs int) chan<- ChannelReader {
	if nprocs < 1 {
		panic("nprocs must be >= 1")
	}

	c := make(chan ChannelReader, nprocs)

	for i := 0; i < nprocs; i++ {
		go func() {
			for info := range c {
				var err error
				ourReader := false
				if info.reader == nil {
					info.reader, err = os.Open(info.path)
					if err != nil {
						info.results <- readResult{nil, 0, err}
						continue
					}
				}

				// Now read until it's done
				var nread int
				for {
					// The buffer will be put back by the one reading from the channel (e.g. in WriteTo()) !
					buf := bufPool.Get().([]byte)
					nread, err = info.reader.Read(buf)
					info.results <- readResult{buf, nread, err}
					// we send all results, but abort if the reader is done for whichever reason
					if err != nil {
						break
					}
				} // read loop
				// Signal the consumer that we are done
				close(info.results)

				if ourReader {
					info.reader.(*os.File).Close()
					info.reader = nil
				}
			}
		}()
	}

	return c
}
