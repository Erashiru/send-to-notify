package dto

type WebhookEventRequest struct {
	Date            string `form:"date"`
	Amount          string `form:"amount"`
	CardMask        string `form:"cardMask"`
	CardHash        string `form:"cardHash"`
	LowerCommission string `form:"lowerCommission"`
	OperationID     string `form:"operationId"`
	Commission      string `form:"commission"`
	Source          string `form:"source"`
	Login           string `form:"login"`
}

type WebhookEventResponse struct {
	Data int `json:"data"`
}
