package types

type HTTPCommonHead struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

type HTTPResponse struct {
	Head HTTPCommonHead `json:"ret"`
	Body interface{}    `json:"body,omitempty"`
}
type GetInfoRequest struct {
	Id int64 `json:"id" form:"id"`
}

type GetInfoResponse struct {
	OsVolume    string   `json:"os_volume"`
	DataVolumes []string `json:"data_volumes"`
	NetInfo     struct {
		Ip      string `json:"ip"`
		Netmask string `json:"netmask"`
		Gateway string `json:"gateway"`
		Dns     string `json:"dns,omitempty"`
	} `json:"net_info"`
	State int `json:"state"`
}
