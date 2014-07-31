// Package testlib provides varous utilities for testing purposes
package testlib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	FirstSubDir = "nothing"
)

// Create a new file at given path and size, without possibly required intermediate directories
func MakeFileOrPanic(path string, size int) string {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	if size != 0 {
		b := [1]byte{0}
		f.WriteAt(b[:], int64(size-1))
	}

	return path
}

// Create a dataset for testing and return the newly created directory
func MakeDatasetOrPanic() (string, string, string) {
	base, err := ioutil.TempDir("", "dataset")
	if err != nil {
		panic(err)
	}

	MakeFileOrPanic(filepath.Join(base, "1mb.ext"), 1024*1024)
	MakeFileOrPanic(filepath.Join(base, "somebytes_noext"), 313)

	subdir := filepath.Join(base, "subdir")
	if err := os.Mkdir(subdir, 0777); err != nil {
		panic(err)
	}
	MakeFileOrPanic(filepath.Join(subdir, "biggie.foo"), 1024*1024+5123)
	MakeFileOrPanic(filepath.Join(subdir, "smallie.blah"), 123)
	MakeFileOrPanic(filepath.Join(subdir, "empty.file"), 0)
	subdir = filepath.Join(base, FirstSubDir, "stillnothing", "ünicod€")
	if err := os.MkdirAll(subdir, 0777); err != nil {
		panic(err)
	}

	file := MakeFileOrPanic(filepath.Join(subdir, "somefile.ext"), 12345)
	symlink := filepath.Join(base, "symlink.ext")
	err = os.Symlink(file, symlink)
	if err != nil {
		symlink = ""
	}
	return base, file, symlink
}

// Delete the given tree entirely. Should only be used in conjunction with makeDataset
// Panics if something is wrong
// Will only do the work if we are not already in panic
func RmTree(tree string) {
	return
	if len(tree) == 0 {
		panic("Invalid tree given")
	}
	res := recover()
	if res != nil {
		fmt.Fprintf(os.Stderr, "Keeping tree for debugging at '%s'", tree)
		panic(res)
	}
	if err := os.RemoveAll(tree); err != nil {
		panic(err)
	}
}
