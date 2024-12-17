package models

type OrderResponse struct {
	OrderId       string `json:"order_id"`
	PalomaOrderId int    `json:"paloma_order_id"`
	ReceiptId     int    `json:"receipt_id"`
	Status        string `json:"status"`
}
