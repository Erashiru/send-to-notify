package models

type OrderItem struct {
	ObjectId      int                 `json:"object_id"`
	Name          string              `json:"name"`
	Count         int                 `json:"count"`
	Price         int                 `json:"price"`
	Modifications []OrderModification `json:"modifications"`
	ComplexItems  []OrderComplexItem  `json:"complex_items"`
}

type OrderModification struct {
	ObjectId int    `json:"object_id"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
	Price    int    `json:"price"`
}

type OrderComplexItem struct {
	ObjectId int    `json:"object_id"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
	Price    int    `json:"price"`
}

type Order struct {
	OrderId        string      `json:"order_id"`
	Date           string      `json:"date"`
	Name           string      `json:"name"`
	Phone          string      `json:"phone"`
	Email          string      `json:"email"`
	Address        string      `json:"address"`
	CoordinateLong string      `json:"coordinate_long,omitempty"`
	CoordinateLat  string      `json:"coordinate_lat,omitempty"`
	Comment        string      `json:"comment,omitempty"`
	PersonAmount   int         `json:"person_amount,omitempty"`
	TotalPrice     int         `json:"total_price"`
	DiscountAmount int         `json:"discount_amount,omitempty"`
	Exchange       int         `json:"exchange,omitempty"`
	DeliveryType   int         `json:"delivery_type"`
	IsCash         bool        `json:"is_cash"`
	IsPayed        bool        `json:"is_payed,"`
	OrderItems     []OrderItem `json:"order_items"`
}
