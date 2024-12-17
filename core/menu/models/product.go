package models

import (
	"time"

	coreModels "github.com/kwaaka-team/orders-core/core/models"
)

type Products []Product

type Product struct {
	ExtID                   string                  `bson:"ext_id" json:"id"`
	MsID                    string                  `bson:"ms_id" json:"ms_id"` // fixme: why we need moysklad ID ?
	PosID                   string                  `bson:"pos_id" json:"pos_id"`
	ProductID               string                  `bson:"product_id" json:"product_id"`
	StarterAppID            string                  `bson:"starter_app_id" json:"starter_app_id"`
	StarterAppOfferID       string                  `bson:"starter_app_offer_id" json:"starter_app_offer_id"`
	IngredientID            string                  `bson:"ingredient_id" json:"ingredient_id"` //poster field- for match stoplist
	IsCombo                 bool                    `bson:"is_combo" json:"is_combo"`
	ChocofoodFoodId         string                  `bson:"chocofood_food_id" json:"chocofood_food_id"`
	SizeID                  string                  `bson:"size_id,omitempty" json:"size_id"`
	GroupID                 string                  `bson:"group_id" json:"group_id"`
	ParentGroupID           string                  `bson:"parent_group_id" json:"parent_group_id"`
	Section                 string                  `bson:"section" json:"section"`
	AttributesGroups        []string                `bson:"attributes_groups,omitempty" json:"attributes_groups"`
	Attributes              []string                `bson:"attributes,omitempty" json:"attributes"`
	Code                    string                  `bson:"code" json:"code"`
	DefaultAttributes       []string                `bson:"default_attributes,omitempty" json:"default_attributes"`
	MenuDefaultAttributes   []MenuDefaultAttributes `bson:"default_attributes_objects,omitempty" json:"menu_default_attributes"`
	Name                    []LanguageDescription   `bson:"name" json:"name"`
	Description             []LanguageDescription   `bson:"description" json:"description"`
	ExtName                 string                  `bson:"ext_name" json:"ext_name"`
	ImageURLs               []string                `bson:"image_urls" json:"images"`
	ImageHash               string                  `bson:"image_hash" json:"image_hash"`
	IsAlcohol               bool                    `bson:"alcohol" json:"alcohol"`
	IsTobacco               bool                    `bson:"tobacco" json:"tobacco"`
	AlcoholPercentage       float32                 `bson:"alcohol_percentage" json:"alcohol_percentage"`
	IsAvailable             bool                    `bson:"available" json:"is_available"`
	IsDeleted               bool                    `bson:"is_deleted" json:"is_deleted"`
	IsIgnored               bool                    `bson:"is_ignored" json:"is_ignored"`
	IsDisabled              bool                    `bson:"is_disabled" json:"is_disabled"` // disable by admin or by sections
	IsFavorite              bool                    `bson:"is_favorite" json:"is_favorite"` // for favorite foods in direct
	IsIncludedInMenu        bool                    `bson:"included_in_menu" json:"is_included_in_menu"`
	HasRequiredAttributes   bool                    `bson:"-" json:"has_required_attributes" `
	IsDeletedReason         string                  `bson:"is_deleted_reason" json:"is_deleted_reason"`
	Price                   []Price                 `bson:"price" json:"price"`
	IsSync                  bool                    `bson:"sync" json:"is_sync"`
	ByAdmin                 bool                    `bson:"by_admin" json:"by_admin"`
	Changes                 []MenuChanges           `bson:"changes,omitempty" json:"changes"`
	Balance                 float64                 `bson:"balance" json:"balance"`
	AverageCookingTime      string                  `bson:"average_cooking_time,omitempty" json:"average_cooking_time,omitempty"`
	FatAmount               float64                 `bson:"fat_amount,omitempty" json:"fat_amount,omitempty"`
	ProteinsAmount          float64                 `bson:"proteins_amount,omitempty" json:"proteins_amount,omitempty"`
	CarbohydratesAmount     float64                 `bson:"carbohydrates_amount,omitempty" json:"carbohydrates_amount,omitempty"`
	EnergyAmount            float64                 `bson:"energy_amount,omitempty" json:"energy_amount,omitempty"`
	FatFullAmount           float64                 `bson:"fat_full_amount,omitempty" json:"fat_full_amount,omitempty"`
	ProteinsFullAmount      float64                 `bson:"proteins_full_amount,omitempty" json:"proteins_full_amount,omitempty"`
	CarbohydratesFullAmount float64                 `bson:"carbohydrates_full_amount,omitempty" json:"carbohydrates_full_amount,omitempty"`
	EnergyFullAmount        float64                 `bson:"energy_full_amount,omitempty" json:"energy_full_amount,omitempty"`
	Weight                  float64                 `bson:"weight" json:"weight"`
	MeasureUnit             string                  `bson:"measure_unit" json:"measure_unit"`
	ProductsCreatedAt       ProductsCreatedAt       `bson:"created_at" json:"products_created_at"`
	UpdatedAt               coreModels.Time         `bson:"updated_at" json:"updated_at"`
	DiscountPrice           DiscountPrice           `bson:"discount_price" json:"discount_price"`
	CookingTime             int32                   `bson:"cooking_time" json:"cooking_time"`
	IsPosIsDeleted          bool                    `bson:"-" json:"pos_is_deleted"`
	Differences             []string                `bson:"differences,omitempty" json:"differences,omitempty"`
	ProductErr              string                  `bson:"product_err" json:"product_err"`
	SpicID                  string                  `bson:"spic_id" json:"spic_id"`
	PackageCode             string                  `json:"package_code" bson:"package_code"`
	ProductInformation      ProductInformation      `bson:"product_information" json:"product_information"`
	Barcode                 Barcode                 `bson:"barcode" json:"barcode"`
	IsCatchWeight           bool                    `bson:"is_catch_weight" json:"is_catch_weight"`
	VendorCode              string                  `bson:"vendor_code" json:"vendor_code"`
	DisabledByValidation    bool                    `bson:"disabled_by_validation" json:"disabled_by_validation"`
	Halal                   bool                    `bson:"halal" json:"halal"`
}
type DiscountPrice struct {
	IsActive bool    `bson:"is_active" json:"is_active"`
	Value    float64 `bson:"value" json:"value"`
}

type ProductRequest struct {
	ID           string   `json:"id"`
	Name         string   `json:"name,omitempty"`
	Price        float64  `json:"price"`
	IsAvailable  *bool    `json:"available"`
	ImageURL     string   `json:"image_url,omitempty"`
	Images       []string `json:"extra_image_urls,omitempty"`
	Restrictions struct {
		IsAlcohol bool `json:"is_alcoholic,omitempty"`
	} `json:"restrictions,omitempty"`
	Description     string   `json:"description,omitempty"`
	AttributeGroups []string `json:"attributes_groups,omitempty"`
	Balance         float64  `json:"balance,omitempty"`
	MSID            string   `json:"ms_id,omitempty"`
	DiscountPrice   float64  `json:"discount_price,omitempty"`
}

type AttributeRequest struct {
	ID          string  `json:"id"`
	Name        string  `json:"name,omitempty"`
	Price       float64 `json:"price"`
	IsAvailable *bool   `json:"available"  binding:"required"`
	Balance     float64 `json:"balance"`
}

func (ar AttributeRequest) ToModel(attributes []Attribute) []AttributeRequest {

	var result []AttributeRequest

	for _, attribute := range attributes {
		temp := AttributeRequest{
			ID:          attribute.ExtID,
			IsAvailable: OfBool(attribute.IsAvailable),
		}
		result = append(result, temp)
	}
	return result
}
func (p *ProductRequest) ToModel(products []Product) []ProductRequest {
	var result []ProductRequest

	for _, product := range products {
		temp := ProductRequest{
			ID:          product.ExtID,
			IsAvailable: OfBool(product.IsAvailable),
			MSID:        product.MsID,
			Balance:     product.Balance,
		}
		result = append(result, temp)
	}
	return result
}

func OfBool(b bool) *bool {
	return &b
}

func (p *ProductRequest) FromModel(products []ProductRequest) StopListProducts {
	var result = make(StopListProducts, 0, len(products))

	for _, product := range products {
		temp := StopListProduct{
			ExtID: product.ID,
			MsID:  product.MSID,
		}
		if product.IsAvailable != nil {
			temp.IsAvailable = *product.IsAvailable
		}
		result = append(result, temp)
	}
	return result
}

func (p *ProductRequest) FromAttribute(attributes []AttributeRequest) []ProductRequest {
	var result = make([]ProductRequest, 0, len(attributes))

	for _, attribute := range attributes {
		temp := ProductRequest{
			ID:          attribute.ID,
			IsAvailable: attribute.IsAvailable,
			Balance:     attribute.Balance,
		}
		result = append(result, temp)
	}
	return result
}

func (p *ProductRequest) RemoveDuplicate(products []ProductRequest) []ProductRequest {
	var unique []ProductRequest
Loop:
	for _, v := range products {
		for i, u := range unique {
			if v.ID == u.ID {
				unique[i] = v
				continue Loop
			}
		}
		unique = append(unique, v)
	}
	return unique
}

func (p *ProductRequest) LessZeroList(products []ProductRequest) []ProductRequest {
	var result []ProductRequest
	for _, item := range products {
		if item.Balance <= 0 {
			result = append(result, item)
		}
	}
	return result
}

func (p Product) FromAttribute(attributes []Attribute) []Product {
	var result = make([]Product, 0, len(attributes))

	for _, attribute := range attributes {
		temp := Product{
			ExtID:       attribute.ExtID,
			IsAvailable: attribute.IsAvailable,
			Balance:     attribute.Balance,
		}
		result = append(result, temp)
	}
	return result
}

func (p *Product) RemoveDuplicate(products []Product) []Product {
	var unique []Product
Loop:
	for _, v := range products {
		for i, u := range unique {
			if v.ExtID == u.ExtID {
				unique[i] = v
				continue Loop
			}
		}
		unique = append(unique, v)
	}
	return unique
}

func (ar AttributeRequest) FromModel(attributes []AttributeRequest) StopListAttributes {

	var result = make(StopListAttributes, 0, len(attributes))

	for _, attribute := range attributes {

		temp := StopListAttribute{
			AttributeID: attribute.ID,
		}
		if attribute.IsAvailable != nil {
			temp.IsAvailable = *attribute.IsAvailable
		}
		result = append(result, temp)
	}
	return result
}

type MoySkladPosition struct {
	//ID           string    `json:"id" bson:"_id"`
	OrderID      string    `bson:"order_id" json:"order_id"`
	RestaurantID string    `bson:"restaurant_id" json:"restaurant_id"`
	ProductID    string    `bson:"product_id" json:"product_id"`
	MsID         string    `bson:"ms_id" json:"ms_id"`
	Code         string    `bson:"code" json:"code"`
	ID           string    `bson:"position_id" json:"id"`
	Available    bool      `bson:"is_available" json:"available"`
	IsDeleted    bool      `bson:"is_deleted" json:"is_deleted"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}

func (p Product) LessZeroList(products Products) Products {
	var result []Product
	for _, item := range products {
		if item.Balance <= 0 {
			result = append(result, item)
		}
	}
	return result
}

type ChocoFoodState string

const (
	Active  ChocoFoodState = "active"
	Soldout ChocoFoodState = "sold_out"
	Removed ChocoFoodState = "removed"
)

func (s ChocoFoodState) String() string {
	return string(s)
}

func (p Products) Unique() Products {

	existProducts := make(map[string]struct{}, len(p))
	results := make(Products, 0, len(p))

	for _, product := range p {

		if _, ok := existProducts[product.ExtID]; ok {
			continue
		}

		results = append(results, product)
		existProducts[product.ExtID] = struct{}{}
	}

	return results
}

func (p Products) Get(productId string) (Product, bool) {

	for _, product := range p {
		if product.ExtID == productId {
			return product, true
		}
	}
	return Product{}, false
}

func (p Products) GetSections() []string {

	res := make([]string, 0, len(p))
	for _, product := range p {
		res = append(res, product.Section)
	}

	return res
}

type ProductsCreatedAt struct {
	Value     coreModels.Time `bson:"value" json:"value"`
	Timezone  string          `bson:"timezone" json:"timezone"`
	TimeStamp int64           `bson:"-" json:"timestamp"`
	UTCOffset float64         `bson:"utc_offset" json:"utc_offset"`
}

type MenuDefaultAttributes struct {
	ExtID         string    `bson:"ext_id" json:"ext_id"`
	Name          string    `bson:"name" json:"name"`
	DefaultAmount int       `bson:"default_amount" json:"default_amount"`
	ByAdmin       bool      `bson:"by_admin" json:"by_admin"`
	Price         int       `bson:"price" json:"price"`
	Attribute     Attribute `bson:"-" json:"attribute,omitempty"`
}

type ProductModifyResponse struct {
	ExtID       string
	ExtName     string
	Price       string
	IsAvailable bool
	Msg         string
}

// UpdateStopListProduct allows you to update product or attribute in aggregator
type UpdateStopListProduct struct {
	ProductID string
	SetToStop bool
	Data      []UpdateStoreData
}

type UpdateStoreData struct {
	ID          string
	Aggregators []Aggregator
}

type UpdateStoreResponse struct {
	ID                string
	AggregatorStoreID string
	Success           bool
	Msg               string
}

type ItemStopList struct {
	ID          string
	Price       float64
	IsAvailable bool
}

type MatchingProducts struct {
	ProductToChange     string
	AggregatorProductID string
	PosProductID        string
	MenuID              string
	IsSync              bool
}

type ProductUpdateRequest struct {
	IsAvailable *bool `json:"is_available"`
	IsDisabled  *bool `json:"is_disabled"`
}

type Aggregators []Aggregator

type Aggregator struct {
	Name     AggregatorName
	IsActive bool
	Success  bool
	Msg      string
}

type ProductWithErr struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	AttributeGroupID string   `json:"attribute_group_id"`
	ProductErr       string   `json:"product_err"`
	Solution         []string `json:"solution"`
}

type UpdateProductImageAndDescription struct {
	PosID       string                `bson:"pos_id"`
	ImageURLs   []string              `bson:"image_urls" json:"images"`
	Weight      float64               `bson:"weight" json:"weight"`
	MeasureUnit string                `bson:"measure_unit" json:"measure_unit"`
	Description []LanguageDescription `bson:"description" json:"description"`
	Price       []Price               `bson:"price" json:"price"`
}

type Barcode struct {
	Type           string `bson:"type" json:"type"`
	Value          string `bson:"value" json:"value"`
	WeightEncoding string `bson:"weight_encoding" json:"weight_encoding"`
	Values         string `bson:"values,omitempty" json:"values,omitempty"`
}

type UpdateAttributePrice struct {
	ExtID string  `bson:"ext_id" json:"ext_id"`
	Price float64 `bson:"price_impact" json:"price_impact"`
}
