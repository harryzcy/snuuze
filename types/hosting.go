package types

import (
	"strings"
	"time"
)

type HostingConfig struct {
	Data    DataConfig    `yaml:"data"`
	Network NetworkConfig `yaml:"network"`
	GitHub  GitHubConfig  `yaml:"github"`
	Gitea   []GiteaConfig `yaml:"gitea"`
}

type DataConfig struct {
	TempDir string `yaml:"tempDir"`
	// Timeout for accessing files and running commands
	Timeout int64 `yaml:"timeout"` // in seconds
}

func (d DataConfig) GetTimeout() time.Duration {
	if d.Timeout <= 0 {
		return 100 * time.Second // default to 100 seconds
	}
	return time.Duration(d.Timeout) * time.Second
}

type NetworkConfig struct {
	Timeout int64 `yaml:"timeout"` // in seconds
}

func (n NetworkConfig) GetTimeout() time.Duration {
	if n.Timeout <= 0 {
		return 10 * time.Second // default to 10 seconds
	}
	return time.Duration(n.Timeout) * time.Second
}

type GitHubConfig struct {
	AuthType string `yaml:"authType"` // token, github-app
	// if auth-type is token
	Token string `yaml:"token"`

	// if auth-type is github-app
	AppID          int64  `yaml:"appID"`
	PEMFile        string `yaml:"pemFile"`
	ClientID       string `yaml:"clientID"`
	InstallationID int64  `yaml:"installationID"`
	AppName        string `yaml:"appName"`
	AppUserID      int64  `yaml:"appUserID"`
}

type GiteaConfig struct {
	Host     string `yaml:"host"`
	AuthType string `yaml:"authType"` // token

	// if auth-type is token
	Token string `yaml:"token"`
}

func (c *GiteaConfig) GetHost() string {
	return strings.TrimSuffix(c.Host, "/")
}

const (
	AuthTypeToken     = "token"
	AuthTypeGithubApp = "github-app"
)
