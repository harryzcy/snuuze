package platform

import "strings"

type Client interface {
	// ListTags returns a sorted list of tags for the given repo
	ListTags(params *ListTagsInput) ([]string, error)

	CreatePullRequest(input *CreatePullRequestInput) error
}

const (
	GitPlatformGitHub = "github.com"
	GitPlatformGitea  = "gitea.com"
)

// NewClient returns a new Client based on Git URL
func NewClient(url string) (Client, error) {
	switch GitPlatform(url) {
	case GitPlatformGitHub:
		return NewGitHubClient()
	case GitPlatformGitea:
		return NewGiteaClient(), nil
	}

	// TODO: check from config after configurations are implemented
	return NewGiteaClient(), nil
}

func GitPlatform(gitURL string) string {
	gitURL = strings.TrimPrefix(gitURL, "git@")
	gitURL = strings.TrimPrefix(gitURL, "https://")
	gitURL = strings.TrimPrefix(gitURL, "http://")

	if strings.HasPrefix(gitURL, "github.com/") {
		return GitPlatformGitHub
	}

	return GitPlatformGitea
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
