package models

type Combos []Combo

type Combo struct {
	ID             string       `bson:"id" json:"id"`
	Name           string       `bson:"name" json:"name"`
	Price          Price        `bson:"price" json:"price"`
	IsActive       bool         `bson:"is_active" json:"is_active"`
	SourceActionID string       `bson:"source_action_id" json:"source_action_id"`
	ProgramID      string       `bson:"program_id" json:"program_id"`
	ComboGroup     []ComboGroup `bson:"groups" json:"combo_group"`
}

type ComboGroup struct {
	Id          string         `bson:"id" json:"id"`
	Name        string         `bson:"name" json:"name"`
	IsMainGroup bool           `bson:"is_main_group" json:"is_main_group"`
	Products    []ComboProduct `bson:"products" json:"products"`
}

type ComboProduct struct {
	ProductId               string `bson:"product_id" json:"product_id"`
	Name                    string `bson:"name" json:"name"`
	PriceModificationAmount Price  `bson:"price_modification_amount" json:"price_modification_amount"`
	IsExistInMenu           bool   `bson:"is_exist_in_menu" json:"is_exist_in_menu"`
}
