package dto

type Category struct {
	ID         int    `json:"id"`
	ExternalID string `json:"externalID"`
	Name       string `json:"name"`
	IsActive   bool   `json:"isActive"`
	Sort       int    `json:"sort"`
	HasSubs    bool   `json:"hasSubs"`
}

type CreateCategoryRequest struct {
	CategoryID string `json:"-"`
	ExternalID string `json:"externalID"`
	Name       string `json:"name"`
	IsActive   bool   `json:"isActive"`
	Sort       int    `json:"sort"`
}

type CreateCategoryResponse struct {
	ID         int    `json:"id"`
	ExternalID string `json:"externalID"`
	Name       string `json:"name"`
	IsActive   bool   `json:"isActive"`
	Sort       int    `json:"sort"`
	HasSubs    bool   `json:"hasSubs"`
}
