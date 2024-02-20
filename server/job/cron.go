package job

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/pkg/errors"
)

func StartCron(state *State) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	_, err = s.NewJob(gocron.DurationJob(time.Hour), gocron.NewTask(func() {
		checkUpdates(state)
	}))
	if err != nil {
		return nil, err
	}

	s.Start()
	return s, nil
}

func StopCron(s gocron.Scheduler) error {
	err := s.Shutdown()
	return errors.Wrap(err, "failed to stop cron")
}
