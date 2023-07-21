package manager

import (
	"strings"

	"github.com/harryzcy/snuuze/types"
)

// matchFile matches a file to a package manager by using manager's Match method
func matchFile(dir, path string) (types.PackageManager, bool) {
	relativePath := strings.TrimPrefix(path, dir)
	relativePath = strings.TrimSuffix(relativePath, "/")

	for _, m := range managers {
		if m.Match(relativePath) {
			return m.Name(), true
		}
	}

	return "", false
}
