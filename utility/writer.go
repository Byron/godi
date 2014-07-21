package utility

import (
	"io"
	"sync"
)

// Similar to MultiWriter, but assumes writes never fail, and provides the same buffer
// to all writers in parallel. However, it will return only when all writes are finished
type UncheckedParallelMultiWriter struct {
	Writers []io.Writer
	wg      sync.WaitGroup
}

func (t *UncheckedParallelMultiWriter) Write(p []byte) (n int, err error) {
	t.wg.Add(len(t.Writers))
	for _, w := range t.Writers {
		go func(w io.Writer) {
			w.Write(p)
			t.wg.Done()
		}(w)
	}
	t.wg.Wait()
	return len(p), nil
}

func NewUncheckedParallelMultiWriter(writers ...io.Writer) UncheckedParallelMultiWriter {
	w := make([]io.Writer, len(writers))
	copy(w, writers)
	return UncheckedParallelMultiWriter{Writers: w}
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

type WriteChannelController struct {
	c chan ChannelWriter
	// if true, we will copy all incoming buffers to assure we own them
	takeOwnership bool

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
		make(chan ChannelWriter, nprocs),
		false,
	}

	// TODO: implementation
	// We will only really need this when we are copying data anyway ... .
	return ctrl
}
