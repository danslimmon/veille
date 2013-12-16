package veille

import (
    "fmt"
    "os/exec"
    "encoding/json"
)

const DEFAULT_PROBES_DIR = "/usr/veille/probes"

type ProbeError struct {
    ErrStr string
    Pr Probe
}
func (e *ProbeError) Error() string {
    return e.ErrStr
}

// This interface defines a pr that checks whether a service is up.
type Probe interface {
    // The service that this probe monitors
    GetService() Service
    // Uniquely names the probe
    GetName() string
    // Interval (in seconds) at which to run the probe when the service is
    // OK
    GetOKInterval() int
    // Interval (in seconds) at which to run the probe when the service is
    // in a problem state.
    GetProblemInterval() int

    // Runs the check and returns its result
    Check() ProbeResult
}

// A pr that simply runs a script with the parameters given.
type ScriptProbe struct {
    // Generic Probe fields
    Srv Service
    Name string
    OKInterval int
    ProblemInterval int

    // The path to the probe script (inside Dir)
    Script string
    // 'path' will override the default path where probe scripts are found.
    Dir string
    // Any additional parameters that need to be passed to the script
    Params map[string]interface{}
}

// Probe interface getters
func (pr *ScriptProbe) GetService() Service { return pr.Srv }
func (pr *ScriptProbe) GetName() string { return pr.Name }
func (pr *ScriptProbe) GetOKInterval() int { return pr.OKInterval }
func (pr *ScriptProbe) GetProblemInterval() int { return pr.ProblemInterval }

// Runs a script to check the status of a service.
func (pr *ScriptProbe) Check() ProbeResult {
    path := pr.scriptPath()
    fmt.Println("Running probe '" + path + "' with params", pr.Params)

    param_blob, _ := json.Marshal(pr.Params)
    output, err := exec.Command("go", "run", path, "--params", string(param_blob)).CombinedOutput()
    if err != nil {
        fmt.Println("Error running probe '" + path + "' with params", pr.Params)
        fmt.Println("OUTPUT:")
        fmt.Println("   ", string(output))
        return ProbeResult{"error", nil, pr}
    }

    var result ProbeResult
    json.Unmarshal(output, &result)
    result.Pr = pr
    fmt.Println("Probe '" + path + "' returned status '" + result.Status + "'")
    return result
}

// Returns the full path to the file containing the pr script.
func (pr *ScriptProbe) scriptPath() string {
    dir := pr.Dir
    if dir == "" {
        dir = DEFAULT_PROBES_DIR
    }

    return dir + "/" + pr.Script + ".go"
}


// The result returned when a Probe is run.
type ProbeResult struct {
    // The status of the probe ("okay" or "problem")
    Status string
    // Any named metrics returned by the pr
    Metrics map[string]interface{}
    // The probe that generated this result
    Pr Probe
}
