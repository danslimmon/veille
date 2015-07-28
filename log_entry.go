package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

// An interface for any log entry from the Nagios log
type LogEntry interface {
	// Timestamp returns the time (UTC is assumed) associated with
	// the log entry.
	Timestamp() time.Time
	// Type returns the type of LogEntry you're looking at. See the
	// individual implementations for the possible Type values.
	Type() string
}

// LogRotationLogEntry represents a Nagios log line that specifies how
// often the Nagios log file is rotated. These lines look like:
//
//	[1438041600] LOG ROTATION: DAILY
type LogRotationLogEntry struct {
	T                time.Time
	RotationInterval string
}

func (ent *LogRotationLogEntry) Timestamp() time.Time { return ent.T }
func (ent *LogRotationLogEntry) Type() string         { return "rotation" }

// ParseLogLine takes a line from the Nagios log and returns a
// LogEntry.
func ParseLogLine(logLine string) (LogEntry, error) {
	var err error
	var re *regexp.Regexp
	var groups []string
	var timestampInt64 int64
	var timestamp time.Time

	re = regexp.MustCompile(`^\[([0-9]+)\] (.*)`)
	groups = re.FindStringSubmatch(logLine)
	if groups == nil {
		return nil, errors.New(fmt.Sprintf("Unable to parse timestamp for log line '%s'", logLine))
	}

	timestampInt64, err = strconv.ParseInt(groups[1], 10, 0)
	if err != nil {
		return nil, err
	}
	timestamp = time.Unix(timestampInt64, 0)
	return ParseLogRotationLogEntry(timestamp, groups[2])
}

// ParseLogRotationLogEntry parses a log line that we already know
// indicates the log rotation interval.
func ParseLogRotationLogEntry(timestamp time.Time, remainder string) (LogEntry, error) {
	var re *regexp.Regexp
	var ent *LogRotationLogEntry
	var groups []string

	re = regexp.MustCompile(`LOG ROTATION: (.*)`)
	groups = re.FindStringSubmatch(remainder)
	if groups == nil {
		return nil, errors.New(fmt.Sprintf("Unable to parse rotation interval for log line '%s'", remainder))
	}

	ent = &LogRotationLogEntry{
		T:                timestamp,
		RotationInterval: groups[1],
	}
	log.WithFields(log.Fields{
		"timestamp":         ent.T,
		"rotation_interval": ent.RotationInterval,
	}).Debug("Parsed log entry")
	return ent, nil
}
