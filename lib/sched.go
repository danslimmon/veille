package veille

import (
    "fmt"
    "time"
    "log"
)

// Starts a scheduler with the given services' tests.
func RunScheduler(services []*Service, cw *ConfigWatcher) error {
    var tests []*Test
    for _, srv := range services {
        for _, t := range srv.Tests {
            tests = append(tests, t)
        }
    }
    sch := &scheduler{
        Services: services,
        Tests: tests,
    }
    ach, e := RunAlerter()
    if e != nil {
        return e
    }
    if e = sch.Run(cw, ach); e != nil {
        return e
    }
    return nil
}

type scheduler struct {
    Services []*Service
    Tests []*Test
}

func (sch *scheduler) Run(cw *ConfigWatcher, alertChan chan AlertEvent) error {
    resultChan := make(chan TestResult)
    confWatchChan := cw.Subscribe()

    for _, t := range sch.Tests {
        go sch.startTest(t, resultChan)
    }

    fmt.Println("Starting test loop")
    for {
        select {
        case rslt := <-resultChan:
            log.Printf("Got status '%s' from test '%s'\n", rslt.Status, rslt.T.Functionality)
            sch.processResult(rslt, alertChan)
        case <- confWatchChan:
            log.Println("Scheduler received notification of config reload")
        }
    }
}

func (sch *scheduler) processResult(rslt TestResult, alertChan chan AlertEvent) {
    if rslt.Status != "ok" {
        rslt.T.RegFailure(rslt)
        rslt.T.Service.RegFailure(rslt)
    } else {
        rslt.T.RegSuccess(rslt)
        rslt.T.Service.RegSuccess(rslt)
    }
    alertChan <- AlertEvent{&rslt}
}

func (sch *scheduler) startTest(t *Test, resultChan chan TestResult) error {
    tkr := time.NewTicker(time.Duration(t.RunEvery) * time.Second)
    for {
        <-tkr.C
        resultChan <-t.Check()
    }
}
