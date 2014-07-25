package utility

// A structure to keep information about what is currently going on.
// It is means to be used as shared resource, used by multiple threads, which is why
// thread-safe counters are used.
// Implementations must keep these numbers up-to-date, while async processors will digest
// and present the data in some form
type Stats struct {
	FilesBeingRead    uint32 // Amount of files currently being read
	FilesBeingWritten uint32 // Amount of files currently being written
	BytesRead         uint64 // Total of bytes read so far, counting all input streams
	BytesWritten      uint64 // Total of bytes written so far, counting all output streams
	BytesHashed       uint64 // Total of bytes hashed so far, counting all active hashers
	NumHashers        uint32 // Amount of hashers running in parallel
}
