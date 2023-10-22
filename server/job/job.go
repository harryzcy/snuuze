package job

import (
	"fmt"
	"os"

	"github.com/harryzcy/snuuze/runner"
)

func checkUpdates(state *State) {
	for _, repo := range state.Repos {
		dependencies, err := runner.GetDependencyForRepo(repo.URL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get dependency for repo", repo.URL, ":", err)
			continue
		}
		state.RepoDependencies[repo] = dependencies

		err = runner.RunForRepo(repo.URL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to run for repo", repo.URL, ":", err)
			continue
		}
	}
}
