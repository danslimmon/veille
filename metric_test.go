package main

import (
	"testing"
	"time"
)

func parseTime(layout, value string) (t time.Time) {
	t, _ = time.Parse(layout, value)
	return
}

func TestCrunchMetrics_SingleOutage(t *testing.T) {
	t.Parallel()
	var timeFmt string
	var end time.Time
	var states []State
	var co CrunchOptions
	var metrics []*Metric
	var err error

	timeFmt = "2006-01-02 15:04:05"
	end = parseTime(timeFmt, "2016-01-01 00:00:00")
	states = []State{
		NewServiceState(
			parseTime(timeFmt, "2015-08-20 13:19:00"),
			"testhost-000", "test Service haha", "CRITICAL", "SOFT", "doesn't matter",
		),
		NewServiceState(
			parseTime(timeFmt, "2015-08-20 13:20:00"),
			"testhost-000", "test Service haha", "CRITICAL", "SOFT", "doesn't matter",
		),
		NewServiceState(
			parseTime(timeFmt, "2015-08-20 13:21:00"),
			"testhost-000", "test Service haha", "CRITICAL", "HARD", "doesn't matter",
		),
		NewServiceState(
			parseTime(timeFmt, "2015-08-20 14:28:00"),
			"testhost-000", "test Service haha", "OK", "HARD", "phew",
		),
	}
	co = CrunchOptions{
		Windows:    []int{7, 14},
		TSInterval: 15,
	}

	metrics, err = CrunchMetrics(states, end, co)
	if err != nil {
		t.Log("Error crunching metrics:", err)
		t.FailNow()
	}

	if len(metrics) != 4 {
		t.Log("Wrong number of metrics: expected 4 but got", len(metrics))
		t.FailNow()
	}
}
