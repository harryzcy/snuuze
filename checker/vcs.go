package checker

import (
	"context"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")

type VSCSource interface {
	// ListTags returns a sorted list of tags for the given repo
	ListTags(params *ListTagsParameters) ([]string, error)
}

type ListTagsParameters struct {
	Owner  string
	Repo   string
	Prefix string // optional
}

type GitHubSource struct{}

func (g *GitHubSource) ListTags(params *ListTagsParameters) ([]string, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	var query struct {
		Repository struct {
			Refs struct {
				Edges []struct {
					Node struct {
						Name string
					}
				}
			} `graphql:"refs(refPrefix: $refPrefix, first: 100, orderBy: {field: TAG_COMMIT_DATE, direction: DESC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	refPrefix := "refs/tags/"
	if params.Prefix != "" {
		refPrefix += params.Prefix
		if !strings.HasSuffix(refPrefix, "/") {
			refPrefix += "/"
		}
	}
	variables := map[string]interface{}{
		"owner":     githubv4.String(params.Owner),
		"name":      githubv4.String(params.Repo),
		"refPrefix": githubv4.String(refPrefix),
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0)
	for _, edge := range query.Repository.Refs.Edges {
		tags = append(tags, edge.Node.Name)
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
	if len(sortedTags) == 0 {
		return ""
	}
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
