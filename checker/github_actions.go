package checker

import (
	"fmt"
	"strings"

	"github.com/harryzcy/snuuze/types"
)

func isUpgradable_GitHubActions(dep types.Dependency) (UpgradeInfo, error) {
	owner, repo, err := parseRepo(dep.Name)
	if err != nil {
		return UpgradeInfo{}, err
	}

	info := UpgradeInfo{
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

	var source VSCSource = &GitHubSource{}
	tags, err := source.ListTags(&ListTagsParameters{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return UpgradeInfo{}, err
	}
	latest := getLatestTag(tags, dep.Version)
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
