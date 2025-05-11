package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/mgumz/speedtest-exporter/internal/pkg/job"
)

func handleSignals(scheduler *cron.Cron) {

	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan, syscall.SIGUSR1)

	for {
		select {
		case sig := <-sChan:
			switch sig {
			case syscall.SIGUSR1:
				job.EntriesToLog(scheduler)
			}
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}
