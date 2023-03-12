package updater

import (
	"fmt"
	"strings"

	"github.com/harryzcy/snuuze/cmdutil"
	"github.com/harryzcy/snuuze/types"
)

func Update(infos []types.UpgradeInfo) {
	groups := groupUpdates(infos)
	for _, group := range groups {
		branchName, ok := generateBranchName(group)
		if ok {
			err := updateDependencies(group.Infos, branchName)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			// create branch for each info
			for _, info := range group.Infos {
				branchName := formateBranchName(info.Dependency.Name)
				infos := []types.UpgradeInfo{info}
				err := updateDependencies(infos, branchName)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
	}
}

func updateDependencies(infos []types.UpgradeInfo, branchName string) error {
	_, err := cmdutil.RunCommand(cmdutil.CommandInputs{
		Command: []string{"git", "checkout", "-b", branchName},
	})
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %s", branchName, err)
	}

	err = delegateUpdate(infos)
	if err != nil {
		return fmt.Errorf("failed to update dependencies: %s", err)
	}

	return nil
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

func delegateUpdate(infos []types.UpgradeInfo) error {
	cache := NewCache()

	errors := []error{}
	for _, info := range infos {
		var err error
		switch info.Dependency.PackageManager {
		case types.PackageManagerGoMod:
			err = upgradeGoMod(cache, info)
		case types.PackageManagerGitHubActions:
			err = upgradeGitHubActions(cache, info)
		default:
			err = fmt.Errorf("unsupported package manager: %s", info.Dependency.PackageManager)
		}
		if err != nil {
			errors = append(errors, err)
		}
	}

	err := cache.Commit()
	if err != nil {
		errors = append(errors, err)
		if len(errors) > 0 {
			return fmt.Errorf("failed to update dependencies: %v", errors)
		}
	}

	return postUpdate(cache)
}

func postUpdate(cache *Cache) error {
	err := postGoMod(cache)
	if err != nil {
		return fmt.Errorf("failed to post update for go.mod: %s", err)
	}
	return nil
}
