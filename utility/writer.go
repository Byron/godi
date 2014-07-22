package utility

import (
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Similar to MultiWriter, but assumes writes never fail, and provides the same buffer
// to all writers in parallel. However, it will return only when all writes are finished
type uncheckedParallelMultiWriter struct {
	writers []io.Writer
	wg      sync.WaitGroup
}

// A writer which dispatches to multiple destinations, collecting errors on the way
// and returning the first one it encounteres.
// If a writer fails, it will not be written anymore until it is closed or reset using SetWriter
type ParallelMultiWriter struct {
	writers []io.Writer
	wg      sync.WaitGroup

	// Pre-allocated array, one slot per writer
	results []error
}

func (t *uncheckedParallelMultiWriter) Write(b []byte) (n int, err error) {
	t.wg.Add(len(t.writers))
	for _, w := range t.writers {
		go func(w io.Writer) {
			w.Write(b)
			t.wg.Done()
		}(w)
	}
	t.wg.Wait()
	return len(b), nil
}

func NewUncheckedParallelMultiWriter(writers ...io.Writer) io.Writer {
	w := make([]io.Writer, len(writers))
	copy(w, writers)
	return &uncheckedParallelMultiWriter{writers: w}
}

func NewParallelMultiWriter(writers []io.Writer) *ParallelMultiWriter {
	w := ParallelMultiWriter{}
	w.writers = writers
	w.results = make([]error, len(writers))
	return &w
}

// Set the given writer to be located at the given index. We don't do bounds checking
func (p *ParallelMultiWriter) SetWriterAtIndex(i int, w io.Writer) {
	p.writers[i] = w
	p.results[i] = nil
}

// Return the writer at the given index, and the first error it might have caused when writing
// to it. We perform no bounds checking
func (p *ParallelMultiWriter) WriterAtIndex(i int) (io.Writer, error) {
	return p.writers[i], p.results[i]
}

// Writes will always succeed, even if individual writers may have failed.
// It's up to our user to check for errors when the write is finished
func (p *ParallelMultiWriter) Write(b []byte) (n int, err error) {
	p.wg.Add(len(p.writers))
	for i, w := range p.writers {
		// continue on writers with errors
		if p.results[i] != nil {
			continue
		}
		go func(i int, w io.Writer) {
			_, p.results[i] = w.Write(b)
			p.wg.Done()
		}(i, w)
	}
	p.wg.Wait()
	return len(b), nil
}

// A simple struct keeping all information we need to make a write and retrieve the results
type writeInfo struct {
	// bytes to write
	b []byte
	// amount of bytes written
	n int
	// error of previous write operation
	e error
}

// Used in conjunction with a WriteChannelController, serving as front-end communicating with
// the actual writer that resides in a separate go-routine
type channelWriter struct {

	// A writer to write to. Must be set if path is nil
	writer io.Writer

	// Our shared write-info instance, similar to the buffer in the channelReader implementation
	wi writeInfo

	// For simplicity, we just use a channel with two-way signaling
	ready chan bool
}

// Like WriteCloser interface, but allows to retrieve more information specific to our usage
type WriteCloser interface {
	io.WriteCloser

	// Writer returns the writer this interface instance contains
	Writer() io.Writer
}

func (c *channelWriter) Writer() io.Writer {
	return c.writer
}

// Send bytes down our channel and wait for the writer on the end to be done, retrieving the result.
func (c *channelWriter) Write(b []byte) (n int, err error) {
	c.wi.b = b
	c.ready <- true
	// this will block until the actual writer is done ...
	<-c.ready
	// ... allowing us to return the actual result safely now
	return c.wi.n, c.wi.e
}

func (c *channelWriter) Close() error {
	// allows writer to break out of his loop and close the actual writer, freeing resources accordingly
	close(c.ready)
	// BUG(st): Right now we can't retrieve the actual error value. Of course there could be one in case
	// of buffered writers that flush and fail because of disk full
	return nil
}

// A writer that will create a new file and intermediate directories on first write.
// You must call the close method to finish the writes and release system resources
type LazyFileWriteCloser struct {

	// The path we should open a writer to on first write. This will fail if the fail already exists.
	Path string

	// A writer we are using to perform the write
	writer *os.File
}

func (l *LazyFileWriteCloser) Write(b []byte) (n int, err error) {
	if l.writer == nil {
		// assure directory exists
		err = os.MkdirAll(filepath.Dir(l.Path), 0777)
		if err != nil {
			return 0, err
		}
		l.writer, err = os.OpenFile(l.Path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return 0, err
		}
	}

	return l.writer.Write(b)
}

// Close our writer if it was initialized already. Therefore it's safe to call this even if Write wasn't called
// beforehand
func (l *LazyFileWriteCloser) Close() error {
	if l.writer != nil {
		return l.writer.Close()
	}
	return nil
}

type WriteChannelController struct {
	c chan channelWriter

	// NOTE: at some point we could hold a pool of large buffers for in-memory write caching
	// However, large buffers could be beneficial for the hashing already as we do less small hash calls
}

// A utility structure to associate a tree with a writer.
// That way, writers can be more easily associated with a device which hosts Tree
type RootedWriteController struct {
	Tree string
	Ctrl WriteChannelController
}

// Create a new controller which deals with writing all incoming requests with nprocs go-routines
func NewWriteChannelController(nprocs int) WriteChannelController {
	ctrl := WriteChannelController{
		make(chan channelWriter, nprocs),
	}

	for i := 0; i < nprocs; i++ {
		go func() {
			for cw := range ctrl.c {
				// read all bytes as long as writer is not closed
				for {
					if _, ok := <-cw.ready; !ok {
						break
					}
					cw.wi.n, cw.wi.e = cw.writer.Write(cw.wi.b)
					// protocol mandates the sender has to listen for our reply, channel must not be closed here ... .
					cw.ready <- true
				} // write endlessly

				// signal the writer we are done, allowing it to cleanup
				if wc, ok := cw.writer.(io.WriteCloser); ok {
					// as ready channel is closed, we can't return any error value here ...
					wc.Close()
				}
			} // for each channel writer
		}()
	} // for each routine to create
	// We will only really need this when we are copying data anyway ... .
	return ctrl
}

// Return a new channel writer, which will write asynchronously to the given writer
func (w *WriteChannelController) NewChannelWriter(writer io.Writer) io.Writer {
	cw := channelWriter{
		writer: writer,
		ready:  make(chan bool),
	}
	w.c <- cw
	return &cw
}
