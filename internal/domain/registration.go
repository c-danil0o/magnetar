package domain

type RegistrationReq struct {
	Labels       []Label
	Resources    map[string]float64
	BindAddress  string
	Capabilities []string
	OsStatus     *OsStatus
}

type RegistrationResp struct {
	NodeId string
}
