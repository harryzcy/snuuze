package common

import "github.com/harryzcy/snuuze/types"

type Manager interface {
	Name() types.PackageManager
	Match(path string) bool
}
