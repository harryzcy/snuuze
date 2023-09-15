package git

import "github.com/harryzcy/snuuze/command"

func PushBranch(repoDir, branchName string) error {
	_, err := command.RunCommand(command.CommandInputs{
		Command: []string{"git", "-C", repoDir, "push", "origin", branchName, "--force"},
	})
	if err != nil {
		return err
	}

	return nil
}
