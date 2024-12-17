package models

type GetCombosRequest struct {
	OrganizationID string `json:"organizationId"`
	ExtraData      string `json:"extraData,omitempty"`
}

type ComboGroup struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	IsMainGroup bool           `json:"isMainGroup"`
	Products    []ComboProduct `json:"products"`
}

type ComboProduct struct {
	ProductId               string   `json:"productId"`
	SizeId                  string   `json:"sizeId"`
	ForbiddenModifiers      []string `json:"forbiddenModifiers"`
	PriceModificationAmount float64  `json:"priceModificationAmount"`
}

type ComboSpecification struct {
	SourceActionId         string       `json:"sourceActionId"`
	CategoryId             string       `json:"categoryId"`
	Name                   string       `json:"name"`
	PriceModificationType  int          `json:"priceModificationType"`
	PriceModification      float64      `json:"priceModification"`
	IsActive               bool         `json:"isActive"`
	StartDate              string       `json:"startDate"`
	ExpirationDate         string       `json:"expirationDate"`
	LackingGroupsToSuggest int          `json:"lackingGroupsToSuggest"`
	IncludeModifiers       bool         `json:"includeModifiers"`
	Groups                 []ComboGroup `json:"groups"`
}

type ComboCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Warning struct {
	Code      string `json:"Code"`
	ErrorCode string `json:"ErrorCode"`
	Message   string `json:"message"`
}

type GetCombosResponse struct {
	ComboSpecifications []ComboSpecification `json:"comboSpecifications"`
	ComboCategories     []ComboCategory      `json:"comboCategories"`
	Warnings            []Warning            `json:"Warnings"`
}
