package models

type OrderEvent struct {
	EventType string         `json:"event"`
	EventBody OrderEventBody `json:"body"`
}

type OrderEventBody struct {
	Order Order `json:"order"`
}

type MenuEvent struct {
	EventType string        `json:"event"`
	EventBody MenuEventBody `json:"body"`
}

type MenuEventBody struct {
	MenuUploadResult MenuUploadResult `json:"menu_upload_result"`
}
