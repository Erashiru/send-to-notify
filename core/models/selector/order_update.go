package selector

type OrderStatusUpdate struct {
	OrderID     string
	OrderStatus string
	StoreID     string
}

func EmptyOrderStatusUpdate() OrderStatusUpdate {
	return OrderStatusUpdate{}
}
func (o OrderStatusUpdate) SetOrderID(id string) OrderStatusUpdate {
	o.OrderID = id
	return o
}

func (o OrderStatusUpdate) SetOrderStatus(status string) OrderStatusUpdate {
	o.OrderStatus = status
	return o
}
func (o OrderStatusUpdate) SetStoreID(id string) OrderStatusUpdate {
	o.StoreID = id
	return o
}
