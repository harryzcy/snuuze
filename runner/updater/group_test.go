package updater

import (
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestGroupUpdates(t *testing.T) {
	originalGetRules := getRules
	defer func() {
		getRules = originalGetRules
	}()

	infos := []*types.UpgradeInfo{
		{
			Dependency: types.Dependency{
				PackageManager: types.PackageManagerGoMod,
				Name:           "github.com/harryzcy/snuuze",
				Indirect:       false,
			},
		},
		{
			Dependency: types.Dependency{
				PackageManager: types.PackageManagerGoMod,
				Name:           "github.com/stretchr/testify",
				Indirect:       false,
			},
		},
	}
	getRules = func() []types.Rule {
		return []types.Rule{
			{
				PackageManagers: []types.PackageManager{types.PackageManagerGoMod},
			},
		}
	}

	groups := groupUpdates(infos)
	assert.Len(t, groups, 1)
	assert.Len(t, groups[0].Infos, 2)
}

func TestMatchRule(t *testing.T) {
	tests := []struct {
		info *types.UpgradeInfo
		rule types.Rule
	}{
		{
			info: &types.UpgradeInfo{
				Dependency: types.Dependency{
					PackageManager: types.PackageManagerGoMod,
					Name:           "github.com/harryzcy/snuuze",
					Indirect:       false,
				},
			},
			rule: types.Rule{
				PackageManagers: []types.PackageManager{types.PackageManagerGoMod},
				PackageNames:    []string{"github.com/harryzcy/snuuze"},
				PackageTypes:    []types.PackageType{types.PackageTypeDirect},
			},
		},
		{
			info: &types.UpgradeInfo{
				Dependency: types.Dependency{
					PackageManager: types.PackageManagerGoMod,
					Name:           "github.com/stretchr/testify",
					Indirect:       false,
				},
			},
			rule: types.Rule{
				PackageManagers: []types.PackageManager{types.PackageManagerGoMod},
			},
		},
	}

	for i, test := range tests {
		assert.True(t, matchRule(test.rule, test.info), "test %d", i)
	}
}
