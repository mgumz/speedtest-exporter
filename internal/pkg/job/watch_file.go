package job

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

type WatchJobsFileInfo struct {
	Name            string
	SpeedtestBinary string
	WatchSchedule   string
}

func WatchJobsFile(wi *WatchJobsFileInfo, scheduler *cron.Cron, collector *Collector) {

	watcher := &jobFileWatch{
		name:         wi.Name,
		speedtestBin: wi.SpeedtestBinary,
		scheduler:    scheduler,
		collector:    collector,
	}

	watcher.Run()

	wscheduler := cron.New(
		cron.WithLocation(time.UTC),
	)

	if _, err := wscheduler.AddJob(wi.WatchSchedule, watcher); err != nil {
		slog.Error("unable to launch watch-jobs scheduler",
			"jobs.fileName", wi.Name,
			"error", err)
		os.Exit(1)
	}
	wscheduler.Start()
}

type jobFileWatch struct {
	name         string
	speedtestBin string
	scheduler    *cron.Cron
	chksum       []byte
	jobs         Jobs
	collector    *Collector
}

func (jw *jobFileWatch) Run() {

	slog.Debug("starting to parse jobFile",
		"jobs.fileName", jw.name)

	jobs, chksum, err := ParseJobFile(jw.name, jw.speedtestBin)
	if err != nil {
		slog.Warn("parsing jobFile failed",
			"jobs.fileName", jw.name,
			"status", "failed",
			"error", err)
		jw.collector.IncMetricJobFileFailed()
		return
	}
	slog.Debug("done parsing jobFile",
		"jobs.fileName", jw.name,
		"jobs.number", len(jobs),
		"jobs.prevSha256", fmt.Sprintf("%x", jw.chksum),
		"jobs.sha256", fmt.Sprintf("%x", chksum),
	)

	if bytes.Equal(jw.chksum, chksum) {
		slog.Debug("watched jobFile is unchanged",
			"jobs.fileName", jw.name,
			"status", "unchanged")
		jw.collector.IncMetricJobFileUnchanged()
		return
	}

	slog.Info("watched jobFile has changed, scheduling jobs",
		"jobs.fileName", jw.name,
		"jobs.number", len(jobs),
		"status", "changed",
	)
	jw.collector.IncMetricJobFileChanged()

	jobs.ReSchedule(jw.scheduler, jw.collector)

	jw.jobs = jobs
	jw.chksum = chksum
}
