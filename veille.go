package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docopt/docopt-go"
)

func ProcessFiles(arguments map[string]interface{}) error {
	var err error
	var logFilesGeneric interface{}
	var logFiles []string
	var entries []LogEntry

	logFilesGeneric = arguments["<logfile>"]
	logFiles, _ = logFilesGeneric.([]string)
	entries, err = ReadFiles(logFiles)
	log.Info(entries)
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
