// Implements the schema proposed by the media-hash-list program. It's an XML based codec

package codec

import (
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/Byron/godi/api"
)

const (
	MHLName      = "mhl"
	MHLExtension = "mhl"
	MHLVersion   = "1.0"
)

type mhlHashList struct {
	XMLName  xml.Name  `xml:"hashlist"`
	Version  string    `xml:"version,attr"`
	HashInfo []mhlHash `xml:"hash"`
}

type mhlHash struct {
	XMLName xml.Name `xml:"hash"`
	File    string   `xml:"file"`
	Size    int64    `xml:"size"`
	// MTimeString string `xml:"lastmodificationdate"`
	Sha1 string `xml:"sha1"`
	Md5  string `xml:"md5"`
	// HashDate    string `xml:"hashdate"`
}

// ToFileInfo copies our parsed XML values into the respective fields of the given FileInfo structure.
// Fields unavailable to h will be reset, and an error is returned if a value could not be parsed/converted to the
// actual type.
func (h *mhlHash) ToFileInfo(f *api.FileInfo) error {
	if len(h.File) == 0 {
		return errors.New("Empty file field")
	}
	f.Path = h.File // This is actually just a relative path, but let's set it anyway
	f.RelaPath = h.File

	if h.Size < 0 {
		return fmt.Errorf("size of '%s' must not be smaller than 0", h.File)
	}
	f.Size = h.Size

	if len(h.Md5) > 0 {
		if md5, err := hex.DecodeString(h.Md5); err != nil {
			return fmt.Errorf("Failed to parse MD5 hash of '%s' with error: %s", h.File, err.Error())
		} else if len(md5) != 16 {
			return fmt.Errorf("Invalid MD5 hash length in '%s'. Expected 16, got %d", h.File, len(md5))
		} else {
			f.MD5 = md5
		}
	} else {
		f.MD5 = nil
	}

	if len(h.Sha1) > 0 {
		if sha1, err := hex.DecodeString(h.Sha1); err != nil {
			return fmt.Errorf("Failed to parse Sha1 hash of '%s' with error: %s", h.File, err.Error())
		} else if len(sha1) != 20 {
			return fmt.Errorf("Invalid Sha1 hash length in '%s'. Expected 20, got %d", h.File, len(sha1))
		} else {
			f.Sha1 = sha1
		}
	} else {
		f.Sha1 = nil
	}

	if len(f.MD5) == 0 && len(f.Sha1) == 0 {
		return fmt.Errorf("Didn't parse a single hash for file '%s'", h.File)
	}

	return nil
}

// Empty type to implement the codec interface
type MHL struct {
}

func (m *MHL) Serialize(in <-chan api.FileInfo, writer io.Writer) (err error) {
	return errors.New("Not implemented")
}

func (m *MHL) Deserialize(reader io.Reader, out chan<- api.FileInfo, predicate func(*api.FileInfo) bool) error {
	dec := xml.NewDecoder(reader)
	hl := mhlHashList{}

	if err := dec.Decode(&hl); err != nil {
		return err
	}

	if hl.Version != MHLVersion {
		return fmt.Errorf("Unsupported MHL version - got %s, want %s", hl.Version, MHLVersion)
	}

	if len(hl.HashInfo) == 0 {
		return errors.New("Didn't find a single hash in media hash list")
	}

	// Otherwise, just stream all the pre-read data
	fi := api.FileInfo{}
	for _, h := range hl.HashInfo {
		if err := h.ToFileInfo(&fi); err != nil {
			return err
		}

		if !predicate(&fi) {
			return nil
		}
		out <- fi
	}

	return nil
}

func (m *MHL) Extension() string {
	return MHLExtension
}
