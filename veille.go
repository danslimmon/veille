package main

import (
	//log "github.com/Sirupsen/logrus"
	"fmt"
	"github.com/docopt/docopt-go"
	"time"
)

func ProcessFiles(arguments map[string]interface{}) error {
	var err error
	var logFilesGeneric interface{}
	var logFiles []string
	var states []State
	var metrics []*Metric

	logFilesGeneric = arguments["<logfile>"]
	logFiles, _ = logFilesGeneric.([]string)

	states, err = ParseFiles(logFiles)
	if err != nil {
		return err
	}
	metrics, err = CrunchMetrics(states, time.Now().Add(time.Duration(120*24)*time.Hour), CrunchOptions{[]int{30, 90}, 15})
	for _, m := range metrics {
		for _, point := range m.Timeseries {
			fmt.Printf("%s %f %d\n", m.Name, point.Value, point.Timestamp.Unix())
		}
	}
	return err
}

func main() {
	var err error
	var usage string

	usage = `veille: determines uptime from Nagios logs.

Usage:
	veille <logfile>...`

	arguments, err := docopt.Parse(usage, nil, true, "veille", false)
	if err != nil {
		panic(err)
	}

	err = ProcessFiles(arguments)
	if err != nil {
		panic(err)
	}
}
