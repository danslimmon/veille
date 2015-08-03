package main

import (
	"bufio"
	"os"
)

// ParseFiles parses all the given log files and returns a sorted slice
// of State structs
func ParseFiles(paths []string) ([]State, error) {
	var err error
	var path string
	var states []State = make([]State, 0)

	for _, path = range paths {
		var newEntries []State = make([]State, 0)
		newEntries, err = ParseFile(path)
		if err != nil {
			return nil, err
		}
		states = append(states, newEntries...)
	}
	return states, nil
}

// ParseFile parses the given Nagios log file and returns the corresponding
// slice of State structs.
func ParseFile(path string) ([]State, error) {
	var err error
	var f *os.File
	var scanner *bufio.Scanner
	var states []State = make([]State, 0)

	f, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		var err error
		var hasState bool
		var state State
		state, hasState, err = ParseLogLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		if hasState {
			states = append(states, state)
		}
	}
	return states, nil
}
