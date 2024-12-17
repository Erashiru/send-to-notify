package models

type ProductCategoryDiscount struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	Percent      float64 `json:"percent"`
}

type DiscountUnit struct {
	ID                       string                    `json:"id"`
	Name                     string                    `json:"name"`
	Percent                  float64                   `json:"percent"`
	IsCategorisedDiscount    bool                      `json:"isCategorisedDiscount"`
	ProductCategoryDiscounts []ProductCategoryDiscount `json:"productCategoryDiscounts"`
	Comment                  string                    `json:"comment"`
	CanBeAppliedSelectively  bool                      `json:"canBeAppliedSelectively"`
	MinOrderSum              int                       `json:"minOrderSum"`
	Mode                     string                    `json:"mode"`
	Sum                      float64                   `json:"sum"`
	CanApplyByCardNumber     bool                      `json:"canApplyByCardNumber"`
	IsManual                 bool                      `json:"isManual"`
	IsCard                   bool                      `json:"isCard"`
	IsAutomatic              bool                      `json:"isAutomatic"`
	IsDeleted                bool                      `json:"isDeleted"`
}

type DiscountByStore struct {
	OrganizationID string         `json:"organizationId"`
	Items          []DiscountUnit `json:"items"`
}

type StoreDiscountsResponse struct {
	CorrelationID string            `json:"correlationId"`
	Discounts     []DiscountByStore `json:"discounts"`
}
