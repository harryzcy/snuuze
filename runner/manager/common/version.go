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
	return getLatestTagMultiParts(input)
}

type MultiVersion struct {
	Original string
	Parts    []string
}

func getLatestTagMultiParts(input *GetLatestTagInput) (string, error) {
	isMajorOnly := !strings.Contains(input.CurrentTag, ".")

	current := MultiVersion{
		Original: input.CurrentTag,
		Parts:    strings.Split(input.CurrentTag, input.Delimiter),
	}

	choices := []MultiVersion{}
	for _, tag := range input.Tags {
		tagParts := strings.Split(tag, input.Delimiter)
		if len(tagParts) != 2 {
			continue
		}
		if len(current.Parts) != len(tagParts) {
			continue
		}
		choices = append(choices, MultiVersion{
			Original: tag,
			Parts:    tagParts,
		})
	}

	latest := []MultiVersion{}

	partsNumber := len(current.Parts)
	for i := 0; i < partsNumber; i++ {
		for _, next := range choices {
			prefix := letterPrefix(current.Parts[i])
			currentPart := strings.TrimPrefix(current.Parts[i], prefix)
			nextPart := strings.TrimPrefix(next.Parts[i], prefix)
			if !containValidPart(currentPart, nextPart) {
				continue
			}

			if len(latest) == 0 {
				latest = []MultiVersion{next}
				continue
			}
			latestPart := strings.TrimPrefix(latest[0].Parts[i], prefix)

			if greater, equal, err := isGreaterAndEqual(latestPart, nextPart, isMajorOnly, input.AllowMajor); err != nil {
				return "", err
			} else if greater {
				latest = []MultiVersion{next}
			} else if equal {
				latest = append(latest, next)
			}
		}

		choices = latest
		if len(choices) == 0 {
			return input.CurrentTag, nil
		}
		latest = []MultiVersion{}
	}

	return choices[0].Original, nil
}

func containValidPart(currentPart, nextPart string) bool {
	if len(currentPart) == 0 && len(nextPart) != 0 {
		return false
	}
	if len(currentPart) != 0 && len(nextPart) == 0 {
		return false
	}
	return true
}

func getLatestTagSinglePart(input *GetLatestTagInput) (string, error) {
	latest := input.CurrentTag
	isMajorOnly := !strings.Contains(latest, ".")

	for _, tag := range input.Tags {
		if tag == "" {
			continue
		}
		if greater, _, err := isGreaterAndEqual(latest, tag, isMajorOnly, input.AllowMajor); err != nil {
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

// isGreaterAndEqual returns true if nextTag is greater than currentTag.
func isGreaterAndEqual(currentTag, nextTag string, isMajorOnly, allowMajor bool) (bool, bool, error) {
	current, err := version.NewVersion(currentTag)
	if err != nil {
		return false, false, err
	}

	next, err := version.NewVersion(nextTag)
	if err != nil {
		fmt.Printf("warning: failed to parse tag (%s), ignoring\n", nextTag)
		return false, false, nil
	}

	if !allowMajor && next.Segments()[0] != current.Segments()[0] {
		return false, false, nil
	}

	// segments of different length
	if !isMajorOnly && segmentLength(nextTag) != segmentLength(currentTag) {
		return false, false, nil
	}

	// ignore weird versions likely it's SemVer vs CalVer.
	// e.g. alpine 3.18.4 vs 20230901 in docker
	if current.Segments()[0] < 200 &&
		next.Segments()[0] > 100000 && next.Segments()[1] == 0 && next.Segments()[2] == 0 {
		return false, false, nil
	}

	// don't upgrade from release to pre-release
	if current.Prerelease() == "" && next.Prerelease() != "" {
		return false, false, nil
	}

	return next.GreaterThan(current), next.Equal(current), nil
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
