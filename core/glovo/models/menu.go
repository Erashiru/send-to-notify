package models

type Collection struct {
	Name     string    `json:"name"`
	Position int       `json:"position"`
	ImageURL string    `json:"image_url"`
	Sections []Section `json:"sections"`
}

type Section struct {
	Name     string   `json:"name"`
	Position int      `json:"position"`
	Products []string `json:"products"`
}
type Supercollections struct {
	Name        string   `json:"name"`
	Position    int      `json:"position"`
	ImageURL    string   `json:"image_url"`
	Collections []string `json:"collections"`
}

type Menu struct {
	Attributes       []Attribute        `json:"attributes"`
	AttributeGroups  []AttributeGroup   `json:"attribute_groups"`
	Products         []Product          `json:"products"`
	Collections      []Collection       `json:"collections"`
	Supercollections []Supercollections `json:"supercollections"`
}
