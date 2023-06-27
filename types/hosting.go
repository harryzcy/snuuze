package types

type HostingConfig struct {
	Data   DataConfig   `yaml:"data"`
	GitHub GitHubConfig `yaml:"github"`
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
}
