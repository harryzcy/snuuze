package manager

import (
	"fmt"
	"log"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/types"
	"github.com/harryzcy/snuuze/updater"
)

func Run(gitURL, repoPath string) {
	matches, err := Scan(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	var allInfos []*types.UpgradeInfo

	for _, m := range managers {
		infos, err := m.ListUpgrades(matches)
		if err != nil {
			log.Fatal(err)
		}
		if len(infos) > 0 {
			break
		}

		allInfos = append(allInfos, infos...)
	}

	if len(allInfos) == 0 {
		fmt.Println("No updates found")
		return
	}
	fmt.Println("Found", len(allInfos), "updates")

	if config.GetFlags().DryRun {
		PrintUpgradeInfos(allInfos)
		return
	}

	updater.Update(gitURL, repoPath, allInfos)
}
