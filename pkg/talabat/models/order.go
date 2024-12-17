package models

type AcceptOrderRequest struct {
	OrderToken     string `json:"-"`
	AcceptanceTime string `json:"acceptanceTime"`
	RemoteOrderId  string `json:"remoteOrderId"`
	Status         string `json:"status"`
}

type RejectOrderRequest struct {
	OrderToken string `json:"-"`
	Message    string `json:"message"`
	Reason     string `json:"reason"`
	Status     string `json:"status"`
}

type OrderPickedUpRequest struct {
	OrderToken string `json:"-"`
	Status     string `json:"status"`
}
