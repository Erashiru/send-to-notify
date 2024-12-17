package dto

import (
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	models3 "github.com/kwaaka-team/orders-core/core/storecore/models"
	models2 "github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

type AuthenticateData struct {
	ClientID     string `json:"client_id" validate:"required" form:"client_id"`
	ClientSecret string `json:"client_secret" validate:"required" form:"client_secret"`
	GrantType    string `json:"grant_type" validate:"required" form:"grant_type"`
	Scope        string `json:"scope" validate:"required" form:"scope"`
}

type SetCredentialsRequest struct {
	RestID  string `json:"restaurant_id"`
	Service string `json:"service"`
	AuthenticateData
}

type ClientIDResponse struct {
	ID       string `json:"id"`
	ClientID string `json:"client_id"`
}

type CredentianRequest struct {
	RestID        string               `json:"restaurant_id"`
	AuthId        string               `json:"auth_id"`
	SendToPos     bool                 `json:"send_to_pos"`
	MenuUrl       string               `json:"menu_url"`
	IsMarketplace bool                 `json:"is_marketplace"`
	PaymentTypes  models2.PaymentTypes `json:"payment_types"`
}

func ToDTO(req []models.AuthClient) []ClientIDResponse {
	var res = make([]ClientIDResponse, 0, len(req))

	for _, v := range req {
		var temp ClientIDResponse
		temp.ID = v.ID.Hex()
		temp.ClientID = v.ClientID
		res = append(res, temp)
	}
	return res
}

func (req CredentianRequest) ToModel() models.SetCredsToStore {
	var res models.SetCredsToStore

	res.RestID = req.RestID
	res.AuthId = req.AuthId
	res.SendToPos = req.SendToPos
	res.MenuUrl = req.MenuUrl
	res.IsMarketplace = req.IsMarketplace
	res.PaymentTypes = models2.PaymentTypes{
		CASH: models3.PaymentType{
			PaymentTypeID:            req.PaymentTypes.CASH.PaymentTypeID,
			PaymentTypeKind:          req.PaymentTypes.CASH.PaymentTypeKind,
			PromotionPaymentTypeID:   req.PaymentTypes.CASH.PromotionPaymentTypeID,
			OrderType:                req.PaymentTypes.CASH.OrderType,
			OrderTypeService:         req.PaymentTypes.CASH.OrderTypeService,
			OrderTypeForVirtualStore: req.PaymentTypes.CASH.OrderTypeForVirtualStore,
			IsProcessedExternally:    req.PaymentTypes.CASH.IsProcessedExternally,
		},
		DELAYED: models3.PaymentType{
			PaymentTypeID:            req.PaymentTypes.DELAYED.PaymentTypeID,
			PaymentTypeKind:          req.PaymentTypes.DELAYED.PaymentTypeKind,
			PromotionPaymentTypeID:   req.PaymentTypes.DELAYED.PromotionPaymentTypeID,
			OrderType:                req.PaymentTypes.DELAYED.OrderType,
			OrderTypeService:         req.PaymentTypes.DELAYED.OrderTypeService,
			OrderTypeForVirtualStore: req.PaymentTypes.DELAYED.OrderTypeForVirtualStore,
			IsProcessedExternally:    req.PaymentTypes.DELAYED.IsProcessedExternally,
		},
	}

	return res
}
