package main

import (
    "fmt"
    "os/exec"
    "encoding/json"
)

var DEFAULT_PROBES_DIR = "/usr/veille/probes"


// A service that we can probe, e.g. "REST API"
type Service struct {
    Name string
    Probe ServiceProbe
}

func (service Service) Check() ProbeResult {
    return service.Probe.Check()
}

// This interface defines a probe that checks whether a service is up.
type ServiceProbe interface {
    Check() ProbeResult
}


// A probe that simply runs a script with the parameters given.
type ScriptProbe struct {
    // Uniquely names the probe and defines its path (unless 'path' is nonempty)
    Script string
    // 'path' will override the default path where probe scripts are found.
    Dir string
    Params map[string]interface{}
}

// Returns the full path to the file containing the probe script.
func (probe ScriptProbe) script_path() string {
    dir := probe.Dir
    if dir == "" {
        dir = DEFAULT_PROBES_DIR
    }

    return dir + "/" + probe.Script + ".go"
}

// Runs a script to check the status of a service.
func (probe ScriptProbe) Check() ProbeResult {
    path := probe.script_path()
    fmt.Println("Running probe '" + path + "' with params", probe.Params)

    param_blob, _ := json.Marshal(probe.Params)
    output, err := exec.Command("go", "run", path, "--params", string(param_blob)).CombinedOutput()
    if err != nil {
        fmt.Println("Error running probe '" + path + "' with params", probe.Params)
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
    // Any named metrics returned by the probe
    Metrics map[string]interface{}
}


func main() {
    probe := ScriptProbe{
        "test_service",
        "./probes",
        map[string]interface{}{
            "port": 80,
        },
    }
        
    srv := Service{"test_service", probe}
    fmt.Println(srv.Check())
}
