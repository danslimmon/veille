package probe

import (
    "fmt"
    "os/exec"
    "encoding/json"
)

const DEFAULT_PROBES_DIR = "/usr/veille/probes"

// This interface defines a pr that checks whether a service is up.
type ServiceProbe interface {
    Check() ProbeResult
}

// A pr that simply runs a script with the parameters given.
type ScriptProbe struct {
    // Uniquely names the pr and defines its path (unless 'path' is nonempty)
    Script string
    // 'path' will override the default path where pr scripts are found.
    Dir string
    Params map[string]interface{}
}

// Returns the full path to the file containing the pr script.
func (pr ScriptProbe) script_path() string {
    dir := pr.Dir
    if dir == "" {
        dir = DEFAULT_PROBES_DIR
    }

    return dir + "/" + pr.Script + ".go"
}

// Runs a script to check the status of a service.
func (pr ScriptProbe) Check() ProbeResult {
    path := pr.script_path()
    fmt.Println("Running probe '" + path + "' with params", pr.Params)

    param_blob, _ := json.Marshal(pr.Params)
    output, err := exec.Command("go", "run", path, "--params", string(param_blob)).CombinedOutput()
    if err != nil {
        fmt.Println("Error running probe '" + path + "' with params", pr.Params)
        fmt.Println("OUTPUT:")
        fmt.Println("   ", string(output))
        return ProbeResult{"error", nil}
    }

    var result ProbeResult
    json.Unmarshal(output, &result)
    fmt.Println("Probe '" + path + "' returned status '" + result.Status + "'")
    return result
}

// The result returned when a ServiceProbe is run.
type ProbeResult struct {
    Status string
    // Any named metrics returned by the pr
    Metrics map[string]interface{}
}
