package seal_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/cli"
	"github.com/Byron/godi/seal"
)

// Create a new file at given path and size, without possibly required intermediate directories
func makeFileOrPanic(path string, size int) string {
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
func makeDatasetOrPanic() (string, string, string) {
	base, err := ioutil.TempDir("", "dataset")
	if err != nil {
		panic(err)
	}

	makeFileOrPanic(filepath.Join(base, "1mb.ext"), 1024*1024)
	makeFileOrPanic(filepath.Join(base, "somebytes_noext"), 313)

	subdir := filepath.Join(base, "subdir")
	if err := os.Mkdir(subdir, 0777); err != nil {
		panic(err)
	}
	makeFileOrPanic(filepath.Join(subdir, "biggie.foo"), 1024*1024+5123)
	makeFileOrPanic(filepath.Join(subdir, "smallie.blah"), 123)
	subdir = filepath.Join(base, "nothing", "stillnothing", "ünicod€")
	if err := os.MkdirAll(subdir, 0777); err != nil {
		panic(err)
	}

	file := makeFileOrPanic(filepath.Join(subdir, "somefile.ext"), 12345)
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
func rmTree(tree string) {
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

func TestSeal(t *testing.T) {
	datasetTree, dataFile, _ := makeDatasetOrPanic()
	defer rmTree(datasetTree)
	var cmd *seal.SealCommand

	scmd, _ := cli.ParseArgs("seal", dataFile)
	cmd = scmd.(*seal.SealCommand)
	if err := cmd.SanitizeArgs(); err == nil {
		t.Error("Expected it to not like files as directory")
	} else {
		t.Log(err)
	}

	scmd, _ = cli.ParseArgs("seal", fmt.Sprintf("--num-readers=%v", runtime.GOMAXPROCS(0)), datasetTree)
	cmd = scmd.(*seal.SealCommand)
	if err := cmd.SanitizeArgs(); err != nil {
		t.Error("Sanitize didn't like existing tree")
	}

	var nprocs uint = uint(runtime.GOMAXPROCS(0))
	if nprocs > cmd.MaxProcs() {
		t.Error("Can't do less than one process here ... ")
	}

	results := make(chan api.Result, nprocs)
	done := make(chan bool)

	// assure we close our done channel on signal
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		close(done)
	}()

	files, generateResult := cmd.Generate(done)
	wg := sync.WaitGroup{}
	for i := 0; uint(i) < nprocs; i++ {
		wg.Add(1)
		go cmd.Gather(files, results, &wg, done)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	accumResult := cmd.Accumulate(results, done)

	// Return true if we should break the loop
	resHandler := func(name string, res api.Result) bool {
		if res == nil {
			// channel closed, have to get out
			t.Log("Channel", name, "is closed")
			return true
		}

		if res.Error() != nil {
			t.Error(res.Error())
		} else {
			t.Log(res.Info())
		}

		return false
	} // end resHandler

infinity:
	for {
		select {
		case r := <-generateResult:
			{
				if resHandler("generator", r) {
					break infinity
				}
			}
		case r := <-accumResult:
			{
				if resHandler("accumulator", r) {
					break infinity
				}
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Didn't get result after timeout")
		} // select
	} // endless loop
}
