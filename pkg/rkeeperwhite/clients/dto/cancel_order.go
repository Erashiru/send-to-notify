package dto

type CancelOrderRequest struct {
	TaskType string           `json:"taskType"`
	Params   CancelOrderParam `json:"params"`
}

type CancelOrderParam struct {
	Async     Sync   `json:"async,omitempty"`
	Sync      Sync   `json:"sync,omitempty"`
	OrderGuid string `json:"orderGuid"`
}

type CancelOrderResponse struct {
	TaskResponse   CancelOrderTaskResponse `json:"taskResponse"`
	ResponseCommon ResponseCommon          `json:"responseCommon"`
	Error          ErrResponse             `json:"error"`
}

type CancelOrderTaskResponse struct {
	Status string `json:"status"`
}
