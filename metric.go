package main

import (
	"fmt"
	"sort"
	"time"
)

type TSPoint struct {
	Timestamp time.Time
	Value     float64
}

type Metric struct {
	Name       string
	Timeseries []TSPoint
}

// CrunchOptions describes the metrics that will be created by CrunchMetrics.
type CrunchOptions struct {
	// A list of windows (in days) for which to generate rolling downtime metrics.
	//
	// For example, []int{30,90,365} would result in rolling 30-, 90-, and 365-day
	// downtime metrics.
	Windows []int
	// Interval (in minutes) between data points for the rolling downtime metrics.
	TSInterval int
}

// CrunchMetrics takes a list of states and returns downtime metrics describing them.
func CrunchMetrics(states []State, end time.Time, co CrunchOptions) ([]*Metric, error) {
	var st State
	var statesByService map[string][]State
	var svcStates []State
	//var svcId string
	var metrics []*Metric

	// Break out states by service.
	statesByService = make(map[string][]State)
	for _, st = range states {
		var exists bool

		// Non-service objects (viz. hosts) are not yet implemented
		if st.ObjectType() != "SERVICE" {
			continue
		}

		svcStates, exists = statesByService[st.ObjIdent()]
		if exists {
			statesByService[st.ObjIdent()] = append(svcStates, st)
		} else {
			statesByService[st.ObjIdent()] = []State{st}
		}
	}

	// Okay, here's the number-crunching part:
	//
	// 1. Divide the whole dataset into (probably minutes-long) steps of length
	//    co.TSInterval.
	// 2. Calculate the downtime contained in each step.
	// 3. Once we have enough steps to span the window, start writing data points.
	// 4. For every new step thereafter, shift the oldest step off the list and
	//    recalculate the rolling downtime total.
	for _, svcStates = range statesByService {
		var wSize int
		for _, wSize = range co.Windows {
			var t time.Time
			// prevStatus keeps track of what the service's status was at the end of
			// the previous step.
			var prevStatus, newStatus string
			var stepDur time.Duration
			var stepsInWindow int
			var stepDowntimes []int
			var downSecsInWindow, downSecsInInterval int
			var met *Metric

			met = &Metric{
				Name:       fmt.Sprintf("veille.hostname.servicename.%d-day", wSize),
				Timeseries: make([]TSPoint, 0),
			}
			prevStatus = svcStates[0].Status()
			stepDur = time.Duration(co.TSInterval) * time.Minute
			stepsInWindow = wSize * 24 * 60 / co.TSInterval
			// stepDowntimes contains the seconds-of-downtime counts for the window.
			// As we step through time, values will fall off the left-hand-side of this
			// slice.
			stepDowntimes = make([]int, 0)

			for t = svcStates[0].Timestamp(); t.Before(end); t = t.Add(stepDur) {
				var statesInStep []State
				statesInStep = statesInInterval(svcStates, t, t.Add(stepDur))
				downSecsInInterval, newStatus = downSecs(statesInStep, prevStatus, t, t.Add(stepDur))
				downSecsInWindow += downSecsInInterval
				stepDowntimes = append(stepDowntimes, downSecsInInterval)
				if len(stepDowntimes) > stepsInWindow {
					// Drop the oldest downtime out of the window
					downSecsInWindow -= stepDowntimes[0]
					stepDowntimes = stepDowntimes[1:]

					// Update the metric
					met.Timeseries = append(met.Timeseries, TSPoint{
						Timestamp: t,
						Value:     float64(downSecsInWindow) / float64(wSize*24*60*60),
					})
				}
				prevStatus = newStatus
			}
			metrics = append(metrics, met)
		}
	}

	return metrics, nil
}

// downSecs calculates the number of seconds of downtime in the interval.
//
// It returns the integer number of seconds of downtime in the interval, as well as the
// service's status at the end of the interval.
func downSecs(statesInStep []State, prevStatus string, start, end time.Time) (int, string) {
	var status string
	var lastChange time.Time
	var st State
	var rslt int

	status = prevStatus
	lastChange = start
	for _, st = range statesInStep {
		if prevStatus == st.Status() {
			continue
		}
		if st.Status() != "CRITICAL" {
			rslt += int(st.Timestamp().Sub(lastChange) / time.Second)
		}
		status = st.Status()
		lastChange = st.Timestamp()
	}

	if status == "CRITICAL" {
		rslt += int(end.Sub(lastChange) / time.Second)
	}
	return rslt, status
}

// statesInInterval filters the given States down to the interval.
//
// It returns a new (possibly empty) slice of State structs with timestamps
// equal to or later than start, but before end.
//
// This function assumes that states is sorted.
func statesInInterval(states []State, start, end time.Time) []State {
	var minInd, maxInd int
	minInd = sort.Search(len(states), func(i int) bool {
		return !states[i].Timestamp().Before(start)
	})
	maxInd = sort.Search(len(states), func(i int) bool {
		return states[i].Timestamp().After(end)
	})
	if maxInd == len(states) {
		// There are no states after the end of the interval
		return states[minInd:]
	}
	return states[minInd:maxInd]
}

// calcStart returns the starting time for rolling downtime metrics.
//
// You pass it the first state in the list of states we're crunching, as well as
// the window size (in days) of the rolling metric you're calculating.
func calcStart(st State, w int) time.Time {
	var dur time.Duration
	dur, _ = time.ParseDuration(fmt.Sprintf("%dh", 24*w))
	return st.Timestamp().Add(dur)
}
