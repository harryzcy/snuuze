package git

import (
	"os"

	"github.com/harryzcy/snuuze/runner/command"
)

func CommitChanges(repoDir, message string) error {
	_, err := command.RunCommand(command.Inputs{
		Command: []string{"git", "-C", repoDir, "add", "."},
	})
	if err != nil {
		return err
	}

	_, err = command.RunCommand(command.Inputs{
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
