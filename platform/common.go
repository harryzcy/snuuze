package platform

import "strings"

type Client interface {
	// ListTags returns a sorted list of tags for the given repo
	ListTags(params *ListTagsInput) ([]string, error)

	CreatePullRequest(input *CreatePullRequestInput) error
}

// NewClient returns a new Client based on Git URL
func NewClient(url string) (Client, error) {
	url = strings.TrimPrefix(url, "git@")
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	if strings.HasPrefix(url, "github.com/") {
		return NewGitHubClient()
	}

	// TODO: check from config after configurations are implemented
	return NewGiteaClient(), nil
}

type ListTagsInput struct {
	Owner  string
	Repo   string
	Prefix string // optional
}

type CreatePullRequestInput struct {
	Title string
	Body  string
	Base  string
	Head  string
	Owner string
	Repo  string
}
