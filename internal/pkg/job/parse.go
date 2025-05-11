package job

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/shlex"

	"github.com/mgumz/speedtest-exporter/internal/pkg/timeshift"
)

// JobFile definition
//
// # comments, ignore everything after #
// ^space*$ - empty lines
// <label> -- <schedule> -- <speedtest-go-flags>

func ParseJobs(r io.Reader, speedtest string) (Jobs, error) {

	var err error
	var jobs = Jobs{}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	n := 0
	for scanner.Scan() {
		line := scanner.Text()
		n++
		job, err2 := parseJobLine(line, n, speedtest)
		if err2 != nil {
			err = err2
			break
		}
		if job != nil {
			jobs = append(jobs, job)
		}
	}

	if err == nil {
		err = scanner.Err()
	}
	if err != nil {
		return Jobs{}, err
	}

	return jobs, nil
}

func parseJobLine(line string, lnr int, speedtest string) (*Job, error) {

	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil, nil
	}

	if strings.HasPrefix(line, "#") {
		return nil, nil
	}

	parts := strings.SplitN(line, " -- ", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid jobLine %d: expect '<label> -- <schedule> -- <speedtest-go-flags>'", lnr)
	}

	label, _ := parseLabel(strings.TrimSpace(parts[0]))
	schedule, tmode, tspec, err := parseSchedule(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid jobline %d: invalid schedule offset in %q", lnr, schedule)
	}
	speedtestArgs, _ := parseSpeedtestArgs(strings.TrimSpace(parts[2]))

	job := NewJob(speedtest, speedtestArgs, schedule, tmode, tspec)
	job.Label = label

	return job, nil
}

func parseLabel(l string) (string, error) {
	return l, nil
}

// schedule s is either
// * "*/5 * * * * " - a regular cron expression
// * @every 5m - a constant delay expression
// optional, there is a random-delay at the end of the
// normal cron/delay expression.
func parseSchedule(s string) (string, timeshift.Mode, string, error) {

	tsMode := timeshift.None
	tsMarker := ""

	if i := strings.IndexAny(s, "~"); i > 0 {
		switch {
		// NOTE: disabled for now, see timeshift.RandomDeviationScheduler for
		// the yet-to-be-solved issues
		//case strings.HasPrefix(s[i:], "±"):
		//	tsMode = timeshift.RandomDeviation
		//	tsMarker = "±"
		case strings.HasPrefix(s[i:], "~"):
			tsMode = timeshift.RandomDelay
			tsMarker = "~"
		}
	}

	if tsMode == timeshift.None {
		return strings.TrimSpace(s), tsMode, "", nil
	}

	// we don't check for "ok": we _know_ the marker is
	// in the string, the above code ensures that
	schedule, tsSpec, _ := strings.Cut(s, tsMarker)
	schedule = strings.TrimSpace(schedule)
	tsSpec = strings.TrimSpace(tsSpec)
	d, err := time.ParseDuration(tsSpec)
	if err != nil {
		return s, tsMode, "", fmt.Errorf("timeshift deviation %q, format error %s", tsSpec, err)
	}
	if d < 0 {
		return s, tsMode, "", fmt.Errorf("timeshift deviation must not be negative %q", tsSpec)
	}
	return schedule, tsMode, tsSpec, nil
}

func parseSpeedtestArgs(s string) ([]string, error) {
	args, err := shlex.Split(s)
	return args, err
}
