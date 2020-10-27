package codisapi

type ClusterResult []string

type CodisInfo struct {
	Closed bool  `json:"closed"`
	Group  Group `json:"group"`
	Proxy  Proxy `json:"proxy"`
}
type Group struct {
	SModels []SModels `json:"models"`
}
type SModels struct {
	Id        int       `json:"id"`
	Servers   []Servers `json:"servers"`
	OutOfSync bool      `json:"out_of_sync"`
}
type Servers struct {
	Server string `json:"server"`
}

type Proxy struct {
	PModels []PModels `json:"models"`
}

type PModels struct {
	ProxyAddr string `json:"proxy_addr"`
	Hostname  string `json:"hostname"`
}
