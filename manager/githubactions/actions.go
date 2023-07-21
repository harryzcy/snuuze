package githubactions

import (
	"strings"

	"github.com/harryzcy/snuuze/manager/common"
)

type GitHubActionsManager struct{}

func New() common.Manager {
	return &GitHubActionsManager{}
}

func (m *GitHubActionsManager) Name() string {
	return "github-actions"
}

func (m *GitHubActionsManager) Match(path string) bool {
	if !strings.HasPrefix(path, ".github/workflows") {
		return false
	}
	return strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")
}
