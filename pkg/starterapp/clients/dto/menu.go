package dto

type CategoryRequest struct {
	Id                int64    `json:"id,omitempty"`
	PosId             string   `json:"posId,omitempty"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Images            []string `json:"images"`
	ParentCategoryIds []int64  `json:"parentCategoryIds"`
	SortIndex         int      `json:"sortIndex"`
	IsActive          bool     `json:"isActive"`
}

type CreateMenuResponse struct {
	Data  []CreateMenuResponseData `json:"data"`
	Count int                      `json:"count"`
}

type CreateMenuResponseData struct {
	PosId string `json:"posId"`
	Id    int64  `json:"id"`
}

type ModifierGroupRequest struct {
	Id        int64                   `json:"id,omitempty"`
	PosId     string                  `json:"posId,omitempty"`
	Name      string                  `json:"name"`
	MaxAmount int                     `json:"maxAmount"`
	MinAmount int                     `json:"minAmount"`
	Required  bool                    `json:"required"`
	Modifiers []ModifierGroupModifier `json:"modifiers"`
}

type ModifierGroupModifier struct {
	Id        int64 `json:"id"`
	MaxAmount int   `json:"maxAmount"`
	MinAmount int   `json:"minAmount"`
	Required  bool  `json:"required"`
}

type ModifiersRequest struct {
	Id        int64    `json:"id,omitempty"`
	PosId     string   `json:"posId,omitempty"`
	Name      string   `json:"name"`
	Price     float64  `json:"price"`
	Images    []string `json:"images"`
	MaxAmount int      `json:"maxAmount"`
	MinAmount int      `json:"minAmount"`
	Required  bool     `json:"required"`
}

type MealRequest struct {
	Id                   int64    `json:"id,omitempty"`
	PosId                string   `json:"posId,omitempty"`
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	Calories             int      `json:"calories"`
	Fats                 int      `json:"fats"`
	Carbohydrates        int      `json:"carbohydrates"`
	Proteins             int      `json:"proteins"`
	Weight               int      `json:"weight"`
	Images               []string `json:"images"`
	ModifierGroups       []int    `json:"modifierGroups"`
	CategoryIds          []int    `json:"categoryIds"`
	DeliveryRestrictions []string `json:"deliveryRestrictions"`
	IsActive             bool     `json:"isActive"`
	SortIndex            int      `json:"sortIndex"`
}

type MealOfferRequest struct {
	ID       int64   `json:"id,omitempty"`
	PosId    string  `json:"posId,omitempty"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	InMenu   bool    `json:"inMenu"`
	MealId   int64   `json:"mealId"`
}

type ModifierOfferRequest struct {
	ID         int64   `json:"id,omitempty"`
	PosId      string  `json:"posId,omitempty"`
	ModifierId int64   `json:"modifierId"`
	ShopId     int     `json:"shopId"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}
