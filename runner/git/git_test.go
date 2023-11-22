package git

import (
	"context"
	"testing"

	"github.com/harryzcy/snuuze/platform"
	"github.com/stretchr/testify/assert"
)

type TestClient struct {
	token string
}

var _ platform.Client = &TestClient{}

func (c *TestClient) Token(ctx context.Context) (string, error) {
	return c.token, nil
}

func (c *TestClient) ListRepos() ([]platform.Repo, error) {
	return nil, nil
}

func (c *TestClient) ListTags(params *platform.ListTagsInput) ([]string, error) {
	return nil, nil
}

func (c *TestClient) CreatePullRequest(input *platform.CreatePullRequestInput) error {
	return nil
}

func TestGetGitURLWithToken(t *testing.T) {
	client := &TestClient{
		token: "test-token",
	}

	url, err := getGitURLWithToken(client, "https://github.com/owner/repo.git")
	assert.NoError(t, err)
	assert.Equal(t, "https://test-token@github.com/owner/repo.git", url)
}
