package codec

import (
	"path/filepath"
	"unicode/utf8"
)

// Finds a codec which can decode the file at the given path.
// We work strictly by name.
func NewByPath(path string) Codec {
	ext := filepath.Ext(path)

	// '.' as extension or no extension
	if len(ext) < 2 {
		return nil
	}

	// I know, we have ascii, even if there are UTF8 characters ... let's just practice sensible string handling
	_, s := utf8.DecodeRuneInString(ext)
	ext = ext[s:]

	switch ext {
	case GobExtension:
		return &Gob{}
	case MHLExtension:
		return &MHL{}
	}

	return nil
}
