package timeshift

import (
	"time"

	"github.com/robfig/cron/v3"
)

type RandomDeviationSchedule struct {
	baseSchedule cron.Schedule
	deviation    time.Duration
}

func NewRandomDeviationSchedule(base cron.Schedule, deviation time.Duration) (*RandomDeviationSchedule, error) {
	sched := RandomDeviationSchedule{
		baseSchedule: base,
		deviation:    deviation,
	}
	return &sched, nil
}

func (rds *RandomDeviationSchedule) Next(t time.Time) time.Time {

	nt := rds.baseSchedule.Next(t)

	if rds.deviation == time.Duration(0) {
		return nt
	}

	off := 2 * randDurationMax(rds.deviation)
	nt = nt.Add(-rds.deviation + off)

	return nt
}
