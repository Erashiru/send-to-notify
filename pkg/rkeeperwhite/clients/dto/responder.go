package dto

type SyncResponse struct {
	TaskResponse   TaskResponse   `json:"taskResponse,omitempty"`
	ResponseCommon ResponseCommon `json:"responseCommon"`
	Error          ErrResponse    `json:"error"`
}

type ResponseCommon struct {
	TaskGUID string `json:"taskGuid"`
	TaskType string `json:"taskType"`
	ObjectID int    `json:"objectId"`
}

type TaskResponse struct {
	Order  TaskResponseOrder `json:"order"`
	Status string            `json:"status"`
}

type Params struct {
	OrderGUID           string `json:"orderGuid,omitempty"`
	TaskGUID            string `json:"taskGuid,omitempty"`
	Sync                *Sync  `json:"sync,omitempty"`
	Async               *Sync  `json:"async,omitempty"`
	PriceTypeID         int    `json:"priceTypeId"`
	FilterByKassPresets bool   `json:"filterByKassPresets"`
}

type Sync struct {
	ObjectID int `json:"objectID,omitempty"`
	Timeout  int `json:"timeout,omitempty"`
}
