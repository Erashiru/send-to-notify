package models

type SendMenuRequest struct {
	Attributes       []Attribute      `json:"attributes"`
	AttributeGroups  []AttributeGroup `json:"attribute_groups"`
	Products         []MenuProduct    `json:"products"`
	Collections      []Collection     `json:"collections"`
	SuperCollections []string         `json:"supercollections"`
	StoreIDs         []string         `json:"store_ids"`
}

type MenuProduct struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	ImageURL       string           `json:"image_url"`
	Price          int              `json:"price"`
	Description    string           `json:"description"`
	Available      bool             `json:"available"`
	AttributeGroup []AttributeGroup `json:"attribute_groups"`
}

type Collection struct {
	Name     string    `json:"name"`
	Position int       `json:"position"`
	ImageURL string    `json:"image_url"`
	Sections []Section `json:"sections"`
}

type Section struct {
	Name     string   `json:"name"`
	Position int      `json:"position"`
	ImageURL string   `json:"image_url"`
	Products []string `json:"products"`
}
