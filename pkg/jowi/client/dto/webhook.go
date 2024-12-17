package dto

type EventType string

const (
	StopListEvent   EventType = "course_count"
	OrderEvent      EventType = "order"
	RestaurantEvent EventType = "restaurant"
)

type JowiEvent struct {
	Status       int       `json:"status"`
	Type         EventType `json:"type"`
	RestaurantId string    `json:"restaurant_id"`
	Data         Data      `json:"data"`
}

type Data struct {
	OrderId      string `json:"order_id"`
	Number       int    `json:"number"`
	Status       int    `json:"status"`
	RestaurantId string `json:"restaurant_id"`
	CourierName  string `json:"courier_name"`
	CourierPhone string `json:"courier_phone"`
	Amount       string `json:"amount"`

	Id       string `json:"id"`
	DeviceId string `json:"device_id"`
	Count    int    `json:"count"`
}

type OrderWebhook struct {
	Status       int       `json:"status"`
	Type         EventType `json:"type"`
	RestaurantId string    `json:"restaurant_id"`
	Data         OrderData `json:"data"`
}

type OrderData struct {
	OrderId      string `json:"order_id"`
	Number       int    `json:"number"`
	Status       int    `json:"status"`
	RestaurantId string `json:"restaurant_id"`
	CourierName  string `json:"courier_name"`
	CourierPhone string `json:"courier_phone"`
	Amount       int    `json:"amount"`
}

type StopListWebhook struct {
	Status       int       `json:"status"`
	Type         EventType `json:"type"`
	RestaurantId string    `json:"restaurant_id"`
	Data         OrderData `json:"data"`
}

type StopListData struct {
	Id       string `json:"id"`
	DeviceId string `json:"device_id"`
	Count    int    `json:"count"`
}
