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
func GetLatestTag(input *GetLatestTagInput) (string, error) {
	isMultiPart := input.Delimiter != "" && strings.Contains(input.CurrentTag, input.Delimiter)
	if !isMultiPart {
		return getLatestTagSinglePart(input)
	}
	return getLatestTagTwoParts(input)
}

type MultiVersion struct {
	Original string
	Parts    []string
}

func getLatestTagTwoParts(input *GetLatestTagInput) (string, error) {
	isMajorOnly := !strings.Contains(input.CurrentTag, ".")

	choices := []MultiVersion{}
	for _, tag := range input.Tags {
		tagParts := strings.Split(tag, input.Delimiter)
		if len(tagParts) != 2 {
			continue
		}
		choices = append(choices, MultiVersion{
			Original: tag,
			Parts:    tagParts,
		})
	}

	latest := MultiVersion{
		Original: input.CurrentTag,
		Parts:    strings.Split(input.CurrentTag, input.Delimiter),
	}
	// get possible versions based on first part
	filtered := []MultiVersion{}
	for _, choice := range choices {
		if len(latest.Parts) != len(choice.Parts) {
			continue
		}

		if greater, err := isGreater(latest.Parts[0], choice.Parts[0], isMajorOnly, input.AllowMajor); err != nil {
			return "", err
		} else if greater {
			filtered = append(filtered, choice)
		}
	}

	for _, current := range filtered {
		latestPart := latest.Parts[1]
		tagPart := current.Parts[1]

		prefix := letterPrefix(latestPart)
		if strings.HasPrefix(tagPart, prefix) {
			latestPart = latestPart[len(prefix):]
			tagPart = tagPart[len(prefix):]
		}

		if latestPart == tagPart {
			latest = current
			continue
		}
		if len(latestPart) == 0 || len(tagPart) == 0 {
			continue
		}

		if greater, err := isGreater(latestPart, tagPart, isMajorOnly, input.AllowMajor); err != nil {
			return "", err
		} else if greater {
			latest = current
		}
	}

	return latest.Original, nil
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
	if !isMajorOnly && segmentLength(nextTag) != segmentLength(currentTag) {
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

func segmentLength(version string) int {
	return strings.Count(version, ".") + 1
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
