package api

import (
	"io"
	"os"
	"sync"
)

// Size for allocated buffers
const bufSize = 64 * 1024
const readChannelSize = 0

// Similar to MultiWriter, but assumes writes never fail, and provides the same buffer
// to all writers in parallel. However, it will return only when all writes are finished
type uncheckedParallelMultiWriter struct {
	writers []io.Writer
	wg      sync.WaitGroup
}

func (t *uncheckedParallelMultiWriter) Write(p []byte) (n int, err error) {
	t.wg.Add(len(t.writers))
	for _, w := range t.writers {
		go func(w io.Writer) {
			w.Write(p)
			t.wg.Done()
		}(w)
	}
	t.wg.Wait()
	return len(p), nil
}

func UncheckedParallelMultiWriter(writers ...io.Writer) io.Writer {
	w := make([]io.Writer, len(writers))
	copy(w, writers)
	return &uncheckedParallelMultiWriter{w, sync.WaitGroup{}}
}

// The result of a read operation, similar to what Reader.Read returns
type readResult struct {
	buf []byte
	n   int
	err error
}

type ReadChannelController struct {
	c chan ChannelReader
}

type WriteChannelController struct {
	c chan ChannelWriter
	// if true, we will copy all incoming buffers to assure we own them
	takeOwnership bool

	// NOTE: at some point we could hold a pool of large buffers for in-memory write caching
	// However, large buffers could be beneficial for the hashing already as we do less small hash calls
}

// Contains all information about a file or reader to be read
type ChannelReader struct {

	// Our controller, providing the pool we deal with
	ctrl *ReadChannelController

	// An optional path, which will be opened for reading when Reader is nil
	path string

	// A Reader interface, in case Path is unset. Use this if you want to open the file or provide your
	// own custom reader
	reader io.Reader

	// The channel to transport read results
	results chan readResult
}

type ChannelWriter struct {

	// Our controller, containing additional information as needed
	ctrl *WriteChannelController

	// A path to write the data to. May be empty, which is when a Writer instance must be set
	// The benefit of using this mechanism is to have file-handles opened only when the operation
	// Actually occours
	path string

	// A writer to write to. Must be set if path is nil
	writer io.Writer

	// A channel through which to receive packets of bytes.
	// They must be owned by us.
	bytes chan []byte

	// A function to call from our go-routine when we are done.
	// It should only do minimal work and route the passed results
	// path is the path of this instance, writer is the writer of this instance (if set initially)
	// nwritten is the amount of written bytes, whereas error denotes the error.
	doneCB func(path string, writer io.Writer, nwritten int64, err error)
}

func (r *ReadChannelController) Channel() chan<- ChannelReader {
	return r.c
}

// Return a new channel reader
// You should set either path
func (r *ReadChannelController) NewChannelReaderFromPath(path string) ChannelReader {
	// NOTE: size of this channel controls how much we can cache into memory before we block
	// as the consumer doesn't keep up
	return ChannelReader{r, path, nil, make(chan readResult, readChannelSize)}
}

func (r *ReadChannelController) NewChannelReaderFromReader(reader io.Reader) ChannelReader {
	return ChannelReader{r, "", reader, make(chan readResult, readChannelSize)}
}

// Allows to use a ChannelReader as source for io.Copy operations
// This should be preferred as it will save a copy operation
// WriteTo will block until a Reader is ready to serve us bytes
// Note that the read operation is performed by N reader routines - we just receive the data
// and pass it on
// Also we assume that write blocks until the operation is finished. If you perform non-blocking writes,
// you must copy the buffer !
func (p *ChannelReader) WriteTo(w io.Writer) (n int64, err error) {
	// We are just consuming, and assume the channel is closed when the reading is finished
	var written int
	for res := range p.results {
		// Write what's possible
		if res.n > 0 {
			written, err = w.Write(res.buf)
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
func NewReadChannelController(nprocs int) ReadChannelController {
	if nprocs < 1 {
		panic("nprocs must be >= 1")
	}

	ctrl := ReadChannelController{
		make(chan ChannelReader, nprocs),
	}

	for i := 0; i < nprocs; i++ {
		go func() {
			for info := range ctrl.c {
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
					buf := make([]byte, bufSize)
					nread, err = info.reader.Read(buf)
					info.results <- readResult{buf[:nread], nread, err}
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

	return ctrl
}

// Create a new controller which deals with writing all incoming requests with nprocs go-routines
func NewWriteChannelController(nprocs int) WriteChannelController {
	ctrl := WriteChannelController{
		make(chan ChannelWriter, nprocs),
		false,
	}

	// TODO: implementation
	// We will only really need this when we are copying data anyway ... .
	return ctrl
}
