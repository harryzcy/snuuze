package job

import (
	"fmt"
	"os"

	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/runner"
	"github.com/harryzcy/snuuze/types"
)

func checkUpdates(state *State) {
	state.Lock()
	repos := state.Repos
	state.Unlock()

	for _, repo := range repos {
		dependencies, err := runner.GetDependencyForRepo(repo.URL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get dependency for repo", repo.URL, ":", err)
			continue
		}
		state.Lock()
		updateDependencyIndex(state, repo, dependencies)
		state.Unlock()

		err = runner.RunForRepo(repo.URL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to run for repo", repo.URL, ":", err)
			continue
		}
	}
}

// updateDependencyIndex updates the dependency index in the state.
// It requires the state to be locked.
func updateDependencyIndex(state *State, repo platform.Repo, dependencies map[types.PackageManager][]*types.Dependency) {
	previousDependencies := flattenDependencies(state.RepoDependencies[repo.URL])
	currentDependencies := flattenDependencies(dependencies)
	currentHashed := make(map[string]bool)

	for _, dep := range currentDependencies {
		if _, ok := state.ReverseDependencyIndex[dep.Hash()]; !ok {
			state.ReverseDependencyIndex[dep.Hash()] = struct {
				Dependency *types.Dependency
				Repos      []platform.Repo
			}{
				Dependency: dep,
				Repos:      make([]platform.Repo, 0),
			}
		}
		if !containRepo(state.ReverseDependencyIndex[dep.Hash()].Repos, repo) {
			state.ReverseDependencyIndex[dep.Hash()] = struct {
				Dependency *types.Dependency
				Repos      []platform.Repo
			}{
				Dependency: state.ReverseDependencyIndex[dep.Hash()].Dependency,
				Repos:      append(state.ReverseDependencyIndex[dep.Hash()].Repos, repo),
			}
		}

		currentHashed[dep.Hash()] = true
	}

	for _, dep := range previousDependencies {
		if _, ok := currentHashed[dep.Hash()]; !ok {
			updatedRepos := removeRepo(state.ReverseDependencyIndex[dep.Hash()].Repos, repo)
			if len(updatedRepos) == 0 {
				delete(state.ReverseDependencyIndex, dep.Hash())
			} else {
				state.ReverseDependencyIndex[dep.Hash()] = struct {
					Dependency *types.Dependency
					Repos      []platform.Repo
				}{
					Dependency: state.ReverseDependencyIndex[dep.Hash()].Dependency,
					Repos:      updatedRepos,
				}
			}
		}
	}

	state.RepoDependencies[repo.URL] = dependencies
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
		result = append(result, deps...)
	}
	return result
}
