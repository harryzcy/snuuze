package parser

import (
	"os"

	"github.com/harryzcy/snuuze/matcher"
	"github.com/harryzcy/snuuze/types"
)

func Parse(match matcher.Match) ([]types.Dependency, error) {
	data, err := os.ReadFile(match.File)
	if err != nil {
		return nil, err
	}

	switch match.PackageManager {
	case types.PackageManagerGoMod:
		return parseGoMod(match.File, data)
	case types.PackageManagerGitHubActions:
		return parseGitHubActions(data)
	}

	return nil, nil
}
