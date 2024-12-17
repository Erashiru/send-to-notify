package models

const (
	StoplistOperationType = "menu_stop_list"
)

type MenuInitiationRequest struct {
	RestaurantID  string `json:"restaurantId"`
	OperationType string `json:"operationType"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
