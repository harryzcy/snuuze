package checker

import (
	"github.com/harryzcy/sailor/matcher"
	"github.com/harryzcy/sailor/parser"
	"github.com/harryzcy/sailor/types"
)

type UpgradeInfo struct {
	Dependency types.Dependency
	Upgradable bool
	ToVersion  string
}

func ListUpgrades(matches []matcher.Match) ([]UpgradeInfo, error) {
	result := []UpgradeInfo{}

	for _, match := range matches {
		dependencies, _ := parser.Parse(match)
		for _, dependency := range dependencies {
			info, err := isUpgradable(dependency)
			if err != nil {
				return nil, err
			}
			if info.Upgradable {
				result = append(result, info)
			}
		}
	}

	return result, nil
}

func isUpgradable(dep types.Dependency) (UpgradeInfo, error) {
	switch dep.PackageManager {
	case matcher.GoMod:
		return isUpgradable_GoMod(dep)
	case matcher.GitHubActions:
		return isUpgradable_GitHubActions(dep)
	default:
		return UpgradeInfo{}, nil
	}
}
