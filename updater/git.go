package updater

import (
	"strings"

	"github.com/harryzcy/snuuze/cmdutil"
)

const (
	DEFAULT_DEFAULT_BRANCH = "main"
)

func getDefaultBranch() string {
	output, err := cmdutil.RunCommand(cmdutil.CommandInputs{
		Command: []string{"git", "rev-parse", "--abbrev-ref", "origin/HEAD"},
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
