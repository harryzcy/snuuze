package updater

import (
	"fmt"
	"strings"

	"github.com/harryzcy/snuuze/cmdutil"
	"github.com/harryzcy/snuuze/types"
)

func Update(gitURL, repoDir string, infos []*types.UpgradeInfo) {
	groups := groupUpdates(infos)
	fmt.Println("Found", len(groups), "groups of updates")
	for _, group := range groups {
		commitInfo, hasGroupName := prepareCommit(group)
		if hasGroupName {
			err := updateDependencies(gitURL, repoDir, group.Infos, commitInfo)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			// create branch for each info
			for _, info := range group.Infos {
				commitInfo := prepareCommitByUpgradeInfo(info)
				err := updateDependencies(gitURL, repoDir, []*types.UpgradeInfo{info}, commitInfo)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
	}
}

func updateDependencies(gitURL, repoDir string, infos []*types.UpgradeInfo, info *commitInfo) error {
	base := getDefaultBranch(repoDir)
	fmt.Println("Creating branch", info.branchName, "from", base)

	output, err := cmdutil.RunCommand(cmdutil.CommandInputs{
		Command: []string{"git", "-C", repoDir, "checkout", "-b", info.branchName, base},
	})
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %s\n%s", info.branchName, err, output.Stderr.String())
	}

	err = delegateUpdate(infos)
	if err != nil {
		return fmt.Errorf("failed to update dependencies: %s", err)
	}

	err = commitChanges(repoDir, info.branchName, info.message)
	if err != nil {
		return fmt.Errorf("failed to commit changes: %s", err)
	}

	err = pushBranch(repoDir, info.branchName)
	if err != nil {
		return fmt.Errorf("failed to push branch: %s", err)
	}

	err = createPullRequest(gitURL, info, base)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %s", err)
	}

	return nil
}

type commitInfo struct {
	branchName string
	message    string
}

// prepareCommit returns commit info and a bool value indicating whether the commit is a group commit
func prepareCommit(group RuleGroup) (*commitInfo, bool) {
	if group.Rule.Name != "" {
		return &commitInfo{
			branchName: formateBranchName(group.Rule.Name),
			message:    "Update " + group.Rule.Name,
		}, true
	}

	if len(group.Infos) == 1 {
		return prepareCommitByUpgradeInfo(group.Infos[0]), true
	}

	return nil, false
}

func prepareCommitByUpgradeInfo(info *types.UpgradeInfo) *commitInfo {
	return &commitInfo{
		branchName: formateBranchName(info.Dependency.Name),
		message:    "Update " + info.Dependency.Name + " to " + info.ToVersion,
	}
}

func formateBranchName(name string) string {
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ToLower(name)
	return fmt.Sprintf("snuuze/%s", name)
}

func delegateUpdate(infos []*types.UpgradeInfo) error {
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
