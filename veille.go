package main

import (
    "fmt"

    "github.com/danslimmon/veille"
)

// A service that we can probe, e.g. "REST API"
type Service struct {
    Name string
    Probe probe.ServiceProbe
}

func (svc Service) Check() probe.ProbeResult {
    return svc.Probe.Check()
}

func main() {
    pr := probe.ScriptProbe{
        "test_service",
        "./probes",
        map[string]interface{}{
            "port": 80,
        },
    }
        
    srv := Service{"test_service", pr}
    fmt.Println(srv.Check())
}
