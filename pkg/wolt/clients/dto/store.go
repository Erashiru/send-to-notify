package dto

type IsStoreOpen struct {
	AvailableStore string
	VenueId        string
}

type IsStoreOpenRequest struct {
	Status string `json:"status"`
}

type Status struct {
	IsOpen     bool `json:"is_open"`
	IsOnline   bool `json:"is_online"`
	IsIpadFree bool `json:"is_ipad_free"`
}

type ContactDetails struct {
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type OpeningTimes struct {
	OpeningDay  string `json:"opening_day"`
	OpeningTime string `json:"opening_time"`
	ClosingDay  string `json:"closing_day"`
	ClosingTime string `json:"closing_time"`
}

type SpecialTimes struct {
	OpeningDate string `json:"opening_date"`
	OpeningTime string `json:"opening_time"`
	ClosingDate string `json:"closing_date"`
	ClosingTime string `json:"closing_time"`
}

type DeliveryArea struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type StoreStatusResponse struct {
	Status                Status         `json:"status"`
	ContactDetails        ContactDetails `json:"contact_details"`
	OpeningTimes          []OpeningTimes `json:"opening_times"`
	SpecialTimes          []SpecialTimes `json:"special_times"`
	DeliveryArea          []DeliveryArea `json:"delivery_area"`
	LastThreeOrdersStatus []string       `json:"last_three_orders_status"`
}
