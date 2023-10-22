package manager

import (
	"fmt"

	"github.com/harryzcy/snuuze/config"
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

	if config.GetCLIConfig().DryRun {
		PrintUpgradeInfos(allInfos)
		return nil, nil
	}

	return allInfos, nil
}

func FindAll(repoPath string) ([]*types.Dependency, error) {
	allMatches, err := Scan(repoPath)
	if err != nil {
		return nil, err
	}

	var result []*types.Dependency

	for _, m := range managers {
		matches := allMatches[m.Name()]
		dependencies, err := m.FindDependencies(matches)
		if err != nil {
			return nil, err
		}
		result = append(result, dependencies...)
	}

	return result, nil
}
