package platform

import (
	"errors"
	"strings"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/types"
)

type Client interface {
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

// NewClient returns a new Client based on Git URL
func NewClient(url string) (Client, error) {
	platform, host := DetermineGitPlatform(url)
	switch platform {
	case GitPlatformGitHub:
		return NewGitHubClient()
	case GitPlatformGitea:
		return NewGiteaClient(host), nil
	}

	return nil, errors.New("unsupported git platform")
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

	return nil, errors.New("invalid git url")
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
