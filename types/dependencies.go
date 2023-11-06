package types

import "fmt"

type Dependency struct {
	File           string // file path of go.mod
	Name           string
	Version        string
	Indirect       bool
	PackageManager PackageManager
	Position       Position // position is only used for some package managers
	Extra          map[string]interface{}
}

func (d Dependency) Hash() string {
	return fmt.Sprintf("%s:%s:%s", d.File, d.Name, d.PackageManager)
}

type Position struct {
	Line      int
	ColStart  int
	ColEnd    int
	ByteStart int
	ByteEnd   int
}

type UpgradeInfo struct {
	Dependency Dependency
	Upgradable bool
	ToVersion  string
}

type PackageManager string

const (
	PackageManagerDocker        PackageManager = "docker"
	PackageManagerGitHubActions PackageManager = "github-actions"
	PackageManagerGoMod         PackageManager = "go-mod"
	PackageManagerPip           PackageManager = "pip"
)
