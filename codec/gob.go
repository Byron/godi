package codec

import (
	"compress/gzip"
	"crypto/sha1"
	"encoding/gob"
	"io"

	"github.com/Byron/godi/api"
)

const (
	Name      = "gob"
	extension = "gobz"
	version   = 1
)

// Used in the serialization format
type gobValueV1 struct {
	RelaPath string
	FileInfo *api.FileInfo
}

// Reads and writes a file structured like so
// - version
// - numEntries
// - gobValue...
// - sha1 (hash of all hashes in prior map)
type Gob struct {
}

func (g *Gob) Extension() string {
	return extension
}

func (g *Gob) Serialize(paths map[string]*api.FileInfo, writer io.Writer) (err error) {
	gzipWriter, _ := gzip.NewWriterLevel(writer, 9)
	defer gzipWriter.Close()
	encoder := gob.NewEncoder(gzipWriter)

	sha1enc := sha1.New()

	if err = encoder.Encode(version); err != nil {
		return
	}

	if err = encoder.Encode(len(paths)); err != nil {
		return
	}

	// NOTE: we re-encode to get rid of the map
	for relaPath, finfo := range paths {
		sha1enc.Write(finfo.Sha1)
		sha1enc.Write(finfo.MD5)
		if err = encoder.Encode(gobValueV1{relaPath, finfo}); err != nil {
			return
		}
	}

	if err = encoder.Encode(paths); err != nil {
		return
	}

	if err = encoder.Encode(sha1enc.Sum(nil)); err != nil {
		return
	}

	return
}
