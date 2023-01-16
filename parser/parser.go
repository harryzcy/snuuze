package parser

import (
	"os"

	"github.com/harryzcy/sailor/matcher"
	"github.com/harryzcy/sailor/types"
)

func Parse(match matcher.Match) ([]types.Dependency, error) {
	data, err := os.ReadFile(match.File)
	if err != nil {
		return nil, err
	}

	switch match.PackageManager {
	case matcher.GoMod:
		return parseGoMod(match.File, data)
	case matcher.GitHubActions:
		return parseGitHubActions(data)
	}

	return nil, nil
}
