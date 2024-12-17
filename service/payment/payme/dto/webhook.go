package dto

type WebhookEvent struct {
	ID     int64         `json:"id,omitempty"`
	Method string        `json:"method"`
	Params WebhookParams `json:"params"`
}

type WebhookParams struct {
	ID      string  `json:"id,omitempty"`
	Time    int64   `json:"time,omitempty"`
	Amount  int     `json:"amount,omitempty"`
	Account Account `json:"account,omitempty"`
}

type WebhookResultResponse struct {
	ID     int64         `json:"id,omitempty"`
	Result WebhookResult `json:"result"`
}

type WebhookResult struct {
	Transaction string `json:"transaction,omitempty"`
	PerformTime int64  `json:"perform_time,omitempty"`
	State       int    `json:"state,omitempty"`
	Allow       bool   `json:"allow,omitempty"`
	CreateTime  int64  `json:"create_time,omitempty"`
	ID          int64  `json:"id,omitempty"`
}
