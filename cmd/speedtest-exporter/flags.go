package main

import "flag"

type steFlags struct {
	speedtestBin string
	jobLabel     string
	bindAddr     string
	jobFile      string
	schedule     string
	timeShift    string
	logLevel     string

	doWatchJobsFile string
	doPrintVersion  bool
	doPrintLicense  bool
	doPrintUsage    bool
}

func newFlags() *steFlags {

	ste := new(steFlags)

	flag.StringVar(&ste.speedtestBin, "speedtest", "speedtest-go", "path to `speedtest-go` binary")
	flag.StringVar(&ste.jobLabel, "label", "speedtest-exporter-cli", "job label")
	flag.StringVar(&ste.bindAddr, "bind", ":8080", "bind address")
	flag.StringVar(&ste.jobFile, "jobs", "", "file containing job definitions")
	flag.StringVar(&ste.schedule, "schedule", "@every 24h", "schedule at which often `speedtest-go` is launched")
	flag.StringVar(&ste.timeShift, "timeshift", "", "timeshift the -schedule a bit")
	flag.StringVar(&ste.doWatchJobsFile, "watch-jobs", "", "re-parse -jobs file to schedule")
	flag.BoolVar(&ste.doPrintVersion, "version", false, "show version")
	flag.BoolVar(&ste.doPrintLicense, "show-license", false, "show license")
	flag.BoolVar(&ste.doPrintUsage, "h", false, "show help")
	flag.StringVar(&ste.logLevel, "log-level", "info", "log level")

	return ste
}
