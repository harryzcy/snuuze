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
		dependencies, err := runner.GetDependencyForRepo(repo.URL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get dependency for repo", repo.URL, ":", err)
			continue
		}
		state.RepoDependencies[repo] = dependencies
		for manager, deps := range dependencies {
			for _, dep := range deps {
				state.ReverseIndex[dep] = struct {
					Repo    platform.Repo
					Manager types.PackageManager
				}{
					Repo:    repo,
					Manager: manager,
				}
			}
		}

		err = runner.RunForRepo(repo.URL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to run for repo", repo.URL, ":", err)
			continue
		}
	}
}
