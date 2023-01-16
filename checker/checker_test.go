package checker

import (
	"strconv"
	"testing"

	"github.com/harryzcy/sailor/matcher"
	"github.com/harryzcy/sailor/types"
	"github.com/stretchr/testify/assert"
)

func TestIsUpgradable(t *testing.T) {
	tests := []struct {
		dep  types.Dependency
		want UpgradeInfo
		err  error
	}{
		{
			dep: types.Dependency{
				Name:           "github.com/gofiber/fiber/v2",
				Version:        "v2.40.0", // use an old version
				PackageManager: matcher.GoMod,
			},
			want: UpgradeInfo{
				Dependency: types.Dependency{
					Name:           "github.com/gofiber/fiber/v2",
					Version:        "v2.40.0",
					PackageManager: matcher.GoMod,
				},
				Upgradable: true,
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			info, err := isUpgradable(test.dep)
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.want.Upgradable, info.Upgradable)
			assert.Equal(t, test.want.Dependency, info.Dependency)
		})
	}
}
