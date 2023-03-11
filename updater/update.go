package updater

import (
	"fmt"

	"github.com/harryzcy/snuuze/types"
)

func Update(infos []types.UpgradeInfo) {
	groups := groupUpdates(infos)
	fmt.Println(groups)
}
