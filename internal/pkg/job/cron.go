package job

import (
	"log/slog"
	"time"
)

// cron.v3 interface
func (job *Job) Run() {

	slog.Info("launching",
		"job.label", job.Label,
		"job.cmdline", job.cmdLine)

	if err := job.Launch(); err != nil {
		slog.Warn("failed launch",
			"job.label", job.Label,
			"error", err)
		return
	}

	logAttrs := []any{
		slog.String("job.label", job.Label),
		slog.Duration("job.duration", job.Duration.Truncate(time.Second)),
	}
	if job.Result.ErrorMsg != "" {
		logAttrs = append(logAttrs, slog.String("error", job.Result.ErrorMsg))
	}

	slog.Info("finished", logAttrs...)

	// The cron instance, just after .Run() was called, calculated
	// the .Next() run. To give a hint to the user, we log the duration
	// until .Next()
	if job.scheduler.instance != nil && job.scheduler.entryID != 0 {
		entry := job.scheduler.instance.Entry(job.scheduler.entryID)
		nextRun := time.Until(entry.Next).Truncate(time.Second)
		slog.Info("next launch",
			"time.when", entry.Next,
			"job.label", job.Label,
			"time.in", nextRun)
	}
}
