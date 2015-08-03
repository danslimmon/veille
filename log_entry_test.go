package main

import (
	"testing"
)

func TestParseLog_LogRotation(t *testing.T) {
	t.Parallel()
	var err error
	var hasState bool
	var st State

	st, hasState, err = ParseLogLine("[1438041600] LOG ROTATION: DAILY")
	if err != nil {
		t.Log("Got error when parsing log line:", err)
		t.FailNow()
	}
	if hasState {
		t.Log("Got state information when parsing stateless log line:", st)
		t.FailNow()
	}
}
