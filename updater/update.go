package updater

import (
	"fmt"
	"strings"

	"github.com/harryzcy/snuuze/types"
)

func Update(infos []types.UpgradeInfo) {
	groups := groupUpdates(infos)
	for _, group := range groups {
		branchName, ok := generateBranchName(group)
		if ok {
			fmt.Println(branchName)
			// TODO create branch
		} else {
			// create branch for each info
			for _, info := range group.Infos {
				branchName := formateBranchName(info.Dependency.Name)
				fmt.Println(branchName)
				// TODO create branch
			}
		}
	}
}

func generateBranchName(group RuleGroup) (string, bool) {
	if group.Rule.Name != "" {
		return formateBranchName(group.Rule.Name), true
	}

	if len(group.Infos) == 1 {
		return formateBranchName(group.Infos[0].Dependency.Name), true
	}

	return "", false
}

func formateBranchName(name string) string {
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ToLower(name)
	return fmt.Sprintf("snuuze/%s", name)
}
