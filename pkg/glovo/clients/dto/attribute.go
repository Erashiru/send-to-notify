package dto

type Attribute struct {
	ID                string  `json:"id"`
	Name              string  `json:"name,omitempty"`
	PriceImpact       float64 `json:"price_impact"`
	Available         bool    `json:"available"`
	SelectedByDefault bool    `json:"selected_by_default"`
}

type AttributeGroup struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Min               int      `json:"min"`
	Max               int      `json:"max"`
	Collapse          bool     `json:"collapse"`
	MultipleSelection bool     `json:"multiple_selection"`
	Attributes        []string `json:"attributes"`
}

type AttributeModifyRequest struct {
	ID          string   `json:"-"`
	StoreID     string   `json:"-"`
	Price       *float64 `json:"price_impact,omitempty"`
	IsAvailable *bool    `json:"available,omitempty"`
}

type AttributeModifyResponse struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Price             string `json:"price"`
	Available         bool   `json:"available"`
	SelectedByDefault bool   `json:"selected_by_default"`
}
