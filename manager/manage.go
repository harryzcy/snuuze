package manager

import (
	"fmt"
	"log"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/manager/checker"
	"github.com/harryzcy/snuuze/manager/matcher"
	"github.com/harryzcy/snuuze/updater"
)

func Run(gitURL, repoPath string) {
	matches, err := matcher.Scan(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	infos, err := checker.ListUpgrades(matches)
	if err != nil {
		log.Fatal(err)
	}

	if len(infos) == 0 {
		fmt.Println("No updates found")
		return
	}
	fmt.Println("Found", len(infos), "updates")

	if config.GetFlags().DryRun {
		checker.PrintUpgradeInfos(infos)
		return
	}

	updater.Update(gitURL, repoPath, infos)
}
