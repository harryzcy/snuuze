package common

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

// GetLatestTag returns the latest tag that is not a pre-release, or the current tag if no such tag exists
func GetLatestTag(depName string, tags []string, currentTag string, includeMajor bool) (string, error) {
	currentVersion, err := version.NewVersion(currentTag)
	if err != nil {
		return "", err
	}
	versions := make([]*version.Version, 0, len(tags))
	for _, tag := range tags {
		v, err := version.NewVersion(tag)
		if err != nil {
			fmt.Printf("warning: failed to parse tag (%s) for %s, ignoring\n", tag, depName)
			continue
		}
		if !includeMajor {
			if v.Segments()[0] != currentVersion.Segments()[0] {
				continue
			}
		}

		// ignore weird versions likely it's SemVer vs CalVer.
		// e.g. alpine 3.18.4 vs 20230901
		if currentVersion.Segments()[0] < 200 && v.Segments()[0] > 100000 {
			continue
		}

		versions = append(versions, v)
	}
	if len(versions) == 0 {
		return currentTag, nil
	}

	fmt.Println(36, currentTag, versions)

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
