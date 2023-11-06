package docker

import (
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestDockerManager_IsUpgradable(t *testing.T) {
	manager := New()
	info, err := manager.IsUpgradable(types.Dependency{
		Name:    "alpine",
		Version: "3.18.0",
	})
	assert.NoError(t, err)
	assert.True(t, info.Upgradable)
}

func TestGetDockerImageTags(t *testing.T) {
	tags, err := getDockerImageTags("library/alpine")
	assert.NoError(t, err)
	assert.NotEmpty(t, tags)
	assert.Contains(t, tags, "3.18.3")
}

func TestParseImageName(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		image    string
	}{
		{
			name:     "alpine",
			endpoint: "index.docker.io",
			image:    "alpine",
		},
		{
			name:     "ghcr.io/harryzcy/snuuze",
			endpoint: "ghcr.io",
			image:    "harryzcy/snuuze",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			endpoint, image := parseImageName(test.name)
			assert.Equal(t, test.endpoint, endpoint)
			assert.Equal(t, test.image, image)
		})
	}
}
