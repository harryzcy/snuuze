package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostingConfig(t *testing.T) {
	if strings.HasSuffix(basepath, "snuuze/config") {
		ConfigFile = filepath.Join(basepath, "..", "testdata", "config.yaml")
	} else {
		ConfigFile = filepath.Join(basepath, "testdata", "config.yaml")
	}

	err := LoadHostingConfig()
	assert.NoError(t, err)

	config := GetHostingConfig()
	assert.Equal(t, int64(12345), config.GitHub.AppID)

	// test Env override
	os.Setenv("SNUUZE_GITHUB_APP_ID", "54321")

	err = LoadHostingConfig()
	assert.NoError(t, err)
	config = GetHostingConfig()
	assert.Equal(t, int64(54321), config.GitHub.AppID)
	assert.Len(t, config.Gitea, 1)
	assert.Equal(t, "https://gitea.com", config.Gitea[0].Host)
}

func TestToEnvName(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", "A"},
		{"tempDir", "TEMP_DIR"},
		{"data.tempDir", "DATA_TEMP_DIR"},
		{"github.authType", "GITHUB_AUTH_TYPE"},
	}

	for _, test := range tests {
		assert.Equal(t, test.out, toEnvName(test.in))
	}
}
