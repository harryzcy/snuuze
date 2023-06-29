package types

import "strings"

type HostingConfig struct {
	Data   DataConfig    `yaml:"data"`
	GitHub GitHubConfig  `yaml:"github"`
	Gitea  []GiteaConfig `yaml:"gitea"`
}

type DataConfig struct {
	TempDir string `yaml:"tempDir"`
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
