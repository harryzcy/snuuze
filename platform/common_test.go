package platform

import (
	"errors"
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestDetermineGitPlatform(t *testing.T) {
	originalGetGiteaConfigs := getGiteaConfigs
	defer func() {
		getGiteaConfigs = originalGetGiteaConfigs
	}()
	getGiteaConfigs = func() []types.GiteaConfig {
		return []types.GiteaConfig{
			{
				Host: "https://gitea.com",
			},
			{
				Host: "invalid",
			},
			{
				Host: "https://git.example.com",
			},
		}
	}

	tests := []struct {
		gitURL   string
		platform GitPlatform
		host     string
	}{
		{
			gitURL:   "https://github.com/owner/repo",
			platform: GitPlatformGitHub,
			host:     "https://github.com",
		},
		{
			gitURL:   "git@github.com:owner/repo",
			platform: GitPlatformGitHub,
			host:     "https://github.com",
		},
		{
			gitURL:   "git@gitea.com:owner/repo",
			platform: GitPlatformGitea,
			host:     "https://gitea.com",
		},
		{
			gitURL:   "https://git.example.com/owner/repo",
			platform: GitPlatformGitea,
			host:     "https://git.example.com",
		},
		{
			gitURL:   "invalid",
			platform: GitPlatformUnknown,
		},
		{
			gitURL:   "http://no.such.host/owner/repo",
			platform: GitPlatformUnknown,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			platform, host := DetermineGitPlatform(test.gitURL)
			assert.Equal(t, test.platform, platform)
			assert.Equal(t, test.host, host)
		})
	}
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		url      string
		host     string
		protocol string
		err      error
	}{
		{
			url:      "https://github.com/owner/repo",
			host:     "github.com",
			protocol: "https",
		},
		{
			url:      "http://git.example.com",
			host:     "git.example.com",
			protocol: "http",
		},
		{
			url:      "git@github.com:owner/repo",
			host:     "github.com",
			protocol: "ssh",
		},
		{
			url: "invalid@github.com",
			err: errors.New("invalid git url"),
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			hostInfo, err := parseURL(test.url)
			assert.Equal(t, test.err, err)
			if err == nil {
				assert.Equal(t, test.host, hostInfo.host)
				assert.Equal(t, test.protocol, hostInfo.protocol)
			}
		})
	}
}
