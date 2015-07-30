package main

import (
	"time"
)

// ReadFiles parses all the given log files and returns a sorted slice
// of LogEntry structs.
func ReadFiles(paths []string) ([]LogEntry, error) {
	var err error
	var path string
	var entries []LogEntry = make([]LogEntry, 0)

	for _, path = range paths {
		var newEntries []LogEntry = make([]LogEntry, 0)
		newEntries, err = ReadFile(path)
		if err != nil {
			return make([]LogEntry, 0), err
		}
		entries = append(entries, newEntries...)
	}
	return entries, nil
}

func ReadFile(path string) ([]LogEntry, error) {
	return []LogEntry{
		&LogRotationLogEntry{time.Now(), "DAILY"},
		&LogVersionLogEntry{time.Now(), "2.0"},
	}, nil
}
