package manager

import (
	"fmt"

	"github.com/harryzcy/snuuze/types"
)

func Run(gitURL, repoPath string) ([]*types.UpgradeInfo, error) {
	allMatches, err := Scan(repoPath)
	if err != nil {
		return nil, err
	}

	var allInfos []*types.UpgradeInfo

	for _, m := range managers {
		matches := allMatches[m.Name()]
		infos, err := m.ListUpgrades(matches)
		if err != nil {
			return nil, err
		}
		allInfos = append(allInfos, infos...)
	}

	if len(allInfos) == 0 {
		fmt.Println("No updates found")
		return nil, nil
	}
	fmt.Println("Found", len(allInfos), "updates")

	return allInfos, nil
}

func FindAll(repoPath string) (map[types.PackageManager][]*types.Dependency, error) {
	allMatches, err := Scan(repoPath)
	if err != nil {
		return nil, err
	}

	result := make(map[types.PackageManager][]*types.Dependency)

	for name, m := range managers {
		matches := allMatches[name]
		dependencies, err := m.FindDependencies(matches)
		if err != nil {
			return nil, err
		}
		result[name] = dependencies
	}

	return result, nil
}
