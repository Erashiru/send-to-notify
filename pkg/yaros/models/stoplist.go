package models

type StopListItem struct {
	Id       string `json:"ID"`
	Quantity int    `json:"QUANTITY"`
}

type StopListResponse struct {
	Status        string         `json:"status"`
	StopListItems []StopListItem `json:"stopList"`
}
