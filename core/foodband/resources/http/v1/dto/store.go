package dto

type Store struct {
	ID               string   `json:"id" example:"test-id"`
	Name             string   `json:"name" example:"test restaurant"`
	PosType          string   `json:"pos_type" example:"test-pos"`
	DeliveryServices []string `json:"delivery_services" example:"glovo,wolt"`
}

type ManageAggregatorStoreRequest struct {
	IsOpen *bool `json:"is_open" binding:"required" example:"true"`
}
