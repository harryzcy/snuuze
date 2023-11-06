package pip

import (
	"path/filepath"
	"strings"

	"github.com/harryzcy/snuuze/runner/manager/common"
	"github.com/harryzcy/snuuze/types"
)

type PipManager struct{}

func New() common.Manager {
	return &PipManager{}
}

func (m *PipManager) Name() types.PackageManager {
	return types.PackageManagerPip
}

func (m *PipManager) Match(path string) bool {
	filename := filepath.Base(path)
	return filename == "requirements.txt"
}

func (m *PipManager) Parse(match types.Match, data []byte) ([]*types.Dependency, error) {
	text := string(data)

	lines := strings.Split(text, "\n")

	dependencies := make([]*types.Dependency, 0)
	for _, line := range lines {
		parser := NewParser(match.File, line)
		dep, err := parser.Parse()
		if err != nil {
			continue
		}
		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}

func (m *PipManager) FindDependencies(matches []types.Match) ([]*types.Dependency, error) {
	return common.FindDependencies(m, matches)
}

func (m *PipManager) ListUpgrades(matches []types.Match) ([]*types.UpgradeInfo, error) {
	return common.ListUpgrades(m, matches)
}

func (m *PipManager) IsUpgradable(dep types.Dependency) (*types.UpgradeInfo, error) {
	return nil, nil
}
