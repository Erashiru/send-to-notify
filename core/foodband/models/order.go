package models

type UpdateOrderStatusReq struct {
	StoreID string `json:"store_id"`
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
