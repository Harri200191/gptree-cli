package internal

import (
	"os"
	"path/filepath"
)

// WalkDir walks the tree rooted at root, skipping any dirs whose
// base name matches ignoreDirs, and skipping any files matching ignoreFileGlobs.
func WalkDir(root string, ignoreDirs []string, ignoreFileGlobs []string) ([]string, error) {
	var files []string

	// Build a set of directory names to ignore
	ignoreDirSet := make(map[string]struct{}, len(ignoreDirs))
	for _, d := range ignoreDirs {
		ignoreDirSet[d] = struct{}{}
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If this is a directory and its name is in ignoreDirSet, skip it
		if d.IsDir() {
			if _, skip := ignoreDirSet[d.Name()]; skip {
				return filepath.SkipDir
			}
			return nil
		}

		// Otherwise, it’s a file—skip if its base matches any ignoreFileGlobs
		base := filepath.Base(path)
		for _, pattern := range ignoreFileGlobs {
			ok, _ := filepath.Match(pattern, base)
			if ok {
				return nil
			}
		}

		files = append(files, path)
		return nil
	})
	return files, err
}
