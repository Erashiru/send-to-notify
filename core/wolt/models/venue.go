package models

type VenueStatus struct {
	Status struct {
		IsOpen     bool `json:"is_open"`
		IsOnline   bool `json:"is_online"`
		IsIpadFree bool `json:"is_ipad_free"`
	} `json:"status"`
}
