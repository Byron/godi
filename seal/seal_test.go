package seal_test

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/seal"
	"github.com/Byron/godi/testlib"
)

func TestSeal(t *testing.T) {
	datasetTree, dataFile, _ := testlib.MakeDatasetOrPanic()
	defer testlib.RmTree(datasetTree)
	var cmd seal.SealCommand
	var err error

	_, err = seal.NewCommand([]string{dataFile}, 1)
	if err == nil {
		t.Error("Expected it to not like files as directory")
	} else {
		t.Log(err)
	}

	cmd, err = seal.NewCommand([]string{datasetTree, filepath.Join(datasetTree, testlib.FirstSubDir, "..", testlib.FirstSubDir)}, 1)
	if err != nil {
		t.Error("Expected to not fail sanitization")
	} else if len(cmd.Trees) != 1 {
		t.Error("Trees should have been pruned, contained one should have been dropped")
	}

	cmd, err = seal.NewCommand([]string{datasetTree}, runtime.GOMAXPROCS(0))
	if err != nil {
		t.Error("Sanitize didn't like existing tree")
	}

	var nprocs uint = uint(runtime.GOMAXPROCS(0))
	if nprocs > cmd.MaxProcs() {
		t.Error("Can't do less than one process here ... ")
	}

	results := make(chan godi.Result, nprocs)
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
	accumResult := cmd.Aggregate(results, done)

	// Return true if we should break the loop
	resHandler := func(name string, res godi.Result) bool {
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
