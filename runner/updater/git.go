package updater

import (
	"fmt"
	"os"
	"strings"

	"github.com/harryzcy/snuuze/command"
	"github.com/harryzcy/snuuze/platform"
)

const (
	DEFAULT_DEFAULT_BRANCH = "main"
)

func getDefaultBranch(repoDir string) string {
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

func commitChanges(repoDir, branchName, message string) error {
	_, err := command.RunCommand(command.CommandInputs{
		Command: []string{"git", "-C", repoDir, "add", "."},
	})
	if err != nil {
		return err
	}

	_, err = command.RunCommand(command.CommandInputs{
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
	_, err := command.RunCommand(command.CommandInputs{
		Command: []string{"git", "-C", repoDir, "push", "origin", branchName, "--force"},
	})
	if err != nil {
		return err
	}

	return nil
}

func createPullRequest(gitURL string, info *commitInfo, base string) error {
	repoInfo, err := parseGitURL(gitURL)
	if err != nil {
		return err
	}

	client, err := platform.NewClient(platform.NewClientOptions{
		URL: gitURL,
	})
	if err != nil {
		return err
	}

	err = client.CreatePullRequest(&platform.CreatePullRequestInput{
		Title: info.message,
		Body:  "",
		Base:  base,
		Head:  info.branchName,
		Owner: repoInfo.owner,
		Repo:  repoInfo.repo,
	})
	if err != nil {
		return err
	}

	return nil
}

type gitInfo struct {
	owner string
	repo  string
}

func parseGitURL(gitURL string) (*gitInfo, error) {
	url := strings.TrimPrefix(gitURL, "git@")
	url = strings.TrimPrefix(url, "https://")
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid git url: %s", gitURL)
	}
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]
	repo = strings.TrimSuffix(repo, ".git")

	return &gitInfo{
		owner: owner,
		repo:  repo,
	}, nil
}
