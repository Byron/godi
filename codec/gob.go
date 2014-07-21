package codec

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/gob"
	"errors"
	"fmt"
	"hash"
	"io"

	"github.com/Byron/godi/api"
)

const (
	GobName      = "gob"
	GobExtension = "gobz"
	Version      = 1
)

// Reads and writes a file structured like so
// - version
// - numEntries
// - gobValue...
// - sha1 (hash of all hashes in prior map)
type Gob struct {
}

func (g *Gob) Extension() string {
	return GobExtension
}

// Take hashes of input arguments in predefined order
// NOTE: If order changes for some reason, we have to change the file version !
func hashInfo(sha1enc hash.Hash, relaPath string, finfo *godi.FileInfo) {
	sha1enc.Write([]byte(relaPath))
	sha1enc.Write([]byte(finfo.Path))
	sha1enc.Write(finfo.Sha1)
	sha1enc.Write(finfo.MD5)
}

func (g *Gob) Serialize(paths map[string]*godi.FileInfo, writer io.Writer) (err error) {
	gzipWriter, _ := gzip.NewWriterLevel(writer, 9)
	defer gzipWriter.Close()
	encoder := gob.NewEncoder(gzipWriter)

	sha1enc := sha1.New()

	if err = encoder.Encode(Version); err != nil {
		return
	}

	if err = encoder.Encode(len(paths)); err != nil {
		return
	}

	// NOTE: we re-encode to get rid of the map
	for relaPath, finfo := range paths {
		hashInfo(sha1enc, relaPath, finfo)
		if err = encoder.Encode(finfo); err != nil {
			return
		}
	}

	if err = encoder.Encode(sha1enc.Sum(nil)); err != nil {
		return
	}

	return
}

func (g *Gob) Deserialize(reader io.Reader) ([]godi.FileInfo, error) {
	gzipReader, _ := gzip.NewReader(reader)
	sha1enc := sha1.New()
	d := gob.NewDecoder(gzipReader)

	// Lets make the fields clear, and not reuse variables even if we could
	fileVersion := 0
	if err := d.Decode(&fileVersion); err != nil {
		return nil, err
	}

	// Of course we would implement reading other formats too
	if fileVersion != Version {
		return nil, fmt.Errorf("Cannot handle index file: invalid header version: %d", fileVersion)
	}

	numValues := 0
	if err := d.Decode(&numValues); err != nil {
		return nil, err
	}

	res := make([]godi.FileInfo, numValues)
	for i := 0; i < numValues; i++ {
		v := &res[i]
		if err := d.Decode(v); err != nil {
			return nil, err
		}

		hashInfo(sha1enc, v.RelaPath, v)
	}

	var signature []byte
	if err := d.Decode(&signature); err != nil {
		return nil, err
	}

	// Finally, compare signature of seal with the one we made ...
	if bytes.Compare(signature, sha1enc.Sum(nil)) != 0 {
		return nil, errors.New("Signature mismatch")
	}

	return res, nil
}
