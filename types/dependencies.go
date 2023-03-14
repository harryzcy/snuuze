package types

type Dependency struct {
	File           string // file path of go.mod
	Name           string
	Version        string
	Indirect       bool
	PackageManager PackageManager
	Position       Position // position is only used for some package managers
	Extra          map[string]interface{}
}

type Position struct {
	Line      int
	StartByte int
	EndByte   int
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
