package utility

import (
	"errors"
	"io"
	"os"
	"sync"
)

// Size for allocated buffers
const bufSize = 32 * 1024

// Actually, this must remain 0 for our sync to work, right now, without pool
const readChannelSize = 0

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

	// A Reader interface, in case Path is unset. Use this if you want to open the file or provide your
	// own custom reader
	reader io.Reader

	// The channel to transport read results
	results chan readResult

	// Protects the buffer from simulateous access
	exclusive sync.Mutex

	// Our buffer
	buf [bufSize]byte
}

// Return amount of streams we handle in parallel
func (r *ReadChannelController) Streams() int {
	return cap(r.c)
}

// Return a new channel reader
// You should set either path
func (r *ReadChannelController) NewChannelReaderFromPath(path string) *ChannelReader {
	// NOTE: size of this channel controls how much we can cache into memory before we block
	// as the consumer doesn't keep up
	cr := ChannelReader{
		path:    path,
		results: make(chan readResult, readChannelSize),
	}

	r.c <- &cr
	return &cr
}

func (r *ReadChannelController) NewChannelReaderFromReader(reader io.Reader) *ChannelReader {
	cr := ChannelReader{
		reader:  reader,
		results: make(chan readResult, readChannelSize),
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
	for res := range p.results {
		// Write what's possible
		if res.n > 0 {
			p.exclusive.Lock()
			written, err = w.Write(res.buf)
			p.exclusive.Unlock()
			n += int64(written)
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
func NewReadChannelController(nprocs int, done <-chan bool) ReadChannelController {
	if nprocs < 1 {
		panic("nprocs must be >= 1")
	}

	ctrl := ReadChannelController{
		make(chan *ChannelReader, nprocs),
	}

	infoHandler := func(info *ChannelReader) {
		// in any case, close the results channel
		defer close(info.results)

		var err error
		ourReader := false
		if info.reader == nil {
			info.reader, err = os.Open(info.path)
			if err != nil {
				// Add one - the client reader will call Done after receiving our result
				info.results <- readResult{nil, 0, err}
				return
			}
		}

		// Now read until it's done
		var nread int
	readForever:
		for {
			select {
			case <-done:
				{
					info.results <- readResult{err: errors.New("Cancelled by user")}
					break readForever
				}
			default:
				{
					// The buffer will be put back by the one reading from the channel (e.g. in WriteTo()) !
					// wait until writer from previous iteration is done using the buffer
					info.exclusive.Lock()
					nread, err = info.reader.Read(info.buf[:])
					info.exclusive.Unlock()
					info.results <- readResult{info.buf[:nread], nread, err}
					// we send all results, but abort if the reader is done for whichever reason
					if err != nil {
						break readForever
					}
				}
			} // end select
		} // read loop

		if ourReader {
			info.reader.(*os.File).Close()
			info.reader = nil
		}
	}

	for i := 0; i < nprocs; i++ {
		go func() {
			for info := range ctrl.c {
				infoHandler(info)
			}
		}()
	}

	return ctrl
}

// NewReadChannelDeviceMap returns a mapping from each of the given trees to a controller which deals with the
// device the tree is on. If all trees are on the same device, you will get a map with len(trees) length, each one
// referring to the same controller
func NewReadChannelDeviceMap(nprocs int, trees []string, done <-chan bool) map[string]*ReadChannelController {
	dm := DeviceMap(trees)
	res := make(map[string]*ReadChannelController, len(dm))

	for _, dirs := range dm {
		rctrl := NewReadChannelController(nprocs, done)
		for _, dir := range dirs {
			res[dir] = &rctrl
		}
	}

	return res
}

// NOTE: Can this be a custom type, with just a function ? I think so !
// Return the number of streams being handled in parallel
func ReadChannelDeviceMapStreams(rm map[string]*ReadChannelController) int {
	if len(rm) == 0 {
		panic("Input map was empty")
	}

	nstreams := 0
	// count unique controllers to figure out stream multiplier
	seen := make([]*ReadChannelController, 0, len(rm))

	for _, ctrl := range rm {
		cseen := false
		for _, c := range seen {
			if c == ctrl {
				cseen = true
				break
			}
		}
		if !cseen {
			seen = append(seen, ctrl)
			nstreams += ctrl.Streams()
		}
	}

	return nstreams
}
