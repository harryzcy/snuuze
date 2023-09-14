package updater

import (
	"bytes"
	"fmt"

	"github.com/harryzcy/snuuze/types"
)

func upgradeGitHubActions(cache *Cache, info *types.UpgradeInfo) error {
	file, err := cache.Get(info.Dependency.File)
	if err != nil {
		return fmt.Errorf("failed to get file %s from cache: %s", info.Dependency.File, err)
	}

	originalText := fmt.Sprintf("%s@%s", info.Dependency.Name, info.Dependency.Version)
	newText := fmt.Sprintf("%s@%s", info.Dependency.Name, info.ToVersion)
	file = bytes.ReplaceAll(file, []byte(originalText), []byte(newText))
	cache.Set(info.Dependency.File, file)

	return nil
}
