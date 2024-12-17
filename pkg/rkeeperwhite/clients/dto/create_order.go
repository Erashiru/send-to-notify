package dto

type CreateOrderRequest struct {
	TaskType string           `json:"taskType"`
	Params   CreateOrderParam `json:"params"`
}

type CreateOrderTaskResponse struct {
	TaskResponse   CreateOrderTaskBody `json:"taskResponse"`
	ResponseCommon ResponseCommon      `json:"responseCommon"`
	Error          ErrResponse         `json:"error"`
}

type CreateOrderTaskBody struct {
	Order TaskResponseOrder `json:"order"`
}

type TaskResponseOrder struct {
	OrderGuid string                  `json:"orderGuid"`
	Status    OrderResponseBodyStatus `json:"status"`
	TableCode int                     `json:"tableCode"`
}
