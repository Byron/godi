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
func hashInfo(sha1enc hash.Hash, relaPath string, finfo *api.FileInfo) {
	sha1enc.Write([]byte(relaPath))
	sha1enc.Write([]byte(finfo.Path))
	sha1enc.Write(finfo.Sha1)
	sha1enc.Write(finfo.MD5)
}

func (g *Gob) Serialize(paths []SerializableFileInfo, writer io.Writer) (err error) {
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
	for _, finfo := range paths {
		hashInfo(sha1enc, finfo.RelaPath, &finfo.FileInfo)
		if err = encoder.Encode(finfo.FileInfo); err != nil {
			return
		}
	}

	if err = encoder.Encode(sha1enc.Sum(nil)); err != nil {
		return
	}

	return
}

func (g *Gob) Deserialize(reader io.Reader, out chan<- api.FileInfo, predicate func(*api.FileInfo) bool) error {
	gzipReader, _ := gzip.NewReader(reader)
	sha1enc := sha1.New()
	d := gob.NewDecoder(gzipReader)

	// Lets make the fields clear, and not reuse variables even if we could
	fileVersion := 0
	if err := d.Decode(&fileVersion); err != nil {
		return err
	}

	// Of course we would implement reading other formats too
	if fileVersion != Version {
		return fmt.Errorf("Cannot handle index file: invalid header version: %d", fileVersion)
	}

	numValues := 0
	if err := d.Decode(&numValues); err != nil {
		return err
	}

	for i := 0; i < numValues; i++ {
		// Yes - we need a fresh one every loop iteration ! Gob doesn't set fields which have the nil value
		v := api.FileInfo{}
		if err := d.Decode(&v); err != nil {
			return err
		}

		// Have to hash it before we hand it to the predicate, as it might alter the data
		hashInfo(sha1enc, v.RelaPath, &v)

		if !predicate(&v) {
			return nil
		}
		out <- v
	}

	var signature []byte
	if err := d.Decode(&signature); err != nil {
		return err
	}

	// Finally, compare signature of seal with the one we made ...
	if bytes.Compare(signature, sha1enc.Sum(nil)) != 0 {
		return errors.New("Signature mismatch")
	}

	return nil
}
