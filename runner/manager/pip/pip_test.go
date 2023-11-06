package pip

import (
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestPipManager_IsUpgradable(t *testing.T) {
	dep := &types.Dependency{
		Name:           "requests",
		Version:        "== 2.25.1",
		PackageManager: types.PackageManagerPip,
		Extra: map[string]interface{}{
			"constraints": [][2]string{
				{"==", "2.25.1"},
			},
		},
	}
	info, err := New().IsUpgradable(*dep)
	assert.NoError(t, err)
	assert.True(t, info.Upgradable)
}

func TestGetPipPackageVersions(t *testing.T) {
	versions, err := getPipPackageVersions("requests")
	assert.NoError(t, err)
	assert.NotEmpty(t, versions)
	assert.Contains(t, versions, "2.31.0")

	_, err = getPipPackageVersions("invalid-package-from-pypi")
	assert.Error(t, err)
}
