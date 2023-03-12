package updater

import (
	"strings"

	"github.com/harryzcy/snuuze/cmdutil"
)

const (
	DEFAULT_DEFAULT_BRANCH = "main"
)

func getDefaultBranch(repoDir string) string {
	output, err := cmdutil.RunCommand(cmdutil.CommandInputs{
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

func pushBranch(repoDir, branchName string) error {
	_, err := cmdutil.RunCommand(cmdutil.CommandInputs{
		Command: []string{"git", "-C", repoDir, "push", "origin", branchName},
	})
	if err != nil {
		return err
	}

	return nil
}
