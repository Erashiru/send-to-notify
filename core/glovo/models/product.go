package models

type Restriction struct {
	IsAlcoholic bool `json:"is_alcoholic"`
}

type Product struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Price            float32     `json:"price"`
	ImageURL         string      `json:"image_url"`
	ExtraImageUrls   []string    `json:"extra_image_urls"`
	Available        bool        `json:"available,omitempty"`
	Description      string      `json:"description,omitempty"`
	AttributesGroups []string    `json:"attributes_groups,omitempty"`
	Restrictions     Restriction `json:"restrictions,omitempty"`
}
