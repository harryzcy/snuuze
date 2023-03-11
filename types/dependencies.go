package types

type Dependency struct {
	Name           string
	Version        string
	Indirect       bool
	PackageManager PackageManager
	Extra          map[string]interface{}
}

type UpgradeInfo struct {
	Dependency Dependency
	Upgradable bool
	ToVersion  string
}

type PackageManager string

const (
	PackageManagerGoMod         PackageManager = "go-mod"
	PackageManagerGitHubActions PackageManager = "github-actions"
)
