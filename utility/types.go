package utility

// A structure to keep information about what is currently going on.
// It is means to be used as shared resource, used by multiple threads, which is why
// thread-safe counters are used.
// Implementations must keep these numbers up-to-date, while async processors will digest
// and present the data in some form
type Stats struct {
	FilesBeingRead    int32 // Amount of files currently being read
	FilesBeingWritten int32 // Amount of files currently being written
	BytesRead         int64 // Total of bytes read so far, counting all input streams
	BytesWritten      int64 // Total of bytes written so far, counting all output streams
	BytesHashed       int64 // Total of bytes hashed so far, counting all active hashers
}
