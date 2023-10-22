package job

import (
	"fmt"

	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/types"
)

type State struct {
	Repos                  []platform.Repo
	RepoDependencies       map[platform.Repo]map[types.PackageManager][]*types.Dependency
	ReverseDependencyIndex map[string]struct {
		Dependency *types.Dependency
		Repos      []platform.Repo
	}
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
		Repos:            repos,
		RepoDependencies: make(map[platform.Repo]map[types.PackageManager][]*types.Dependency),
		ReverseDependencyIndex: make(map[string]struct {
			Dependency *types.Dependency
			Repos      []platform.Repo
		}),
	}, nil
}
