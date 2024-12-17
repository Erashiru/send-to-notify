package dto

type RestaurantReport struct {
	ID       string
	Name     string
	Quantity int
	Report   []OrderStatusBody
}

type OrderStatusBody struct {
	ID              string
	OrderID         string
	OrderCode       string
	PosOrderID      string
	DeliveryService string
	Message         string
}
