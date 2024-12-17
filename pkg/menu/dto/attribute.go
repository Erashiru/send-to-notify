package dto

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
)

type AttributesUpdate struct {
	Attributes      models.Attributes
	AttributeGroups models.AttributeGroups
	Total           int
}

type AttributeSelector struct {
	Page   int64
	Limit  int64
	ID     string
	MenuID string
	Sorting
}

type Sorting struct {
	Param     string
	Direction int8
}

type UpdateMenuName struct {
	MenuID   string
	MenuName string
}

func (um *UpdateMenuName) ToModel() models.UpdateMenuName {
	var res models.UpdateMenuName
	res.MenuID = um.MenuID
	res.MenuName = um.MenuName
	return res
}
