package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// setupTailTest prepares for a test of the tail input mode. It returns
// a temporary file and a ConfStruct configured to tail that file.
func setupTailTest() (tmpFile *os.File, inp *TailInput) {
	var conf *ConfStruct
	tmpFile, _ = ioutil.TempFile("", "veille_test_")
	conf = &ConfStruct{
		InputMode: "tail",
		TailPath:  tmpFile.Name(),
	}
	inp = NewTailInput(conf)
	return
}

// teardownTailTest gets rid of the fixture we created for a test of
// the tail input mode.
func teardownTailTest(tmpFile *os.File) {
	tmpFile.Close()
	os.Remove(tmpFile.Name())
}

// Tests basic functionality of the "tail" input mode.
func TestTail(t *testing.T) {
	t.Parallel()
	var line string
	var tmpFile *os.File
	var inp *TailInput
	var lineChan chan string
	var err error

	tmpFile, inp = setupTailTest()
	defer teardownTailTest(tmpFile)

	tmpFile.Write([]byte("line 1\n"))
	tmpFile.Sync()
	err = inp.Start()
	if err != nil {
		t.Log("Error starting TailInput:", err)
		t.FailNow()
	}

	lineChan = inp.LineChan()
	line = <-lineChan
	if line != "line 1" {
		t.Log(fmt.Sprintf("Expected 'line 1' from TailInput but got '%s'", line))
		t.Fail()
	}
}
