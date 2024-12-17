package dto

type MenuSyncReq struct {
	Categories    []MenuSyncCategoryReq      `json:"categories"`
	SubCategories []MenuSyncSubCategoryReq   `json:"subCategories"`
	ModifierItems []MenuSyncModifierItemsReq `json:"modifierItems"`
	Modifiers     []MenuSyncModifierReq      `json:"modifiers"`
	Products      []MenuSyncProductReq       `json:"products"`
}

type MenuSyncCategoryReq struct {
	ExternalID    string   `json:"externalID"`
	Name          string   `json:"name"`
	IsActive      bool     `json:"isActive"`
	Sort          int      `json:"sort"`
	SubCategories []string `json:"subCategories"`
}

type MenuSyncSubCategoryReq struct {
	ExternalID string `json:"externalID"`
	Name       string `json:"name"`
	IsActive   bool   `json:"isActive"`
	Sort       int    `json:"sort"`
}

type MenuSyncModifierItemsReq struct {
	Name       string `json:"name"`
	ExternalID string `json:"externalID"`
	Price      int    `json:"price"`
}

type MenuSyncModifierReq struct {
	Name       string   `json:"name"`
	ExternalID string   `json:"externalID"`
	Items      []string `json:"items"`
}

type MenuSyncProductReq struct {
	ExternalID    string        `json:"externalID"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Price         int           `json:"price"`
	CategoryID    string        `json:"categoryID"`
	Modifiers     []string      `json:"modifiers"`
	Fiscalization Fiscalization `json:"fiscalization"`
	Vat           int           `json:"vat"`
	Images        []Images      `json:"images"`
}

type MenuSyncResp struct {
	Message     string `json:"message"`
	Transaction string `json:"transaction"`
}
