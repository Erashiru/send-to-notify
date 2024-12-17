package dto

type CreateProductRequest struct {
	ExternalID       string           `json:"externalID"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Price            int              `json:"price"`
	CategoryID       int              `json:"categoryID"`
	Fiscalization    Fiscalization    `json:"fiscalization"`
	Vat              int              `json:"vat"`
	AttachedBranches []AttachedBranch `json:"attachedBranches"`
	Images           []Images         `json:"images"`
	AttributeGroups  []int            `json:"modifiers"`
}

type Fiscalization struct {
	SpicID      string `json:"spicID"`
	PackageCode int    `json:"packageCode"`
}
type AttachedBranch struct {
	ID          int    `json:"id"`
	ExternalID  string `json:"externalID,omitempty"`
	IsActive    bool   `json:"isActive"`
	IsAvailable bool   `json:"isAvailable,omitempty"`
	Qty         int    `json:"qty,omitempty"`
}
type Images struct {
	URL       string `json:"url"`
	IsPreview bool   `json:"isPreview"`
}

type CreateProductResponse struct {
	ID               int              `json:"id"`
	ExternalID       string           `json:"externalID"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Price            int              `json:"price"`
	CategoryID       int              `json:"categoryID"`
	SubCategoryID    int              `json:"subCategoryID"`
	Fiscalization    Fiscalization    `json:"fiscalization"`
	Vat              int              `json:"vat"`
	AttachedBranches []AttachedBranch `json:"attachedBranches"`
	Images           []Images         `json:"images"`
	AttributeGroups  []ModifierGroup  `json:"modifiers"`
}

type UpdateProductRequest struct {
	ID               string            `json:"-"`
	ExternalID       *string           `json:"externalID,omitempty"`
	Name             *string           `json:"name,omitempty"`
	Description      *string           `json:"description,omitempty"`
	Price            *int              `json:"price,omitempty"`
	CategoryID       *int              `json:"categoryID,omitempty"`
	Fiscalization    *Fiscalization    `json:"fiscalization,omitempty"`
	Vat              *int              `json:"vat,omitempty"`
	AttachedBranches *[]AttachedBranch `json:"attachedBranches,omitempty"`
	Images           *[]Images         `json:"images,omitempty"`
	AttributeGroups  *[]int            `json:"modifiers,omitempty"`
}

type UpdateProductResponse struct {
	Id               int              `json:"id"`
	ExternalID       string           `json:"externalID"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Price            int              `json:"price"`
	Fiscalization    Fiscalization    `json:"fiscalization"`
	AttachedBranches []AttachedBranch `json:"attachedBranches"`
	Images           []Images         `json:"images"`
	AttributeGroups  []ModifierGroup  `json:"modifiers"`
	CategoryID       int              `json:"categoryID"`
	Vat              int              `json:"vat"`
}

type GetCategoryProductsResponse struct {
	ID               int              `json:"id"`
	ExternalID       string           `json:"externalID"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Price            int              `json:"price"`
	CategoryID       int              `json:"categoryID"`
	SubCategoryID    int              `json:"subCategoryID"`
	Fiscalization    Fiscalization    `json:"fiscalization"`
	Vat              int              `json:"vat"`
	AttachedBranches []AttachedBranch `json:"attachedBranches"`
	Images           []Images         `json:"images"`
	AttributeGroups  []ModifierGroup  `json:"modifiers"`
}

type ModifierGroup struct {
	ID    int        `json:"id"`
	Name  string     `json:"name"`
	Items []Modifier `json:"items"`
}

type Modifier struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ExternalID string `json:"externalID"`
}
