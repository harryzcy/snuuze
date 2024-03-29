package gomod

import (
	"path/filepath"

	"github.com/harryzcy/snuuze/runner/manager/common"
	"github.com/harryzcy/snuuze/thirdparty/gomajor"
	"github.com/harryzcy/snuuze/types"
	"golang.org/x/mod/modfile"
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

func (m *GolangManager) Parse(match types.Match, data []byte) ([]*types.Dependency, error) {
	filename := filepath.Base(match.File)
	file, err := modfile.Parse(filename, data, nil)
	if err != nil {
		return nil, err
	}

	dependencies := make([]*types.Dependency, 0, len(file.Require))
	for _, require := range file.Require {
		dependencies = append(dependencies, &types.Dependency{
			File:           match.File,
			Name:           require.Mod.Path,
			Version:        require.Mod.Version,
			Indirect:       require.Indirect,
			PackageManager: types.PackageManagerGoMod,
			Position: types.Position{
				Line:      require.Syntax.Start.Line,
				ColStart:  require.Syntax.Start.LineRune,
				ColEnd:    require.Syntax.End.LineRune,
				ByteStart: require.Syntax.Start.Byte,
				ByteEnd:   require.Syntax.End.Byte,
			},
		})
	}

	return dependencies, nil
}

func (m *GolangManager) FindDependencies(matches []types.Match) ([]*types.Dependency, error) {
	return common.FindDependencies(m, matches)
}

func (m *GolangManager) ListUpgrades(matches []types.Match) ([]*types.UpgradeInfo, error) {
	return common.ListUpgrades(m, matches)
}

func (m *GolangManager) IsUpgradable(dep types.Dependency) (*types.UpgradeInfo, error) {
	mod, err := gomajor.Query(dep.Name, false)
	if err != nil {
		return nil, err
	}

	latest, err := common.GetLatestTag(&common.GetLatestTagInput{
		DepName:    dep.Name,
		Tags:       mod.Versions(),
		CurrentTag: dep.Version,
		AllowMajor: !dep.Indirect,
	})
	if err != nil {
		return nil, err
	}
	if gomajor.IsNewerVersion(dep.Version, latest, false) {
		return &types.UpgradeInfo{
			Dependency: dep,
			Upgradable: true,
			ToVersion:  latest,
		}, nil
	}

	return nil, nil
}
