package veille

// A service that we can probe, e.g. "REST API"
type Service struct {
    Name string
}

func (srv *Service) PopFromConf(srvConf ServiceConfig) {
    srv.Name = srvConf.Service_Name
}
