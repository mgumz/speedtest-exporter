// *speedtest-exporter* periodically executes *speedtest-go* to a given host and
// provides the measured results as prometheus metrics.

package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/mgumz/speedtest-exporter/internal/pkg/job"
	"github.com/mgumz/speedtest-exporter/internal/pkg/timeshift"
)

func main() {

	stef := newFlags()
	flag.Usage = usage
	flag.Parse()

	if stef.doPrintVersion {
		printVersion()
		return
	}
	if stef.doPrintUsage {
		flag.Usage()
		return
	}
	if stef.doPrintLicense {
		printLicense()
		return
	}

	// logging
	logLevel, err := logLevel(stef.logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// ok, lets go:
	jobsCollector := job.NewCollector()
	jobs := job.Jobs{}

	if len(flag.Args()) > 0 {
		tmode := timeshift.None
		if stef.timeShift != "" {
			tmode = timeshift.RandomDeviation
		}
		j := job.NewJob(stef.speedtestBin, flag.Args(), stef.schedule, tmode, stef.timeShift)
		j.Label = stef.jobLabel
		jobs = append(jobs, j)
	}

	jobsScheduler := cron.New(
		cron.WithLocation(time.UTC),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DiscardLogger),
		),
	)

	jobsArePossible := !jobs.Empty()

	if stef.jobFile != "" {
		if stef.doWatchJobsFile != "" {

			wi := job.WatchJobsFileInfo{
				Name:            stef.jobFile,
				SpeedtestBinary: stef.speedtestBin,
				WatchSchedule:   stef.doWatchJobsFile,
			}
			job.WatchJobsFile(&wi, jobsScheduler, jobsCollector)
			jobsArePossible = true

		} else {

			jobsFromFile, _, err := job.ParseJobFile(stef.jobFile, stef.speedtestBin)
			if err != nil {
				slog.Error("parsing jobs file failed", "fileName", stef.jobFile, "error", err)
				os.Exit(1)
			}
			if !jobsFromFile.Empty() {
				jobs = append(jobs, jobsFromFile...)
				jobsArePossible = true
			}
		}
	}

	if !jobsArePossible {
		slog.Error("no speedtest jobs defined - provide at least one via -file or via arguments")
		os.Exit(1)
	}

	if !jobs.Empty() { // jobs are only filled when not -watch-jobs is active
		if err := jobs.ReSchedule(jobsScheduler, jobsCollector); err != nil {
			slog.Error("", "error", err)
			os.Exit(1)
		}
	}

	go handleSignals(jobsScheduler)

	http.Handle("/metrics", jobsCollector)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/", handleRoot)

	slog.Info("serving ...",
		"http.path", "/metrics",
		"http.bindAddr", stef.bindAddr)

	http.ListenAndServe(stef.bindAddr, nil)
}
