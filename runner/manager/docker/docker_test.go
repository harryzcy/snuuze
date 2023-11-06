package docker

import (
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestDockerManager_Parse(t *testing.T) {
	manager := New()
	match := types.Match{
		File: "Dockerfile",
	}
	data := []byte(`FROM alpine:3.18.0
RUN echo "hello world"
`)

	dependencies, err := manager.Parse(match, data)
	assert.NoError(t, err)
	assert.Len(t, dependencies, 1)
	assert.Equal(t, "alpine", dependencies[0].Name)
	assert.Equal(t, "3.18.0", dependencies[0].Version)
	assert.Equal(t, types.PackageManagerDocker, dependencies[0].PackageManager)
	assert.Equal(t, 1, dependencies[0].Position.Line)
	assert.Equal(t, map[string]interface{}{
		"versionType": "tag",
	}, dependencies[0].Extra)

	data = []byte(`FROM alpine@sha256:48d9183eb12a05c99bcc0bf44a003607b8e941e1d4f41f9ad12bdcc4b5672f86
RUN echo "hello world"
`)
	dependencies, err = manager.Parse(match, data)
	assert.NoError(t, err)
	assert.Len(t, dependencies, 1)
	assert.Equal(t, "sha256:48d9183eb12a05c99bcc0bf44a003607b8e941e1d4f41f9ad12bdcc4b5672f86", dependencies[0].Version)
	assert.Equal(t, map[string]interface{}{
		"versionType": "digest",
	}, dependencies[0].Extra)

	data = []byte(`FROM alpine
RUN echo "hello world"
`)
	dependencies, err = manager.Parse(match, data)
	assert.NoError(t, err)
	assert.Len(t, dependencies, 1)
	assert.Equal(t, "", dependencies[0].Version)
}

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
	assert.Contains(t, tags, "3.18.4")
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
			image:    "library/alpine",
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
