package main

import (
	"io/ioutil"
	"os"
	"testing"
)

// setupInputTest prepares for a test of input functionality. It copies
// the contents of the given test file into a temporary file and returns
// the latter.
//
// If testLogFile is "", then the returned temporary file will be empty.
func setupInputTest(testLogFile string) (tmpFile *os.File) {
	tmpFile, _ = ioutil.TempFile("", "veille_test_")
	if testLogFile != "" {
		var err error
		var inputF *os.File
		var inputContents []byte

		inputF, err = os.Open(testLogFile)
		if err != nil {
			panic(err)
		}
		inputContents, err = ioutil.ReadAll(inputF)
		if err != nil {
			panic(err)
		}
		_, err = tmpFile.Write(inputContents)
		if err != nil {
			panic(err)
		}
		tmpFile.Sync()
		tmpFile.Seek(0, 0)
	}
	return
}

// teardownInputTest gets rid of the fixture we created for a test of
// the input functionality.
func teardownInputTest(tmpFile *os.File) {
	tmpFile.Close()
	os.Remove(tmpFile.Name())
}

// Tests basic functionality of ParseFile().
func TestParseFile(t *testing.T) {
	t.Parallel()
	var err error
	var entries []LogEntry
	var f *os.File

	f = setupInputTest("test_data/empty.log")
	defer teardownInputTest(f)

	entries, err = ParseFile(f.Name())
	if err != nil {
		t.Log("Error reading almost-empty Nagios log file:", err)
		t.FailNow()
	}
	if len(entries) != 2 {
		t.Log("Wrong number of entries from almost-empty Nagios log file; expected 2 but got", len(entries))
		t.Fail()
	}
	if entries[0].StrVal() != "DAILY" {
		t.Log("Wrong rotation interval from almost-empty Nagios log file:", entries[0].StrVal())
		t.Fail()
	}
	if entries[1].StrVal() != "2.0" {
		t.Log("Wrong log version from almost-empty Nagios log file:", entries[1].StrVal())
		t.Fail()
	}
}
