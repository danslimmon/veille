package veille

import (
    "log"
)

type AlertEvent struct {
    TR *TestResult
}

func RunAlerter() (chan AlertEvent, error) {
    ch := make(chan AlertEvent)
    al := &Alerter{}
    go al.Run(ch)
    return ch, nil
}

type Alerter struct {}

func (al *Alerter) Run(eventChan chan AlertEvent) {
    for {
        select {
        case ev := <-eventChan:
            log.Printf("Alerter received event for test '%s'\n", ev.TR.T.Functionality)
        }
    }
}
