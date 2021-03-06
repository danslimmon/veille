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

// ParseLogLine takes a line from the Nagios log and returns the state
// information contained therein.
func ParseLogLine(logLine string) (st State, hasState bool, err error) {
	var re *regexp.Regexp
	var groups []string
	var remainder, beforeColon string
	var timestampInt64 int64
	var timestamp time.Time
	var colonInd int

	re = regexp.MustCompile(`^\[([0-9]+)\] (.*)`)
	groups = re.FindStringSubmatch(logLine)
	if groups == nil {
		return nil, false, errors.New(fmt.Sprintf("Unable to parse timestamp for log line '%s'", logLine))
	}

	timestampInt64, err = strconv.ParseInt(groups[1], 10, 0)
	if err != nil {
		return nil, false, err
	}
	timestamp = time.Unix(timestampInt64, 0)

	remainder = groups[2]
	colonInd = strings.IndexByte(remainder, ':')
	if colonInd == -1 {
		return nil, false, nil
	}

	beforeColon = remainder[:colonInd]
	switch beforeColon {
	default:
		return nil, false, nil
	case "CURRENT HOST STATE":
		return ParseHostStateLogLine(timestamp, remainder)
	case "CURRENT SERVICE STATE":
		return ParseServiceStateLogLine(timestamp, remainder)
	case "SERVICE ALERT":
		return ParseServiceAlertLogLine(timestamp, remainder)
	}
	return nil, false, errors.New(fmt.Sprintf("Unable to parse log entry beginning with '%s'", beforeColon))
}

// ParseCurrentHostStateLogLine parses a log line that we already know
// indicates the current status of a given host.
func ParseHostStateLogLine(timestamp time.Time, remainder string) (State, bool, error) {
	var re *regexp.Regexp
	var st State
	var groups []string

	re = regexp.MustCompile(`CURRENT HOST STATE: ([^;]+);([A-Z]+);([A-Z]+);[^;]+;(.*)$`)
	groups = re.FindStringSubmatch(remainder)
	if groups == nil {
		return nil, false, errors.New(fmt.Sprintf("Unable to parse current host status for log line '%s'", remainder))
	}

	st = NewHostState(timestamp, groups[1], groups[2], groups[3], groups[4])
	log.WithFields(log.Fields{
		"timestamp": st.Timestamp(),
		"hostname":  st.Hostname(),
		"status":    st.Status(),
	}).Debug("Parsed log entry")
	return st, true, nil
}

// ParseCurrentServiceStateLogLine parses a log line that we already know
// indicates the current status of a given host.
func ParseServiceStateLogLine(timestamp time.Time, remainder string) (State, bool, error) {
	var re *regexp.Regexp
	var st State
	var groups []string

	re = regexp.MustCompile(`CURRENT SERVICE STATE: ([^;]+);([^;]+);([A-Z]+);([A-Z]+);[^;]+;(.*)$`)
	groups = re.FindStringSubmatch(remainder)
	if groups == nil {
		return nil, false, errors.New(fmt.Sprintf("Unable to parse current service status for log line '%s'", remainder))
	}

	st = NewServiceState(timestamp, groups[1], groups[2], groups[3], groups[4], groups[5])
	log.WithFields(log.Fields{
		"timestamp":   st.Timestamp(),
		"hostname":    st.Hostname(),
		"servicename": st.Servicename(),
		"status":      st.Status(),
	}).Debug("Parsed log entry")
	return st, true, nil
}

// ParseCurrentServiceAlertLogLine parses a log line that we already know
// indicates the current status of a given host.
func ParseServiceAlertLogLine(timestamp time.Time, remainder string) (State, bool, error) {
	var re *regexp.Regexp
	var st State
	var groups []string

	re = regexp.MustCompile(`SERVICE ALERT: ([^;]+);([^;]+);([A-Z]+);([A-Z]+);[0-9]+;(.*)$`)
	groups = re.FindStringSubmatch(remainder)
	if groups == nil {
		return nil, false, errors.New(fmt.Sprintf("Unable to parse service alert for log line '%s'", remainder))
	}

	st = NewServiceState(timestamp, groups[1], groups[2], groups[3], groups[4], groups[5])
	log.WithFields(log.Fields{
		"timestamp":   st.Timestamp(),
		"hostname":    st.Hostname(),
		"servicename": st.Servicename(),
		"status":      st.Status(),
	}).Debug("Parsed log entry")
	return st, true, nil
}
