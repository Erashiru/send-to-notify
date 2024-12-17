package models

type Event struct {
	Status       int    `json:"status"`
	Type         string `json:"type"`
	RestaurantId string `json:"restaurant_id"`
	Data         Data   `json:"data"`
}

type Data struct {
	// integration_api fields
	RestaurantId string `json:"restaurant_id"`

	// order fields
	OrderId      string `json:"order_id"`
	Number       int    `json:"number"`
	Status       int    `json:"status"`
	CourierName  string `json:"courier_name"`
	CourierPhone string `json:"courier_phone"`

	// stoplist fields
	DeviceId string `json:"device_id"`
	Id       string `json:"id"`
	Count    string `json:"count"`
}

type WebhookType string

const (
	Order    WebhookType = "order"
	StopList WebhookType = "course_count"
)

func (wt WebhookType) String() string {
	return string(wt)
}
