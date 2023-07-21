package checker

import (
	"fmt"

	"github.com/harryzcy/snuuze/matcher"
	"github.com/harryzcy/snuuze/parser"
	"github.com/harryzcy/snuuze/types"
)

func ListUpgrades(matches []matcher.Match) ([]*types.UpgradeInfo, error) {
	result := []*types.UpgradeInfo{}

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

func isUpgradable(dep types.Dependency) (*types.UpgradeInfo, error) {
	switch dep.PackageManager {
	case types.PackageManagerGoMod:
		return isUpgradable_GoMod(dep)
	case types.PackageManagerGitHubActions:
		return isUpgradable_GitHubActions(dep)
	default:
		return nil, nil
	}
}

func PrintUpgradeInfos(infos []*types.UpgradeInfo) {
	for _, info := range infos {
		fmt.Println(info.Dependency.PackageManager, info.Dependency.Name, info.Dependency.Version, info.ToVersion)
	}
}
