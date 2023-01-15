package matcher

import (
	"path/filepath"
	"strings"
)

type Match struct {
	file           string
	packageManager PackageManager
}

type PackageManager string

const (
	GoMod         PackageManager = "gomod"
	GitHubActions PackageManager = "github-actions"
)

func matchFile(path string) (PackageManager, bool) {
	filename := filepath.Base(path)
	if filename == "go.mod" {
		return GoMod, true
	}

	if strings.HasPrefix(path, ".github/workflows") && (strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")) {
		return GitHubActions, true
	}

	return "", false
}
