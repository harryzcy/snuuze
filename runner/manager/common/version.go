package common

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
)

type GetLatestTagInput struct {
	DepName    string
	Tags       []string
	CurrentTag string
	// AllowMajor specifies whether to allow major version upgrade.
	AllowMajor bool
	// Delimiter is optional, used to split the tag, i.e. docker image tag go:1.20-alpine3.18
	Delimiter string
}

// GetLatestTag returns the latest tag that is not a pre-release,
// or the current tag if no such tag exists.
// If a delimiter is provided, it will be used to split the tag
func GetLatestTag(input *GetLatestTagInput) (string, error) {
	latest := input.CurrentTag
	for _, tag := range input.Tags {
		if greater, err := isGreater(latest, tag, input.AllowMajor); err != nil {
			return "", err
		} else if greater {
			latest = tag
		}
	}

	if isMajorOnly(input.CurrentTag) {
		fullVersion, err := version.NewVersion(latest)
		if err != nil {
			return "", err
		}
		major := fullVersion.Segments()[0]
		if fullVersion.Original()[0] == 'v' {
			return fmt.Sprintf("v%d", major), nil
		}
		return fmt.Sprintf("%d", major), nil
	}

	return latest, nil
}

func isGreater(currentTag, nextTag string, allowMajor bool) (bool, error) {
	current, err := version.NewVersion(currentTag)
	if err != nil {
		return false, err
	}

	next, err := version.NewVersion(nextTag)
	if err != nil {
		fmt.Printf("warning: failed to parse tag (%s), ignoring\n", next)
		return false, nil
	}

	if !allowMajor && next.Segments()[0] != current.Segments()[0] {
		return false, nil
	}

	// ignore weird versions likely it's SemVer vs CalVer.
	// e.g. alpine 3.18.4 vs 20230901
	if current.Segments()[0] < 200 &&
		next.Segments()[0] > 100000 && next.Segments()[1] == 0 && next.Segments()[2] == 0 {
		return false, nil
	}

	// don't upgrade from release to pre-release
	if current.Prerelease() == "" && next.Prerelease() != "" {
		return false, nil
	}

	return next.GreaterThan(current), nil
}

func isMajorOnly(v string) bool {
	return !strings.Contains(v, ".")
}
