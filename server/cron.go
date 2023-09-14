package server

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func startCron() error {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(5).Seconds().Do(func() {
		fmt.Println("Running every 5 seconds")
	})
	if err != nil {
		return err
	}

	s.StartAsync()
	return nil
}
