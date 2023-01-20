package parser

import (
	"os"

	"github.com/harryzcy/latte/matcher"
	"github.com/harryzcy/latte/types"
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
