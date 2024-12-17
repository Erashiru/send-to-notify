package models

type PosOrder struct {
	Id              string      `json:"id"`
	Type            string      `json:"type"`
	InfoSystem      string      `json:"infosystem"`
	Department      string      `json:"department"`
	Date            string      `json:"date,omitempty"`
	DeliveryService string      `json:"deliveryService"`
	OrderCode       string      `json:"orderCode"`
	PickUpCode      string      `json:"pickUpCode"`
	Status          string      `json:"status,omitempty"`
	PayMethod       string      `json:"pay_method,omitempty"`
	Change          string      `json:"change,omitempty"`
	Total           string      `json:"total,omitempty"`
	User            OrderUser   `json:"user,omitempty"`
	Address         string      `json:"address,omitempty"`
	Comment         string      `json:"comment,omitempty"`
	Items           []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ProductId string `json:"product_id"`
	Quantity  string `json:"quantity"`
	Price     string `json:"price"`
	Amount    string `json:"amount"`
}

type OrderUser struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
type OrderRequest struct {
	Orders []PosOrder `json:"orders"`
}

type OrderResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Orders  []Order `json:"orders"`
}

type Order struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
