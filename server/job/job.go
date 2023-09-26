package job

import "github.com/harryzcy/snuuze/runner"

func checkUpdates(state *State) error {
	for _, repo := range state.Repos {
		err := runner.RunForRepo(repo.URL)
		if err != nil {
			return err
		}
	}
	return nil
}
