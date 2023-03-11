package types

type Config struct {
	Version string
	Presets []string
	Rules   []Rule
}

type Rule struct {
	PackageManagers []PackageManager
	PackageNames    []string
	PackageTypes    []PackageType // direct, indirect, all
	Labels          []string
}

type PackageType string

const (
	PackageTypeDirect   PackageType = "direct"
	PackageTypeIndirect PackageType = "indirect"
	PackageTypeAll      PackageType = "all"
)
