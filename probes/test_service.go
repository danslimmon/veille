package main

import (
    "fmt"
    "flag"
    "encoding/json"

    "github.com/danslimmon/veille"
)

var params_blob = flag.String("params", "", "The parameters with which the probe should be run")

func main() {
    flag.Parse()

    var params probe.ProbeResult
    json.Unmarshal([]byte(*params_blob), &params)

    fmt.Print(`{"status":"ok","metrics":{}}`)
}
