package dto

type Product struct {
	StoreID   string  `json:"store_id"`
	ProductID string  `json:"item_id"`
	Price     float64 `json:"price,omitempty"`
	Available bool    `json:"available"`
}

type Modifier struct {
	StoreID    string  `json:"store_id"`
	ModifierID string  `json:"item_id"`
	Price      float64 `json:"price,omitempty"`
	Available  bool    `json:"available"`
}
