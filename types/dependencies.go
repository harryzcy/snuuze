package types

import "github.com/harryzcy/sailor/matcher"

type Dependency struct {
	Name           string
	Version        string
	Indirect       bool
	PackageManager matcher.PackageManager
	Extra          map[string]interface{}
}
