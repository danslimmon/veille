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
    if e := sch.Run(cw); e != nil {
        return e
    }
    return nil
}

type scheduler struct {
    Services []*Service
    Tests []*Test
}

func (sch *scheduler) Run(cw *ConfigWatcher) error {
    resultChan := make(chan TestResult)
    errorChan := make(chan TestError)
    confWatchChan := cw.Subscribe()

    for _, t := range sch.Tests {
        go sch.startTest(t, resultChan, errorChan)
    }

    fmt.Println("Starting test loop")
    for {
        select {
        case rslt := <-resultChan:
            log.Printf("Got status '%s' from test '%s'\n", rslt.Status, rslt.T.Functionality)
        case e := <-errorChan: 
            log.Printf("Got error '%s' from test '%s'\n", e, e.T.Functionality)
        case <- confWatchChan:
            log.Println("Scheduler received notification of config reload")
        }
    }
}

func (sch *scheduler) startTest(t *Test, resultChan chan TestResult,
                                 errorChan chan TestError) error {
    tkr := time.NewTicker(time.Duration(t.RunEvery) * time.Second)
    for {
        <-tkr.C
        resultChan <-t.Check()
    }
}
