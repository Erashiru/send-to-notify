package models

import (
	"github.com/google/uuid"
)

type GetMenuResponse struct {
	Groups            []Group    `json:"groups"`
	ProductCategories []Category `json:"productCategories"`
	Products          []Product  `json:"products"`
	Sizes             []Size     `json:"sizes"`
	Revision          int        `json:"revision"`
}

type Modifier struct {
	ID                 string `json:"id"`
	DefaultAmount      int    `json:"defaultAmount"`
	MinAmount          int    `json:"minAmount"`
	MaxAmount          int    `json:"maxAmount"`
	IsRequired         bool   `json:"required"`
	FreeOfChargeAmount int    `json:"freeOfChargeAmount"`
}

type GroupModifiers struct {
	Modifier
	ChildModifiers []Modifier `json:"childModifiers"`
	Restriction    bool       `json:"childModifiersHaveMinMaxRestrictions,omitempty"`
}

type Price struct {
	ID    string    `json:"sizeId"`
	Price PriceInfo `json:"price"`
}

type PriceInfo struct {
	Value    float64 `json:"currentPrice"`
	IsInMenu bool    `json:"isIncludedInMenu"`
}

type Group struct {
	ID              uuid.UUID     `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Images          []string      `json:"imageLinks"`
	ParentGroup     uuid.NullUUID `json:"parentGroup"`
	Order           int           `json:"order"`
	InMenu          bool          `json:"isIncludedInMenu"`
	IsGroupModifier bool          `json:"isGroupModifier"`
}

type Nomenclature struct {
	Groups     []Group   `json:"groups"`
	Products   []Product `json:"products"`
	Revision   int64     `json:"revision"`
	UploadDate DateTime  `json:"uploadDate"`
}

type Size struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Priority  int    `json:"priority"`
	IsDefault bool   `json:"isDefault"`
}

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsDeleted bool   `json:"isDeleted"`
}

type Product struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	IsDeleted      bool             `json:"isDeleted"`
	GroupID        string           `json:"groupId"`
	CategoryID     string           `json:"productCategoryId"`
	Type           Type             `json:"type"`
	ParentGroup    string           `json:"parentGroup"`
	Images         []string         `json:"imageLinks"`
	Modifiers      []Modifier       `json:"modifiers"`
	GroupModifiers []GroupModifiers `json:"groupModifiers"`
	Article        string           `json:"code"`
	Prices         []Price          `json:"sizePrices"`

	Weight              float64 `json:"weight"`
	MeasureUnit         string  `json:"measureUnit"`
	FatAmount           float64 `json:"fatAmount"`
	ProteinsAmount      float64 `json:"proteinsAmount"`
	CarbohydratesAmount float64 `json:"carbohydratesAmount"`
	EnergyAmount        float64 `json:"energyAmount"`

	FatFullAmount           float64 `json:"fatFullAmount"`
	ProteinsFullAmount      float64 `json:"proteinsFullAmount"`
	CarbohydratesFullAmount float64 `json:"carbohydratesFullAmount"`
	EnergyFullAmount        float64 `json:"energyFullAmount"`
}

func (p Product) IsInMenu() bool {
	return p.Prices[0].Price.IsInMenu
}

func (p Product) Price() float64 {
	return p.Prices[0].Price.Value
}

type GetExternalMenuResponse struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	ItemCategories []ItemCategories `json:"itemCategories"`
}
type Prices struct {
	OrganizationID string  `json:"organizationId"`
	Price          float64 `json:"price"`
}
type Restrictions struct {
	MinQuantity  int `json:"minQuantity"`
	MaxQuantity  int `json:"maxQuantity"`
	FreeQuantity int `json:"freeQuantity"`
	ByDefault    int `json:"byDefault"`
}
type AllergenGroups struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
type NutritionPerHundredGrams struct {
	Fats                float64 `json:"fats"`     // check on existing
	Proteins            float64 `json:"proteins"` // check on existing
	Carbs               float64 `json:"carbs"`    // check on existing
	Energy              float64 `json:"energy"`   // check on existing
	CarbohydratesAmount float64 `json:"carbohydratesAmount"`
	ProteinsAmount      float64 `json:"proteinsAmount"`
	EnergyAmount        float64 `json:"energyAmount"`
	FatAmount           float64 `json:"fatAmount"`
}

type Tags struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type ModifierGroupItems struct {
	Prices                   []Prices                 `json:"prices"`
	Sku                      string                   `json:"sku"`
	Name                     string                   `json:"name"`
	Description              string                   `json:"description"`
	ButtonImage              string                   `json:"buttonImage"`
	Restrictions             Restrictions             `json:"restrictions"`
	AllergenGroups           []AllergenGroups         `json:"allergenGroups"`
	NutritionPerHundredGrams NutritionPerHundredGrams `json:"nutritionPerHundredGrams"`
	PortionWeightGrams       int                      `json:"portionWeightGrams"`
	Tags                     []Tags                   `json:"tags"`
	ItemID                   string                   `json:"itemId"`
}
type ItemModifierGroups struct {
	Items                                []ModifierGroupItems `json:"items"`
	Name                                 string               `json:"name"`
	Description                          string               `json:"description"`
	Restrictions                         Restrictions         `json:"restrictions"`
	CanBeDivided                         bool                 `json:"canBeDivided"`
	ItemGroupID                          string               `json:"itemGroupId"`
	ChildModifiersHaveMinMaxRestrictions bool                 `json:"childModifiersHaveMinMaxRestrictions"`
	Sku                                  string               `json:"sku"`
	IsDefault                            bool
}
type ItemSizes struct {
	Prices                   []Prices                 `json:"prices"`
	ItemModifierGroups       []ItemModifierGroups     `json:"itemModifierGroups"`
	Sku                      string                   `json:"sku"`
	SizeCode                 string                   `json:"sizeCode"`
	SizeName                 string                   `json:"sizeName"`
	IsDefault                bool                     `json:"isDefault"`
	PortionWeightGrams       float64                  `json:"portionWeightGrams"`
	MeasureUnitType          string                   `json:"measureUnitType"`
	SizeID                   string                   `json:"sizeId"`
	NutritionPerHundredGrams NutritionPerHundredGrams `json:"nutritionPerHundredGrams"`
	ButtonImageURL           string                   `json:"buttonImageUrl"`
}

type TaxCategory struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Percentage int    `json:"percentage"`
}
type Items struct {
	ItemSizes        []ItemSizes      `json:"itemSizes"`
	Sku              string           `json:"sku"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	AllergenGroups   []AllergenGroups `json:"allergenGroups"`
	ItemID           string           `json:"itemId"`
	ModifierSchemaID string           `json:"modifierSchemaId"`
	TaxCategory      TaxCategory      `json:"taxCategory"`
	OrderItemType    string           `json:"orderItemType"`
}
type ItemCategories struct {
	Items          []Items `json:"items"`
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	ButtonImageURL string  `json:"buttonImageUrl"`
	HeaderImageURL string  `json:"headerImageUrl"`
}

type GetExternalMenuRequest struct {
	ExternalMenuID  string   `json:"externalMenuId"`
	OrganizationIDS []string `json:"organizationIds"`
	PriceCategoryId string   `json:"priceCategoryId,omitempty"`
}
