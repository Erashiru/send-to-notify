package dto

type ErrResponse struct {
	WsError    WsError `json:"wsError,omitempty"`
	AgentError WsError `json:"agentError,omitempty"`
}

type WsError struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}
