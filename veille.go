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

    var services []*veille.Service
    for _, srvConf := range conf.Services {
        srv := new(veille.Service)
        srv.PopFromConf(srvConf)
        services = append(services, srv)
    }

    cw := new(veille.ConfigWatcher)
    signal.Notify(cw.PublishOnSignals(), syscall.SIGHUP)
    cw.Loader = loader
    err = veille.RunScheduler(services, cw)
    if err != nil {
        log.Println("Error running scheduler:", err)
    }
}
