package main

import (
	"testing"
)

// Tests the parsing of a couple log lines that don't contain state info
func TestParseLog_NonStateful(t *testing.T) {
	t.Parallel()
	var err error
	var hasState bool
	var st State
	var line string
	var lines []string

	lines = []string{
		"[1438041600] LOG ROTATION: DAILY",
		"[1438044812] Auto-save of retention data completed successfully.",
	}
	for _, line = range lines {
		st, hasState, err = ParseLogLine(line)
		if err != nil {
			t.Log("Got error when parsing log line:", err)
			t.FailNow()
		}
		if hasState {
			t.Log("Got state information when parsing stateless log line:", st)
			t.FailNow()
		}
	}
}

// Tests the parsing of a log line that contains initial host state info.
func TestParseLog_CurrentHostState(t *testing.T) {
	t.Parallel()
	var err error
	var hasState bool
	var st State

	st, hasState, err = ParseLogLine("[1438041600] CURRENT HOST STATE: api-old;UP;HARD;1;SSH OK - Totally SSHed to host, no problem")
	if err != nil {
		t.Log("Got error when parsing log line:", err)
		t.FailNow()
	}
	if !hasState {
		t.Log("Missing state information when parsing CURRENT HOST STATE line")
		t.FailNow()
	}
	if st.Timestamp().Unix() != 1438041600 {
		t.Log("Wrong timestamp from CURRENT HOST STATE line; expected 1438041600 but got", st.Timestamp().Unix())
		t.Fail()
	}
	if st.ObjectType() != "HOST" {
		t.Log("Wrong object type from CURRENT HOST STATE line; expected HOST but got", st.ObjectType())
		t.Fail()
	}
	if st.Hostname() != "api-old" {
		t.Log("Wrong hostname from CURRENT HOST STATE line; expected api-old but got", st.Hostname())
		t.Fail()
	}
	if st.Servicename() != "" {
		t.Log("Non-empty servicename from CURRENT HOST STATE line:", st.Servicename())
		t.Fail()
	}
	if st.Status() != "UP" {
		t.Log("Non-UP status from CURRENT HOST STATE line:", st.Status())
		t.Fail()
	}
	if st.Hardness() != "HARD" {
		t.Log("Non-HARD hardness from CURRENT HOST STATE line:", st.Hardness())
		t.Fail()
	}
	if st.PluginOutput() != "SSH OK - Totally SSHed to host, no problem" {
		t.Log("Wrong plugin output from CURRENT HOST STATE line:", st.PluginOutput())
		t.Fail()
	}
}

// Tests the parsing of a log line that contains initial service state info.
func TestParseLog_CurrentServiceState(t *testing.T) {
	t.Parallel()
	var err error
	var hasState bool
	var st State

	st, hasState, err = ParseLogLine("[1438041600] CURRENT SERVICE STATE: site-000;site-000 HTTPS;OK;HARD;1;TCP OK - 0.058 second response time on port 443")
	if err != nil {
		t.Log("Got error when parsing log line:", err)
		t.FailNow()
	}
	if !hasState {
		t.Log("Missing state information when parsing CURRENT SERVICE STATE line")
		t.FailNow()
	}
	if st.Timestamp().Unix() != 1438041600 {
		t.Log("Wrong timestamp from CURRENT SERVICE STATE line; expected 1438041600 but got", st.Timestamp().Unix())
		t.Fail()
	}
	if st.ObjectType() != "SERVICE" {
		t.Log("Wrong object type from CURRENT SERVICE STATE line; expected SERVICE but got", st.ObjectType())
		t.Fail()
	}
	if st.Hostname() != "site-000" {
		t.Log("Wrong hostname from CURRENT SERVICE STATE line; expected api-old but got", st.Hostname())
		t.Fail()
	}
	if st.Servicename() != "site-000 HTTPS" {
		t.Log("Wrong servicename from CURRENT SERVICE STATE line:", st.Servicename())
		t.Fail()
	}
	if st.Status() != "OK" {
		t.Log("Non-OK status from CURRENT SERVICE STATE line:", st.Status())
		t.Fail()
	}
	if st.Hardness() != "HARD" {
		t.Log("Non-HARD hardness from CURRENT SERVICE STATE line:", st.Hardness())
		t.Fail()
	}
	if st.PluginOutput() != "TCP OK - 0.058 second response time on port 443" {
		t.Log("Wrong plugin output from CURRENT SERVICE STATE line:", st.PluginOutput())
		t.Fail()
	}
}

// Tests the parsing of a log line that contains a service alert.
func TestParseLog_ServiceAlert(t *testing.T) {
	t.Parallel()
	var err error
	var hasState bool
	var st State

	st, hasState, err = ParseLogLine("[1438041722] SERVICE ALERT: api-dev;Old API without SSL;CRITICAL;SOFT;1;CRITICAL: Failed to write value to API within timeout")
	if err != nil {
		t.Log("Got error when parsing log line:", err)
		t.FailNow()
	}
	if !hasState {
		t.Log("Missing state information when parsing SERVICE ALERT line")
		t.FailNow()
	}
	if st.Timestamp().Unix() != 1438041722 {
		t.Log("Wrong timestamp from SERVICE ALERT line; expected 1438041722 but got", st.Timestamp().Unix())
		t.Fail()
	}
	if st.ObjectType() != "SERVICE" {
		t.Log("Wrong object type from SERVICE ALERT line; expected SERVICE but got", st.ObjectType())
		t.Fail()
	}
	if st.Hostname() != "api-dev" {
		t.Log("Wrong hostname from SERVICE ALERT line; expected api-old but got", st.Hostname())
		t.Fail()
	}
	if st.Servicename() != "Old API without SSL" {
		t.Log("Wrong servicename from SERVICE ALERT line:", st.Servicename())
		t.Fail()
	}
	if st.Status() != "CRITICAL" {
		t.Log("Non-CRITICAL status from SERVICE ALERT line:", st.Status())
		t.Fail()
	}
	if st.Hardness() != "SOFT" {
		t.Log("Non-SOFT hardness from SERVICE ALERT line:", st.Hardness())
		t.Fail()
	}
	if st.PluginOutput() != "CRITICAL: Failed to write value to API within timeout" {
		t.Log("Wrong plugin output from SERVICE ALERT line:", st.PluginOutput())
		t.Fail()
	}
}
