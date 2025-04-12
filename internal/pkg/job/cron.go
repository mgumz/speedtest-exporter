package job

import (
	"fmt"
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

	errMsg := ""
	if job.Result.ErrorMsg != "" {
		errMsg = fmt.Sprintf("(err: %q)", job.Result.ErrorMsg)
	}

	slog.Info("done", "job", job.Label, "duration", job.Duration, "errorMsg", errMsg)

	if job.scheduler.instance != nil && job.scheduler.entryID != 0 {
		entry := job.scheduler.instance.Entry(job.scheduler.entryID)
		nextRun := time.Until(entry.Next).Truncate(time.Second)
		slog.Info("next launch", "job", job.Label, "nextRun", nextRun)
	}
}
