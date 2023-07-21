package manager

import (
	"fmt"

	"github.com/harryzcy/snuuze/types"
)

func PrintUpgradeInfos(infos []*types.UpgradeInfo) {
	for _, info := range infos {
		fmt.Println(info.Dependency.PackageManager, info.Dependency.Name, info.Dependency.Version, info.ToVersion)
	}
}
