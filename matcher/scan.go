package matcher

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

	dir, err := filepath.Abs(rootDirectory)
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if shouldIgnore(path) {
				return filepath.SkipDir
			}
			return nil
		}

		if packageManager, ok := matchFile(path); ok {
			matches = append(matches, Match{
				File:           path,
				PackageManager: packageManager,
			})
		}
		return nil
	})

	return matches, err
}

func shouldIgnore(path string) bool {
	return path == ".git"
}
