package common

import "github.com/harryzcy/snuuze/types"

// Manager is the interface that all package managers must implement
type Manager interface {
	Name() types.PackageManager
	Match(path string) bool
}
