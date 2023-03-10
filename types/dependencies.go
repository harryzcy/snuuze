package types

import "github.com/harryzcy/snuuze/matcher"

type Dependency struct {
	Name           string
	Version        string
	Indirect       bool
	PackageManager matcher.PackageManager
	Extra          map[string]interface{}
}
