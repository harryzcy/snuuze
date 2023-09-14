package updater

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	if strings.HasSuffix(basepath, "snuuze/updater") {
		_ = os.Chdir("../")
	}

	_ = os.Chdir("./testdata")

	m.Run()
}

func TestNewCache(t *testing.T) {
	cache := NewCache()
	assert.NotNil(t, cache)
}

func TestCache_Get(t *testing.T) {
	cache := NewCache()
	data, err := cache.Get("cache_test")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(data))
}

func TestCache_Set(t *testing.T) {
	cache := NewCache()
	cache.Set("cache_test", []byte("test"))
	data, err := cache.Get("cache_test")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(data))
}

func TestCache_Read(t *testing.T) {
	cache := NewCache()
	data, err := cache.Read("cache_test")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(data))
}

func TestCache_Commit(t *testing.T) {
	cache := NewCache()
	cache.Set("cache_test", []byte("test"))
	err := cache.Commit()
	assert.NoError(t, err)
}
