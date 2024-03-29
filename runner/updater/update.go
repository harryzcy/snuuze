package updater

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/harryzcy/snuuze/runner/command"
	"github.com/harryzcy/snuuze/runner/git"
	"github.com/harryzcy/snuuze/types"
)

func Update(gitURL, repoDir string, infos []*types.UpgradeInfo, newCommit bool) error {
	if !newCommit {
		err := delegateUpdate(infos)
		if err != nil {
			return fmt.Errorf("failed to update dependencies: %s", err)
		}
		return nil
	}

	groups := groupUpdates(infos)
	fmt.Println("Found", len(groups), "groups of updates")

	hasError := false
	for _, group := range groups {
		commitInfo, hasGroupName := prepareCommit(group)
		if hasGroupName {
			err := updateDependencies(gitURL, repoDir, group.Infos, commitInfo)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				hasError = true
				continue
			}
		} else {
			// create branch for each info
			for _, info := range group.Infos {
				commitInfo := prepareCommitByUpgradeInfo(info)
				err := updateDependencies(gitURL, repoDir, []*types.UpgradeInfo{info}, commitInfo)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					hasError = true
					continue
				}
			}
		}
	}
	if hasError {
		return errors.New("error occurred during update")
	}
	return nil
}

func updateDependencies(gitURL, repoDir string, infos []*types.UpgradeInfo, info *commitInfo) error {
	base := git.GetDefaultBranch(repoDir)
	fmt.Println("Creating branch", info.branchName, "from", base)

	output, err := command.RunCommand(command.Inputs{
		Command: []string{"git", "-C", repoDir, "checkout", "-b", info.branchName, base},
	})
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %s\n%s", info.branchName, err, output.Stderr.String())
	}

	err = delegateUpdate(infos)
	if err != nil {
		return fmt.Errorf("failed to update dependencies: %s", err)
	}

	err = git.CommitChanges(repoDir, info.message)
	if err != nil {
		return fmt.Errorf("failed to commit changes to branch [%s]: %s", info.branchName, err)
	}

	err = git.PushBranch(repoDir, info.branchName)
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
			branchName: formatBranchName(group.Rule.Name),
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
		branchName: formatBranchName(info.Dependency.Name),
		message:    "Update " + info.Dependency.Name + " to " + info.ToVersion,
	}
}

func formatBranchName(name string) string {
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ToLower(name)
	return fmt.Sprintf("snuuze/%s", name)
}

func delegateUpdate(infos []*types.UpgradeInfo) error {
	cache := NewCache()

	goReplaceItems := []*ReplaceItem{}

	var errs []error
	for _, info := range infos {
		var err error
		switch info.Dependency.PackageManager {
		case types.PackageManagerDocker:
			err = upgradeDocker(cache, info)
		case types.PackageManagerGitHubActions:
			err = upgradeGitHubActions(cache, info)
		case types.PackageManagerGoMod:
			var replace *ReplaceItem
			replace, err = upgradeGoMod(cache, info)
			if err == nil && replace != nil {
				goReplaceItems = append(goReplaceItems, replace)
			}
		case types.PackageManagerPip:
			err = upgradePip(cache, info)
		default:
			err = fmt.Errorf("unsupported package manager: %s", info.Dependency.PackageManager)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}

	err := postUpdate(cache, goReplaceItems)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to update dependencies: %v", errs)
	}
	return nil
}

func postUpdate(cache *Cache, goReplaceItems []*ReplaceItem) error {
	goModFiles, err := cache.ListGoMod()
	if err != nil {
		return fmt.Errorf("postGoMod: failed to list go.mod files: %s", err)
	}

	err = cache.Commit()
	if err != nil {
		return err
	}

	err = postGoMod(cache, goModFiles, goReplaceItems)
	if err != nil {
		return fmt.Errorf("failed to post update for go.mod: %s", err)
	}
	return nil
}
