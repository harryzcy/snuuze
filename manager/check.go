package manager

import (
	"fmt"

	"github.com/harryzcy/snuuze/types"
)

// func ListUpgrades(matches []types.Match) ([]*types.UpgradeInfo, error) {
// 	matchesByManager := make(map[types.PackageManager][]types.Match)
// 	for _, match := range matches {
// 		matchesByManager[match.PackageManager] = append(matchesByManager[match.PackageManager], match)
// 	}

// 	for _, m := range managers {
// 		matches, ok := matchesByManager[m.Name()]
// 		if !ok {
// 			continue
// 		}

// 		infos, err := m.ListUpgrades(matches)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return infos, nil
// 	}

// 	result := []*types.UpgradeInfo{}

// 	for _, match := range matches {

// 		for _, dependency := range dependencies {
// 			info, err := isUpgradable(dependency)
// 			if err != nil {
// 				// log error and continue
// 				fmt.Println(err)
// 				continue
// 			}
// 			if info != nil && info.Upgradable {
// 				result = append(result, info)
// 			}
// 		}
// 	}

// 	return result, nil
// }

func PrintUpgradeInfos(infos []*types.UpgradeInfo) {
	for _, info := range infos {
		fmt.Println(info.Dependency.PackageManager, info.Dependency.Name, info.Dependency.Version, info.ToVersion)
	}
}
