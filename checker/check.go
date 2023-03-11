package checker

import (
	"fmt"

	"github.com/harryzcy/snuuze/matcher"
	"github.com/harryzcy/snuuze/parser"
	"github.com/harryzcy/snuuze/types"
)

func ListUpgrades(matches []matcher.Match) ([]types.UpgradeInfo, error) {
	result := []types.UpgradeInfo{}

	for _, match := range matches {
		dependencies, _ := parser.Parse(match)
		for _, dependency := range dependencies {
			info, err := isUpgradable(dependency)
			if err != nil {
				// log error and continue
				fmt.Println(err)
				continue
			}
			if info.Upgradable {
				result = append(result, info)
			}
		}
	}

	return result, nil
}

func isUpgradable(dep types.Dependency) (types.UpgradeInfo, error) {
	switch dep.PackageManager {
	case types.GoMod:
		return isUpgradable_GoMod(dep)
	case types.GitHubActions:
		return isUpgradable_GitHubActions(dep)
	default:
		return types.UpgradeInfo{}, nil
	}
}
