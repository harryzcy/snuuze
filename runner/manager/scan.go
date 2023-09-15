package manager

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/harryzcy/snuuze/types"
)

// Scan scans all the files and returns those containing dependencies
func Scan(dir string) (map[types.PackageManager][]types.Match, error) {
	fmt.Println("Scanning files in", dir)
	matches := make(map[types.PackageManager][]types.Match)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if shouldIgnore(path) {
				return filepath.SkipDir
			}
			return nil
		}

		if packageManager, ok := matchFile(dir, path); ok {
			matches[packageManager] = append(matches[packageManager], types.Match{
				File:           path,
				PackageManager: packageManager,
			})
		}
		return nil
	})

	return matches, err
}

func shouldIgnore(path string) bool {
	return strings.HasSuffix(path, ".git")
}
