package veille

import (
    "fmt"
    "time"
)

// Starts a scheduler with the given services' tests.
func RunScheduler(tests []*Test) error {
    sch := &scheduler{tests}
    if e := sch.Run(); e != nil {
        return e
    }
    return nil
}

type scheduler struct {
    Tests []*Test
}

func (sch *scheduler) Run() error {
    resultChan := make(chan TestResult)
    errorChan := make(chan TestError)
    confWatchChan := ConfigSubscribe()

    for _, t := range sch.Tests {
        go sch.startTest(t, resultChan, errorChan)
    }

    fmt.Println("Starting test loop")
    for {
        select {
        case rslt := <-resultChan:
            fmt.Printf("Got status '%s' from test '%s'\n", rslt.Status, rslt.T.Functionality)
        case e := <-errorChan: 
            fmt.Printf("Got error '%s' from test '%s'\n", e, e.T.Functionality)
        case <- confWatchChan:
            fmt.Println("Scheduler received notification of config reload")
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
