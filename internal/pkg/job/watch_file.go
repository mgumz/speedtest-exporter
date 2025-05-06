package job

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func WatchJobsFile(name, speedtestBin, watchSchedule string, collector *Collector) {

	watcher := &jobFileWatch{
		name:         name,
		speedtestBin: speedtestBin,
		scheduler: cron.New(
			cron.WithLocation(time.UTC),
			cron.WithChain(
				cron.SkipIfStillRunning(cron.DiscardLogger),
			),
		),
		collector: collector,
	}

	watcher.Run()

	// check `name` according to `watchSchedule`
	scheduler := cron.New(
		cron.WithLocation(time.UTC),
	)

	if _, err := scheduler.AddJob(watchSchedule, watcher); err != nil {
		slog.Error("unable to launch watch-jobs scheduler", "error", err)
		os.Exit(1)
	}
	scheduler.Start()
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

	slog.Debug("starting to parse jobFile", "fileName", jw.name)
	jobs, chksum, err := ParseJobFile(jw.name, jw.speedtestBin)
	if err != nil {
		slog.Warn("parsing jobFile failed", "fileName", jw.name, "error", err)
		return
	}
	slog.Debug("done parsing jobFile",
		"fileName", jw.name,
		"numberJobs", len(jobs),
		"previousSha256", fmt.Sprintf("%x", jw.chksum),
		"currentSha256", fmt.Sprintf("%x", chksum),
	)

	if bytes.Equal(jw.chksum, chksum) {
		slog.Debug("watched jobFile is unchanged", "fileName", jw.name)
		return
	}

	slog.Info("watched jobFile has changed, scheduling jobs",
		"fileName", jw.name,
		"numberJobs", len(jobs),
	)

	jobs.ReSchedule(jw.scheduler, jw.collector)

	jw.jobs = jobs
	jw.chksum = chksum
}
