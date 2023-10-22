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
		updateDependencyIndex(state, repo, dependencies)

		err = runner.RunForRepo(repo.URL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to run for repo", repo.URL, ":", err)
			continue
		}
	}
}

func updateDependencyIndex(state *State, repo platform.Repo, dependencies map[types.PackageManager][]*types.Dependency) {
	state.RepoDependencies[repo] = dependencies

	previousDependencies := flattenDependencies(state.RepoDependencies[repo])
	currentDependencies := flattenDependencies(dependencies)
	currentHashed := make(map[string]bool)

	for _, dep := range currentDependencies {
		if !containRepo(state.ReverseDependencyIndex[dep], repo) {
			state.ReverseDependencyIndex[dep] = append(state.ReverseDependencyIndex[dep], repo)
		}

		currentHashed[dep.Hash()] = true
	}

	for _, dep := range previousDependencies {
		if _, ok := currentHashed[dep.Hash()]; !ok {
			state.ReverseDependencyIndex[dep] = removeRepo(state.ReverseDependencyIndex[dep], repo)
		}
	}
}

func containRepo(repos []platform.Repo, repo platform.Repo) bool {
	for _, r := range repos {
		if r.URL == repo.URL {
			return true
		}
	}
	return false
}

func removeRepo(repos []platform.Repo, repo platform.Repo) []platform.Repo {
	for i, r := range repos {
		if r.URL == repo.URL {
			return append(repos[:i], repos[i+1:]...)
		}
	}
	return repos
}

func flattenDependencies(dependencies map[types.PackageManager][]*types.Dependency) []*types.Dependency {
	result := make([]*types.Dependency, 0)
	for _, deps := range dependencies {
		for _, dep := range deps {
			result = append(result, dep)
		}
	}
	return result
}
