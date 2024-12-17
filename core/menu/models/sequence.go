package models

const (
	MenuSequence        = "menus"
	StoreSequence       = "stores"
	ProductMenuSequence = "products"
)

type Sequence struct {
	Value int `bson:"value" json:"value"`
}
