package dto

type YTimesWebhookRequest struct {
	Body string `json:"body"`
}

type YTimesUpdateOrderStatusBody struct {
	EventId       string `json:"eventId"`
	Guid          string `json:"guid"`
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
}

type YTimesMenuUpdatesRequestBody struct {
	EventId string `json:"eventId"`
	Type    string `json:"type"`
}
