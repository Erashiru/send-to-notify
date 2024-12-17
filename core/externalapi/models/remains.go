package models

type RemainsResponse struct {
	Items []RemainsItem `json:"items"`
}

type RemainsItem struct {
	Id    string  `json:"id"`
	Stock float64 `json:"stock"`
}
