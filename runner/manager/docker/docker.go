package docker

import (
	"path/filepath"

	"github.com/harryzcy/snuuze/runner/manager/common"
	"github.com/harryzcy/snuuze/types"
)

type DockerManager struct{}

func New() common.Manager {
	return &DockerManager{}
}

func (m *DockerManager) Name() types.PackageManager {
	return types.PackageManagerDocker
}

func (m *DockerManager) Match(path string) bool {
	filename := filepath.Base(path)
	return filename == "Dockerfile"
}

func (m *DockerManager) Parse(match types.Match, data []byte) ([]*types.Dependency, error) {
	return nil, nil
}

func (m *DockerManager) FindDependencies(matches []types.Match) ([]*types.Dependency, error) {
	return common.FindDependencies(m, matches)
}

func (m *DockerManager) ListUpgrades(matches []types.Match) ([]*types.UpgradeInfo, error) {
	return common.ListUpgrades(m, matches)
}

func (m *DockerManager) IsUpgradable(dep types.Dependency) (*types.UpgradeInfo, error) {
	return nil, nil
}
