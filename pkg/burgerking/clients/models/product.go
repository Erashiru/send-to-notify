package models

type Product struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Quantity           int         `json:"quantity"`
	Price              int         `json:"price"`
	Attributes         []Attribute `json:"attributes"`
	PurchasedProductID string      `json:"purchased_product_id"`
}

type Attribute struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Name     string `json:"name"`
}

type AttributeGroup struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Min               int         `json:"min"`
	Max               int         `json:"max"`
	Collapse          bool        `json:"collapse"`
	MultipleSelection bool        `json:"multiple_selection"`
	Attributes        []Attribute `json:"attributes"`
}
