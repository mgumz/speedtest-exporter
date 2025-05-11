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
	collector := job.NewCollector()
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

	jobsAvailable := !jobs.Empty()

	if stef.jobFile != "" {
		if stef.doWatchJobsFile != "" {
			slog.Info("watching -jobs-file", "fileName", stef.jobFile, "schedule", stef.doWatchJobsFile)
			job.WatchJobsFile(stef.jobFile, stef.speedtestBin, stef.doWatchJobsFile, collector)
			jobsAvailable = true
		} else {
			jobsFromFile, _, err := job.ParseJobFile(stef.jobFile, stef.speedtestBin)
			if err != nil {
				slog.Error("parsing jobs file failed", "fileName", stef.jobFile, "error", err)
				os.Exit(1)
			}
			if !jobsFromFile.Empty() {
				jobs = append(jobs, jobsFromFile...)
				jobsAvailable = true
			}
		}
	}

	if !jobsAvailable {
		slog.Error("no speedtest jobs defined - provide at least one via -file or via arguments")
		os.Exit(1)
	}

	scheduler := cron.New(
		cron.WithLocation(time.UTC),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DiscardLogger),
		),
	)

	if err := jobs.ReSchedule(scheduler, collector); err != nil {
		slog.Error("", "error", err)
		os.Exit(1)
	}

	http.Handle("/metrics", collector)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/", handleRoot)

	slog.Info("serving ...", "path", "/metrics", "bindAddr", stef.bindAddr)
	log.Fatal(http.ListenAndServe(stef.bindAddr, nil))
}
