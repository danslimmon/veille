package main

import (
	"testing"
	"time"
)

func TestParseLog_LogRotation(t *testing.T) {
	t.Parallel()
	var err error
	var ent LogEntry
	var expTime time.Time
	var expType string

	ent, err = ParseLogLine("[1438041600] LOG ROTATION: DAILY")
	if err != nil {
		t.Log("Got error when parsing log line:", err)
		t.Fail()
	}

	expTime = time.Unix(1438041600, 0)
	expType = "rotation"
	if !ent.Timestamp().Equal(expTime) {
		t.Log("Wrong timestamp on log rotation entry. Expected", expTime, "but got", ent.Timestamp())
		t.Fail()
	}
	if ent.Type() != expType {
		t.Log("Wrong type on log rotation entry. Expected", expType, "but got", ent.Type())
		t.Fail()
	}
}
