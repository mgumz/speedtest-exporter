package job

import (
	"errors"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/mgumz/speedtest-exporter/internal/pkg/timeshift"
)

type Jobs []*Job

func (jobs Jobs) Count() int { return len(jobs) }

func (jobs Jobs) Empty() bool { return len(jobs) == 0 }

func (jobs Jobs) CollectedResults() int {
	n := 0
	for _, job := range jobs {
		if !job.Result.Empty() {
			n++
		}
	}
	return n
}

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
	n := 0
	for _, j := range jobs {
		if !collector.AddJob(j.JobMeta) {
			slog.Error("unable to add job to collector", "job", j.Label)
			if err != nil {
				err = errors.New("collector error")
			}
			continue
		}
		slog.Info("schedule job", "job", j.Label, "schedule", j.scheduler.spec, "timeshift", &j.Timeshift)
		j.UpdateFn = func(meta JobMeta) bool { return collector.UpdateJob(meta) }

		schedule, err2 := timeshift.NewSchedule(j.Timeshift.Mode, j.scheduler.spec, j.Timeshift.Spec)
		if err2 != nil {
			slog.Error("unable to add job to scheduler", "job", j.Label, "error", err2)
			if err != nil {
				err = errors.New("schedule error")
			}
			continue
		}

		j.scheduler.instance = scheduler
		j.scheduler.entryID = scheduler.Schedule(schedule, j)

		n++
	}

	if n > 0 {
		slog.Info("restart scheduler", "njobs", n)
		scheduler.Start()

		for _, entry := range scheduler.Entries() {
			j := entry.Job.(*Job)
			nextRun := time.Until(entry.Next).Truncate(time.Second)
			slog.Info("next run", "job", j.Label, "duration", nextRun)
		}
	}

	return err
}
