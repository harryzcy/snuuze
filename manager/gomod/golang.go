package gomod

import (
	"path/filepath"

	"github.com/harryzcy/snuuze/manager/common"
)

type GolangManager struct{}

func New() common.Manager {
	return &GolangManager{}
}

func (m *GolangManager) Name() string {
	return "golang"
}

func (m *GolangManager) Match(path string) bool {
	filename := filepath.Base(path)
	return filename == "go.mod"
}
