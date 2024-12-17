package dto

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/core/menu/models"
)

type MenuGroupRequest struct {
	MenuID  string
	StoreID string
	Token   string
	GroupID string
}

type MenuGroups []MenuGroup

type MenuGroup struct {
	ID              string   `json:"id" bson:"id"`
	Name            string   `json:"name" bson:"name"`
	Description     string   `json:"description" bson:"description,omitempty"`
	Images          []string `json:"imageLinks" bson:"image_links,omitempty"`
	ParentGroup     string   `json:"parentGroup" bson:"parent_group,omitempty"`
	Order           int      `json:"order" bson:"order,omitempty"`
	InMenu          bool     `json:"isIncludedInMenu" bson:"is_included_in_menu,omitempty"`
	IsGroupModifier bool     `json:"isGroupModifier" bson:"is_group_modifier,omitempty"`
}

func FromMenuGroups(req models.Groups) MenuGroups {
	res := make(MenuGroups, 0, len(req))

	for _, group := range req {
		res = append(res, fromMenuGroup(group))
	}
	return res
}

func fromMenuGroup(req models.Group) MenuGroup {
	return MenuGroup{
		ID:              req.ID,
		Name:            req.Name,
		Description:     req.Description,
		Images:          req.Images,
		ParentGroup:     req.ParentGroup,
		Order:           req.Order,
		InMenu:          req.InMenu,
		IsGroupModifier: req.IsGroupModifier,
	}
}

type MenuUploadRequest struct {
	StoreId      string
	MenuId       string
	DeliveryName string
	Sv3          *s3.S3
	UserRole     string
	UserName     string
}

type MenuUploadVerifyRequest struct {
	TransactionId string
}

type GetMenuUploadTransactions struct {
	DeliveryService string
	StoreId         string
	Pagination
}

type MenuValidateRequest struct {
	ID      string
	StoreID string
}

type MenuSuperCollections []MenuSuperCollection

type MenuSuperCollection struct {
	ExtId                string `json:"ext_id"`
	Name                 string `json:"name"`
	ImgUrl               string `json:"img_url"`
	SuperCollectionOrder int    `json:"order"`
}

type Pagination struct {
	Page  int64
	Limit int64
}
