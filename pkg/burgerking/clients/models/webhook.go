package models

type StopListRequest struct {
	StoreIDs   []string            `json:"store_ids"`
	Attributes []StopListAttribute `json:"attributes"`
	Products   []StopListProduct   `json:"products"`
}

type StopListAttribute struct {
	AttributeID string `json:"attributeId"`
	Available   bool   `json:"available"`
}

type StopListProduct struct {
	ProductID string `json:"productId"`
	Available bool   `json:"available"`
}
