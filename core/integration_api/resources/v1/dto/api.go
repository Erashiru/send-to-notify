package dto

type RequestGenerateNewAggregatorMenu struct {
	StoreID          string `json:"store_id"`
	Delivery         string `json:"delivery"`
	AggregatorMenuId string `json:"aggregator_menu_id"`
}

type RequestAutoUpdateAggregatorMenu struct {
	StoreID          string `json:"store_id"`
	Delivery         string `json:"delivery"`
	AggregatorMenuId string `json:"aggregator_menu_id"`
}

type RequestSetMarkUpToAggregatorMenu struct {
	StoreId string `json:"store_id"`
	MenuId  string `json:"menu_id"`
}

type RequestGenerateAggregatorMenuFromPos struct {
	StoreId  string `json:"store_id"`
	Delivery string `json:"delivery"`
}

type ErrorResponse struct {
	Status  int    `json:"status,omitempty"`
	Error   int    `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
