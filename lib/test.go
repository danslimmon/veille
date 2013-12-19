package veille

import (
    "fmt"
    "os/exec"
    "encoding/json"
)

const DEFAULT_TESTS_DIR = "./tests"


// Tests a given service.
type Test struct {
    Service *Service
    Functionality string
    Script string
    RunEvery int
    AlertAfter int
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
    fmt.Println("Running test '" + path + "'")

    output, err := exec.Command(DEFAULT_TESTS_DIR + "/" + t.Script).CombinedOutput()
    if err != nil {
        fmt.Println("Error running script '" + t.Script + "'")
        fmt.Println("OUTPUT:")
        fmt.Println("   ", string(output))
        return TestResult{"error", nil, t}
    }

    var result TestResult
    json.Unmarshal(output, &result)
    result.T = t
    fmt.Println("Test '" + path + "' returned status '" + result.Status + "'")
    return result
}

// Returns the full path to the file containing the test script.
func (t *Test) scriptPath() string {
    return DEFAULT_TESTS_DIR + "/" + t.Script + ".go"
}


// The result returned when a Test is run.
type TestResult struct {
    // The status of the probe ("okay" or "problem")
    Status string
    // Any named metrics returned by the pr
    Metrics map[string]interface{}
    // The probe that generated this result
    T *Test
}
