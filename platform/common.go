package platform

import (
	"context"
	"errors"
	"strings"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/types"
)

var (
	ErrInvalidGitPlatform = errors.New("unsupported git platform")
	ErrInvalidGitURL      = errors.New("invalid git url")
)

type Client interface {
	Token(ctx context.Context) (string, error)

	ListRepos() ([]Repo, error)
	// ListTags returns a sorted list of tags for the given repo
	ListTags(params *ListTagsInput) ([]string, error)

	CreatePullRequest(input *CreatePullRequestInput) error
}

type GitPlatform int

const (
	GitPlatformUnknown GitPlatform = iota
	GitPlatformGitHub
	GitPlatformGitea
)

type NewClientOptions struct {
	Platform GitPlatform
	URL      string
}

// NewClient returns a new Client based on Git URL
func NewClient(options NewClientOptions) (Client, error) {
	platform := options.Platform
	var host string
	if platform == GitPlatformUnknown || platform == GitPlatformGitea {
		platform, host = DetermineGitPlatform(options.URL)
	}

	switch platform {
	case GitPlatformGitHub:
		return NewGitHubClient()
	case GitPlatformGitea:
		return NewGiteaClient(host), nil
	}

	return nil, ErrInvalidGitPlatform
}

var getGiteaConfigs = func() []types.GiteaConfig {
	return config.GetHostingConfig().Gitea
}

// DetermineGitPlatform returns the GitPlatform and the host of the given git URL
func DetermineGitPlatform(gitURL string) (GitPlatform, string) {
	urlInfo, err := parseURL(gitURL)
	if err != nil {
		return GitPlatformUnknown, ""
	}

	if urlInfo.host == "github.com" {
		return GitPlatformGitHub, "https://github.com"
	}

	for _, giteaConfig := range getGiteaConfigs() {
		configuredHost := giteaConfig.GetHost()
		configuredInfo, err := parseURL(configuredHost)
		if err != nil {
			continue
		}

		if urlInfo.host == configuredInfo.host {
			return GitPlatformGitea, configuredHost
		}
	}

	return GitPlatformUnknown, ""
}

type hostInfo struct {
	protocol string
	host     string
}

func parseURL(url string) (*hostInfo, error) {
	for _, protocol := range []string{"https", "http"} {
		prefix := protocol + "://"
		if strings.HasPrefix(url, prefix) {
			url = strings.TrimPrefix(url, prefix)
			url = strings.SplitN(url, "/", 2)[0]
			return &hostInfo{
				protocol: protocol,
				host:     url,
			}, nil
		}
	}

	if strings.HasPrefix(url, "git@") {
		url = strings.TrimPrefix(url, "git@")
		url = strings.SplitN(url, ":", 2)[0]
		return &hostInfo{
			protocol: "ssh",
			host:     url,
		}, nil
	}

	return nil, ErrInvalidGitURL
}

type Repo struct {
	Server        string `json:"server"`
	Owner         string `json:"owner"`
	Repo          string `json:"repo"`
	URL           string `json:"url"`
	IsPrivate     bool   `json:"isPrivate"`
	DefaultBranch string `json:"defaultBranch"`
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
