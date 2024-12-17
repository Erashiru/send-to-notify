package dto

type SendOrderErrorNotificationRequest struct {
	IsOrderSent bool                            `json:"isOrderSent"`
	IsPosError  bool                            `json:"isPosError"`
	Error       SendOrderErrorNotificationError `json:"error"`
}

type SendOrderErrorNotificationError struct {
	Message  string `json:"message"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

type ChangeOrderStatusRequest struct {
	Status     string `json:"status"`
	PosNumber  string `json:"posNumber"`
	TerminalId string `json:"terminalId"`
}
