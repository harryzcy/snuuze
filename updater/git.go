package updater

import (
	"os"
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

func commitChanges(repoDir, branchName, message string) error {
	_, err := cmdutil.RunCommand(cmdutil.CommandInputs{
		Command: []string{"git", "-C", repoDir, "add", "."},
	})
	if err != nil {
		return err
	}

	_, err = cmdutil.RunCommand(cmdutil.CommandInputs{
		Command: []string{"git", "-C", repoDir, "commit", "-m", message},
		Env: map[string]string{
			"HOME": os.Getenv("HOME"), // required for git to find the user's config
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func pushBranch(repoDir, branchName string) error {
	_, err := cmdutil.RunCommand(cmdutil.CommandInputs{
		Command: []string{"git", "-C", repoDir, "push", "origin", branchName, "--force"},
	})
	if err != nil {
		return err
	}

	return nil
}
