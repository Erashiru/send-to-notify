package dto

type Product struct {
	ID               string       `json:"id"`
	Name             string       `json:"name,omitempty"`
	Price            float64      `json:"price,omitempty"`
	ImageURL         string       `json:"image_url,omitempty"`
	ExtraImageUrls   []string     `json:"extra_image_urls"`
	Available        *bool        `json:"available"`
	Description      string       `json:"description"`
	AttributesGroups []string     `json:"attributes_groups"`
	Restrictions     *Restriction `json:"restrictions,omitempty"`
}

type Restriction struct {
	IsAlcoholic bool `json:"is_alcoholic"`
	IsTobacco   bool `json:"is_tobacco"`
}

type ProductModifyRequest struct {
	ID          string   `json:"-"`
	StoreID     string   `json:"-"`
	Price       *float64 `json:"price,omitempty"`
	IsAvailable *bool    `json:"available,omitempty"`
}

type ProductModifyResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Price       string `json:"price"`
	IsAvailable bool   `json:"available"`
	Msg         string `json:"message"`
}
