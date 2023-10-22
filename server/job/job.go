package job

import (
	"fmt"
	"os"

	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/runner"
	"github.com/harryzcy/snuuze/types"
)

func checkUpdates(state *State) {
	for _, repo := range state.Repos {
		dependencies, err := getDependencyForRepo(repo)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get dependency for repo", repo.URL, ":", err)
			continue
		}
		state.Dependencies = append(state.Dependencies, dependencies...)

		err = runner.RunForRepo(repo.URL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to run for repo", repo.URL, ":", err)
			continue
		}
	}
}

func getDependencyForRepo(repo platform.Repo) ([]*types.Dependency, error) {
	dependencies, err := runner.GetDependencyForRepo(repo.URL)
	if err != nil {
		return nil, err
	}
	return dependencies, nil
}
