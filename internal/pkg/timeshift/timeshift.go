package timeshift

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Mode int

const (
	None Mode = iota
	RandomDeviation
	RandomDelay
)

func NewSchedule(mode Mode, spec, timeshift string) (cron.Schedule, error) {

	sched, err := cron.ParseStandard(spec)
	if err != nil {
		return nil, err
	}

	switch mode {
	case None:
		return sched, nil

	case RandomDeviation:
		deviation, err := time.ParseDuration(timeshift)
		if err != nil {
			return nil, fmt.Errorf("unparseable timeshift deviation %q", timeshift)
		}
		return NewRandomDeviationSchedule(sched, deviation)

	case RandomDelay:
		delay, err := time.ParseDuration(timeshift)
		if err != nil {
			return nil, fmt.Errorf("unparseable timeshift delay %q", timeshift)
		}
		return NewRandomDelaySchedule(sched, delay)
	}

	return nil, fmt.Errorf("unknown timeshift.Mode %d", mode)
}
