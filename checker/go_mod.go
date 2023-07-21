package checker

import (
	"errors"

	"github.com/harryzcy/snuuze/thirdparty/gomajor"
	"github.com/harryzcy/snuuze/types"
)

var (
	ErrRequestFailed = errors.New("request failed")
)

type GoListOutput struct {
	Path    string
	Version string
	Time    string
	Update  *struct {
		Path    string
		Version string
		Time    string
	}
	Dir   string
	GoMod string
}

func isUpgradable_GoMod(dep types.Dependency) (*types.UpgradeInfo, error) {
	mod, err := gomajor.Latest(dep.Name, false)
	if err != nil {
		return nil, err
	}

	latestVersion := mod.MaxVersion("", true)
	if gomajor.IsNewerVersion(dep.Version, latestVersion, false) {
		return &types.UpgradeInfo{
			Dependency: dep,
			Upgradable: true,
			ToVersion:  latestVersion,
		}, nil
	}

	return nil, nil
}
