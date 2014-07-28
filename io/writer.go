package io

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
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
type ChannelWriter struct {
	ctrl *WriteChannelController

	// A writer to write to. Must be set if path is nil
	writer io.Writer

	// bytes to write - it's just a temporary
	b []byte
	// amount of bytes written
	n int
	// error of previous write operation
	e error

	// will let us know when reomte is done
	// NOTE: Would a channel be faster ?
	wg sync.WaitGroup
}

// Like WriteCloser interface, but allows to retrieve more information specific to our usage
type WriteCloser interface {
	io.WriteCloser

	// Writer returns the writer this interface instance contains
	Writer() io.Writer
}

func (c *ChannelWriter) Writer() io.Writer {
	return c.writer
}

// Set our writer to be the given one. Allows to reuse ChannelWriters
func (c *ChannelWriter) SetWriter(w io.Writer) {
	c.writer = w
}

// Send bytes down our channel and wait for the writer on the end to be done, retrieving the result.
func (c *ChannelWriter) Write(b []byte) (n int, err error) {
	c.b = b
	c.wg.Add(1)
	c.ctrl.c <- c
	c.wg.Wait()

	// ... allowing us to return the actual result safely now
	return c.n, c.e
}

func (c *ChannelWriter) Close() error {
	atomic.AddUint32(&c.ctrl.stats.TotalFilesWritten, 1)
	if w, ok := c.writer.(io.Closer); ok {
		return w.Close()
	}
	return nil
}

// A writer that will create a new file and intermediate directories on first write.
// You must call the close method to finish the writes and release system resources
type LazyFileWriteCloser struct {

	// The path we should open a writer to on first write. This will fail if the fail already exists.
	path string

	// The mode the destination file should have when done writing
	mode os.FileMode

	// A writer we are using to perform the write
	writer *os.File
}

// Path returns the currently set path
func (l *LazyFileWriteCloser) Path() string {
	return l.path
}

// SetPath changes the path to the given one.
// It's an error to set a new path while the previous writer wasn't closed yet
func (l *LazyFileWriteCloser) SetPath(p string, mode os.FileMode) {
	if l.writer != nil {
		panic("Previous writer wasn't close - can't set new path")
	}
	l.path = p
	l.mode = mode
}

func (l *LazyFileWriteCloser) Write(b []byte) (n int, err error) {
	if l.writer == nil {
		// assure directory exists
		err = os.MkdirAll(filepath.Dir(l.path), 0777)
		if err != nil {
			return 0, err
		}

		// NOTE: We may rightfully assume we see only one write in case of symlinks !
		// This is because the read-buffer is large enough to hold any symlink.
		// If not, the reader will panic

		// Symlinks are created right away
		if l.mode&os.ModeSymlink == os.ModeSymlink {
			err = os.Symlink(string(b), l.path)
			return len(b), err
		} else {
			l.writer, err = os.OpenFile(l.path, os.O_EXCL|os.O_WRONLY|os.O_CREATE, l.mode)
			if err != nil {
				return 0, err
			}
		}
	}

	return l.writer.Write(b)
}

// Close our writer if it was initialized already. Therefore it's safe to call this even if Write wasn't called
// beforehand
func (l *LazyFileWriteCloser) Close() error {
	if l.writer != nil {
		err := l.writer.Close()
		l.writer = nil
		return err
	}
	return nil
}

// A utility to help control how parallel we try to write
type WriteChannelController struct {
	// Keeps all write requests, which contain all information we could possibly want to write something.
	// As the ChannelWriter is keeping all information, we serves as request right away
	c chan *ChannelWriter

	// Allows to track amount of written files
	stats *Stats
}

// A utility structure to associate a tree with a writer.
// That way, writers can be more easily associated with a device which hosts a particular Tree
type RootedWriteController struct {
	// The trees the controller should write to
	Trees []string

	// A possibly shared controller which may write to the given tree
	Ctrl WriteChannelController
}

// Create a new controller which deals with writing all incoming requests with nprocs go-routines.
// Use the channel capacity to assure less blocking will occur. A good value is depending heavily on your
// algorithm's patterns. Should at least be nprocs, or larger.
func NewWriteChannelController(nprocs, channelCap int, stats *Stats) WriteChannelController {
	ctrl := WriteChannelController{
		make(chan *ChannelWriter, channelCap),
		stats,
	}
	if nprocs < 1 {
		panic("Need at least one go routine to process work")
	}

	for i := 0; i < nprocs; i++ {
		go func() {
			for cw := range ctrl.c {
				atomic.AddUint32(&stats.FilesBeingWritten, 1)
				cw.n, cw.e = cw.writer.Write(cw.b)
				atomic.AddUint64(&stats.BytesWritten, uint64(cw.n))
				atomic.AddUint32(&stats.FilesBeingWritten, ^uint32(0))
				// protocol mandates the sender has to listen for our reply, channel must not be closed here ... .
				cw.wg.Done()
			} // for each channel writer
		}()
	} // for each routine to create

	// We will only really need this when we are copying data anyway ... .
	return ctrl
}

// Initialize as many new ChannelWriters as fit into the given slice of writers
// You will have to set it's writer before using it
func (w *WriteChannelController) InitChannelWriters(out []ChannelWriter) {
	// create one writer per
	for i := range out {
		out[i] = ChannelWriter{
			ctrl: w,
		}
	}
}

// Return amount of streams we handle in parallel
func (w *WriteChannelController) Streams() int {
	return cap(w.c)
}

// Returns the amount of Trees/Destinations we can write to in total
// TODO(st): objectify
func WriteChannelDeviceMapTrees(wm []RootedWriteController) (n int) {
	for _, rctrl := range wm {
		n += len(rctrl.Trees)
	}
	return
}
