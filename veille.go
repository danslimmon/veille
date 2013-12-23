package main

import (
    "log"
    "syscall"
    "os/signal"

    "github.com/danslimmon/veille/lib"
)

func main() {
    log.Println("Loading config")

    loader := new(veille.YamlFileConfigLoader)
    loader.Path = "./etc/example.yaml"
    conf, err := loader.GetConfig()
    if err != nil { panic(err) }

    var tests []*veille.Test
    for _, srvConf := range conf.Services {
        srv := new(veille.Service)
        srv.PopFromConf(srvConf)
        log.Printf("Loaded service \"%s\"\n", srv.Name)

        for _, testConf := range srvConf.Tests {
            t := new(veille.Test)
            t.PopFromConf(testConf, srv)
            tests = append(tests, t)
            log.Printf("Loaded test \"%s\" of service \"%s\"\n", t.Functionality, srv.Name)
        }
    }

    cw := new(veille.ConfigWatcher)
    signal.Notify(cw.PublishOnSignals(), syscall.SIGHUP)
    cw.Loader = loader
    err = veille.RunScheduler(tests, cw)
    if err != nil {
        log.Println("Error running scheduler:", err)
    }
}
