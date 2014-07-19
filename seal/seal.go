package seal

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
)

const (
	IndexBaseName = "godi"
	Name          = "seal"
)

// A type representing all arguments required to drive a Seal operation
type SealCommand struct {

	// One or more trees to seal
	Trees []string

	// Amount of readers to use
	nReaders int

	// parallel reader
	pCtrl api.ReadChannelController
}

// REVIEW:
func NewCommand(trees []string, nReaders int) SealCommand {
	c := SealCommand{}
	c.Trees = trees
	c.nReaders = nReaders
	return c
}

// Implements information about a seal operation
type SealResult struct {
	finfo *api.FileInfo
	msg   string
	err   error
	prio  api.Priority
}

func (s *SealResult) Info() (string, api.Priority) {
	if s.err != nil {
		return s.err.Error(), api.Error
	}
	return s.msg, s.prio
}

func (s *SealResult) Error() error {
	return s.err
}

func (s *SealResult) FileInformation() *api.FileInfo {
	return s.finfo
}

func (s *SealCommand) SetUnparsedArgs(args []string) error {
	s.Trees = args
	return nil
}

func (s *SealCommand) MaxProcs() uint {
	return uint(math.MaxUint32)
}

func (s *SealCommand) SanitizeArgs() (err error) {
	if len(s.Trees) == 0 {
		return errors.New("Please provide at least one tree to work on")
	}

	invalidTrees := make([]string, 0, len(s.Trees))
	noTrees := make([]string, 0, len(s.Trees))
	for i, tree := range s.Trees {
		if stat, err := os.Stat(tree); err != nil {
			invalidTrees = append(invalidTrees, tree)
		} else if !stat.IsDir() {
			noTrees = append(noTrees, tree)
		}
		s.Trees[i] = path.Clean(tree)
	}

	if len(invalidTrees) > 0 {
		return errors.New("Coulnd't read at least one of the given trees to verify: " + strings.Join(invalidTrees, ", "))
	}
	if len(noTrees) > 0 {
		return errors.New("The following trees are no directory: " + strings.Join(noTrees, ", "))
	}
	if s.nReaders < 1 {
		return errors.New("--num-readers must not be smaller than 1")
	}

	// drop trees which are a sub-tree of another
	if len(s.Trees) > 1 {
		validTrees := make([]string, 0, len(s.Trees))
		for i, ltree := range s.Trees {
			for x, rtree := range s.Trees {
				if i == x || strings.HasPrefix(ltree, rtree) {
					continue
				}
				validTrees = append(validTrees, ltree)
			}
		}
		if len(validTrees) == 0 {
			panic("Didn't find a single valid tree")
		}

		s.Trees = validTrees
	}

	return err
}

func (s *SealCommand) SetupParser(parser *flag.FlagSet) error {
	parser.IntVar(&s.nReaders, "num-readers", 1, "Amount of parallel read streams we can use")
	return nil
}

func (s *SealCommand) Generate(done <-chan bool) (<-chan api.FileInfo, <-chan api.Result) {
	files := make(chan api.FileInfo)
	results := make(chan api.Result)

	s.pCtrl = api.NewReadChannelController(s.nReaders)

	go func() {
		for _, tree := range s.Trees {
			if !s.traverseFilesRecursively(files, results, done, tree) {
				// interrupted usually, or there was an error
				break
			}
		}
		defer close(files)
	}()

	return files, results
}

// Traverse recursively, return false if the caller should stop traversing due to an error
func (s *SealCommand) traverseFilesRecursively(files chan<- api.FileInfo, results chan<- api.Result, done <-chan bool, tree string) bool {
	select {
	case <-done:
		return false
	default:
		{
			// read dir and, build file info, and recurse into subdirectories
			f, err := os.Open(tree)
			if err != nil {
				results <- &SealResult{nil, "", err, api.Error}
				return false
			}
			dirInfos, err := f.Readdir(-1)
			f.Close()
			if err != nil {
				results <- &SealResult{nil, "", err, api.Error}
				return false
			}

			// first generate infos
			for _, fi := range dirInfos {
				if !fi.IsDir() {
					path := filepath.Join(tree, fi.Name())
					if !fi.Mode().IsRegular() {
						results <- &SealResult{nil, fmt.Sprintf("Ignoring symbolic link: '%s'", path), nil, api.Warn}
						continue
					}
					if fi.Name()[0] == '.' {
						results <- &SealResult{nil, fmt.Sprintf("Ignoring hidden file: '%s'", path), nil, api.Warn}
						continue
					}
					files <- api.FileInfo{
						Path: path,
						Size: fi.Size(),
					}
				}
			}

			// then recurse
			for _, fi := range dirInfos {
				if fi.IsDir() {
					if !s.traverseFilesRecursively(files, results, done, filepath.Join(tree, fi.Name())) {
						return false
					}
				}
			}
		} //  default
	} // selcect

	return true
}

func (s *SealCommand) Gather(files <-chan api.FileInfo, results chan<- api.Result, wg *sync.WaitGroup, done <-chan bool) {
	defer wg.Done()
	sha1gen := sha1.New()
	md5gen := md5.New()
	// This makes the write as slow as the slowest hash, instead of hash+hash
	allHashes := api.UncheckedParallelMultiWriter(sha1gen, md5gen)

	// This MUST be a copy of f here, otherwise we will be in trouble thanks to the user of defer in handleHash
	// we will get f overwritten by the next iteration variable ... it's kind of special, might
	// be intersting for the mailing list.
	handleHash := func(f api.FileInfo) {
		res := SealResult{&f, "", nil, api.Progress}
		err := &res.err
		defer func(res *SealResult) { results <- res }(&res)

		// let the other end open the file and close it as well
		reader := s.pCtrl.NewChannelReaderFromPath(f.Path)
		s.pCtrl.Channel() <- reader

		sha1gen.Reset()
		md5gen.Reset()
		var written int64
		written, *err = reader.WriteTo(allHashes)
		if *err != nil {
			return
		}
		f.Sha1 = sha1gen.Sum(nil)
		f.MD5 = md5gen.Sum(nil)
		if written != f.Size {
			*err = fmt.Errorf("Filesize of '%s' reported as %d, yet only %d bytes were hashed", f.Path, f.Size, written)
			return
		}
	}

	for f := range files {
		select {
		case <-done:
			return
		default:
			handleHash(f)
		}
	}
}

// return a path to an index file residing at tree
func (s *SealCommand) IndexPath(tree string, extension string) string {
	n := time.Now()
	return filepath.Join(tree, fmt.Sprintf("%s_%04d-%02d-%02d_%02d%02d%02d.%s",
		IndexBaseName,
		n.Year(),
		n.Month(),
		n.Day(),
		n.Hour(),
		n.Minute(),
		n.Second(),
		extension))
}

// When called, we have seen no error and everything can be assumed to be in order
// Returns error in case we can't write, and all index files written so far.
// It's up to the caller to remove existing files on error
func (s *SealCommand) writeIndex(treeMap map[string]map[string]*api.FileInfo) ([]string, error) {
	// Serialize all fileinfo structures
	// NOTE: As the parallel writer will send results only when writing finished, we can just operate serially here ...
	// For this there is also no need to optimize performance
	// However, we could use a parallel writer if we are so inclined
	// For now, we do only gob
	encoder := codec.Gob{}
	indices := make([]string, 0, len(treeMap))
	for tree, paths := range treeMap {
		// This will and should fail if the file already exists
		fd, err := os.OpenFile(s.IndexPath(tree, encoder.Extension()), os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return indices, err
		}

		err = encoder.Serialize(paths, fd)
		fd.Close()
		if err != nil {
			// remove possibly half-written file
			os.Remove(fd.Name())
			return indices, err
		}
		indices = append(indices, fd.Name())
	}

	return indices, nil
}

func (s *SealCommand) Accumulate(results <-chan api.Result, done <-chan bool) <-chan api.Result {
	accumResult := make(chan api.Result)

	go func() {
		defer close(accumResult)
		treePathmap := make(map[string]map[string]*api.FileInfo)

		// Presort all paths by their root
		for _, tree := range s.Trees {
			treePathmap[tree] = make(map[string]*api.FileInfo)
		}

		var count uint = 0
		var errCount uint = 0
		var size uint64 = 0
		st := time.Now()
		wasCancelled := false

		// ACCUMULATE PATHS INFO
		/////////////////////////
		for r := range results {
			if r.Error() != nil {
				errCount += 1
				accumResult <- r
			}

			// Be sure we take note of cancellation.
			// If this happens, soon our results will be drained and we leave naturally
			select {
			case <-done:
				wasCancelled = true
			default:
				{
					sr := r.(*SealResult)
					// find root
					var pathmap map[string]*api.FileInfo
					var pathTree string
					for _, tree := range s.Trees {
						if strings.HasPrefix(sr.finfo.Path, tree) {
							pathTree = tree
							pathmap = treePathmap[tree]
							break
						}
					}
					if pathmap == nil {
						panic(fmt.Sprintf("Couldn't determine root of path '%s', roots are %v", sr.finfo.Path, s.Trees))
					}
					// we store only relative paths in the map - this is all we want to serialize
					relaPath := sr.finfo.Path[len(pathTree)+1:]

					_, ok := pathmap[relaPath]
					if ok {
						// duplicate path - shouldn't happen, but we deal with it
						sr.err = fmt.Errorf("Path '%s' was handled multiple times - ignoring occurrence", sr.finfo.Path)
						errCount += 1
						accumResult <- sr
					} else {
						pathmap[relaPath] = sr.finfo
						count += 1
						size += uint64(sr.finfo.Size)
						accumResult <- &SealResult{nil, fmt.Sprintf("DONE ...%s", relaPath), nil, api.Progress}
					}
				} // default
			} // select
		} // range results
		elapsed := time.Now().Sub(st).Seconds()
		sizeMB := float32(size) / (1024.0 * 1024.0)

		// INDEX HANDLING
		//////////////////
		// Serialize all fileinfo structures
		// NOTE: As the parallel writer will send results only when writing finished, we can just operate serially here ...
		// For this there is also no need to optimize performance
		// However, we could use a parallel writer if we are so inclined
		if !wasCancelled {
			var indices []string
			var err error
			if indices, err = s.writeIndex(treePathmap); err != nil {
				// NOTE: We keep successfully written indices as each tree is consistent in itself
				accumResult <- &SealResult{nil, "", err, api.Error}
			}

			// Inform about successfully written indices
			for _, index := range indices {
				accumResult <- &SealResult{nil, fmt.Sprintf("Wrote seal at '%s'", index), err, api.Info}
			}
		}

		accumResult <- &SealResult{nil, fmt.Sprintf("Sealed %d files with total size of %#vMB in %vs (%#v MB/s, %d errors, cancelled=%v)", count, sizeMB, elapsed, float64(sizeMB)/elapsed, errCount, wasCancelled), nil, api.Info}
	}()

	return accumResult
}
