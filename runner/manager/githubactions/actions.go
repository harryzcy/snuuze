package githubactions

import (
	"fmt"
	"strings"

	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/runner/manager/common"
	"github.com/harryzcy/snuuze/types"
)

type GitHubActionsManager struct{}

func New() common.Manager {
	return &GitHubActionsManager{}
}

func (m *GitHubActionsManager) Name() types.PackageManager {
	return types.PackageManagerGitHubActions
}

func (m *GitHubActionsManager) Match(path string) bool {
	if !strings.HasPrefix(path, ".github/workflows") && !strings.HasPrefix(path, ".gitea/workflows") {
		return false
	}
	return strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")
}

func (m *GitHubActionsManager) Parse(match types.Match, data []byte) ([]*types.Dependency, error) {
	return parseGitHubActions(match.File, data)
}

func (m *GitHubActionsManager) FindDependencies(matches []types.Match) ([]*types.Dependency, error) {
	return common.FindDependencies(m, matches)
}

func (m *GitHubActionsManager) ListUpgrades(matches []types.Match) ([]*types.UpgradeInfo, error) {
	return common.ListUpgrades(m, matches)
}

func (m *GitHubActionsManager) IsUpgradable(dep types.Dependency) (*types.UpgradeInfo, error) {
	owner, repo, err := parseRepo(dep.Name)
	if err != nil {
		return nil, err
	}

	info := &types.UpgradeInfo{
		Dependency: dep,
	}
	if isSha(dep.Version) {
		// don't check if sha is upgradable
		return info, nil
	}
	if !strings.HasPrefix(dep.Version, "v") {
		// not a versioned tag
		return info, nil
	}

	client, err := platform.NewGitHubClient()
	if err != nil {
		return nil, err
	}

	tags, err := client.ListTags(&platform.ListTagsInput{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return nil, err
	}

	latest, err := common.GetLatestTag(&common.GetLatestTagInput{
		DepName:    dep.Name,
		Tags:       tags,
		CurrentTag: dep.Version,
		AllowMajor: true,
	})
	if err != nil {
		return nil, err
	}

	if latest != dep.Version {
		info.Upgradable = true
		info.ToVersion = latest
	}
	return info, nil
}

func parseRepo(uses string) (string, string, error) {
	uses = strings.Split(uses, "@")[0]
	parts := strings.Split(uses, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid uses in github workflow file: %s", uses)
	}
	owner := parts[0]
	repo := parts[1]
	return owner, repo, nil
}

func isSha(version string) bool {
	return len(version) == 40
}
