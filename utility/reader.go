package utility

import (
	"io"
	"os"
	"sync"
)

// Size for allocated buffers
const bufSize = 64 * 1024

// Actually, this must remain 0 for our sync to work, right now, without pool
const readChannelSize = 0

// The result of a read operation, similar to what Reader.Read returns
type readResult struct {
	buf []byte
	n   int
	err error
}

type ReadChannelController struct {
	c chan ChannelReader
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

	// Sync availability of our buffers, as we cam't use a sync.Pool due tot he slice truncation we have to do
	// p.Put(p.Get()[:5]) will put in a different slice, even though they point to the same array. It's not possible
	// to retrieve the underlying array.
	// When sending an array pointer directly, the same issue occurred. Don't know why ... hangs returned ...
	// Hanging occurred most of the time, but not always, which points at some async issue, maybe a bug in Pool ?
	wg *sync.WaitGroup

	buf *[bufSize]byte
}

// Return amount of streams we handle in parallel
func (r *ReadChannelController) Streams() int {
	return cap(r.c)
}

// Return a new channel reader
// You should set either path
func (r *ReadChannelController) NewChannelReaderFromPath(path string) ChannelReader {
	// NOTE: size of this channel controls how much we can cache into memory before we block
	// as the consumer doesn't keep up
	cr := ChannelReader{path, nil, make(chan readResult, readChannelSize),
		&sync.WaitGroup{},
		new([bufSize]byte),
	}

	r.c <- cr
	return cr
}

func (r *ReadChannelController) NewChannelReaderFromReader(reader io.Reader) ChannelReader {
	cr := ChannelReader{"", reader, make(chan readResult, readChannelSize),
		&sync.WaitGroup{},
		new([bufSize]byte),
	}

	r.c <- cr
	return cr
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
		// in any case, claim we are done with the result !
		p.wg.Done()
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
						info.wg.Add(1)
						info.results <- readResult{nil, 0, err}
						close(info.results)
						continue
					}
				}

				// Now read until it's done
				var nread int
				for {
					// The buffer will be put back by the one reading from the channel (e.g. in WriteTo()) !
					// wait until writer from previous iteration is done using the buffer
					info.wg.Wait()
					nread, err = info.reader.Read(info.buf[:])
					info.wg.Add(1)
					info.results <- readResult{info.buf[:nread], nread, err}
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

// NewReadChannelDeviceMap returns a mapping from each of the given trees to a controller which deals with the
// device the tree is on. If all trees are on the same device, you will get a map with len(trees) length, each one
// referring to the same controller
func NewReadChannelDeviceMap(nprocs int, trees []string) map[string]*ReadChannelController {
	dm := DeviceMap(trees)
	res := make(map[string]*ReadChannelController, len(dm))

	for _, dirs := range dm {
		rctrl := NewReadChannelController(nprocs)
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
