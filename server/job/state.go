package job

import (
	"fmt"

	"github.com/harryzcy/snuuze/platform"
)

type State struct {
	Repos []platform.Repo
}

// InitState loads the state for the server.
func InitState() (*State, error) {
	client, err := platform.NewClient(platform.NewClientOptions{
		Platform: platform.GitPlatformGitHub,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create platform client: %w", err)
	}

	repos, err := client.ListRepos()
	if err != nil {
		return nil, fmt.Errorf("failed to list repos: %w", err)
	}

	return &State{
		Repos: repos,
	}, nil
}
