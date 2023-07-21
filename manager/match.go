package manager

import (
	"path/filepath"
	"strings"

	"github.com/harryzcy/snuuze/types"
)

func matchFile(dir, path string) (types.PackageManager, bool) {
	filename := filepath.Base(path)

	relativePath := strings.TrimPrefix(path, dir)
	relativePath = strings.TrimPrefix(relativePath, "/")

	if filename == "go.mod" {
		return types.PackageManagerGoMod, true
	}

	if strings.HasPrefix(relativePath, ".github/workflows") && (strings.HasSuffix(relativePath, ".yml") || strings.HasSuffix(relativePath, ".yaml")) {
		return types.PackageManagerGitHubActions, true
	}

	return "", false
}
