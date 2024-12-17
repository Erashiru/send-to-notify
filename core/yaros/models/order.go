package models

type OrderUpdateRequestBody struct {
	PosOrderID       string `json:"order_id"`
	Status           string `json:"status"`
	ErrorDescription string `json:"error_description,omitempty"`
	Synchronized     bool   `json:"synch"`
}
