package common

import (
	"fmt"
	"strings"
	"unicode"

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
	isMultiPart := input.Delimiter != "" && strings.Contains(input.CurrentTag, input.Delimiter)
	if !isMultiPart {
		return getLatestTagSinglePart(input)
	}
	return getLatestTagTwoParts(input)
}

func getLatestTagTwoParts(input *GetLatestTagInput) (string, error) {
	isMajorOnly := !strings.Contains(input.CurrentTag, ".")

	// get possible versions based on first part
	possibleVersions := []string{}
	for _, tag := range input.Tags {
		latestParts := strings.Split(input.CurrentTag, input.Delimiter)
		tagParts := strings.Split(tag, input.Delimiter)

		if len(latestParts) != len(tagParts) {
			continue
		}

		if greater, err := isGreater(latestParts[0], tagParts[0], isMajorOnly, input.AllowMajor); err != nil {
			return "", err
		} else if greater {
			possibleVersions = append(possibleVersions, tag)
		}
	}

	latest := input.CurrentTag
	for _, tag := range possibleVersions {
		latestPart := strings.Split(latest, input.Delimiter)[1]
		tagPart := strings.Split(tag, input.Delimiter)[1]

		prefix := letterPrefix(latestPart)
		if strings.HasPrefix(tagPart, prefix) {
			latestPart = latestPart[len(prefix):]
			tagPart = tagPart[len(prefix):]
		}

		if latestPart == tagPart {
			latest = tag
			continue
		}
		if len(latestPart) == 0 || len(tagPart) == 0 {
			continue
		}

		if greater, err := isGreater(latestPart, tagPart, isMajorOnly, input.AllowMajor); err != nil {
			return "", err
		} else if greater {
			latest = tag
		}
	}

	return latest, nil
}

func getLatestTagSinglePart(input *GetLatestTagInput) (string, error) {
	latest := input.CurrentTag
	isMajorOnly := !strings.Contains(latest, ".")

	for _, tag := range input.Tags {
		if greater, err := isGreater(latest, tag, isMajorOnly, input.AllowMajor); err != nil {
			return "", err
		} else if greater {
			latest = tag
		}
	}

	if isMajorOnly {
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

// isGreater returns true if nextTag is greater than currentTag.
func isGreater(currentTag, nextTag string, isMajorOnly, allowMajor bool) (bool, error) {
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

	// segments of different length
	if !isMajorOnly && segmentLength(next.Segments()) != segmentLength(current.Segments()) {
		return false, nil
	}

	// ignore weird versions likely it's SemVer vs CalVer.
	// e.g. alpine 3.18.4 vs 20230901 in docker
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

func segmentLength(segments []int) int {
	l := 0
	for _, s := range segments {
		if s == 0 {
			break
		}
		l++
	}
	return l
}

func letterPrefix(version string) string {
	prefix := ""
	for _, c := range version {
		if unicode.IsLetter(c) {
			prefix += string(c)
		} else {
			break
		}
	}
	return prefix
}
