package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
	// StrVal returns any uncategorized information in the log entry.
	// This is used for miscellaneous entry types like log rotation
	// interval and log version.
	StrVal() string
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
func (ent *LogRotationLogEntry) StrVal() string       { return ent.RotationInterval }

// LogVersionLogEntry represents a Nagios log line that specifies the
// version of the log file. These lines look like:
//
//	[1438041600] LOG VERSION: 2.0
type LogVersionLogEntry struct {
	T       time.Time
	Version string
}

func (ent *LogVersionLogEntry) Timestamp() time.Time { return ent.T }
func (ent *LogVersionLogEntry) Type() string         { return "version" }
func (ent *LogVersionLogEntry) StrVal() string       { return ent.Version }

// ParseLogLine takes a line from the Nagios log and returns a
// LogEntry.
func ParseLogLine(logLine string) (LogEntry, error) {
	var err error
	var re *regexp.Regexp
	var groups []string
	var remainder, beforeColon string
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

	remainder = groups[2]
	beforeColon = remainder[:strings.IndexByte(remainder, ':')]
	switch beforeColon {
	case "LOG ROTATION":
		return ParseLogRotationLogEntry(timestamp, remainder)
	case "LOG VERSION":
		return ParseLogVersionLogEntry(timestamp, remainder)
	}
	return nil, errors.New(fmt.Sprintf("Unable to parse log entry beginning with '%s'", beforeColon))
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

// ParseLogVersionLogEntry parses a log line that we already know
// indicates the logfile version.
func ParseLogVersionLogEntry(timestamp time.Time, remainder string) (LogEntry, error) {
	var re *regexp.Regexp
	var ent *LogVersionLogEntry
	var groups []string

	re = regexp.MustCompile(`LOG VERSION: (.*)`)
	groups = re.FindStringSubmatch(remainder)
	if groups == nil {
		return nil, errors.New(fmt.Sprintf("Unable to parse version for log line '%s'", remainder))
	}

	ent = &LogVersionLogEntry{
		T:       timestamp,
		Version: groups[1],
	}
	log.WithFields(log.Fields{
		"timestamp":   ent.T,
		"log_version": ent.Version,
	}).Debug("Parsed log entry")
	return ent, nil
}
