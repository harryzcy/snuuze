package checker

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/v49/github"

	"github.com/harryzcy/sailor/types"
)

func isUpgradable_GitHubActions(dep types.Dependency) (UpgradeInfo, error) {
	owner, repo, err := getRepo(dep.Name)
	if err != nil {
		return UpgradeInfo{}, err
	}

	info := UpgradeInfo{
		Dependency: dep,
	}

	client := github.NewClient(nil)
	if isSha(dep.Version) {
		// don't check if sha is upgradable
		return info, nil
	}

	if !strings.HasPrefix(dep.Version, "v") {
		// not a versioned tag
		return info, nil
	}

	repoTags, _, err := client.Repositories.ListTags(context.Background(), owner, repo, nil)
	if err != nil {
		return UpgradeInfo{}, err
	}
	tags, err := getSortedTags(repoTags)
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

func getRepo(uses string) (string, string, error) {
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

func getSortedTags(repoTags []*github.RepositoryTag) ([]string, error) {
	tags := make([]string, len(repoTags))
	for i, tag := range repoTags {
		tags[i] = tag.GetName()
	}
	sort.SliceStable(tags, func(i, j int) bool {
		tag1, err1 := parseTag(tags[i])
		tag2, err2 := parseTag(tags[j])
		if err1 != nil || err2 != nil {
			// if one of the tags is not a valid tag, then compare the string
			return tags[i] > tags[j]
		}
		minLen := len(tag1)
		if len(tag2) < minLen {
			minLen = len(tag2)
		}

		for i := 0; i < minLen; i++ {
			if tag1[i] > tag2[i] {
				return true
			}
			if tag1[i] < tag2[i] {
				return false
			}
		}
		return len(tag1) > len(tag2)
	})

	return tags, nil
}

func getLatestTag(sortedTags []string, currentTag string) string {
	latest := sortedTags[0]
	parts := strings.Split(latest, ".")
	dotCount := strings.Count(currentTag, ".")
	return strings.Join(parts[:dotCount+1], ".")
}

func parseTag(tag string) ([]int, error) {
	tag = strings.TrimPrefix(tag, "v")
	parts := strings.Split(tag, ".")
	intParts := make([]int, len(parts))
	var err error
	for i, part := range parts {
		intParts[i], err = strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
	}
	return intParts, nil
}
