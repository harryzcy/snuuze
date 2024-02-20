package git

import (
	"strings"

	"github.com/harryzcy/snuuze/runner/command"
)

const (
	DefaultDefaultBranch = "main"
)

func GetDefaultBranch(repoDir string) string {
	output, err := command.RunCommand(command.CommandInputs{
		Command: []string{"git", "-C", repoDir, "rev-parse", "--abbrev-ref", "origin/HEAD"},
	})
	if err != nil {
		return DefaultDefaultBranch
	}

	branch := output.Stdout.String()
	if branch == "" {
		return DefaultDefaultBranch
	}

	branch = strings.TrimSpace(branch)
	branch = strings.TrimPrefix(branch, "origin/")
	return branch
}
