package dto

type AttributeGroup struct {
	ID         int    `json:"id"`
	ExternalID string `json:"externalID"`
	Name       string `json:"name"`
}

type GetAttributeGroupResponse struct {
	Count  int              `json:"count"`
	Result []AttributeGroup `json:"result"`
}
type CreateProductsAttributeGroupRequest struct {
	ProductID  string `json:"-"`
	ExternalID string `json:"externalID"`
	Name       string `json:"name"`
}

type AttributeGroupItem struct {
	ID               int              `json:"id"`
	ExternalID       string           `json:"externalID"`
	Name             string           `json:"name"`
	Price            int              `json:"price"`
	AttachedBranches []AttachedBranch `json:"attachedBranches"`
}

type CreateAttributeGroupsItemRequest struct {
	AttributeGroupID string `json:"-"`
	ExternalID       string `json:"externalID"`
	Name             string `json:"name"`
	Price            int    `json:"price"`
}

type GetAttributeGroupsItemResponse struct {
	ID         int              `json:"id"`
	Name       string           `json:"name"`
	Price      int              `json:"price"`
	ModifierID int              `json:"modifierID"`
	ExternalID string           `json:"externalID"`
	Branches   []AttachedBranch `json:"branches"`
}
