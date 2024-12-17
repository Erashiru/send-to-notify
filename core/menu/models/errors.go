package models

import "github.com/pkg/errors"

var (
	ErrNotFound                 = errors.New("not found")
	ErrNoAnyStopListTransaction = errors.New("not updated in any delivery service")
)

type ProductStatus struct {
	Name              string
	IsDeleted         bool
	IsIncludedInMenu  bool
	AttributeGroupIDs []string
	Defaults          []string
}

type AttributeStatus struct {
	IsDeleted        bool
	IsIncludedInMenu bool
}

type AttributeGroupStatus struct {
	Position   int
	Name       string
	Min        int
	Max        int
	Attributes []string
}

type ValidateReport struct {
	ID                   string                 `json:"restaurant_id"`
	RestaurantName       string                 `json:"name"`
	Delivery             string                 `json:"delivery"`
	MenuID               string                 `json:"menu_id"`
	Products             []ProductReport        `json:"products"`
	AttributeGroups      []AttributeGroupReport `json:"attribute_groups"`
	AttributeGroupMinMax []ProductMinMaxReport  `json:"attribute_group_min_max"`
}

type ProductReport struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Position int      `json:"position"`
	Status   string   `json:"status"`
	Solution []string `json:"solution"`
}

type AttributeReport struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Position int      `json:"position"`
	Status   string   `json:"status"`
	Solution []string `json:"solution"`
}

type AttributeGroupReport struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Min        int               `json:"min"`
	Max        int               `json:"max"`
	Position   int               `json:"position"`
	Status     string            `json:"status"`
	Solution   []string          `json:"solution"`
	Attributes []AttributeReport `json:"not_exist_attributes"`
	Products   []ProductReport   `json:"in_products"`
}

type ProductMinMaxReport struct {
	ID                     string                  `json:"product_id"`
	Name                   string                  `json:"product_name"`
	Position               int                     `json:"position"`
	AttributeMinMaxReports []AttributeMinMaxReport `json:"attribute_groups"`
}

type AttributeMinMaxReport struct {
	ID                      string   `json:"pos_group_id"`
	Name                    string   `json:"pos_group_name"`
	Min                     int      `json:"pos_min"`
	Max                     int      `json:"pos_max"`
	PosAttributeName        []string `json:"pos_attributes_name"`
	AggregatorAttributeName []string `json:"aggregator_attributes_name"`
	CurrentMin              int      `json:"aggregator_min"`
}
