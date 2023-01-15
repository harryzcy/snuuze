package manager

import (
	"io/fs"
	"path/filepath"
)

const (
	rootDirectory = "."
)

// Scan scans all the files and returns those containing dependencies
func Scan() ([]Match, error) {
	matches := make([]Match, 0)

	err := filepath.WalkDir(rootDirectory, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if shouldIgnore(path) {
				return filepath.SkipDir
			}
			return nil
		}

		if packageManager, ok := matchFile(path); ok {
			matches = append(matches, Match{
				file:           path,
				packageManager: packageManager,
			})
		}
		return nil
	})

	return matches, err
}

func shouldIgnore(path string) bool {
	return path == ".git"
}
