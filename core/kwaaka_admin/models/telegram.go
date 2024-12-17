package models

type SendTelegramMessageRequest struct {
	Message          string `json:"message"`
	NotificationType string `json:"notification_type"`
	OrderID          string `json:"order_id"`
}

type SendCompensationMessageRequest struct {
	OrderID            string `json:"order_id"`
	CompensationID     string `json:"compensation_id"`
	CompensationNumber int    `json:"compensation_number"`
	CompensationText   string `json:"compensation_text"`
}
