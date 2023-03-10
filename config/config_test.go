package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func TestMain(m *testing.M) {
	if strings.HasSuffix(basepath, "snuuze/config") {
		_ = os.Chdir("../")
	}

	_ = os.Chdir("./testdata")
	os.Exit(m.Run())
}

func TestConfig(t *testing.T) {
	err := Load()
	if err != nil {
		t.Error(err)
	}

	config := GetConfig()
	assert.Equal(t, "1", config.Version)

	rules := GetRules()
	assert.Len(t, rules, 1)
}
