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

type NodeUpdateReq struct {
	NodeId       NodeId
	Capabilities []string
	OsStatus     *OsStatus
}

type NodeUpdateResp struct {
	NotFound bool
}
