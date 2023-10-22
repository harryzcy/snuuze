package job

import (
	"time"

	"github.com/go-co-op/gocron"
)

func StartCron(state *State) (*gocron.Scheduler, error) {
	s := gocron.NewScheduler(time.UTC)
	s.WaitForScheduleAll()

	_, err := s.Every(1).Hour().Do(func() {
		checkUpdates(state)
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
