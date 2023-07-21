package common

import (
	"fmt"
	"os"

	"github.com/harryzcy/snuuze/types"
)

// Manager is the interface that all package managers must implement
type Manager interface {
	Name() types.PackageManager
	Match(path string) bool
	Parse(match types.Match, data []byte) ([]types.Dependency, error)

	// ListUpgrades returns a list of upgrades for the given matches.
	// This could be implemented by ListUpgrades in this package, which calls IsUpgradable in a loop.
	// Or it could be implemented by the package manager itself, which may be more efficient.
	ListUpgrades(matches []types.Match) ([]*types.UpgradeInfo, error)
	IsUpgradable(dep types.Dependency) (*types.UpgradeInfo, error)
}

// ListUpgrades provides a common implementation that lists all the upgrades in a loop
func ListUpgrades(m Manager, matches []types.Match) ([]*types.UpgradeInfo, error) {
	result := []*types.UpgradeInfo{}

	for _, match := range matches {
		data, err := os.ReadFile(match.File)
		if err != nil {
			return nil, err
		}

		dependencies, _ := m.Parse(match, data)
		for _, dependency := range dependencies {
			info, err := m.IsUpgradable(dependency)
			if err != nil {
				// log error and continue
				fmt.Println(err)
				continue
			}
			if info != nil && info.Upgradable {
				result = append(result, info)
			}
		}
	}

	return result, nil
}
