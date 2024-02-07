package platform

import "context"

type GitClient struct {
	server string
}

var _ Client = &GitClient{}

// NewGitClient creates a new NewGitClient for a git server.
func NewGitClient(server string) (Client, error) {
	client := &GitClient{
		server: server,
	}

	return client, nil
}

func (c *GitClient) Token(ctx context.Context) (string, error) {
	return "", ErrUnimplemented
}

func (c *GitClient) ListRepos() ([]Repo, error) {
	return nil, ErrUnimplemented
}

func (c *GitClient) ListTags(params *ListTagsInput) ([]string, error) {
	return nil, ErrUnimplemented
}

func (c *GitClient) CreatePullRequest(input *CreatePullRequestInput) error {
	return ErrUnimplemented
}
