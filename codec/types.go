package codec

import (
	"io"

	"github.com/Byron/godi/api"
)

type SerializableFileInfo struct {
	godi.FileInfo

	// The error associated with the file, usually read errors
	// This can be expected to be unset if the structure should be written
	Err error
}

// Helper to sort by longest path, descending
type ByLongestPathDescending []SerializableFileInfo

func (a ByLongestPathDescending) Len() int           { return len(a) }
func (a ByLongestPathDescending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLongestPathDescending) Less(i, j int) bool { return len(a[i].Path) > len(a[j].Path) }

// Represents a codec's standard capabilities.
// A codec is a specialized implementation able to read and write indices of file hash information
// NOTE: Even though it would be more idiomatic to have two interfaces for read and write respectively,
// we just don't need it here
type Codec interface {
	// Write the given FileInfo structure to the given writer.
	// The codec must protect the written data against modification, usually by hashing the contained information
	Serialize(paths []SerializableFileInfo, writer io.Writer) (err error)

	// Read a FileInfo slice from the given reader. The fileinfo Paths must be relative to the index file
	// An error must be returned if the data read could not be verified.
	Deserialize(reader io.Reader) ([]godi.FileInfo, error)

	// Extension returns the file extension of the codec, without the '.' prefix
	Extension() string
}
