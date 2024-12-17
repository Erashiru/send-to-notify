package dto

type OrderUpdateRequest struct {
	ID      int64  `json:"-"`
	Status  string `json:"status"`
	StoreID string `json:"-"`
}

type OrderUpdateResponse struct {
	Code       string `json:"code"`
	RequestID  string `json:"requestId"`
	Domain     string `json:"domain"`
	Message    string `json:"message"`
	StaticCode int    `json:"staticCode"`
}
