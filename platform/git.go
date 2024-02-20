package platform

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GitClient struct {
	// server is the git server URL (e.g. https://git.example.com)
	server string
}

var _ Client = &GitClient{}

// NewGitClient creates a new NewGitClient for a git server.
func NewGitClient(server string) (Client, error) {
	if server == "" {
		return nil, ErrServerRequired
	}

	if strings.HasPrefix(server, "http://") {
		return nil, ErrNoInsecureServer
	}

	if !strings.HasPrefix(server, "https://") {
		return nil, ErrInvalidServerURL
	}

	client := &GitClient{
		server: server,
	}

	return client, nil
}

func (c *GitClient) Token(_ context.Context) (string, error) {
	return "", ErrUnimplemented
}

func (c *GitClient) ListRepos() ([]Repo, error) {
	return nil, ErrUnimplemented
}

func (c *GitClient) ListTags(params *ListTagsInput) ([]string, error) {
	s := memory.NewStorage()
	r, err := git.Clone(s, nil, &git.CloneOptions{
		URL: strings.Join([]string{c.server, params.Owner, params.Repo}, "/"),
	})
	if err != nil {
		return nil, err
	}

	// tags := make(map[plumbing.Hash]string)
	tags := make([]string, 0)

	iter, err := r.Tags()
	if err != nil {
		return nil, err
	}

	for {
		ref, err := iter.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		// tags[ref.Hash()] = ref.Name().Short()
		tags = append(tags, ref.Name().Short())

	}
	return tags, nil
}

func (c *GitClient) CreatePullRequest(_ *CreatePullRequestInput) error {
	return ErrUnimplemented
}
