package git

import "github.com/harryzcy/snuuze/runner/command"

func PushBranch(repoDir, branchName string) error {
	_, err := command.RunCommand(command.Inputs{
		Command: []string{"git", "-C", repoDir, "push", "origin", branchName, "--force"},
	})
	if err != nil {
		return err
	}

	return nil
}
