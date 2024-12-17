package dto

type PayOrderRequest struct {
	TaskType string         `json:"taskType"`
	Params   PayOrderParams `json:"params"`
}

type PayOrderParams struct {
	Async               PayOrderAsync     `json:"async"`
	OrderGuid           string            `json:"orderGuid"`
	IsFullOrderRequired bool              `json:"isFullOrderRequired"`
	Payments            []PayOrderPayment `json:"payments"`
}

type PayOrderAsync struct {
	ObjectId int `json:"objectId"`
	Timeout  int `json:"timeout"`
}

type PayOrderPayment struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
