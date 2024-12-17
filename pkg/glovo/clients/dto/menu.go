package dto

import "time"

type UploadMenuRequest struct {
	StoreId string `json:"-"`
	MenuURL string `json:"menuUrl"`
}

type UploadMenuResponse struct {
	TransactionID string    `json:"transaction_id"`
	Status        Status    `json:"status,omitempty"`
	UpdatedAt     time.Time `json:"last_updated_at,omitempty"`
	Details       []string  `json:"details,omitempty"`
}

type ValidateMenuRequest struct {
	Attributes       []Attribute       `json:"attributes"`
	AttributeGroups  []AttributeGroup  `json:"attribute_groups"`
	Products         []Product         `json:"products"`
	Collections      []Collection      `json:"collections"`
	SuperCollections []SuperCollection `json:"supercollections"`
}

type ValidateMenuResponse struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}
