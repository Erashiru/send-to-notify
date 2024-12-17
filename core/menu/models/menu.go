package models

import (
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"time"
)

type Menu struct {
	ID             string  `bson:"_id,omitempty" json:"id"`
	Name           string  `bson:"name" json:"name"`
	ExtName        string  `bson:"ext_name,omitempty" json:"ext_name"`
	Description    string  `bson:"description,omitempty" json:"description"`
	Delivery       string  `bson:"delivery" json:"delivery"`
	Comment        Comment `bson:"comment,omitempty" json:"comment"`
	IsActive       bool    `bson:"is_active" json:"is_active"`
	IsDeleted      bool    `bson:"is_deleted" json:"is_deleted"`
	IsSync         bool    `bson:"is_sync" json:"is_sync"`
	SyncAttributes bool    `bson:"sync_attributes" json:"sync_attributes"`

	StopLists        []string        `bson:"stoplist" json:"stoplists"`
	Attributes       Attributes      `bson:"attributes" json:"attributes"`
	AttributesGroups AttributeGroups `bson:"attributes_groups" json:"attributes_groups"`
	Sections         Sections        `bson:"sections" json:"sections"`
	Products         Products        `bson:"products" json:"products"`
	Combos           Combos          `bson:"combos" json:"combos"`

	Collections      MenuCollections      `bson:"collections,omitempty" json:"collections"`
	SuperCollections MenuSuperCollections `bson:"super_collections,omitempty" json:"super_collections"`

	AggregatorName string `bson:"aggregator_name,omitempty" json:"aggregator_name"`

	Groups Groups `bson:"groups,omitempty" json:"groups"`

	UpdatedAt          coreModels.Time `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt          coreModels.Time `bson:"created_at,omitempty" json:"created_at"`
	HasMoreThanOneLang bool            `bson:"has_more_than_one_lang" json:"has_more_than_one_lang"`
	IsDiscount         bool            `bson:"is_discount,omitempty" json:"is_discount"`
	HasWoltPromo       bool            `bson:"has_wolt_promo" json:"has_wolt_promo"`
	QRPromo            []string        `bson:"qr_promo" json:"qr_promo"`
	IsProductOnStop    bool            `bson:"-" json:"is_product_on_stop"`
	CreationSource     string          `bson:"creation_source,omitempty" json:"creation_source"`
}

type MenuChanges struct {
	TxID         string          `bson:"transaction_id" json:"tx_id"`
	MenuID       string          `bson:"menu_id" json:"menu_id"`
	RestaurantID string          `bson:"restaurant_id" json:"restaurant_id"`
	Delivery     string          `bson:"delivery" json:"delivery"`
	CreatedAt    coreModels.Time `bson:"created_at" json:"created_at"`
	Products     Products        `bson:"products" json:"products"`
	Attributes   Attributes      `bson:"attributes" json:"attributes"`
}

type MenuCollections []MenuCollection

type MenuCollection struct {
	ExtID           string    `bson:"ext_id" json:"id"`
	StarterAppID    string    `bson:"starter_app_id" json:"starter_app_id"`
	Name            string    `bson:"name" json:"name"`
	ImageURL        string    `bson:"image_url" json:"image_url"`
	LogoURL         string    `bson:"logo_url" json:"logo_url"`
	Description     string    `bson:"description,omitempty" json:"description,omitempty"`
	ImageUpdatedAt  time.Time `bson:"image_updated_at" json:"image_updated_at"`
	CollectionOrder int       `bson:"collection_order" json:"collection_order"`
	SuperCollection string    `bson:"super_collection" json:"super_collection"`
	Sections        []Section `bson:"sections" json:"sections"`
	IsDeleted       bool      `bson:"is_deleted" json:"is_deleted"`
	Amount          int       `bson:"amount" json:"amount"`
	Schedule        Schedule  `bson:"schedule" json:"schedule"`
}

func (c MenuCollections) Unique() MenuCollections {

	existCollections := make(map[string]struct{}, len(c))
	result := make(MenuCollections, 0, len(c))

	for _, collection := range c {
		if _, ok := existCollections[collection.ExtID]; ok {
			continue
		}
		if collection.ExtID == "" {
			continue
		}
		result = append(result, collection)
		existCollections[collection.ExtID] = struct{}{}
	}
	return result
}

type Schedule struct {
	ID             string         `bson:"id" json:"id"`
	Name           string         `bson:"name" json:"name"`
	Availabilities []Availability `bson:"availabilities" json:"availabilities"`
}

type Availability struct {
	Day       string     `bson:"day" json:"day"`
	TimeSlots []TimeSlot `bson:"time_slots" json:"time_slots"`
}

type TimeSlot struct {
	Start string `bson:"start" json:"start"`
	End   string `bson:"end" json:"end"`
}

type MenuSuperCollections []MenuSuperCollection

type MenuSuperCollection struct {
	ExtID                string           `bson:"ext_id" json:"ext_id"`
	StarterAppID         string           `bson:"starter_app_id" json:"starter_app_id"`
	Name                 string           `bson:"name" json:"name"`
	SuperCollectionOrder int              `bson:"supercollection_order" json:"position"`
	ImageUrl             string           `bson:"image_url" json:"image_url"`
	Collections          []MenuCollection `bson:"collections,omitempty" json:"collections"`
	Amount               int              `bson:"amount" json:"amount"`
}

type Groups []Group

type Group struct {
	ID              string   `bson:"id" json:"id"`
	Name            string   `bson:"name" json:"name"`
	Description     string   `bson:"description,omitempty" json:"description"`
	Images          []string `bson:"image_links,omitempty" json:"images"`
	ParentGroup     string   `bson:"parent_group,omitempty" json:"parent_group"`
	Order           int      `bson:"order,omitempty" json:"order"`
	InMenu          bool     `bson:"is_included_in_menu,omitempty" json:"in_menu"`
	IsGroupModifier bool     `bson:"is_group_modifier,omitempty" json:"is_group_modifier"`
}

type ProductInformation struct {
	RegulatoryInformation []RegulatoryInformationValues `bson:"regulatory_information" json:"regulatory_information"`
}
type RegulatoryInformationValues struct {
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}

func (g Groups) IsExist(groupID string) bool {
	isExist := false
	for _, group := range g {
		if group.ID == groupID {
			isExist = true
		}
	}

	return isExist
}

func (g Groups) Group(groupID string) Group {
	for _, group := range g {
		if group.ID == groupID {
			return group
		}
	}

	return Group{}
}

func (g Group) Image() string {

	if len(g.Images) != 0 {
		return g.Images[0]
	}

	return ""
}

type UpdateMenu struct {
	Name           string  `bson:"name" json:"name"`
	ExtName        string  `bson:"ext_name,omitempty" json:"ext_name"`
	Description    string  `bson:"description,omitempty" json:"description"`
	Delivery       string  `bson:"delivery" json:"delivery"`
	Comment        Comment `bson:"comment,omitempty" json:"comment"`
	IsActive       bool    `bson:"is_active" json:"is_active"`
	IsDeleted      bool    `bson:"is_deleted" json:"is_deleted"` // ?
	IsSync         bool    `bson:"is_sync" json:"is_sync"`
	SyncAttributes bool    `bson:"sync_attributes" json:"sync_attributes"`

	StopLists        []string             `bson:"stoplist" json:"stoplists"`
	Attributes       Attributes           `bson:"attributes" json:"attributes"`
	AttributesGroups AttributeGroups      `bson:"attributes_groups" json:"attributes_groups"`
	Sections         Sections             `bson:"sections" json:"sections"`
	Products         Products             `bson:"products" json:"products"`
	Combos           []Combo              `bson:"combos" json:"combos"`
	Collections      MenuCollections      `bson:"collections,omitempty" json:"collections"`
	SuperCollections MenuSuperCollections `bson:"super_collections,omitempty" json:"super_collections"`

	AggregatorName string `bson:"aggregator_name,omitempty" json:"aggregator_name"`

	Groups Groups `bson:"groups,omitempty" json:"groups"`

	UpdatedAt          coreModels.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt          coreModels.Time `bson:"created_at" json:"created_at"`
	HasWoltPromo       bool            `bson:"has_wolt_promo" json:"has_wolt_promo"`
	HasMoreThanOneLang bool            `bson:"has_more_than_one_lang" json:"has_more_than_one_lang"`
}

func (m Menu) ToUpdate() UpdateMenu {
	return UpdateMenu{
		Name:             m.Name,
		ExtName:          m.ExtName,
		Delivery:         m.Delivery,
		Comment:          m.Comment,
		Description:      m.Description,
		IsDeleted:        m.IsDeleted,
		IsActive:         m.IsActive,
		IsSync:           m.IsSync,
		SyncAttributes:   m.SyncAttributes,
		StopLists:        m.StopLists,
		Attributes:       m.Attributes,
		AttributesGroups: m.AttributesGroups,
		Sections:         m.Sections,
		Products:         m.Products,
		Combos:           m.Combos,
		Collections:      m.Collections,
		SuperCollections: m.SuperCollections,
		AggregatorName:   m.AggregatorName,
		Groups:           m.Groups,
		UpdatedAt:        m.UpdatedAt,
		CreatedAt:        m.CreatedAt,
		HasWoltPromo:     m.HasWoltPromo,
	}
}

type MenuValidateRequest struct {
	MenuUploadTransaction MenuUploadTransaction
	Menu                  Menu
	OffersBK              []BkOffers
}

type UpdateMenuName struct {
	MenuID   string
	MenuName string
}

type Comment struct {
	Comment string `bson:"comment,omitempty" json:"comment"`
	Active  bool   `bson:"active" json:"active"`
}

type AddLanguageDescriptionRequest struct {
	MenuID    string `json:"menu_id"`
	PosMenuID string `json:"pos_menu_id"`
	ObjectID  string `json:"object_id"`
	Request   LanguageDescription
}

type RegulatoryInformationRequest struct {
	MenuID                string `json:"menu_id"`
	PosMenuID             string `json:"pos_menu_id"`
	ProductID             string `json:"product_id"`
	RegulatoryInformation RegulatoryInformationValues
}

type UpdateFields struct {
	ProductName          bool `json:"product_name"`
	ProductPrice         bool `json:"product_price"`
	ProductDescription   bool `json:"product_description"`
	ProductImage         bool `json:"product_image"`
	AttributeGroupName   bool `json:"attribute_group_name"`
	AttributeGroupMinMax bool `json:"attribute_group_min_max"`
	AttributeName        bool `json:"attribute_name"`
	AttributePrice       bool `json:"attribute_price"`
}

type UpdateFieldsAggregators struct {
	Wolt   bool
	Glovo  bool
	Yandex bool
}

type StoreUpdateMenuByAgg struct {
	RestaurantID string `json:"restaurant_id"`
	Wolt         bool   `json:"wolt"`
	Glovo        bool   `json:"glovo"`
	Yandex       bool   `json:"yandex"`
}

type UpsertMenuFields struct {
	StoreAndAggs []StoreUpdateMenuByAgg `json:"store_and_aggs"`
	Fields       UpdateFields           `json:"fields"`
}
