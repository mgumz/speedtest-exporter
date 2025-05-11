package job

import (
	"errors"
	"log/slog"

	"github.com/robfig/cron/v3"

	"github.com/mgumz/speedtest-exporter/internal/pkg/timeshift"
)

func (jobs Jobs) ReSchedule(scheduler *cron.Cron, collector *Collector) error {

	scheduler.Stop()

	// step 1: clean out currently scheduled jobs
	entries := []cron.EntryID{}
	for _, entry := range scheduler.Entries() {
		entries = append(entries, entry.ID)
	}
	for _, entry := range entries {
		scheduler.Remove(entry)
	}

	// step 2: unregister previous jobs
	for i := range jobs {
		collector.RemoveJob(jobs[i].Label)
	}

	// step 3: launch current set of jobs and register them
	// at the collector
	err := error(nil)
	for _, j := range jobs {
		if !collector.AddJob(j.JobMeta) {
			slog.Error("unable to add job to collector",
				"job.label", j.Label)
			if err != nil {
				err = errors.New("collector error")
			}
			continue
		}
		j.UpdateFn = func(meta JobMeta) bool { return collector.UpdateJob(meta) }

		s, err2 := timeshift.NewSchedule(j.Timeshift.Mode, j.scheduler.spec, j.Timeshift.Spec)
		if err2 != nil {
			slog.Error("unable to add job to scheduler",
				"job.label", j.Label,
				"error", err2)
			if err != nil {
				err = errors.New("schedule error")
			}
			continue
		}

		j.scheduler.entryID = scheduler.Schedule(s, j)
		j.scheduler.instance = scheduler

	}

	// step 4: restart jobs if needed
	if collector.NumberJobs() > 0 {
		slog.Info("restart scheduler", "status", "started")
		scheduler.Start()
		EntriesToLog(scheduler)
		slog.Info("restart scheduler", "status", "done")
	}

	return err
}
