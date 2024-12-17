package models

type RestGroupMenu struct {
	ID            string          `bson:"_id,omitempty" json:"id"`
	RestGroupId   string          `bson:"rest_group_id" json:"rest_group_id"`
	Category      []Category      `bson:"category" json:"category"`
	SubCategory   []SubCategory   `bson:"sub_category" json:"sub_category"`
	ModifierItems []ModifierItems `bson:"modifier_items" json:"modifier_items"`
	MenuProduct   []MenuProduct   `bson:"menu_product" json:"menu_product"`
	Modifier      []Modifier      `bson:"modifier" json:"modifier"`
}

type Category struct {
	ExternalID    string   `json:"externalID"`
	Name          string   `json:"name"`
	IsActive      bool     `json:"isActive"`
	Sort          int      `json:"sort"`
	SubCategories []string `json:"subCategories"`
}

type SubCategory struct {
	ExternalID string `json:"externalID"`
	Name       string `json:"name"`
	IsActive   bool   `json:"isActive"`
	Sort       int    `json:"sort"`
}

type ModifierItems struct {
	Name       string `json:"name"`
	ExternalID string `json:"externalID"`
	Price      int    `json:"price"`
}

type Modifier struct {
	Name       string   `json:"name"`
	ExternalID string   `json:"externalID"`
	Items      []string `json:"items"`
}

type MenuProduct struct {
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

type Fiscalization struct {
	SpicID      string `json:"spicID"`
	PackageCode int    `json:"packageCode"`
}

type Images struct {
	URL       string `json:"url"`
	IsPreview bool   `json:"isPreview"`
}
