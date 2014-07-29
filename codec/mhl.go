// Implements the schema proposed by the media-hash-list program. It's an XML based codec

package codec

import (
	"bytes"
	"crypto/sha1"
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
	XMLName   xml.Name     `xml:"hashlist"`
	Version   string       `xml:"version,attr"`
	HashInfo  []mhlHash    `xml:"hash"`
	Signature mhlSignature `xml:"signature"`
}

type mhlSignature struct {
	XMLName xml.Name `xml:"signature"`
	Sha1    string   `xml:"sha1"`
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

// The inverse method of toFileInfo()
func (m *mhlHash) fromFileInfo(f *api.FileInfo) {
	m.File = f.RelaPath
	m.Size = f.Size
	m.Sha1 = fmt.Sprintf("%x", f.Sha1)
	m.Md5 = fmt.Sprintf("%x", f.MD5)
}

// toFileInfo copies our parsed XML values into the respective fields of the given FileInfo structure.
// Fields unavailable to h will be reset, and an error is returned if a value could not be parsed/converted to the
// actual type.
func (h *mhlHash) toFileInfo(f *api.FileInfo) error {
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
	enc := xml.NewEncoder(writer)
	enc.Indent("  ", "    ")
	sha1enc := sha1.New()
	writer.Write([]byte(xml.Header))
	hl := mhlHashList{
		Version: MHLVersion,
	}

	h := mhlHash{}
	for fi := range in {
		// Have to flatten the Path - after all, mhl has no support for absolute paths
		fi.Path = fi.RelaPath
		hashInfo(sha1enc, &fi)
		h.fromFileInfo(&fi)
		hl.HashInfo = append(hl.HashInfo, h)
	}

	hl.Signature.Sha1 = fmt.Sprintf("%x", sha1enc.Sum(nil))

	return enc.Encode(hl)
}

func (m *MHL) Deserialize(reader io.Reader, out chan<- api.FileInfo, predicate func(*api.FileInfo) bool) error {
	dec := xml.NewDecoder(reader)
	sha1enc := sha1.New()
	hl := mhlHashList{}

	if err := dec.Decode(&hl); err != nil {
		return &DecodeError{Msg: err.Error()}
	}

	if hl.Version != MHLVersion {
		return &DecodeError{Msg: fmt.Sprintf("Unsupported MHL version - got %s, want %s", hl.Version, MHLVersion)}
	}

	if len(hl.HashInfo) == 0 {
		return &DecodeError{Msg: "Didn't find a single hash in media hash list"}
	}

	// Otherwise, just stream all the pre-read data
	fi := api.FileInfo{}
	for _, h := range hl.HashInfo {
		if err := h.toFileInfo(&fi); err != nil {
			return err
		}
		// Bring back the Path, which is unset in XML, for the hashing to have something useful
		fi.Path = fi.RelaPath
		hashInfo(sha1enc, &fi)

		if !predicate(&fi) {
			return nil
		}
		out <- fi
	}

	// Yes, we do the check last, this way the user can at least see what might be wrong ... even though
	// the verify operations fails in the end ... .
	// This is disputable - if we know the file changed, the seal is broken and we have no reason to assume
	// we could find out anything different ... .
	if len(hl.Signature.Sha1) > 0 {
		if hls, err := hex.DecodeString(hl.Signature.Sha1); err != nil {
			return &DecodeError{Msg: fmt.Sprintf("Invalid signature fomat: %s", hl.Signature.Sha1)}
		} else if bytes.Compare(hls, sha1enc.Sum(nil)) != 0 {
			return &SignatureMismatchError{}
		}
	}

	return nil
}

func (m *MHL) Extension() string {
	return MHLExtension
}
