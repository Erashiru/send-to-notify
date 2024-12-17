package dto

import "time"

type TaskRequest struct {
	TaskType string `json:"taskType"`
	Params   Params `json:"params"`
}

type MenuResponse struct {
	TaskResponse   MenuTaskResponse `json:"taskResponse"`
	ResponseCommon ResponseCommon   `json:"responseCommon"`
	ErrResponse    ErrResponse      `json:"error,omitempty"`
}

type MenuTaskResponse struct {
	Menu          Menu      `json:"menu"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

type Menu struct {
	Categories         Categories         `json:"categories"`
	Products           Products           `json:"products"`
	IngredientsSchemes IngredientsSchemes `json:"ingredientsSchemes"`
	IngredientsGroups  IngredientsGroups  `json:"ingredientsGroups"`
	Ingredients        Ingredients        `json:"ingredients"`
	LastUpdatedAt      time.Time          `json:"lastUpdatedAt"`
	Version            int                `json:"version"`
}

type Categories []Category

type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentId string `json:"parentId,omitempty"`
}

type Products []Product

type Product struct {
	ID          string   `json:"id"`
	CategoryID  string   `json:"categoryId"`
	Name        string   `json:"name"`
	Price       string   `json:"price"`
	SchemeId    string   `json:"schemeId,omitempty"` // вариант блюда с модификаторами или комбо
	Description string   `json:"description"`
	ImageUrls   []string `json:"imageUrls"`
	Measure     Measure  `json:"measure"`
}

type Measure struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type IngredientsSchemes []IngredientsScheme

type IngredientsScheme struct {
	ID                string                  `json:"id"`
	IngredientsGroups IngredientsSchemeGroups `json:"ingredientsGroups,omitempty"`
}

type IngredientsSchemeGroups []IngredientsSchemeGroup

type IngredientsSchemeGroup struct {
	ID       string `json:"id"`
	MinCount int    `json:"minCount"`
	MaxCount int    `json:"maxCount"`
}

type IngredientsGroups []IngredientsGroup

type IngredientsGroup struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
}

type Ingredients []Ingredient

type Ingredient struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Price            string   `json:"price"`
	Description      string   `json:"description"`
	MaxAmountForDish int      `json:"maxAmountForDish"`
	ImageUrls        []string `json:"imageUrls"`
	Measure          Measure  `json:"measure"`
}
