package server

import (
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/harryzcy/snuuze/runner"
)

func startCron(state *State) (*gocron.Scheduler, error) {
	s := gocron.NewScheduler(time.UTC)
	s.WaitForScheduleAll()

	_, err := s.Every(5).Seconds().Do(func() {
		err := checkUpdates(state)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to check for updates:", err)
		}
	})
	if err != nil {
		return nil, err
	}

	s.StartAsync()
	return s, nil
}

func stopCron(s *gocron.Scheduler) {
	s.Stop()
}

func checkUpdates(state *State) error {
	for _, repo := range state.Repos {
		err := runner.RunForRepo(repo.URL)
		if err != nil {
			return err
		}
	}
	return nil
}
