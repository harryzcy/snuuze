package config

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostingConfig(t *testing.T) {
	if strings.HasSuffix(basepath, "snuuze/config") {
		CONFIG_FILE = filepath.Join(basepath, "..", "testdata", "config.yaml")
	} else {
		CONFIG_FILE = filepath.Join(basepath, "testdata", "config.yaml")
	}

	err := LoadHostingConfig()
	assert.NoError(t, err)
}
