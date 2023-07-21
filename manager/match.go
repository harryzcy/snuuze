package manager

import (
	"strings"

	"github.com/harryzcy/snuuze/types"
)

func matchFile(dir, path string) (types.PackageManager, bool) {
	relativePath := strings.TrimPrefix(path, dir)
	relativePath = strings.TrimPrefix(relativePath, "/")

	for _, m := range managers {
		if m.Match(relativePath) {
			return m.Name(), true
		}
	}

	return "", false
}
