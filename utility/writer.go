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
	for i, w := range p.writers {
		// continue on writers with errors
		if p.results[i] != nil || p.writers[i] == nil {
			continue
		}
		p.wg.Add(1)
		go func(i int, w io.Writer) {
			_, p.results[i] = w.Write(b)
			p.wg.Done()
		}(i, w)
	}
	p.wg.Wait()
	return len(b), nil
}

// Used in conjunction with a WriteChannelController, serving as front-end communicating with
// the actual writer that resides in a separate go-routine
type channelWriter struct {
	// The controller owning us
	ctrl *WriteChannelController

	// A writer to write to. Must be set if path is nil
	writer io.Writer

	// Our shared write-information, similar to the buffer in the channelReader implementation
	// bytes to write
	b []byte
	// amount of bytes written
	n int
	// error of previous write operation
	e error

	// Helps us to wait until the destination writer is done with our bytes
	// NOTE: Would it be faster to use a channel ?
	wg sync.WaitGroup
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
	c.b = b
	c.wg.Add(1)
	c.ctrl.c <- c
	c.wg.Wait()

	// ... allowing us to return the actual result safely now
	return c.n, c.e
}

func (c *channelWriter) Close() error {
	if w, ok := c.writer.(io.Closer); ok {
		return w.Close()
	}
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
		l.writer, err = os.OpenFile(l.Path, os.O_EXCL|os.O_WRONLY|os.O_CREATE, 0666)
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

// A utility to help control how parallel we try to write
type WriteChannelController struct {
	// Keeps all write requests, which contain all information we could possibly want to write something.
	// As the channelWriter is keeping all information, we serves as request right away
	c chan *channelWriter

	// NOTE: at some point we could hold a pool of large buffers for in-memory write caching
	// However, large buffers could be beneficial for the hashing already as we do less small hash calls
}

// A utility structure to associate a tree with a writer.
// That way, writers can be more easily associated with a device which hosts Tree
type RootedWriteController struct {
	// The tree the controller should write to
	Tree string

	// A possibly shared controller which may write to the given tree
	Ctrl *WriteChannelController
}

// Create a new controller which deals with writing all incoming requests with nprocs go-routines
func NewWriteChannelController(nprocs int) WriteChannelController {
	ctrl := WriteChannelController{
		make(chan *channelWriter, nprocs),
	}

	for i := 0; i < nprocs; i++ {
		go func() {
			for cw := range ctrl.c {
				cw.n, cw.e = cw.writer.Write(cw.b)
				cw.wg.Done()
			} // for each channel writer
		}()
	} // for each routine to create
	// We will only really need this when we are copying data anyway ... .
	return ctrl
}

// Return a new channel writer, which will write asynchronously to the given writer
func (w *WriteChannelController) NewChannelWriter(writer io.Writer) io.Writer {
	cw := channelWriter{
		ctrl:   w,
		writer: writer,
	}
	return &cw
}
