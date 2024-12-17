package dto

type BulkUpdateRequest struct {
	Products   []Product   `json:"products"`
	Attributes []Attribute `json:"attributes"`
}

type BulkUpdateResponse struct {
	TransactionID string `json:"transaction_id"`
}
