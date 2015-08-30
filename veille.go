package main

import (
	//log "github.com/Sirupsen/logrus"
	"github.com/docopt/docopt-go"
)

func ProcessFiles(arguments map[string]interface{}) error {
	return nil
}

/*
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
	metrics, err = CrunchMetrics(states, []int{30})
	log.Info(metrics)
	return err
}
*/

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
