package dto

type UpdateProduct struct {
	//ID              string   `json:"sku"`
	ExtID           string   `json:"external_id"`
	IsAvailable     *bool    `json:"enabled,omitempty"`
	DiscountedPrice *int     `json:"discounted_price,omitempty"`
	VatPercentage   *float64 `json:"vat_percentage,omitempty"`
}

type UpdateProducts struct {
	Product []UpdateProduct `json:"data"`
}

type UpdateAttribute struct {
	ExtID       string `json:"external_id"`
	IsAvailable *bool  `json:"enabled"`
}

type UpdateAttributes struct {
	Attribute []UpdateAttribute `json:"data"`
}

type Item struct {
	ExtID     string `json:"external_id"`
	Inventory int    `json:"inventory"`
}

type WoltInventory struct {
	Data []Item `json:"data"`
}
