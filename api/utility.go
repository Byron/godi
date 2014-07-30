package api

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Append elm if it is not yet on dest
func AppendUniqueString(dest []string, elm string) []string {
	for _, d := range dest {
		if d == elm {
			return dest
		}
	}
	return append(dest, elm)
}

// Parse all valid source items from the given list.
// May either be files or directories. The returned list may be shorter, as contained paths are
// skipped automatically. Paths will be normalized.
func ParseSources(items []string, allowFiles bool) (res []string, err error) {
	var invalidTrees, noTrees, noRegularFiles []string
	res = make([]string, len(items))
	copy(res, items)

	for i, tree := range res {
		if stat, err := os.Stat(tree); err != nil {
			invalidTrees = append(invalidTrees, tree)
			continue
		} else if !stat.IsDir() {
			if !allowFiles {
				noTrees = append(noTrees, tree)
				continue
			} else if !stat.Mode().IsRegular() {
				noRegularFiles = append(noRegularFiles, tree)
			}
			// otherwise it's a regular file
		}
		tree = path.Clean(tree)
		if !filepath.IsAbs(tree) {
			tree, err = filepath.Abs(tree)
			if err != nil {
				return nil, err
			}
		}
		res[i] = tree
	}

	if len(invalidTrees) > 0 {
		return nil, errors.New("None of the following items exists: " + strings.Join(invalidTrees, ", "))
	}
	if len(noTrees) > 0 {
		return nil, errors.New("The following items are no directory: " + strings.Join(noTrees, ", "))
	}
	if len(noRegularFiles) > 0 {
		return nil, errors.New("The following items are no regular files: " + strings.Join(noRegularFiles, ", "))
	}

	// drop trees which are a sub-tree of another, or which are equal !
	if len(res) > 1 {
		validTrees := make([]string, 0, len(res))
		for i, ltree := range res {
			for x, rtree := range res {
				if i == x || strings.HasPrefix(ltree, rtree+string(os.PathSeparator)) {
					continue
				}
				validTrees = AppendUniqueString(validTrees, ltree)
			}
		}
		if len(validTrees) == 0 {
			panic("Didn't find a single valid tree after pruning - shouldn't happen")
		}

		res = validTrees
	}

	return res, nil
}
