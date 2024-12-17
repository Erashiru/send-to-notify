package dto

type UpdateProductsRequest struct {
	Data UpdateProductData `json:"data"`
}
type UpdateProductData struct {
	Branches []int     `json:"branches"`
	Products []Product `json:"products"`
}
type Product struct {
	ExternalId  string `json:"externalId"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price,omitempty"`
	IsAvailable int    `json:"isAvailable"`
	SpicId      string `json:"spicId,omitempty"`
	PackageCode int    `json:"packageCode,omitempty"`
}

type UpdateProductsResponse struct {
	Failed  []Failed  `json:"failed"`
	Updated []Updated `json:"updated"`
}

type Failed struct {
	ExternalId string `json:"external_id"`
	Message    string `json:"message"`
}

type Updated struct {
	ExternalId  string `json:"externalId"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	IsAvailable int    `json:"isAvailable"`
	SpicId      int    `json:"spicId"`
	PackageCode int    `json:"packageCode"`
}

type ProductsError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Code   int    `json:"code"`
}

type UpdateOffersRequest struct {
	Data UpdateOfferData `json:"data"`
}
type UpdateOfferData struct {
	Branches []int     `json:"branches"`
	Options  []Options `json:"options"`
}

type Options struct {
	ExternalId string `json:"externalId"`
	IsActive   int    `json:"isActive"`
}

type UpdateOffersResponse struct {
	Failed  []Failed              `json:"failed"`
	Updated []UpdateOffersUpdated `json:"updated"`
}

type UpdateOffersUpdated struct {
	ExternalId string `json:"externalId"`
	IsActive   int    `json:"isActive"`
}
