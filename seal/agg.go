package seal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync/atomic"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/codec"
	"github.com/Byron/godi/io"
)

// Will setup a go-routine which writes a seal for the given tree, continuously as new files come in
func setupIndexWriter(commonTree string, encoder codec.Codec) (chan<- api.FileInfo, <-chan indexWriterResult) {
	if encoder == nil {
		panic("No encoder provided")
	}

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

		// This will and should fail if the file already exists
		indexPath := api.IndexPath(commonTree, encoder.Extension())
		fd, err := os.OpenFile(indexPath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			// We assume the serializer deals with buffering if he needs it.
			// MHL caches in memory, and gob uses zip, which allocates a big buffer itself
			err = encoder.Serialize(sealFiles, fd)
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

		// If we have no file information, this one is likely to be sent by the generator.
		// For now we just pass it on.
		// NOTE(st): right now we are trying to seal as much as possible, but count errors on the way
		// We could use the gen info to abort early, by marking entire trees as failed.
		if sr.FromGenerator() {
			accumResult <- r
			return sr.Err == nil
		}

		deleteResultSafely := func(tree, path string) {
			if !isWriting {
				panic("Shouldn't ever try to delete a file we have not written ... ")
			}

			br := api.BasicResult{Prio: api.Info}

			// If the file doesn't exist anymore, we don't care either
			err := os.Remove(path)
			if err == nil {
				s.Stats.NumUndoneFiles += 1
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
			treeInfo.sealFInfos, treeInfo.sealResult = setupIndexWriter(treeRoot, codec.NewByName(s.format))
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

		if !hasError && treeInfo.lsr.err == nil {
			// Provide some informational logging
			sr.Prio = api.Info
			if len(sr.source) == 0 {
				sr.Msg = fmt.Sprintf("# %s", sr.Finfo.Path)
			} else {
				sr.Msg = fmt.Sprintf("CP %s -> %s", sr.source, sr.Finfo.Path)
			}

			// The seal can fail anytime, for instance on permission issues or when there
			// This select will block until something happens, usually this means
			// TODO(st): As we don't expect to inrease performance with this async operation,
			// it would be good to use a begin/do[,...]/end style serialization pattern
			select {
			case lsr, ok := <-treeInfo.sealResult:
				if ok {
					// Be sore we stop sending !
					// Results are only sent once - if it sends now, it must be an error !
					treeInfo.lsr = lsr
					close(treeInfo.sealFInfos)

					// Mark the tree early - I would always expect the error to be set here ...
					if treeInfo.lsr.err != nil {
						// error is counted where it is handled
						treeInfo.hasError = true
						// Check if all trees have failures, and provide feedback to the generators !
						// If we can't do anything useful anymore, we should stop
						nft := 0 // num failed trees
						for _, ti := range treeInfoMap {
							if ti.hasError {
								nft += 1
							}
						}

						maxTrees := 0
						if isWriting {
							maxTrees = io.WriteChannelDeviceMapTrees(s.rootedWriters)
						} else {
							for _, rctrl := range s.RootedReaders {
								maxTrees += len(rctrl.Trees)
							}
						}

						if nft == maxTrees {
							atomic.AddUint32(&s.Stats.StopTheEngines, 1)
						}
					} else {
						panic("Didn't expect seal writer to stop early, but not provie an error")
					}
				}
			// Send the valid file to the sealer
			case treeInfo.sealFInfos <- sr.Finfo:
			}
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
		accumResult chan<- api.Result) {

		// All we have to do is to stop the sealers and gather their result, possibly deleting
		// incomplete seals (created because there was some error on the way)
		for tree, treeInfo := range treeInfoMap {
			// Free now unused memory, just in case we do a verify or something else afterwards.
			// Of course, the gc will decide when to actually free it.
			treeInfo.writtenFiles = nil

			if treeInfo.lsr.err == nil {
				close(treeInfo.sealFInfos)
				treeInfo.lsr = <-treeInfo.sealResult
			}

			br := api.BasicResult{}
			if treeInfo.lsr.err == nil {
				// Can we have an error here ? Just be sure we don't, otherwise we say to have
				// written the file, and remove it in the next step
				if !treeInfo.hasError {
					br.Msg = fmt.Sprintf("Wrote seal file to '%s'", treeInfo.lsr.path)
					// special marker, to allow others to easily retrieve seal files from the result
					// No normal file has a size of -1
					br.Finfo = api.FileInfo{Path: treeInfo.lsr.path, Size: -1}
					br.Prio = api.Valuable
				}
			} else {
				// Just forward the error, hoping it is informative enough
				// The error will be handled
				s.Stats.ErrCount += 1
				treeInfo.hasError = true
			}

			if treeInfo.hasError {
				if len(treeInfo.lsr.path) > 0 {
					os.Remove(treeInfo.lsr.path)
				}
				if treeInfo.lsr.err != nil {
					br.Msg = fmt.Sprintln(treeInfo.lsr.err.Error())
				}
				br.Msg += fmt.Sprintf("Did not write seal for '%s' due to preceeding errors", tree)
				br.Prio = api.Error
			}

			accumResult <- &br
		} // end for each tree/treeInfo tuple

		prefix := fmt.Sprintf("SEAL %s", SymbolSuccess)
		if s.Stats.ErrCount > 0 {
			prefix = fmt.Sprintf("SEAL %s", SymbolFail)
		}

		// Final seal result !
		accumResult <- &api.BasicResult{
			Msg: fmt.Sprintf(
				"%s: %s",
				prefix,
				s.Stats.DeltaString(&s.Stats, s.Stats.Elapsed(), io.StatsClientSep),
			) + s.Stats.String(),
			Prio: api.Valuable,
		}

	} // end finalizer()

	return api.Aggregate(results, s.Done, resultHandler, finalizer, &s.Stats)
}
