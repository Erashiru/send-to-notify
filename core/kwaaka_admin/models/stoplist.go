package models

type StopListByProductIDRequest struct {
	IsAvailabe bool   `json:"is_availabe"`
	StoreID    string `json:"store_id"`
	ProductID  string `json:"product_id"`
}

type StopListByAttributeIDRequest struct {
	IsAvailabe  bool   `json:"is_availabe"`
	StoreID     string `json:"store_id"`
	AttributeID string `json:"product_id"`
}
