package main

import "fmt"

func usage() {

	const usage string = `Usage: speedtest-exporter [FLAGS] -- [SPEEDTEST-GO-FLAGS]

FLAGS:
-bind       <bind-address>
            bind address (default ":8080")
-h
            show help
-jobs       <path-to-jobsfile>
            file describing multiple speedtest-jobs. syntax is given below.
-label      <job-label>
            use <job-label> in prometheus-metrics (default: "speedtest-exporter-cli")
-log-level  <log-level>
            either of "debug", "info" (default), "warn", "error"
-speedtest  <path-to-binary>
            path to speedtest-go binary (default: "speedtest-go")
-schedule   <schedule>
            schedule at which often speedtest-go is launched (default: "@every 24h")
            examples:
               @every <dur>  - example "@every 24h"
               @hourly       - run once per hour
               10 * * * *    - execute 10 minutes after the full hour
            see https://en.wikipedia.org/wiki/Cron
-timeshift  <timeshift> 
            timeshift around the point in time when -schedule would trigger otherwise
            (default: "" - no timeshift)
-watch-jobs <schedule>
            periodically watch the file defined via -jobs (default: "")
            if it has changed stop previously running speedtest-jobs and apply
            all jobs defined in -jobs.
-show-license
            show license
-show-version
            show version

SPEEDTEST-GO-FLAGS:
see "speedtest-go" for valid flags.

Examples:

$> speedtest-exporter -- -s 00001
# probe every minute the speedtest-net server "00001"

Example Job File:

    # comments are ignored
    job1 -- @every 24h ±10m -- -s 00001
    job2 -- @daily ±15m -- -s 00002`

	fmt.Println(usage)
}
