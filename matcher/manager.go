package matcher

import (
	"path/filepath"
	"strings"

	"github.com/harryzcy/snuuze/types"
)

type Match struct {
	File           string
	PackageManager types.PackageManager
}

func matchFile(path string) (types.PackageManager, bool) {
	filename := filepath.Base(path)
	if filename == "go.mod" {
		return types.PackageManagerGoMod, true
	}

	if strings.HasPrefix(path, ".github/workflows") && (strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")) {
		return types.PackageManagerGitHubActions, true
	}

	return "", false
}
