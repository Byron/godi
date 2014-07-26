package seal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/utility"
)

// Must be kept in sync with indexPath generator
var reIsIndexPath = regexp.MustCompile(fmt.Sprintf(`%s_\d{4}-\d{2}-\d{2}_\d{2}\d{2}\d{2}\.%s`, IndexBaseName, codec.GobExtension))

// return a path to an index file residing at tree
func indexPath(tree string, extension string) string {
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

// Will setup a go-routine which writes a seal for the given tree, continuously as new files come in
func setupIndexWriter(commonTree string) (chan<- api.FileInfo, <-chan indexWriterResult) {

	sealFiles := make(chan api.FileInfo)
	// we may always put in one item, which allows this go-routine to go down without having to wait
	results := make(chan indexWriterResult, 1)

	go func() {
		defer close(results)

		// Serialize all fileinfo structures
		// For now, we just use a standard writer, without using our parallel writers, which would at least
		// assure we don't write with more streams than defined.
		// Reason is that // we want this to be as fast as possible without blocking, which is also why it is cached
		// reasonably well, allowing it to only write on larger chunks.
		encoder := codec.Gob{}

		// This will and should fail if the file already exists
		indexPath := indexPath(commonTree, encoder.Extension())
		fd, err := os.OpenFile(indexPath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			bfd := bufio.NewWriterSize(fd, 512*1024)
			err = encoder.Serialize(sealFiles, fd)
			bfd.Flush()
			fd.Close()
			if err != nil {
				// Remove intermediate results
				os.Remove(indexPath)
			}
		}

		results <- indexWriterResult{path: indexPath, err: err}
	}()

	return sealFiles, results
}

func (s *SealCommand) Aggregate(results <-chan api.Result) <-chan api.Result {
	treeInfoMap := make(map[string]*aggregationTreeInfo)
	isWriting := len(s.rootedWriters) > 0

	resultHandler := func(r api.Result, accumResult chan<- api.Result) bool {
		sr := r.(*SealResult)

		deleteResultSafely := func(tree, path string) {
			if !isWriting {
				panic("Shouldn't ever try to delete a file we have not written ... ")
			}

			br := api.BasicResult{Prio: api.Info}

			// If the file doesn't exist anymore, we don't care either
			err := os.Remove(path)
			if err == nil {
				br.Msg = fmt.Sprintf("Removed '%s'", path)
				accumResult <- &br

				// try to remove the directory - will fail if non-empty.
				// only do that if we wouldn't remove the tree.
				// Also crawl upwards
				var derr error
				for dir := filepath.Dir(path); dir != tree && derr == nil; dir = filepath.Dir(dir) {
					derr = os.Remove(dir)
				}
			}
		}

		// Keep the file-information, in any case
		treeRoot := sr.Finfo.Root()
		treeInfo, didExist := treeInfoMap[treeRoot]
		if !didExist {
			// Initialize this root
			// Create a new go-routine which will take care of streaming file-information straight to file
			treeInfo = &aggregationTreeInfo{}
			treeInfo.sealFInfos, treeInfo.sealResult = setupIndexWriter(treeRoot)
			treeInfoMap[treeRoot] = treeInfo
		}

		// We will keep track of the file even if it reported an error.
		// That way, we can later determine what to cleanup
		hasError := r.Error() != nil || treeInfo.hasError

		// In any case, remember the file we have written in some way (may be partial write)
		// However, don't remember the file if we didn't actually write it in any way
		// and failed to write because it existed
		if isWriting && !os.IsExist(sr.Err) {
			treeInfo.writtenFiles = append(treeInfo.writtenFiles, sr.Finfo.Path)
		}

		if !hasError {
			// Provide some informational logging
			sr.Prio = api.Info
			if len(sr.source) == 0 {
				sr.Msg = fmt.Sprintf("# %s", sr.Finfo.Path)
			} else {
				sr.Msg = fmt.Sprintf("CP %s -> %s", sr.source, sr.Finfo.Path)
			}

			// Send the valid file to the sealer
			treeInfo.sealFInfos <- sr.Finfo
		}

		// Send previous error first, before error handling
		accumResult <- sr

		if hasError {
			// mark the entire tree as having errors
			treeInfo.hasError = true
			hasError = true

			if isWriting {
				// Remove all files we have created so far
				sort.Sort(byLongestPathDescending(treeInfo.writtenFiles))
				for _, path := range treeInfo.writtenFiles {
					deleteResultSafely(treeRoot, path)
				}
				// Clear it - next time we have to remove whatever happened so far
				treeInfo.writtenFiles = nil
			}
		}

		return !hasError
	} // end resultHandler()

	finalizer := func(
		accumResult chan<- api.Result,
		st *api.AggregateFinalizerState) {

		// All we have to do is to stop the sealers and gather their result, possibly deleting
		// incomplete seals (created because there was some error on the way)
		for tree, treeInfo := range treeInfoMap {
			close(treeInfo.sealFInfos)
			sres := <-treeInfo.sealResult

			// Free now unused memory, just in case we do a verify or something else afterwards.
			// Of course, the gc will decide when to actually free it.
			treeInfo.writtenFiles = nil

			br := api.BasicResult{}

			if treeInfo.hasError {
				os.Remove(sres.path)
				br.Msg = fmt.Sprintf("Did not write seal for '%s' due to preceeding errors", tree)
				br.Prio = api.Error
			} else {
				if sres.err == nil {
					br.Msg = fmt.Sprintf("Wrote seal file to '%s'", sres.path)
					// special marker, to allow others to easily retrieve seal files from the result
					// No normal file has a size of -1
					br.Finfo = api.FileInfo{Path: sres.path, Size: -1}
					br.Prio = api.Valuable
				} else {
					// Just forward the error, hoping it is informative enough
					st.ErrCount += 1
					br.Err = sres.err
				}
			}

			accumResult <- &br
		} // end for each tree/treeInfo tuple

		prefix := "SEAL DONE"
		if st.ErrCount > 0 {
			prefix = "SEAL FAILED"
		}

		// Final seal result !
		accumResult <- &api.BasicResult{
			Msg: fmt.Sprintf(
				"%s: %s",
				prefix,
				s.Stats.DeltaString(&s.Stats, st.Elapsed, utility.StatsClientSep),
			) + st.String(),
			Prio: api.Valuable,
		}

	} // end finalizer()

	return api.Aggregate(results, s.Done, resultHandler, finalizer)
}
