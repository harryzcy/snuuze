package job

import (
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
)

func StartCron(state *State) (*gocron.Scheduler, error) {
	s := gocron.NewScheduler(time.UTC)
	s.WaitForScheduleAll()

	_, err := s.Every(1).Hour().Do(func() {
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

func StopCron(s *gocron.Scheduler) {
	s.Stop()
}
