package dto

import (
	models2 "github.com/kwaaka-team/orders-core/core/storecore/models"
	glovoModels "github.com/kwaaka-team/orders-core/pkg/glovo/clients/dto"
	"time"
)

type UpdateStore struct {
	ID                *string                       `json:"id"`
	Token             *string                       `json:"token"`
	Name              *string                       `json:"name"`
	MenuID            *string                       `json:"menu_id"`
	PosType           *string                       `json:"pos_type"`
	LegalEntityId     *string                       `json:"legal_entity_id"`
	AccountManagerId  *string                       `json:"account_manager_id"`
	SalesManagerId    *string                       `json:"sales_manager_id"`
	Address           *UpdateStoreAddress           `json:"address"`
	Glovo             *UpdateStoreGlovoConfig       `json:"glovo"`
	Wolt              *UpdateStoreWoltConfig        `json:"wolt"`
	KwaakaAdmin       *UpdateStoreKwaakaAdminConfig `json:"kwaaka_admin"`
	Chocofood         *UpdateStoreChocofoodConfig   `json:"chocofood"`
	RKeeper           *UpdateStoreRKeeperConfig     `json:"rkeeper"`
	RKeeper7XML       *UpdateStoreRKeeper7XMLConfig `json:"rkeeper7_xml"`
	Paloma            *UpdateStorePalomaConfig      `json:"paloma"`
	Yandex            *UpdateStoreYandexConfig      `json:"yandex"`
	External          []UpdateStoreExternalConfig   `json:"external"`
	QRMenu            *UpdateStoreQRMenuConfig      `json:"qr_menu"`
	MoySklad          *UpdateStoreMoySkladConfig    `json:"moysklad"`
	IikoCloud         *UpdateStoreIikoConfig        `json:"iiko_cloud"`
	Telegram          *UpdateStoreTelegramConfig    `json:"telegram"`
	Notification      *UpdateNotification           `json:"notification"`
	CallCenter        *UpdateStoreCallcenter        `json:"callcenter"`
	Delivery          []UpdateStoreDelivery         `json:"delivery"`
	Menus             []UpdateStoreDSMenu           `json:"menus"`
	Contacts          []UpdateContact               `json:"contacts"`
	Links             []UpdateLinks                 `json:"links"`
	Settings          *UpdateSettings               `json:"settings"`
	IntegrationDate   *time.Time                    `json:"integration_date"`
	UpdatedAt         *time.Time                    `json:"updated_at"`
	CreatedAt         *time.Time                    `json:"created_at"`
	Payments          []UpdatePayment               `json:"payments"`
	RestaurantGroupID *string                       `json:"restaurant_group_id"`
	BillParameter     *UpdateParameters             `json:"bill_parameter"`
	StoreSchedule     *StoreSchedule                `json:"store_schedule"`
	IsDeleted         *bool                         `json:"is_deleted"`
	SocialMediaLinks  []UpdateSocialMediaLinks      `json:"social_media_links"`
	CompensationCount *int                          `json:"compensation_count"`
}

type UpdateLinks struct {
	Name      string `json:"name"`
	Url       string `json:"url" bson:"url"`
	ImageLink string `json:"image_link" bson:"image_link"`
}

type UpdateSocialMediaLinks struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	URL  string `bson:"url" json:"url"`
	Logo string `bson:"logo,omitempty" json:"logo,omitempty"`
}

type UpdateContact struct {
	FullName string `json:"full_name"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Comment  string `json:"comment"`
}

type StoreSchedule struct {
	GlovoSchedule  *glovoModels.StoreScheduleResponse `json:"glovo_schedule,omitempty"`
	WoltSchedule   *glovoModels.StoreScheduleResponse `json:"wolt_schedule,omitempty"`
	DirectSchedule []DirectSchedule                   `json:"direct_schedule"`
}

type UpdateStoreKwaakaAdminConfig struct {
	IsIntegrated *bool    `json:"is_integrated"`
	IsActive     *bool    `json:"is_active"`
	CookingTime  *int32   `json:"cooking_time"`
	StoreID      []string `json:"store_id"`
	SendToPos    *bool    `json:"send_to_pos"`
}

type UpdateStoreRKeeper7XMLConfig struct {
	Domain              *string `json:"domain"`
	Username            *string `json:"username"`
	Password            *string `json:"password"`
	UCSUsername         *string `json:"ucs_username"`
	UCSPassword         *string `json:"ucs_password"`
	Token               *string `json:"token"`
	ObjectID            *string `json:"object_id"`
	Anchor              *string `json:"anchor"`
	LicenseInstanceGUID *string `json:"license_instance_guid"`
	SeqNumber           *int    `json:"seq_number"`
}

type Parameters struct {
	IsActive       bool           `json:"is_active"`
	BillParameters BillParameters `json:"bill_parameters"`
}
type BillParameters struct {
	AddPaymentType     bool `json:"add_payment_type"`
	AddOrderCode       bool `json:"add_order_code"`
	AddComments        bool `json:"add-comments"`
	AddDelivery        bool `json:"add_delivery"`
	AddAddress         bool `json:"add_address"`
	AddQuantityPersons bool `json:"add_quantity_persons"`
}
type UpdateParameters struct {
	IsActive             *bool                `json:"is_active"`
	UpdateBillParameters UpdateBillParameters `json:"update_bill_parameters"`
}

type UpdateBillParameters struct {
	AddPaymentType     *bool `json:"add_payment_type"`
	AddOrderCode       *bool `json:"add_order_code"`
	AddComments        *bool `json:"add-comments"`
	AddDelivery        *bool `json:"add_delivery"`
	AddAddress         *bool `json:"add_address"`
	AddQuantityPersons *bool `json:"add_quantity_persons"`
}

type UpdateStoreAddress struct {
	City        *string            `json:"city"`
	Street      *string            `json:"street"`
	Entrance    *string            `json:"entrance,omitempty"`
	Coordinates *UpdateCoordinates `json:"coordinates"`
}

type UpdateCoordinates struct {
	Longitude *float64 `json:"longitude,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
}

type UpdateStoreGlovoConfig struct {
	StoreID                            []string                          `json:"store_id"`
	MenuUrl                            *string                           `json:"menu_url"`
	SendToPos                          *bool                             `json:"send_to_pos"`
	IsMarketplace                      *bool                             `json:"is_marketplace"`
	PaymentTypes                       *UpdateDeliveryServicePaymentType `json:"payment_types"`
	PurchaseTypes                      *UpdatePurchaseTypes              `json:"purchase_types"`
	AdditionalPreparationTimeInMinutes *int                              `json:"additional_preparation_time_in_minutes"`
}

type UpdateDeliveryServicePaymentType struct {
	CASH    *UpdateIIKOPaymentType `json:"CASH"`
	DELAYED *UpdateIIKOPaymentType `json:"DELAYED"`
}

type UpdateIIKOPaymentType struct {
	IikoPaymentTypeID        *string `json:"iiko_payment_type_id"`
	IikoPaymentTypeKind      *string `json:"iiko_payment_type_kind"`
	OrderType                *string `json:"order_type"`
	PromotionPaymentTypeID   *string `json:"promotion_payment_type_id"`
	OrderTypeService         *string `json:"order_type_service"`
	OrderTypeForVirtualStore *string `json:"order_type_for_virtual_store"`
	IsProcessedExternally    *bool   `json:"is_processed_externally"`
}

type UpdateStoreWoltConfig struct {
	StoreID               []string                          `json:"store_id"`
	MenuUsername          *string                           `json:"menu_username"`
	MenuPassword          *string                           `json:"menu_password"`
	ApiKey                *string                           `json:"api_key"`
	AdjustedPickupMinutes *int                              `json:"adjusted_pickup_minutes"`
	MenuUrl               *string                           `json:"menu_url"`
	SendToPos             *bool                             `json:"send_to_pos"`
	IsMarketplace         *bool                             `json:"is_marketplace"`
	PaymentTypes          *UpdateDeliveryServicePaymentType `json:"payment_types"`
	PurchaseTypes         *UpdatePurchaseTypes              `json:"purchase_types"`
	IgnoreStatusUpdate    *bool                             `json:"ignore_status_update"`
	AutoAcceptOn          *bool                             `json:"auto_accept_on"`
}

type UpdatePurchaseTypes struct {
	Instant  []UpdateStatus `json:"instant"`
	Preorder []UpdateStatus `json:"preorder"`
	TakeAway []UpdateStatus `json:"takeaway"`
}

type UpdateStatus struct {
	PosStatus *string `json:"pos_status"`
	Status    *string `json:"status"`
}

type UpdateStoreChocofoodConfig struct {
	StoreID       []string                          `json:"store_id"`
	MenuUrl       *string                           `json:"menu_url"`
	SendToPos     *bool                             `json:"send_to_pos"`
	IsMarketplace *bool                             `json:"is_marketplace"`
	PaymentTypes  *UpdateDeliveryServicePaymentType `json:"payment_types"`
}

type UpdateStoreRKeeperConfig struct {
	ObjectId *int `json:"object_id"`
}

type UpdateStorePalomaConfig struct {
	PointID *string `json:"point_id"`
	ApiKey  *string `json:"api_key"`
}

type UpdateStoreYandexConfig struct {
	StoreID       []string                          `json:"store_id"`
	MenuUrl       *string                           `json:"menu_url"`
	SendToPos     *bool                             `json:"send_to_pos"`
	IsMarketplace *bool                             `json:"is_marketplace"`
	PaymentTypes  *UpdateDeliveryServicePaymentType `json:"payment_types"`
}

type UpdateStoreExternalConfig struct {
	StoreID                  []string                          `json:"store_id"`
	Type                     *string                           `json:"type"`
	MenuUrl                  *string                           `json:"menu_url"`
	SendToPos                *bool                             `json:"send_to_pos"`
	IsMarketplace            *bool                             `json:"is_marketplace"`
	PaymentTypes             *UpdateDeliveryServicePaymentType `json:"payment_types"`
	ClientSecret             *string                           `json:"client_secret"`
	WebhookURL               *string                           `json:"webhook_url"`
	AuthToken                *string                           `json:"auth_token"`
	WebhookProductStoplist   *string                           `json:"webhook_product_stoplist"`
	WebhookAttributeStoplist *string                           `json:"webhook_attribute_stoplist"`
}

type UpdateStoreQRMenuConfig struct {
	StoreID               []string                          `json:"store_id"`
	URL                   *string                           `json:"url"`
	IsIntegrated          *bool                             `json:"is_integrated"`
	PaymentTypes          *UpdateDeliveryServicePaymentType `json:"payment_types"`
	Hash                  *string                           `json:"hash"`
	CookingTime           *int                              `json:"cooking_time"`
	DeliveryTime          *int                              `json:"delivery_time"`
	NoTable               *bool                             `json:"no_table"`
	Theme                 *string                           `json:"theme"`
	IsMarketplace         *bool                             `json:"is_marketplace"`
	SendToPos             *bool                             `json:"send_to_pos"`
	IgnoreStatusUpdate    *bool                             `json:"ignore_status_update"`
	AdjustedPickupMinutes *int                              `json:"adjusted_pickup_minutes" bson:"adjusted_pickup_minutes"`
	BusyMode              *bool                             `json:"busy_mode" bson:"busy_mode"`
}

type UpdateStoreMoySkladConfig struct {
	UserName       *string                           `json:"username"`
	Password       *string                           `json:"password"`
	OrderID        *string                           `json:"order_id"`
	OrganizationID *string                           `json:"organization_id"`
	Status         *UpdateMoySkladStatus             `json:"status"`
	SendToPos      *bool                             `json:"send_to_pos"`
	IsMarketPlace  *bool                             `json:"is_marketplace"`
	PaymentTypes   *UpdateDeliveryServicePaymentType `json:"payment_types"`
}

type UpdateMoySkladStatus struct {
	ID         *string `json:"id"`
	Name       *string `json:"name"`
	StatusType *string `json:"status_type"`
}

type UpdateStoreIikoConfig struct {
	OrganizationID       *string `json:"organization_id"`
	TerminalID           *string `json:"terminal_id"`
	Key                  *string `json:"key"`
	StopListByBalance    *bool   `json:"stoplist_by_balance,omitempty"` // temporary field for Traveler`s menu
	StopListBalanceLimit *int    `json:"stoplist_balance_limit"`
	IsExternalMenu       *bool   `json:"is_external_menu"`
	ExternalMenuID       *string `json:"external_menu_id"`
	PriceCategory        *string `json:"price_category"`
}

type UpdateStoreTelegramConfig struct {
	GroupChatID    *string `json:"group_chat_id"`
	StopListChatID *string `json:"stoplist_chat_id"`
	CancelChatID   *string `json:"cancel_chat_id"`
}

type UpdateNotification struct {
	Whatsapp *UpdateWhatsapp `json:"whatsapp"`
}

type UpdateWhatsapp struct {
	Receivers []UpdateWhatsappReceiver `json:"receivers"`
}

type UpdateWhatsappReceiver struct {
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phone_number"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateStoreCallcenter struct {
	Restaurants []string `json:"restaurants"`
}

type UpdateStoreDelivery struct {
	ID       *string `json:"id"`
	Code     *string `json:"code"`
	Price    *int    `json:"price"`
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
	Service  *string `json:"service"`
}

type UpdateStoreDSMenu struct {
	MenuID         *string    `json:"menu_id"`
	Name           *string    `json:"name"`
	IsActive       *bool      `json:"is_active"`
	IsDeleted      *bool      `json:"is_deleted"`
	IsSync         *bool      `json:"is_sync"`
	SyncAttributes *bool      `json:"sync_attributes"`
	Delivery       *string    `json:"delivery"`
	Timestamp      *int       `json:"timestamp"`
	UpdatedAt      *time.Time `json:"updated_at"`
	HasWoltPromo   *bool      `json:"has_wolt_promo"`
	CreationSource *string    `json:"creation_source"`
}

type UpdateSettings struct {
	TimeZone               *UpdateTimeZone               `json:"timezone"`
	Currency               *string                       `json:"currency"`
	LanguageCode           *string                       `json:"language_code"`
	SendToPos              *bool                         `json:"send_to_pos"`
	IsMarketplace          *bool                         `json:"is_marketplace"`
	PriceSource            *string                       `json:"price_source"`
	OrderDestination       *OrderDestination             `json:"order_destination"`
	IsAutoUpdate           *bool                         `json:"is_auto_update"`
	IsDeleted              *bool                         `json:"is_deleted"`
	Group                  *UpdateMenuGroup              `json:"group"`
	HasVirtualStore        *bool                         `json:"has_virtual_store"`
	StopListClosingActions []UpdateStopListClosingAction `json:"stoplist_closing_actions"`
	Email                  *string                       `json:"email"`
}

type UpdateStopListClosingAction struct {
	RestaurantID *string `json:"restaurant_id"`
	Opening      *string `json:"opening"` // format 15:04:05
	Closing      *string `json:"closing"` // format 15:04:05
	Status       *string `json:"status"`
}

type UpdateTimeZone struct {
	TZ        *string  `json:"tz"`
	UTCOffset *float64 `json:"utc_offset"`
}

type UpdateMenuGroup struct {
	MenuID  *string `json:"menu_id"`
	GroupID *string `json:"group_id"`
}

type UpdatePayment struct {
	Type     *string `json:"type"`
	Service  *string `json:"service"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func toUpdatePurchaseTypes(req *UpdatePurchaseTypes) *models2.UpdatePurchaseTypes {
	var purchaseTypes models2.UpdatePurchaseTypes

	for i := 0; i < len(req.Instant); i++ {
		purchaseTypes.Instant = append(purchaseTypes.Instant, models2.UpdateStatus{
			PosStatus: req.Instant[i].PosStatus,
			Status:    req.Instant[i].Status,
		})
	}

	for i := 0; i < len(req.Preorder); i++ {
		purchaseTypes.Preorder = append(purchaseTypes.Preorder, models2.UpdateStatus{
			PosStatus: req.Preorder[i].PosStatus,
			Status:    req.Preorder[i].Status,
		})
	}

	for i := 0; i < len(req.TakeAway); i++ {
		purchaseTypes.TakeAway = append(purchaseTypes.TakeAway, models2.UpdateStatus{
			PosStatus: req.TakeAway[i].PosStatus,
			Status:    req.TakeAway[i].Status,
		})
	}

	return &purchaseTypes
}

func toUpdatePaymentTypes(req *UpdateDeliveryServicePaymentType) *models2.UpdateDeliveryServicePaymentType {
	var paymentType models2.UpdateDeliveryServicePaymentType

	if req.CASH != nil {
		paymentType.CASH = &models2.UpdateIIKOPaymentType{
			IikoPaymentTypeID:        req.CASH.IikoPaymentTypeID,
			IikoPaymentTypeKind:      req.CASH.IikoPaymentTypeKind,
			PromotionPaymentTypeID:   req.CASH.PromotionPaymentTypeID,
			OrderType:                req.CASH.OrderType,
			OrderTypeService:         req.CASH.OrderTypeService,
			OrderTypeForVirtualStore: req.CASH.OrderTypeForVirtualStore,
			IsProcessedExternally:    req.CASH.IsProcessedExternally,
		}
	}

	if req.DELAYED != nil {
		paymentType.DELAYED = &models2.UpdateIIKOPaymentType{
			IikoPaymentTypeID:        req.DELAYED.IikoPaymentTypeID,
			IikoPaymentTypeKind:      req.DELAYED.IikoPaymentTypeKind,
			PromotionPaymentTypeID:   req.DELAYED.PromotionPaymentTypeID,
			OrderType:                req.DELAYED.OrderType,
			OrderTypeService:         req.DELAYED.OrderTypeService,
			OrderTypeForVirtualStore: req.DELAYED.OrderTypeForVirtualStore,
			IsProcessedExternally:    req.DELAYED.IsProcessedExternally,
		}
	}

	return &paymentType
}

func toUpdateExternalConfig(req []UpdateStoreExternalConfig) []models2.UpdateStoreExternalConfig {
	var external = make([]models2.UpdateStoreExternalConfig, 0, len(req))

	for i := 0; i < len(req); i++ {
		config := models2.UpdateStoreExternalConfig{
			StoreID:                  req[i].StoreID,
			Type:                     req[i].Type,
			MenuUrl:                  req[i].MenuUrl,
			SendToPos:                req[i].SendToPos,
			IsMarketplace:            req[i].IsMarketplace,
			ClientSecret:             req[i].ClientSecret,
			AuthToken:                req[i].AuthToken,
			WebhookURL:               req[i].WebhookURL,
			WebhookProductStoplist:   req[i].WebhookProductStoplist,
			WebhookAttributeStoplist: req[i].WebhookAttributeStoplist,
		}

		if req[i].PaymentTypes != nil {
			config.PaymentTypes = toUpdatePaymentTypes(req[i].PaymentTypes)
		}

		external = append(external, config)
	}

	return external
}

func (s UpdateStore) ToModel() models2.UpdateStore {
	store := models2.UpdateStore{
		ID:                s.ID,
		Token:             s.Token,
		Name:              s.Name,
		MenuID:            s.MenuID,
		PosType:           s.PosType,
		IntegrationDate:   s.IntegrationDate,
		UpdatedAt:         s.UpdatedAt,
		CreatedAt:         s.CreatedAt,
		RestaurantGroupID: s.RestaurantGroupID,
		IsDeleted:         s.IsDeleted,
		LegalEntityId:     s.LegalEntityId,
		AccountManagerId:  s.AccountManagerId,
		SalesManagerId:    s.SalesManagerId,
		CompensationCount: s.CompensationCount,
	}

	if s.Address != nil {
		store.Address = &models2.UpdateStoreAddress{
			City:     s.Address.City,
			Street:   s.Address.Street,
			Entrance: s.Address.Entrance,
		}
		if store.Address.UpdateCoordinates != nil {
			store.Address = &models2.UpdateStoreAddress{
				UpdateCoordinates: &models2.UpdateCoordinates{
					Longitude: store.Address.UpdateCoordinates.Longitude,
					Latitude:  store.Address.UpdateCoordinates.Latitude,
				},
			}
		}
	}

	if s.RKeeper7XML != nil {
		rkeeper7XmlConfig := &models2.UpdateStoreRKeeper7XMLConfig{
			Domain:              s.RKeeper7XML.Domain,
			Username:            s.RKeeper7XML.Username,
			Password:            s.RKeeper7XML.Password,
			UCSUsername:         s.RKeeper7XML.UCSUsername,
			UCSPassword:         s.RKeeper7XML.UCSPassword,
			Token:               s.RKeeper7XML.Token,
			ObjectID:            s.RKeeper7XML.ObjectID,
			Anchor:              s.RKeeper7XML.Anchor,
			LicenseInstanceGUID: s.RKeeper7XML.LicenseInstanceGUID,
			SeqNumber:           s.RKeeper7XML.SeqNumber,
		}

		store.RKeeper7XML = rkeeper7XmlConfig
	}

	if s.Glovo != nil {
		glovoCfg := &models2.UpdateStoreGlovoConfig{
			StoreID:                            s.Glovo.StoreID,
			MenuUrl:                            s.Glovo.MenuUrl,
			SendToPos:                          s.Glovo.SendToPos,
			IsMarketplace:                      s.Glovo.IsMarketplace,
			AdditionalPreparationTimeInMinutes: s.Glovo.AdditionalPreparationTimeInMinutes,
		}

		if s.Glovo.PaymentTypes != nil {
			glovoCfg.PaymentTypes = toUpdatePaymentTypes(s.Glovo.PaymentTypes)
		}

		if s.Glovo.PurchaseTypes != nil {
			glovoCfg.PurchaseTypes = toUpdatePurchaseTypes(s.Glovo.PurchaseTypes)
		}

		store.Glovo = glovoCfg
	}

	if s.Wolt != nil {
		woltCfg := &models2.UpdateStoreWoltConfig{
			StoreID:               s.Wolt.StoreID,
			MenuUsername:          s.Wolt.MenuUsername,
			MenuPassword:          s.Wolt.MenuPassword,
			ApiKey:                s.Wolt.ApiKey,
			AdjustedPickupMinutes: s.Wolt.AdjustedPickupMinutes,
			MenuUrl:               s.Wolt.MenuUrl,
			SendToPos:             s.Wolt.SendToPos,
			IsMarketplace:         s.Wolt.IsMarketplace,
			IgnoreStatusUpdate:    s.Wolt.IgnoreStatusUpdate,
			AutoAcceptOn:          s.Wolt.AutoAcceptOn,
		}

		if s.Wolt.PaymentTypes != nil {
			woltCfg.PaymentTypes = toUpdatePaymentTypes(s.Wolt.PaymentTypes)
		}

		if s.Wolt.PurchaseTypes != nil {
			woltCfg.PurchaseTypes = toUpdatePurchaseTypes(s.Wolt.PurchaseTypes)
		}

		store.Wolt = woltCfg
	}

	if s.Chocofood != nil {
		chocoCfg := &models2.UpdateStoreChocofoodConfig{
			StoreID:       s.Chocofood.StoreID,
			MenuUrl:       s.Chocofood.MenuUrl,
			SendToPos:     s.Chocofood.SendToPos,
			IsMarketplace: s.Chocofood.IsMarketplace,
		}

		if s.Chocofood.PaymentTypes != nil {
			chocoCfg.PaymentTypes = toUpdatePaymentTypes(s.Chocofood.PaymentTypes)
		}

		store.Chocofood = chocoCfg
	}

	if s.RKeeper != nil {
		rkeeperCfg := &models2.UpdateStoreRKeeperConfig{
			ObjectId: s.RKeeper.ObjectId,
		}

		store.RKeeper = rkeeperCfg
	}

	if s.Paloma != nil {
		palomaCfg := &models2.UpdateStorePalomaConfig{
			PointID: s.Paloma.PointID,
			ApiKey:  s.Paloma.ApiKey,
		}

		store.Paloma = palomaCfg
	}

	if s.External != nil {
		store.External = toUpdateExternalConfig(s.External)
	}

	if s.QRMenu != nil {
		qrCfg := &models2.UpdateStoreQRMenuConfig{
			StoreID:               s.QRMenu.StoreID,
			URL:                   s.QRMenu.URL,
			IsIntegrated:          s.QRMenu.IsIntegrated,
			Hash:                  s.QRMenu.Hash,
			CookingTime:           s.QRMenu.CookingTime,
			DeliveryTime:          s.QRMenu.DeliveryTime,
			NoTable:               s.QRMenu.NoTable,
			Theme:                 s.QRMenu.Theme,
			IsMarketplace:         s.QRMenu.IsMarketplace,
			SendToPos:             s.QRMenu.SendToPos,
			IgnoreStatusUpdate:    s.QRMenu.IgnoreStatusUpdate,
			BusyMode:              s.QRMenu.BusyMode,
			AdjustedPickupMinutes: s.QRMenu.AdjustedPickupMinutes,
		}

		if s.QRMenu.PaymentTypes != nil {
			qrCfg.PaymentTypes = toUpdatePaymentTypes(s.QRMenu.PaymentTypes)
		}
		store.QRMenu = qrCfg
	}

	if s.KwaakaAdmin != nil {
		kwaakaAdminCfg := &models2.UpdateStoreKwaakaAdminConfig{
			IsIntegrated: s.KwaakaAdmin.IsIntegrated,
			IsActive:     s.KwaakaAdmin.IsActive,
			CookingTime:  s.KwaakaAdmin.CookingTime,
			StoreID:      s.KwaakaAdmin.StoreID,
			SendToPos:    s.KwaakaAdmin.SendToPos,
		}

		store.KwaakaAdmin = kwaakaAdminCfg
	}

	if s.MoySklad != nil {
		moyskladCfg := &models2.UpdateStoreMoySkladConfig{
			UserName:       s.MoySklad.UserName,
			Password:       s.MoySklad.Password,
			OrderID:        s.MoySklad.OrderID,
			OrganizationID: s.MoySklad.OrganizationID,
			SendToPos:      s.MoySklad.SendToPos,
			IsMarketPlace:  s.MoySklad.IsMarketPlace,
		}

		if s.MoySklad.PaymentTypes != nil {
			moyskladCfg.PaymentTypes = toUpdatePaymentTypes(s.MoySklad.PaymentTypes)
		}

		if s.MoySklad.Status != nil {
			moyskladCfg.Status = &models2.UpdateMoySkladStatus{
				ID:         s.MoySklad.Status.ID,
				Name:       s.MoySklad.Status.Name,
				StatusType: s.MoySklad.Status.StatusType,
			}
		}

		store.MoySklad = moyskladCfg
	}

	if s.IikoCloud != nil {
		iikoCfg := &models2.UpdateStoreIikoConfig{
			OrganizationID:       s.IikoCloud.OrganizationID,
			TerminalID:           s.IikoCloud.TerminalID,
			Key:                  s.IikoCloud.Key,
			StopListByBalance:    s.IikoCloud.StopListByBalance,
			StopListBalanceLimit: s.IikoCloud.StopListBalanceLimit,
			IsExternalMenu:       s.IikoCloud.IsExternalMenu,
			ExternalMenuID:       s.IikoCloud.ExternalMenuID,
			PriceCategory:        s.IikoCloud.PriceCategory,
		}

		store.IikoCloud = iikoCfg
	}

	if s.Telegram != nil {
		telegramCfg := &models2.UpdateStoreTelegramConfig{
			GroupChatID:    s.Telegram.GroupChatID,
			StopListChatID: s.Telegram.StopListChatID,
			CancelChatID:   s.Telegram.CancelChatID,
		}

		store.Telegram = telegramCfg
	}

	if s.Notification != nil {
		notificationCfg := &models2.UpdateNotification{}

		if s.Notification.Whatsapp != nil {
			whatsapp := &models2.UpdateWhatsapp{}

			if s.Notification.Whatsapp.Receivers != nil {
				receivers := make([]models2.UpdateWhatsappReceiver, 0, len(s.Notification.Whatsapp.Receivers))

				for i := 0; i < len(s.Notification.Whatsapp.Receivers); i++ {
					receivers = append(receivers, models2.UpdateWhatsappReceiver{
						Name:        s.Notification.Whatsapp.Receivers[i].Name,
						PhoneNumber: s.Notification.Whatsapp.Receivers[i].PhoneNumber,
						IsActive:    s.Notification.Whatsapp.Receivers[i].IsActive,
					})
				}

				whatsapp.Receivers = receivers
			}

			notificationCfg.Whatsapp = whatsapp
		}

		store.Notification = notificationCfg
	}

	if s.CallCenter != nil {
		callCenterCfg := &models2.UpdateStoreCallcenter{}

		if s.CallCenter.Restaurants != nil {
			callCenterCfg.Restaurants = s.CallCenter.Restaurants
		}

		store.CallCenter = callCenterCfg
	}

	if s.Delivery != nil {
		var deliveries = make([]models2.UpdateStoreDelivery, 0, len(s.Delivery))

		for i := 0; i < len(s.Delivery); i++ {
			deliveries = append(deliveries, models2.UpdateStoreDelivery{
				ID:       s.Delivery[i].ID,
				Code:     s.Delivery[i].Code,
				Price:    s.Delivery[i].Price,
				Name:     s.Delivery[i].Name,
				IsActive: s.Delivery[i].IsActive,
				Service:  s.Delivery[i].Service,
			})
		}

		store.Delivery = deliveries
	}

	if s.Links != nil {
		links := make([]models2.UpdateLinks, 0, len(s.Links))

		for i := 0; i < len(s.Links); i++ {
			links = append(links, models2.UpdateLinks{
				Name:      s.Links[i].Name,
				Url:       s.Links[i].Url,
				ImageLink: s.Links[i].ImageLink,
			})
		}

		store.Links = links
	}

	if s.SocialMediaLinks != nil {
		socialLinks := make([]models2.UpdateSocialMediaLinks, 0, len(s.SocialMediaLinks))
		for i := 0; i < len(s.SocialMediaLinks); i++ {
			socialLinks = append(socialLinks, models2.UpdateSocialMediaLinks{
				Name: s.SocialMediaLinks[i].Name,
				URL:  s.SocialMediaLinks[i].URL,
				Logo: s.SocialMediaLinks[i].Logo,
			})
		}
		store.SocialMediaLinks = socialLinks
	}

	if s.Contacts != nil {
		contacts := make([]models2.UpdateContact, 0, len(s.Contacts))

		for _, contact := range s.Contacts {
			contacts = append(contacts, models2.UpdateContact{
				FullName: contact.FullName,
				Position: contact.Position,
				Phone:    contact.Phone,
				Comment:  contact.Comment,
			})
		}

		store.Contacts = contacts
	}

	if s.Menus != nil {
		var menus = make([]models2.UpdateStoreDSMenu, 0, len(s.Menus))

		for i := 0; i < len(s.Menus); i++ {
			menus = append(menus, models2.UpdateStoreDSMenu{
				MenuID:         s.Menus[i].MenuID,
				Name:           s.Menus[i].Name,
				IsActive:       s.Menus[i].IsActive,
				IsDeleted:      s.Menus[i].IsDeleted,
				IsSync:         s.Menus[i].IsSync,
				SyncAttributes: s.Menus[i].SyncAttributes,
				Delivery:       s.Menus[i].Delivery,
				Timestamp:      s.Menus[i].Timestamp,
				UpdatedAt:      s.Menus[i].UpdatedAt,
				HasWoltPromo:   s.Menus[i].HasWoltPromo,
				CreationSource: s.Menus[i].CreationSource,
			})
		}

		store.Menus = menus
	}

	if s.Settings != nil {
		var settings = &models2.UpdateSettings{
			Currency:        s.Settings.Currency,
			LanguageCode:    s.Settings.LanguageCode,
			SendToPos:       s.Settings.SendToPos,
			IsMarketplace:   s.Settings.IsMarketplace,
			PriceSource:     s.Settings.PriceSource,
			IsAutoUpdate:    s.Settings.IsAutoUpdate,
			IsDeleted:       s.Settings.IsDeleted,
			HasVirtualStore: s.Settings.HasVirtualStore,
			Email:           s.Settings.Email,
		}

		if s.Settings.TimeZone != nil {
			settings.TimeZone = &models2.UpdateTimeZone{
				TZ:        s.Settings.TimeZone.TZ,
				UTCOffset: s.Settings.TimeZone.UTCOffset,
			}
		}

		if s.Settings.OrderDestination != nil {
			dest := models2.OrderDestination(*s.Settings.OrderDestination)
			settings.OrderDestination = &dest
		}

		if s.Settings.Group != nil {
			settings.Group = &models2.UpdateMenuGroup{
				MenuID:  s.Settings.Group.MenuID,
				GroupID: s.Settings.Group.GroupID,
			}
		}

		if s.Settings.StopListClosingActions != nil {
			var actions = make([]models2.UpdateStopListClosingAction, 0, len(s.Settings.StopListClosingActions))

			for i := 0; i < len(s.Settings.StopListClosingActions); i++ {
				actions = append(actions, models2.UpdateStopListClosingAction{
					RestaurantID: s.Settings.StopListClosingActions[i].RestaurantID,
					Opening:      s.Settings.StopListClosingActions[i].Opening,
					Closing:      s.Settings.StopListClosingActions[i].Closing,
					Status:       s.Settings.StopListClosingActions[i].Status,
				})
			}

			settings.StopListClosingActions = actions
		}

		store.Settings = settings
	}

	if s.Payments != nil {
		var payments = make([]models2.UpdatePayment, 0, len(s.Payments))

		for i := 0; i < len(s.Payments); i++ {
			payments = append(payments, models2.UpdatePayment{
				Type:     s.Payments[i].Type,
				Service:  s.Payments[i].Service,
				Username: s.Payments[i].Username,
				Password: s.Payments[i].Password,
			})
		}

		store.Payments = payments
	}

	if s.BillParameter != nil {
		var billParameter = &models2.UpdateParameters{
			IsActive: s.BillParameter.IsActive,
			UpdateBillParameters: models2.UpdateBillParameters{
				AddPaymentType:     s.BillParameter.UpdateBillParameters.AddPaymentType,
				AddOrderCode:       s.BillParameter.UpdateBillParameters.AddOrderCode,
				AddComments:        s.BillParameter.UpdateBillParameters.AddComments,
				AddDelivery:        s.BillParameter.UpdateBillParameters.AddDelivery,
				AddAddress:         s.BillParameter.UpdateBillParameters.AddAddress,
				AddQuantityPersons: s.BillParameter.UpdateBillParameters.AddQuantityPersons,
			},
		}
		store.BillParameter = billParameter
	}
	if s.StoreSchedule != nil {
		glovoSchedule := models2.AggregatorSchedule{}
		if s.StoreSchedule.GlovoSchedule != nil {
			glovoSchedule = models2.AggregatorSchedule{
				Timezone: s.StoreSchedule.GlovoSchedule.Timezone,
				Schedule: fillSchedule(s.StoreSchedule.GlovoSchedule.Schedule),
			}
		}
		woltSchedule := models2.AggregatorSchedule{}
		if s.StoreSchedule.WoltSchedule != nil {
			woltSchedule = models2.AggregatorSchedule{
				Timezone: s.StoreSchedule.WoltSchedule.Timezone,
				Schedule: fillSchedule(s.StoreSchedule.WoltSchedule.Schedule),
			}
		}

		directSchedules := make([]models2.DirectSchedule, 0, len(s.StoreSchedule.DirectSchedule))
		if s.StoreSchedule.DirectSchedule != nil {
			directSchedule := models2.DirectSchedule{}
			for _, storeSch := range s.StoreSchedule.DirectSchedule {
				directSchedule.DayOfWeek = storeSch.DayOfWeek
				directSchedule.TimeSlots.Opening = storeSch.TimeSlots.Opening
				directSchedule.TimeSlots.Closing = storeSch.TimeSlots.Closing
				directSchedules = append(directSchedules, directSchedule)
			}
		}

		var storeSchedule = &models2.UpdateAggregatorSchedule{}

		if s.StoreSchedule.GlovoSchedule != nil {
			storeSchedule.GlovoSchedule = &glovoSchedule
		}
		if s.StoreSchedule.WoltSchedule != nil {
			storeSchedule.WoltSchedule = &woltSchedule
		}
		if s.StoreSchedule.DirectSchedule != nil {
			storeSchedule.DirectSchedule = directSchedules
		}
		store.StoreSchedule = storeSchedule
	}

	return store
}

func fillSchedule(req []glovoModels.StoreSchedule) []models2.Schedule {
	if len(req) == 0 {
		return []models2.Schedule{}
	}
	schedule := make([]models2.Schedule, 0, len(req))
	for _, v := range req {
		timeSlots := make([]models2.TimeSlot, 0, len(v.TimeSlots))
		for _, val := range v.TimeSlots {
			timeSlots = append(timeSlots, models2.TimeSlot{
				Opening: val.Opening,
				Closing: val.Closing,
			})
		}
		schedule = append(schedule, models2.Schedule{
			DayOfWeek: v.DayOfWeek,
			TimeSlots: timeSlots,
		})
	}
	return schedule
}
