package job

import (
	"fmt"

	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/types"
)

type State struct {
	Repos []platform.Repo

	// RepoDependencies maps a repo URL to its dependencies.
	RepoDependencies map[string]map[types.PackageManager][]*types.Dependency

	// ReverseDependencyIndex maps a dependency hash to its dependency and the repos that depend on it.
	// The dependency hash is obtained by calling the Hash() method on the dependency.
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
		RepoDependencies: make(map[string]map[types.PackageManager][]*types.Dependency),
		ReverseDependencyIndex: make(map[string]struct {
			Dependency *types.Dependency
			Repos      []platform.Repo
		}),
	}, nil
}
