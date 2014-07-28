package codec

import (
	"path/filepath"
	"unicode/utf8"
)

func Names() []string {
	// I believe I have seen this somewhere - maybe it can be optimized to be constant ?
	names := [...]string{GobName, MHLName}
	return names[:]
}

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

// Find a codec matching the given name, and return it. Retuns nil otherwise
func NewByName(name string) Codec {
	switch {
	case name == GobName:
		return &Gob{}
	case name == MHLName:
		return &MHL{}
	}
	return nil
}
