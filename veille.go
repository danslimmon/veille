package main

import (
    "fmt"
    "github.com/danslimmon/veille/lib"
)

func main() {
    fmt.Println("Loading config")
    if err := veille.LoadConfig("./etc/example.yaml"); err != nil {
        panic(err)
    }

    conf := veille.GetConfig()
    var tests []*veille.Test
    for _, srvConf := range conf.Services {
        srv := new(veille.Service)
        srv.PopFromConf(srvConf)
        fmt.Printf("Loaded service \"%s\"\n", srv.Name)

        for _, testConf := range srvConf.Tests {
            t := new(veille.Test)
            t.PopFromConf(testConf, srv)
            tests = append(tests, t)
            fmt.Printf("Loaded test \"%s\" of service \"%s\"\n", t.Functionality, srv.Name)
        }
    }

    err := veille.RunScheduler(tests)
    if err != nil {
        fmt.Println("Error running scheduler:", err)
    }
}
