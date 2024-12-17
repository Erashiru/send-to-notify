package dto

type SendNewsletterRequest struct {
	Name        string `json:"name"`
	RestGroupId string `json:"rest_group_id"`
	Text        string `json:"text"`
	SendTime    string `json:"send_time"`
}

type SendMessage struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
	StoreId string `json:"store_id"`
}
