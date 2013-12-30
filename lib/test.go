package veille

import (
    "strings"
    "log"
    "os/exec"
    "encoding/json"
)

const DEFAULT_TESTS_DIR = "./tests"


// Tests a given service.
type Test struct {
    // Configured
    Service *Service
    Functionality string
    Script string
    RunEvery int
    AlertAfter int

    // Dynamic
    FailCount int
}

func (t *Test) PopFromConf(testConf TestConfig, s *Service) {
    t.Service = s
    t.Functionality = testConf.Functionality
    t.Script = testConf.Script
    t.RunEvery = testConf.Run_Every
    t.AlertAfter = testConf.Alert_After
}


type TestError struct {
    E string
    T Test
}
func (e *TestError) Error() string {
    return e.E
}

// Runs a script to check the status of a service.
func (t *Test) Check() TestResult {
    path := t.scriptPath()
    log.Printf("Running test '%s'\n", path)

    output, err := exec.Command(path).CombinedOutput()
    if err != nil {
        log.Printf("Error running script '%s'\n", t.Script)
        log.Printf("OUTPUT:\n")
        log.Printf("    %s\n", string(output))
        return TestResult{"error", err.Error(), nil, t}
    }

    var result TestResult
    err = json.Unmarshal(output, &result)
    if err != nil {
        log.Printf("Script '%s' didn't return valid JSON\n", t.Script)
        log.Printf("Error: %s\n", err.Error())
        log.Printf("OUTPUT:\n")
        log.Printf("    %s\n", string(output))
        return TestResult{"error", err.Error(), nil, t}
    }
    result.T = t
    return result
}

func (t *Test) RegSuccess(rslt TestResult) {
    t.FailCount = 0
}

func (t *Test) RegFailure(rslt TestResult) {
    t.FailCount += 1
}

// Returns the full path to the file containing the test script.
func (t *Test) scriptPath() string {
    if strings.HasPrefix(t.Script, "/") {
        return t.Script
    }
    return DEFAULT_TESTS_DIR + "/" + t.Script
}


// The result returned when a Test is run.
type TestResult struct {
    // The status of the probe ("ok" or "problem")
    Status string
    // A message providing details on the probe's result
    Message string
    // Any named metrics returned by the pr
    Metrics map[string]interface{}
    // The probe that generated this result
    T *Test
}
