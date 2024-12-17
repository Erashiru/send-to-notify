package models

import (
	"time"
)

type UpdateStore struct {
	ID                *string                       `bson:"id"`
	Token             *string                       `bson:"token"`
	Name              *string                       `bson:"name"`
	MenuID            *string                       `bson:"menu_id"`
	PosType           *string                       `bson:"pos_type"`
	LegalEntityId     *string                       `bson:"legal_entity_id"`
	AccountManagerId  *string                       `bson:"account_manager_id"`
	SalesManagerId    *string                       `bson:"sales_manager_id"`
	Address           *UpdateStoreAddress           `bson:"address"`
	Glovo             *UpdateStoreGlovoConfig       `bson:"glovo"`
	Wolt              *UpdateStoreWoltConfig        `bson:"wolt"`
	Chocofood         *UpdateStoreChocofoodConfig   `bson:"chocofood"`
	RKeeper           *UpdateStoreRKeeperConfig     `bson:"rkeeper"`
	RKeeper7XML       *UpdateStoreRKeeper7XMLConfig `bson:"rkeeper7_xml"`
	Paloma            *UpdateStorePalomaConfig      `bson:"paloma"`
	Yandex            *UpdateStoreYandexConfig      `bson:"yandex"`
	External          []UpdateStoreExternalConfig   `bson:"external"`
	QRMenu            *UpdateStoreQRMenuConfig      `bson:"qr_menu"`
	KwaakaAdmin       *UpdateStoreKwaakaAdminConfig `bson:"kwaaka_admin"`
	MoySklad          *UpdateStoreMoySkladConfig    `bson:"moysklad"`
	IikoCloud         *UpdateStoreIikoConfig        `bson:"iiko_cloud"`
	Telegram          *UpdateStoreTelegramConfig    `bson:"telegram"`
	Notification      *UpdateNotification           `bson:"notification"`
	CallCenter        *UpdateStoreCallcenter        `bson:"callcenter"`
	Delivery          []UpdateStoreDelivery         `bson:"delivery"`
	Menus             []UpdateStoreDSMenu           `bson:"menus"`
	Contacts          []UpdateContact               `bson:"contacts"`
	Links             []UpdateLinks                 `bson:"links"`
	Settings          *UpdateSettings               `bson:"settings"`
	IntegrationDate   *time.Time                    `bson:"integration_date"`
	UpdatedAt         *time.Time                    `bson:"updated_at"`
	CreatedAt         *time.Time                    `bson:"created_at"`
	Payments          []UpdatePayment               `bson:"payments"`
	RestaurantGroupID *string                       `bson:"restaurant_group_id"`
	BillParameter     *UpdateParameters             `bson:"bill_parameter"`
	StoreSchedule     *UpdateAggregatorSchedule     `bson:"store_schedule"`
	IsDeleted         *bool                         `bson:"is_deleted"`
	SocialMediaLinks  []UpdateSocialMediaLinks      `bson:"social_media_links"`
	CompensationCount *int                          `bson:"compensation_count"`
}

type UpdateLinks struct {
	Name      string `bson:"name"`
	Url       string `bson:"url"`
	ImageLink string `bson:"image_link"`
}

type UpdateSocialMediaLinks struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	URL  string `bson:"url" json:"url"`
	Logo string `bson:"logo,omitempty" json:"logo,omitempty"`
}

type UpdateContact struct {
	FullName string `bson:"full_name"`
	Position string `bson:"position"`
	Phone    string `bson:"phone"`
	Comment  string `bson:"comment"`
	IsMain   bool   `bson:"is_main"`
}

type UpdateAggregatorSchedule struct {
	GlovoSchedule  *AggregatorSchedule `bson:"glovo_schedule"`
	WoltSchedule   *AggregatorSchedule `bson:"wolt_schedule"`
	DirectSchedule []DirectSchedule    `bson:"direct_schedule"`
}

type AggregatorSchedule struct {
	Timezone string     `bson:"timezone"`
	Schedule []Schedule `bson:"schedule"`
}

type Schedule struct {
	DayOfWeek int        `bson:"day_of_week"`
	TimeSlots []TimeSlot `bson:"time_slots"`
}
type TimeSlot struct {
	Opening string `bson:"opening"`
	Closing string `bson:"closing"`
}

type DirectSchedule struct {
	DayOfWeek int                    `json:"day_of_week" bson:"day_of_week"`
	TimeSlots DirectScheduleTimeSlot `json:"time_slots" bson:"time_slots"`
}

type DirectScheduleTimeSlot struct {
	Opening string `json:"opening" bson:"opening"`
	Closing string `json:"closing" bson:"closing"`
}

type UpdateStoreRKeeper7XMLConfig struct {
	Domain              *string `bson:"domain"`
	Username            *string `bson:"username"`
	Password            *string `bson:"password"`
	UCSUsername         *string `bson:"ucs_username"`
	UCSPassword         *string `bson:"ucs_password"`
	Token               *string `bson:"token"`
	ObjectID            *string `bson:"object_id"`
	Anchor              *string `bson:"anchor"`
	LicenseInstanceGUID *string `bson:"license_instance_guid"`
	SeqNumber           *int    `bson:"seq_number"`
}

type UpdateParameters struct {
	IsActive             *bool                `bson:"is_active"`
	UpdateBillParameters UpdateBillParameters `bson:"bill_params"`
}

type UpdateBillParameters struct {
	AddPaymentType     *bool `bson:"add_payment_type"`
	AddOrderCode       *bool `bson:"add_order_code"`
	AddComments        *bool `bson:"add-comments"`
	AddDelivery        *bool `bson:"add_delivery"`
	AddAddress         *bool `bson:"add_address"`
	AddQuantityPersons *bool `bson:"add_quantity_persons"`
}

type UpdateStoreAddress struct {
	City              *string            `bson:"city,omitempty"`
	Street            *string            `bson:"street,omitempty"`
	Entrance          *string            `bson:"entrance,omitempty"`
	UpdateCoordinates *UpdateCoordinates `bson:"coordinates"`
}

type UpdateCoordinates struct {
	Latitude  *float64 `bson:"latitude"`
	Longitude *float64 `bson:"longitude"`
}

type UpdateStoreGlovoConfig struct {
	StoreID                            []string                          `bson:"store_id"`
	MenuUrl                            *string                           `bson:"menu_url"`
	SendToPos                          *bool                             `bson:"send_to_pos"`
	IsMarketplace                      *bool                             `bson:"is_marketplace"`
	PaymentTypes                       *UpdateDeliveryServicePaymentType `bson:"payment_types"`
	PurchaseTypes                      *UpdatePurchaseTypes              `bson:"purchase_types"`
	IsOpen                             *bool                             `bson:"is_open"`
	AdditionalPreparationTimeInMinutes *int                              `bson:"additional_preparation_time_in_minutes"`
}

type UpdateDeliveryServicePaymentType struct {
	CASH    *UpdateIIKOPaymentType `bson:"CASH"`
	DELAYED *UpdateIIKOPaymentType `bson:"DELAYED"`
}

type UpdateIIKOPaymentType struct {
	IikoPaymentTypeID        *string `bson:"iiko_payment_type_id"`
	IikoPaymentTypeKind      *string `bson:"iiko_payment_type_kind"`
	PromotionPaymentTypeID   *string `bson:"promotion_payment_type_id"`
	OrderType                *string `bson:"order_type"`
	OrderTypeService         *string `bson:"order_type_service"`
	OrderTypeForVirtualStore *string `bson:"order_type_for_virtual_store"`
	IsProcessedExternally    *bool   `bson:"is_processed_externally"`
}

type UpdateStoreWoltConfig struct {
	StoreID               []string                          `bson:"store_id"`
	MenuUsername          *string                           `bson:"menu_username"`
	MenuPassword          *string                           `bson:"menu_password"`
	ApiKey                *string                           `bson:"api_key"`
	AdjustedPickupMinutes *int                              `bson:"adjusted_pickup_minutes"`
	MenuUrl               *string                           `bson:"menu_url"`
	SendToPos             *bool                             `bson:"send_to_pos"`
	IsMarketplace         *bool                             `bson:"is_marketplace"`
	PaymentTypes          *UpdateDeliveryServicePaymentType `bson:"payment_types"`
	PurchaseTypes         *UpdatePurchaseTypes              `bson:"purchase_types"`
	IgnoreStatusUpdate    *bool                             `bson:"ignore_status_update"`
	AutoAcceptOn          *bool                             `bson:"auto_accept_on"`
	IsOpen                *bool                             `bson:"is_open"`
}

type UpdatePurchaseTypes struct {
	Instant  []UpdateStatus `bson:"instant" json:"instant"`
	Preorder []UpdateStatus `bson:"preorder" json:"preorder"`
	TakeAway []UpdateStatus `bson:"takeaway" json:"takeaway"`
}

type UpdateStatus struct {
	PosStatus *string `bson:"pos_status" json:"pos_status"`
	Status    *string `bson:"status" json:"status"`
}

type UpdateStoreChocofoodConfig struct {
	StoreID       []string                          `bson:"store_id" json:"store_id"`
	MenuUrl       *string                           `bson:"menu_url" json:"menu_url"`
	SendToPos     *bool                             `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketplace *bool                             `bson:"is_marketplace" json:"is_marketplace"`
	PaymentTypes  *UpdateDeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
}

type UpdateStoreRKeeperConfig struct {
	ObjectId *int `bson:"object_id" json:"object_id"`
}

type UpdateStorePalomaConfig struct {
	PointID *string `bson:"point_id" json:"point_id"`
	ApiKey  *string `bson:"api_key" json:"api_key"`
}

type UpdateStoreYandexConfig struct {
	StoreID       []string                          `bson:"store_id" json:"store_id"`
	MenuUrl       *string                           `bson:"menu_url" json:"menu_url"`
	SendToPos     *bool                             `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketplace *bool                             `bson:"is_marketplace" json:"is_marketplace"`
	PaymentTypes  *UpdateDeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	ClientSecret  *string                           `bson:"client_secret" json:"client_secret"`
}

type UpdateStoreExternalConfig struct {
	StoreID                  []string                          `bson:"store_id" json:"store_id"`
	Type                     *string                           `bson:"type" json:"type"`
	MenuUrl                  *string                           `bson:"menu_url" json:"menu_url"`
	SendToPos                *bool                             `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketplace            *bool                             `bson:"is_marketplace" json:"is_marketplace"`
	PaymentTypes             *UpdateDeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	ClientSecret             *string                           `bson:"client_secret" json:"client_secret"`
	WebhookURL               *string                           `bson:"webhook_url" json:"webhook_url"`
	AuthToken                *string                           `bson:"auth_token" json:"auth_token"`
	WebhookProductStoplist   *string                           `bson:"webhook_product_stoplist" json:"webhook_product_stoplist"`
	WebhookAttributeStoplist *string                           `bson:"webhook_attribute_stoplist" json:"webhook_attribute_stoplist"`
}

type UpdateStoreQRMenuConfig struct {
	StoreID               []string                          `bson:"store_id" json:"store_id"`
	URL                   *string                           `bson:"url" json:"url"`
	IsIntegrated          *bool                             `bson:"is_integrated" json:"is_integrated"`
	PaymentTypes          *UpdateDeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	Hash                  *string                           `bson:"hash" json:"hash"`
	CookingTime           *int                              `bson:"cooking_time" json:"cooking_time"`
	DeliveryTime          *int                              `bson:"delivery_time" json:"delivery_time"`
	NoTable               *bool                             `bson:"no_table" json:"no_table"`
	Theme                 *string                           `bson:"theme" json:"theme"`
	IsMarketplace         *bool                             `bson:"is_marketplace" json:"is_marketplace"`
	SendToPos             *bool                             `bson:"send_to_pos" json:"send_to_pos"`
	IgnoreStatusUpdate    *bool                             `bson:"ignore_status_update" json:"ignore_status_update"`
	AdjustedPickupMinutes *int                              `bson:"adjusted_pickup_minutes" json:"adjusted_pickup_minutes"`
	BusyMode              *bool                             `bson:"busy_mode" json:"busy_mode"`
}

type UpdateStoreKwaakaAdminConfig struct {
	IsIntegrated *bool    `bson:"is_integrated" json:"is_integrated"`
	IsActive     *bool    `bson:"is_active" json:"is_active"`
	CookingTime  *int32   `bson:"cooking_time" json:"cooking_time"`
	StoreID      []string `bson:"store_id" json:"store_id"`
	SendToPos    *bool    `bson:"send_to_pos" json:"send_to_pos"`
}

type UpdateStoreMoySkladConfig struct {
	UserName       *string                           `bson:"username" json:"user_name"`
	Password       *string                           `bson:"password" json:"password"`
	OrderID        *string                           `bson:"order_id" json:"order_id"`
	OrganizationID *string                           `bson:"organization_id" json:"organization_id"`
	Status         *UpdateMoySkladStatus             `bson:"status" json:"status"`
	SendToPos      *bool                             `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketPlace  *bool                             `bson:"is_marketplace" json:"is_marketplace"`
	PaymentTypes   *UpdateDeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
}

type UpdateMoySkladStatus struct {
	ID         *string `bson:"id" json:"id"`
	Name       *string `bson:"name" json:"name"`
	StatusType *string `bson:"status_type" json:"status_type"`
}

type UpdateStoreIikoConfig struct {
	OrganizationID       *string `bson:"organization_id" json:"organization_id"`
	TerminalID           *string `bson:"terminal_id" json:"terminal_id"`
	Key                  *string `bson:"key" json:"key"`
	StopListByBalance    *bool   `bson:"stoplist_by_balance,omitempty" json:"stoplist_by_balance"` // temporary field for Traveler`s menu
	StopListBalanceLimit *int    `bson:"stoplist_balance_limit" json:"stoplist_balance_limit"`
	IsExternalMenu       *bool   `bson:"is_external_menu" json:"is_external_menu"`
	ExternalMenuID       *string `bson:"external_menu_id" json:"external_menu_id"`
	PriceCategory        *string `bson:"price_category" json:"price_category"`
}

type UpdateStoreTelegramConfig struct {
	GroupChatID    *string `bson:"group_chat_id" json:"group_chat_id"`
	StopListChatID *string `bson:"stoplist_chat_id" json:"stop_list_chat_id"`
	CancelChatID   *string `bson:"cancel_chat_id" json:"cancel_chat_id"`
}

type UpdateNotification struct {
	Whatsapp *UpdateWhatsapp `bson:"whatsapp" json:"whatsapp"`
}

type UpdateWhatsapp struct {
	Receivers []UpdateWhatsappReceiver `bson:"receivers" json:"receivers"`
}

type UpdateWhatsappReceiver struct {
	Name        *string `bson:"name" json:"name"`
	PhoneNumber *string `bson:"phone_number" json:"phone_number"`
	IsActive    *bool   `bson:"is_active" json:"is_active"`
}

type UpdateStoreCallcenter struct {
	Restaurants []string `bson:"restaurants" json:"restaurants"`
}

type UpdateStoreDelivery struct {
	ID       *string `bson:"id" json:"id"`
	Code     *string `bson:"code" json:"code"`
	Price    *int    `bson:"price" json:"price"`
	Name     *string `bson:"name" json:"name"`
	IsActive *bool   `bson:"is_active" json:"is_active"`
	Service  *string `bson:"service" json:"service"`
}

type UpdateStoreDSMenu struct {
	MenuID         *string    `bson:"menu_id" json:"menu_id"`
	Name           *string    `bson:"name" json:"name"`
	IsActive       *bool      `bson:"is_active" json:"is_active"`
	IsDeleted      *bool      `bson:"is_deleted" json:"is_deleted"`
	IsSync         *bool      `bson:"is_sync" json:"is_sync"`
	SyncAttributes *bool      `bson:"sync_attributes" json:"sync_attributes"`
	Delivery       *string    `bson:"delivery" json:"delivery"`
	Timestamp      *int       `bson:"timestamp" json:"timestamp"`
	UpdatedAt      *time.Time `bson:"updated_at" json:"updated_at"`
	HasWoltPromo   *bool      `bson:"has_wolt_promo" json:"has_wolt_promo"`
	CreationSource *string    `bson:"creation_source,omitempty" json:"creation_source"`
}

type UpdateSettings struct {
	TimeZone               *UpdateTimeZone               `bson:"timezone" json:"timezone"`
	Currency               *string                       `bson:"currency" json:"currency"`
	LanguageCode           *string                       `bson:"language_code" json:"language_code"`
	SendToPos              *bool                         `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketplace          *bool                         `bson:"is_marketplace" json:"is_marketplace"`
	PriceSource            *string                       `bson:"price_source" json:"price_source"`
	OrderDestination       *OrderDestination             `bson:"order_destination" json:"order_destination"`
	IsAutoUpdate           *bool                         `bson:"is_auto_update" json:"is_auto_update"`
	IsDeleted              *bool                         `bson:"is_deleted" json:"is_deleted"`
	Group                  *UpdateMenuGroup              `bson:"group" json:"group"`
	HasVirtualStore        *bool                         `bson:"has_virtual_store" json:"has_virtual_store"`
	StopListClosingActions []UpdateStopListClosingAction `bson:"stoplist_closing_actions" json:"stoplist_closing_actions"`
	Email                  *string                       `bson:"email" json:"email"`
}

type UpdateStopListClosingAction struct {
	RestaurantID *string `bson:"restaurant_id" json:"restaurant_id"`
	Opening      *string `bson:"opening" json:"opening"` // format 15:04:05
	Closing      *string `bson:"closing" json:"closing"` // format 15:04:05
	Status       *string `bson:"status" json:"status"`
}

type UpdateTimeZone struct {
	TZ        *string  `bson:"tz" json:"tz"`
	UTCOffset *float64 `bson:"utc_offset" json:"utc_offset"`
}

type UpdateMenuGroup struct {
	MenuID  *string `bson:"menu_id" json:"menu_id"`
	GroupID *string `bson:"group_id" json:"group_id"`
}

type UpdatePayment struct {
	Type     *string `bson:"type" json:"type"`
	Service  *string `bson:"service" json:"service"`
	Username *string `bson:"username" json:"username"`
	Password *string `bson:"password" json:"password"`
}

type Update3plRestaurantStatus struct {
	RestaurantId *string `bson:"restaurant_id" json:"restaurant_id"`
	Is3pl        *bool   `bson:"is_3pl" json:"is_3pl"`
}

type UpdateDispatchDeliveryAvailable struct {
	RestaurantID    string `bson:"restaurant_id" json:"restaurant_id"`
	DeliveryService string `bson:"delivery_service" json:"delivery_service"`
	Available       bool   `bson:"available" json:"available"`
	Is3pl           bool   `bson:"is_3pl" json:"is_3pl"`
}

type UpdateRestaurantCharge struct {
	IsRestaurantChargeOn    *bool    `json:"is_service_fee_on"`
	MinRestaurantCharge     *float64 `json:"min_service_fee"`
	MaxRestaurantCharge     *float64 `json:"max_service_fee"`
	RestaurantChargePercent *float64 `json:"service_fee_percent"`
}
