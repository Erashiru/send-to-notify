package dto

type StopListByProductsRequest struct {
	Products  Products `json:"products"`
	BranchIDs []int    `json:"branches"`
}

type StopListBulkRequest struct {
	Products  Products  `json:"products,omitempty"`
	Modifiers Modifiers `json:"modifiers,omitempty"`
	BranchIDs []int     `json:"branches"`
}

type StopListByAttributesRequest struct {
	Modifiers Modifiers `json:"modifiers"`
	BranchIDs []int     `json:"branches"`
}

type ProductItem struct {
	ExternalID  string `json:"externalID"`
	IsAvailable bool   `json:"isAvailable"`
	Quantity    int    `json:"quantity,omitempty"`
}

type AttributeItem struct {
	ExternalID  string `json:"externalID"`
	IsAvailable bool   `json:"isAvailable"`
}

type Products struct {
	Items                   []ProductItem `json:"items"`
	MakeAvailableOtherItems bool          `json:"makeAvailableOtherItems"`
}

type Modifiers struct {
	Items                   []AttributeItem `json:"items"`
	MakeAvailableOtherItems bool            `json:"makeAvailableOtherItems"`
}
