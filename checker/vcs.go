package checker

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

type VSCSource interface {
}

type GitHubSource struct{}

// getLatestTag returns the latest tag that is not a pre-release, or the current tag if no such tag exists
func getLatestTag(tags []string, currentTag string) (string, error) {
	currentVersion, err := version.NewVersion(currentTag)
	if err != nil {
		return "", err
	}
	versions := make([]*version.Version, len(tags))
	for i, tag := range tags {
		if versions[i], err = version.NewVersion(tag); err != nil {
			return "", err
		}
	}

	sort.Sort(sort.Reverse(version.Collection(versions)))

	var latest *version.Version
	for _, v := range versions {
		if v.Prerelease() != "" {
			continue
		}

		if v.GreaterThan(currentVersion) {
			latest = v
			break
		}
	}

	if latest == nil {
		return currentTag, nil
	}

	if isMajorOnly(currentTag) {
		major := latest.Segments()[0]
		if latest.Original()[0] == 'v' {
			return fmt.Sprintf("v%d", major), nil
		}
		return fmt.Sprintf("%d", major), nil
	}

	return latest.Original(), nil
}

func isMajorOnly(v string) bool {
	return !strings.Contains(v, ".")
}
