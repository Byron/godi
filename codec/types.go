package codec

import (
	"io"

	"github.com/Byron/godi/api"
)

// Represents a codec's standard capabilities.
// A codec is a specialized implementation able to read and write indices of file hash information
// NOTE: Even though it would be more idiomatic to have two interfaces for read and write respectively,
// we just don't need it here
type Codec interface {
	// Write the file-info structures received from in channel to the given writer.
	// When the channel is closed, the read-operation naturally finishes.
	// The codec must protect the written data against modification, usually by hashing the contained information
	Serialize(in <-chan api.FileInfo, writer io.Writer) error

	// Stream a FileInfo obtained from the given reader.
	// An error must be returned if the data read could not be verified.
	// You must pass each read fileinfo structure to the given predicate - it might modify it before
	// sending it down the channel. It returns false in case you should stop reading and return without error
	// Check the done-channel and cancel the operation
	// This function doesn't close the stream
	Deserialize(reader io.Reader, out chan<- api.FileInfo, predicate func(*api.FileInfo) bool) error

	// Extension returns the file extension of the codec, without the '.' prefix
	Extension() string
}
