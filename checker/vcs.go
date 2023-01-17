package checker

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/v49/github"
)

type VSCSource interface {
	// ListTags returns a sorted list of tags for the given repo
	ListTags(owner, repo string) ([]string, error)
}

type GitHubSource struct{}

func (g *GitHubSource) ListTags(owner, repo string) ([]string, error) {
	client := github.NewClient(nil)
	repoTags, _, err := client.Repositories.ListTags(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}
	tags := make([]string, len(repoTags))
	for i, tag := range repoTags {
		tags[i] = tag.GetName()
	}

	return sortTags(tags), nil
}

func sortTags(tags []string) []string {
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
	return tags
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
