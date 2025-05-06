package job

import (
	"log/slog"
	"time"
)

// cron.v3 interface
func (job *Job) Run() {

	slog.Info("launching", "job", job.Label, "cmdline", job.cmdLine)
	if err := job.Launch(); err != nil {
		slog.Warn("failed launch", "job", job.Label, "error", err)
		return
	}

	logAttrs := []any{
		slog.String("job", job.Label),
		slog.Duration("duration", job.Duration),
	}
	if job.Result.ErrorMsg != "" {
		logAttrs = append(logAttrs, slog.String("error", job.Result.ErrorMsg))
	}

	slog.Info("job done", logAttrs...)

	if job.scheduler.instance != nil && job.scheduler.entryID != 0 {
		entry := job.scheduler.instance.Entry(job.scheduler.entryID)
		nextRun := time.Until(entry.Next).Truncate(time.Second)
		slog.Info("next launch", "job", job.Label, "nextRun", nextRun)
	}
}
