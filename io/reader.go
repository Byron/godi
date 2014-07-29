package io

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

// Size for allocated buffers
const bufSize = 32 * 1024

// Actually, this must remain 0 for our sync to work, right now, without pool
const readChannelSize = 0

// A utility structure to associate trees with a reader.
// NOTE: Very similar to RootedReadController !
type RootedReadController struct {
	// The trees the controller should write to
	Trees []string

	// A possibly shared controller which may write to the given tree
	Ctrl ReadChannelController
}

// The result of a read operation, similar to what Reader.Read returns
type readResult struct {
	buf []byte
	n   int
	err error
}

type ReadChannelController struct {
	c chan *ChannelReader
}

// Contains all information about a file or reader to be read
type ChannelReader struct {
	// An optional path, which will be opened for reading when Reader is nil
	path string

	// The mode of the file to read
	mode os.FileMode

	// A Reader interface, in case Path is unset. Use this if you want to open the file or provide your
	// own custom reader
	reader io.Reader

	// The channel to transport read results
	results chan readResult

	// Protects the buffer from simulateous access
	ready chan bool

	// Our buffer
	buf []byte
}

// Return amount of streams we handle in parallel
func (r *ReadChannelController) Streams() int {
	return cap(r.c)
}

// Return a new channel reader
// You should set either path
// The buffer must not be shared among multiple channel readers !
func (r *ReadChannelController) NewChannelReaderFromPath(path string, mode os.FileMode, buf []byte) *ChannelReader {
	// NOTE: size of this channel controls how much we can cache into memory before we block
	// as the consumer doesn't keep up
	cr := ChannelReader{
		path:    path,
		mode:    mode,
		buf:     buf,
		results: make(chan readResult, readChannelSize),
		ready:   make(chan bool),
	}

	r.c <- &cr
	return &cr
}

func (r *ReadChannelController) NewChannelReaderFromReader(reader io.Reader, buf []byte) *ChannelReader {
	cr := ChannelReader{
		reader:  reader,
		buf:     buf,
		results: make(chan readResult, readChannelSize),
		ready:   make(chan bool),
	}

	r.c <- &cr
	return &cr
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

	// initial ready indicator - now remote reader produces result
	p.ready <- true
	// We will receive results until the other end is done reading
	for res := range p.results {
		// Write what's possible - don't check for 0, as we also have to deal with empty files
		// Without the write call, they wouldn't be created after all.
		written, err = w.Write(res.buf)
		n += int64(written)

		// now we are ready for the next one

		// This would block as the remote will stop sending results on error
		if res.err == nil {
			p.ready <- true
		} else {
			if res.err != io.EOF {
				err = res.err
			}
		}

		// in any case, claim we are done with the result !
		if res.n == 0 && res.err == nil {
			panic("If 0 bytes have been read, there should at least be an EOF (in case of empty files)")
		}
	} // for each read result

	// whatever is held in n, err, we return
	return
}

// Create a new parallel reader with nprocs go-routines and return a channel to it.
// Feed the channel with ChannelReader structures and listen on it's channel to read bytes until EOF, which
// is when the channel will be closed by the reader
// done will allow long reads to be interrupted by closing the channel
func NewReadChannelController(nprocs int, stats *Stats, done <-chan bool) ReadChannelController {
	if nprocs < 1 {
		panic("nprocs must be >= 1")
	}

	ctrl := ReadChannelController{
		make(chan *ChannelReader, nprocs),
	}

	reader := func(info *ChannelReader) {
		// in any case, close the results channel
		defer close(info.results)
		defer close(info.ready)

		sendError := func(err error) {
			// Add one - the client reader will call Done after receiving our result
			// We are always required to signal ready before we send
			<-info.ready
			info.results <- readResult{nil, 0, err}
		}

		var err error
		ourReader := false
		if info.reader == nil {
			if info.mode&os.ModeSymlink == os.ModeSymlink {
				ldest, err := os.Readlink(info.path)
				if err != nil {
					sendError(err)
					return
				} else {
					// The contents of the link is our result - therefore, we finish it here
					<-info.ready
					atomic.AddUint64(&stats.BytesRead, uint64(len(ldest)))

					if n := copy(info.buf, []byte(ldest)); n != len(ldest) {
						panic("Couldn't copy symlink into buffer - was it larger than our buffer ??")
					}

					info.results <- readResult{info.buf[:len(ldest)], len(ldest), io.EOF}
					return
				}
			} else {
				ourReader = true
				info.reader, err = os.Open(info.path)
				if err != nil {
					sendError(err)
					return
				}
			}
		}

		// Now read until it's done
		var nread int

	readForever:
		for {
			// The buffer will be put back by the one reading from the channel (e.g. in WriteTo()) !
			// wait until writer from previous iteration is done using the buffer
			// Have to ask for it in any case - if we quit this loop, the receiver may stall otherwise
			<-info.ready
			select {
			case <-done:
				{
					var err error
					if ourReader {
						err = fmt.Errorf("Reading of '%s' cancelled", info.path)
					} else {
						err = errors.New("Reading cancelled by user")
					}
					info.results <- readResult{err: err}
					break readForever
				}
			default:
				{
					nread, err = info.reader.Read(info.buf)
					atomic.AddUint64(&stats.BytesRead, uint64(nread))
					info.results <- readResult{info.buf[:nread], nread, err}
					// we send all results, but abort if the reader is done for whichever reason
					if err != nil {
						break readForever
					}
				}
			} // end select
		} // readForever

		if ourReader {
			info.reader.(*os.File).Close()
			info.reader = nil
		}
	}

	for i := 0; i < nprocs; i++ {
		go func() {
			for info := range ctrl.c {
				atomic.AddUint32(&stats.FilesBeingRead, uint32(1))
				reader(info)
				atomic.AddUint32(&stats.FilesBeingRead, ^uint32(0))
				atomic.AddUint32(&stats.TotalFilesRead, uint32(1))
			}
		}()
	}

	return ctrl
}

// A new list of Controllers, one per device it handles, which is associated with the tree's it can handle
func NewDeviceReadControllers(nprocs int, trees []string, stats *Stats, done <-chan bool) []RootedReadController {
	dm := DeviceMap(trees)
	res := make([]RootedReadController, len(dm))

	for did, trees := range dm {
		// each device as so and so many sources. Each source uses the same read controller
		res[did] = RootedReadController{
			Trees: trees,
			Ctrl:  NewReadChannelController(nprocs, stats, done),
		}
	} // for each tree set in deviceMap

	return res
}

// NOTE: Can this be a custom type, with just a function ? I think so !
// Return the number of streams being handled in parallel
// TODO(st) objectify
func ReadChannelDeviceMapStreams(rctrls []RootedReadController) int {
	if len(rctrls) == 0 {
		panic("Input map was empty")
	}

	nstreams := 0
	for _, rctrl := range rctrls {
		nstreams += rctrl.Ctrl.Streams()
	}

	return nstreams
}
