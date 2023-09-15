package git

import (
	"strings"

	"github.com/harryzcy/snuuze/command"
)

const (
	DEFAULT_DEFAULT_BRANCH = "main"
)

func GetDefaultBranch(repoDir string) string {
	output, err := command.RunCommand(command.CommandInputs{
		Command: []string{"git", "-C", repoDir, "rev-parse", "--abbrev-ref", "origin/HEAD"},
	})
	if err != nil {
		return DEFAULT_DEFAULT_BRANCH
	}

	branch := output.Stdout.String()
	if branch == "" {
		return DEFAULT_DEFAULT_BRANCH
	}

	branch = strings.TrimSpace(branch)
	branch = strings.TrimPrefix(branch, "origin/")
	return branch
}
