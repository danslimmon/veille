package veille

import (
    "log"
)

// A service that we can probe, e.g. "REST API"
type Service struct {
    Name string
    Tests []*Test
}

func (srv *Service) PopFromConf(srvConf ServiceConfig) {
    srv.Name = srvConf.Service_Name
    for _, testConf := range srvConf.Tests {
        t := new(Test)
        t.PopFromConf(testConf, srv)
        srv.Tests = append(srv.Tests, t)
        log.Printf("Loaded test \"%s\" of service \"%s\"\n", t.Functionality, srv.Name)
    }
    log.Printf("Loaded service \"%s\"\n", srv.Name)
}

func (s *Service) RegSuccess(rslt TestResult) {
}

func (s *Service) RegFailure(rslt TestResult) {
}
