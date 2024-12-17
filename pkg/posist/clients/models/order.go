package models

type Source struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	OrderId string `json:"order_id"`
}

type Payments struct {
	Type string `json:"type"`
}

type Discount struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type Charge struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Customer struct {
	Firstname string `json:"firstname"`
	Mobile    string `json:"mobile"`
	AddType   string `json:"addType"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	City      string `json:"city"`
}

type AddOn struct {
	Id       string `json:"_id"`
	Quantity int    `json:"quantity"`
}

type OrderItem struct {
	Id            string      `json:"id"`
	Quantity      int         `json:"quantity"`
	Rate          int         `json:"rate"`
	Discounts     []Discount  `json:"discounts,omitempty"`
	AddOns        []AddOn     `json:"addOns,omitempty"`
	MapComboItems []ComboItem `json:"mapComboItems,omitempty"`
}

type ComboItem struct {
	Id       string  `json:"id"`
	Quantity int     `json:"quantity"`
	AddOns   []AddOn `json:"addOns,omitempty"`
}

type Order struct {
	Source   Source      `json:"source"`
	Payments Payments    `json:"payments"`
	Discount *Discount   `json:"discount,omitempty"`
	Charges  []Charge    `json:"charges,omitempty"`
	Customer Customer    `json:"customer"`
	TabType  string      `json:"tabType"`
	Items    []OrderItem `json:"items"`
}

type OrderStatusResponse struct {
	Id     string `json:"_id"`
	Status string `json:"status"`
}
