package githubactions

import (
	"strings"

	"github.com/harryzcy/snuuze/manager/common"
	"github.com/harryzcy/snuuze/types"
)

type GitHubActionsManager struct{}

func New() common.Manager {
	return &GitHubActionsManager{}
}

func (m *GitHubActionsManager) Name() types.PackageManager {
	return types.PackageManagerGitHubActions
}

func (m *GitHubActionsManager) Match(path string) bool {
	if !strings.HasPrefix(path, ".github/workflows") {
		return false
	}
	return strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")
}
