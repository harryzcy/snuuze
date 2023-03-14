package checker

import (
	"strings"
)

type VSCSource interface {
}

type GitHubSource struct{}

func getLatestTag(sortedTags []string, currentTag string) string {
	if len(sortedTags) == 0 {
		return ""
	}
	latest := sortedTags[0]
	parts := strings.Split(latest, ".")
	dotCount := strings.Count(currentTag, ".")
	return strings.Join(parts[:dotCount+1], ".")
}
