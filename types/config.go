package types

type Config struct {
	Version string   `mapstructure:"version"`
	Presets []string `mapstructure:"presets"`
	Rules   []Rule   `mapstructure:"rules"`
}

type Rule struct {
	PackageManagers []PackageManager `mapstructure:"package-managers"`
	PackageNames    []string         `mapstructure:"package-names"`
	PackageTypes    []PackageType    `mapstructure:"package-types"` // direct, indirect, all
	Labels          []string         `mapstructure:"labels"`        // github labels
}

type PackageType string

const (
	PackageTypeDirect   PackageType = "direct"
	PackageTypeIndirect PackageType = "indirect"
	PackageTypeAll      PackageType = "all"
)
