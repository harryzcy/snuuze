package updater

import (
	"fmt"
	"strings"

	"github.com/harryzcy/snuuze/platform"
)

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
