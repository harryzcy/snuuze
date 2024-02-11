package job

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

func StartCron(state *State) (gocron.Scheduler, error) {
	s, _ := gocron.NewScheduler()

	_, err := s.NewJob(gocron.DurationJob(time.Hour), gocron.NewTask(func() {
		checkUpdates(state)
	}))
	if err != nil {
		return nil, err
	}

	s.Start()
	return s, nil
}

func StopCron(s gocron.Scheduler) {
	_ = s.Shutdown()
}
