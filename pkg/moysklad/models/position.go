package models

type Position struct {
	OrderID    string     `json:"order_id"`
	ProductID  string     `json:"product_id"`
	Quantity   int        `json:"quantity"`
	Price      float64    `json:"price,omitempty"`
	Assortment Assortment `json:"assortment"`
}

type Assortment struct {
	Meta Meta `json:"meta"`
}
