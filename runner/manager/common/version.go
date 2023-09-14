package common

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

// GetLatestTag returns the latest tag that is not a pre-release, or the current tag if no such tag exists
func GetLatestTag(tags []string, currentTag string, includeMajor bool) (string, error) {
	currentVersion, err := version.NewVersion(currentTag)
	if err != nil {
		return "", err
	}
	versions := make([]*version.Version, 0, len(tags))
	for _, tag := range tags {
		v, err := version.NewVersion(tag)
		if err != nil {
			return "", err
		}
		if !includeMajor {
			if v.Segments()[0] != currentVersion.Segments()[0] {
				continue
			}
		}

		versions = append(versions, v)
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