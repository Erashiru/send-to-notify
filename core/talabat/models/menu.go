package models

type MenuUploadCallbackRequest struct {
	RequestId   string `json:"requestId"`
	Status      string `json:"status"`
	Description string `json:"description"`
}
