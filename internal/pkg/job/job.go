package job

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/mgumz/speedtest-exporter/internal/pkg/speedtest"
	"github.com/mgumz/speedtest-exporter/internal/pkg/timeshift"
)

type JobMeta struct {
	Result    speedtest.Result
	Launched  time.Time
	Duration  time.Duration
	Timeshift tsMeta
	Label     string
	CmdLine   string

	Runs map[string]int64
}

func (jm *JobMeta) DataAvailable() bool { return len(jm.Runs) > 0 }

type Job struct {
	JobMeta

	scheduler struct {
		spec     string
		instance *cron.Cron
		entryID  cron.EntryID
	}

	speedtestBinary string
	args            []string
	cmdLine         string

	UpdateFn func(JobMeta) bool
}

func NewJob(speedtest string, args []string, schedule string, tmode timeshift.Mode, tspec string) *Job {
	extra := []string{
		"--json",
		"--unit", "decimal-bytes",
	}
	args = append(extra, args...)
	job := Job{
		args:            args,
		speedtestBinary: speedtest,
		cmdLine:         strings.Join(append([]string{speedtest}, args...), " "),
	}
	job.scheduler.spec = schedule
	job.JobMeta.Runs = map[string]int64{}
	job.JobMeta.Timeshift = tsMeta{tmode, tspec}
	job.JobMeta.CmdLine = job.cmdLine
	return &job
}

func (job *Job) Launch() error {

	// TODO: maybe use CommandContext to have an upper limit in the execution

	cmd := exec.Command(job.speedtestBinary, job.args...)

	// launch speedtest-go
	bufStdout, bufStderr := bytes.Buffer{}, bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = &bufStdout, &bufStderr
	launched := time.Now()
	cmd.Run()
	duration := time.Since(launched)

	errMsg := normalizeSpeedtestErrorMsg(bufStderr.String())
	if val, exists := job.Runs[errMsg]; exists {
		job.Runs[errMsg] = val + 1
	} else {
		job.Runs[errMsg] = 1
	}

	// decode the report
	result := speedtest.Result{}
	if err := result.Decode(&bufStdout); err != nil {
		result.ErrorMsg = "error-decoding-speedtest-json"
	} else {
		result.ErrorMsg = errMsg
	}

	// copy the report into the job
	job.JobMeta.Result = result
	job.JobMeta.Launched = launched
	job.JobMeta.Duration = duration

	if job.UpdateFn != nil {
		job.UpdateFn(job.JobMeta)
	}

	return nil
}

func normalizeSpeedtestErrorMsg(msg string) string {
	mf := func(r rune) rune {
		switch {
		case r == ' ' || r == ':':
			return r
		case r >= '0' && r <= '9':
			return r
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return 'a' + (r - 'A')
		}
		return '-'
	}
	return strings.Map(mf, strings.TrimSpace(msg))
}
