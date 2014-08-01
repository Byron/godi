// The cosmic tool helps to corrupt data.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"

	"code.google.com/p/go.crypto/ssh/terminal"
)

const (
	usage = `cosmic [--force] file [file ...]

Flip a single bit randombly in the given file.
Unless force is specified, it will ask the user for verification, or fail if we are not a tty
`
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Flip a single bit in the given file, changing it without backup.
// Return any error we encounter on the way.
func cosmicRayOn(file string) (err error) {
	fd, err := os.OpenFile(file, os.O_RDWR, 0)
	if err != nil {
		return
	}
	flen, err := fd.Seek(0, os.SEEK_END)
	if err != nil {
		return
	}

	// can't change empty files
	if flen == 0 {
		return
	}

	// jump somewhere
	lpos, err := fd.Seek(rand.Int63n(flen-1), os.SEEK_SET)
	if err != nil {
		return
	}

	// read a byte
	var b [1]byte
	if _, serr := fd.Read(b[:]); err != nil {
		return serr
	}

	// flip the bit and write it back
	b[0] = b[0] ^ (1 << uint(rand.Intn(8)))
	if _, werr := fd.WriteAt(b[:], lpos); werr != nil {
		return werr
	}

	fd.Close()
	return
}

func usageAndDie() {
	os.Stderr.WriteString(usage)
	os.Exit(1)
}

func main() {
	force := false

	flag.CommandLine.Init("command-line", flag.ContinueOnError)
	flag.BoolVar(&force, "force", force, "--force is required to flip a bit without a tty, or to prevent asking")
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		usageAndDie()
	}

	if len(flag.Args()) == 0 {
		usageAndDie()
	}

	if !force {
		if !terminal.IsTerminal(syscall.Stdin) {
			fmt.Fprint(os.Stderr, "Without a tty, --force must be specified")
			os.Exit(3)
		}

		// ask if we should really do it
		answer := ""
		fmt.Printf("The given file(s) %s will be changed randomly - type 'yes' to continue: ", strings.Join(flag.Args(), ", "))
		fmt.Scan(&answer)
		if answer != "yes" {
			fmt.Fprintln(os.Stderr, "Aborted by user")
			os.Exit(2)
		}
	}

	errCount := 0
	for _, file := range flag.Args() {
		if err := cosmicRayOn(file); err == nil {
			fmt.Printf("A cosmic ray hit '%s' and flipped a bit\n", file)
		} else {
			errCount += 1
			fmt.Fprintln(os.Stderr, err)
		}
	}

	os.Exit(errCount)

}
