package veille

import (
    "fmt"
    "time"
)

// Starts a scheduler with the given services' probes.
func RunScheduler(probes []*Probe) error {
    sch := &scheduler{probes}
    if e := sch.Run(); e != nil {
        return e
    }
    return nil
}

type scheduler struct {
    Probes []*Probe
}

func (sch *scheduler) Run() error {
    resultChan := make(chan ProbeResult)
    errorChan := make(chan ProbeError)
    confWatchChan := ConfigSubscribe()

    for _, pr := range sch.Probes {
        go sch.startProbe(pr, resultChan, errorChan)
    }

    fmt.Println("Starting probe loop")
    for {
        select {
        case rslt := <-resultChan:
            fmt.Printf("Got status '%s' from probe '%s'\n", rslt.Status, rslt.Pr.Name)
        case e := <-errorChan: 
            fmt.Printf("Got error '%s' from probe '%s'\n", e, e.Pr.Name)
        case <- confWatchChan:
            fmt.Println("Scheduler received notification of config reload")
        }
    }
}

func (sch *scheduler) startProbe(pr *Probe, resultChan chan ProbeResult,
                                 errorChan chan ProbeError) error {
    tkr := time.NewTicker(time.Duration(pr.OKInterval) * time.Second)
    for {
        <-tkr.C
        resultChan <-pr.Check()
    }
}
