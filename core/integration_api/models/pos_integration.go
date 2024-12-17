package models

// StopListRequest Update stoplist
type StopListRequest struct {
	RestaurantID string            `json:"restaurant_id"`
	Positions    StopListPositions `json:"positions"`
}

type StopListPositions []StopListPosition

func (positions StopListPositions) ToSliceOfString() []string {
	ids := make([]string, 0, len(positions))

	for _, position := range positions {
		ids = append(ids, position.ID)
	}

	return ids
}

type StopListPosition struct {
	ID      string `json:"id"`
	Balance int    `json:"balance"`
}

// UpdateOrderStatusRequest Update order status
type UpdateOrderStatusRequest struct {
	RestaurantID string `json:"restaurant_id"`
	OrderID      string `json:"order_id"`
	Status       string `json:"status"`
	Reason       string `json:"reason"`
}

// GetOrdersResponse GetOrders
type GetOrdersResponse struct {
	RestaurantID string         `json:"restaurant_id"`
	Orders       []GetOrderBody `json:"orders"`
}

type GetOrderBody struct {
	ID                  string                  `json:"id"`
	Status              string                  `json:"status"`
	OrderCode           string                  `json:"order_code"`
	PickUpCode          string                  `json:"pick_up_code"`
	Products            []GetOrderBodyProduct   `json:"products"`
	OrderTime           string                  `json:"order_time"`
	EstimatedPickUpTime string                  `json:"estimated_pickup_time"`
	Customer            GetOrderBodyCustomer    `json:"customer"`
	Comment             string                  `json:"comment"`
	Address             GetOrderBodyAddress     `json:"address"`
	PeopleCount         int                     `json:"people_count"`
	OrderType           string                  `json:"order_type"`
	DeliveryService     string                  `json:"delivery_service"`
	PaymentInfo         GetOrderBodyPaymentInfo `json:"payment_info"`
}

type GetOrderBodyProduct struct {
	ID                   string                 `json:"product_id"`
	Name                 string                 `json:"name"`
	Price                int                    `json:"price"`
	PriceWithoutDiscount int                    `json:"price_without_discount"`
	Quantity             int                    `json:"quantity"`
	Modifiers            []GetOrderBodyModifier `json:"modifiers"`
}

type GetOrderBodyModifier struct {
	ID       string `json:"modifier_id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type GetOrderBodyPaymentInfo struct {
	Sum  int    `json:"sum"`
	Type string `json:"payment_type"`
}

type GetOrderBodyCustomer struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type GetOrderBodyAddress struct {
	Label string `json:"label"`
}
