package models

type Attribute struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	PriceImpact       float32 `json:"price_impact"`
	Available         bool    `json:"available,omitempty"`
	SelectedByDefault bool    `json:"selected_by_default,omitempty"`
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
