package main

import (
	"bufio"
	"os"
)

// ParseFiles parses all the given log files and returns a sorted slice
// of LogEntry structs.
func ParseFiles(paths []string) ([]LogEntry, error) {
	var err error
	var path string
	var entries []LogEntry = make([]LogEntry, 0)

	for _, path = range paths {
		var newEntries []LogEntry = make([]LogEntry, 0)
		newEntries, err = ParseFile(path)
		if err != nil {
			return nil, err
		}
		entries = append(entries, newEntries...)
	}
	return entries, nil
}

// ParseFile parses the given Nagios log file and returns the corresponding
// slice of LogEntry structs.
func ParseFile(path string) ([]LogEntry, error) {
	var err error
	var f *os.File
	var scanner *bufio.Scanner
	var entries []LogEntry = make([]LogEntry, 0)

	f, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		var err error
		var ent LogEntry
		ent, err = ParseLogLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		entries = append(entries, ent)
	}
	return entries, nil
}
