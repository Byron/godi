package codec

import (
	"io"

	"github.com/Byron/godi/api"
)

// Used as structure when deserializing information
// It contains both absolute and relative paths to allow dealing with relocated trees, or relocated seal files
type SerializedFileInfo struct {
	FileInfo godi.FileInfo
	RelaPath string
}

// Represents a codec's standard capabilities.
// A codec is a specialized implementation able to read and write indices of file hash information
// NOTE: Even though it would be more idiomatic to have two interfaces for read and write respectively,
// we just don't need it here
type Codec interface {
	// Write the given FileInfo structure to the given writer.
	// The codec must protect the written data against modification, usually by hashing the contained information
	Serialize(paths map[string]*godi.FileInfo, writer io.Writer) (err error)

	// Read a FileInfo slice from the given reader. The fileinfo Paths must be relative to the index file
	// An error must be returned if the data read could not be verified.
	Deserialize(reader io.Reader) ([]SerializedFileInfo, error)

	// Extension returns the file extension of the codec, without the '.' prefix
	Extension() string
}
