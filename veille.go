package main

import (
    "fmt"
    "github.com/danslimmon/veille/lib"
)

func main() {
    fmt.Println("Reading config")
    if err := veille.LoadConfig("./etc/example.yaml"); err != nil {
        panic(err)
    }
    srv := veille.Service{"test_service"}
    pr := &veille.Probe{
        srv,            //Srv
        "test_service", //Name
        2,              //OKInterval
        1,              //ProblemInterval

        "test_service", //Script
        "./probes",     //Dir
        map[string]interface{}{
            "port": 80,
        },
    }

    probes := make([]*veille.Probe, 0, 256)
    probes = append(probes, pr)
    err := veille.RunScheduler(probes)
    if err != nil {
        fmt.Println("Error running scheduler:", err)
    }
}
