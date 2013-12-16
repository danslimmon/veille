package veille

import (
    "fmt"
    "time"
)

// Starts a scheduler with the given services' probes.
func RunScheduler(probes []Probe) error {
    sch := &scheduler{probes}
    if e := sch.Run(); e != nil {
        return e
    }
    return nil
}

type scheduler struct {
    Probes []Probe
}

func (sch *scheduler) Run() error {
    result_chan := make(chan ProbeResult)
    error_chan := make(chan ProbeError)

    for _, pr := range sch.Probes {
        go sch.startProbe(pr, result_chan, error_chan)
    }

    fmt.Println("Starting probe loop")
    for {
        select {
        case rslt := <-result_chan:
            fmt.Printf("Got status '%s' from probe '%s'\n", rslt.Status, rslt.Pr.GetName())
        case e := <-error_chan: 
            fmt.Printf("Got error '%s' from probe '%s'\n", e, e.Pr.GetName())
        }
    }
}

func (sch *scheduler) startProbe(pr Probe, result_chan chan ProbeResult,
                                 error_chan chan ProbeError) error {
    tkr := time.NewTicker(time.Duration(pr.GetOKInterval()) * time.Second)
    for {
        <-tkr.C
        result_chan <-pr.Check()
    }
}
