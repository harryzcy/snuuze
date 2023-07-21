package gomod

import (
	"path/filepath"

	"github.com/harryzcy/snuuze/manager/common"
	"github.com/harryzcy/snuuze/types"
)

type GolangManager struct{}

func New() common.Manager {
	return &GolangManager{}
}

func (m *GolangManager) Name() types.PackageManager {
	return types.PackageManagerGoMod
}

func (m *GolangManager) Match(path string) bool {
	filename := filepath.Base(path)
	return filename == "go.mod"
}
