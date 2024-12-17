package models

import (
	"sort"
	"time"

	"github.com/pkg/errors"
)

const (
	ONLINE           = "ONLINE"
	OFFLINE          = "OFFLINE"
	DELIVERY_SERVICE = "DELIVERY_SERVICE"
)

type OrderDestination string

type StoreTelegramConfig struct {
	GroupChatID       string `bson:"group_chat_id" json:"group_chat_id"`
	StopListChatID    string `bson:"stoplist_chat_id" json:"stoplist_chat_id"`
	CancelChatID      string `bson:"cancel_chat_id" json:"cancel_chat_id"`
	CreateOrderChatID string `bson:"create_order_chat_id" json:"create_order_chat_id"`
	CheckInChatID     string `bson:"check_in_chat_id" json:"check_in_chat_id"`
	StoreStatusChatId string `bson:"store_status_chat_id"`
	TelegramBotToken  string `bson:"telegram_bot_token"`
}

type TimezoneSetting struct {
	TZ        string `bson:"tz" json:"tz"`
	UTCOffset int    `bson:"utc_offset" json:"utc_offset"`
}

type StoreSettings struct {
	Currency     string          `bson:"currency" json:"currency"`
	LanguageCode string          `bson:"language_code" json:"language_code"`
	PriceSource  string          `bson:"price_source" json:"price_source"`
	IsAutoUpdate bool            `bson:"is_auto_update" json:"is_auto_update"`
	Timezone     TimezoneSetting `bson:"timezone" json:"timezone"`
}

type CommentSetting struct {
	HasCommentSetting   bool                `bson:"has_comment_setting" json:"has_comment_setting"`
	OrderCodeName       string              `bson:"order_code_name" json:"order_code_name"`
	DeliveryName        string              `bson:"delivery_name" json:"delivery_name"`
	CommentName         string              `bson:"comment_name" json:"comment_name"`
	CutleryName         string              `bson:"cutlery_name" json:"cutlery_name"`
	Allergy             string              `bson:"allergy_name" json:"allergy"`
	CourierPhoneName    string              `bson:"courier_phone_name" json:"courier_phone_name"`
	AddressName         string              `bson:"address_name" json:"address_name"`
	DefaultCourierPhone string              `bson:"default_courier_phone" json:"default_courier_phone"`
	DelayedPaymentName  string              `bson:"delayed_payment_name" json:"delayed_payment_name"`
	CashPaymentName     string              `bson:"cash_payment_name" json:"cash_payment_name"`
	PaymentTypeName     string              `bson:"payment_type_name" json:"payment_type_name"`
	PickUpToName        string              `bson:"pick_up_to_name" json:"pick_up_to_name"`
	QuantityPerson      string              `bson:"quantity_person" json:"quantity_person"`
	CommentDynamicName  CommentDynamicNames `bson:"comment_dynamic_names" json:"comment_dynamic_name"`
}

type CommentDynamicNames struct {
	HasAllergy                string `bson:"has_allergy" json:"has_allergy"`
	HasNotAllergy             string `bson:"has_not_allergy" json:"has_not_allergy"`
	HasCutlery                string `bson:"has_cutlery" json:"has_cutlery"`
	HasNotCutlery             string `bson:"has_not_cutlery" json:"has_not_cutlery"`
	HasNotCourierPhone        string `bson:"has_courier_phone" json:"has_courier_phone"`
	HasNotAddress             string `bson:"has_not_address" json:"has_not_address"`
	HasNotSpecialRequirements string `bson:"has_not_special_requirements" json:"has_not_special_requirements"`
}

type StoreDSMenu struct {
	ID                     string    `bson:"menu_id,omitempty" json:"menu_id"` // json is menu_id as in ~/pkg/store/dto/update_store.go UpdateStoreDSMenu
	Name                   string    `bson:"name" json:"name"`
	IsActive               bool      `bson:"is_active" json:"is_active"`
	IsDeleted              bool      `bson:"is_deleted" json:"is_deleted"`
	IsSync                 bool      `bson:"is_sync" json:"is_sync"`
	SyncAttributes         bool      `bson:"sync_attributes" json:"sync_attributes"`
	Delivery               string    `bson:"delivery" json:"delivery"`
	Timestamp              int       `bson:"timestamp" json:"timestamp"`
	UpdatedAt              time.Time `bson:"updated_at" json:"updated_at"`
	Status                 string    `bson:"status" json:"status"`
	IsDiscount             bool      `bson:"is_discount" json:"is_discount"`
	IsProductOnStop        bool      `bson:"is_product_on_stop" json:"is_product_on_stop"`
	HasWoltPromo           bool      `bson:"has_wolt_promo" json:"has_wolt_promo"`
	EmptyProductPercentage int       `bson:"empty_product_percentage" json:"empty_product_percentage"`
	MarkupPercent          int       `bson:"markup_percent" json:"markup_percent"`
	CreationSource         string    `bson:"creation_source,omitempty" json:"creation_source"`
}

func (s StoreDSMenus) SetActiveMenu() {

	sort.Slice(s, func(i, j int) bool {
		return s[i].UpdatedAt.Before(s[j].UpdatedAt)
	})

	m := make(map[string]struct{}, len(s))

	for _, menu := range s {
		m[menu.Delivery] = struct{}{}
	}

	for i := len(s) - 1; i >= 0; i-- {
		if _, ok := m[s[i].Delivery]; ok {

			s[i].IsActive = true
			delete(m, s[i].Delivery)
			continue
		}
		s[i].IsActive = false
	}

}

type StoreCallcenter struct {
	Restaurants []string `bson:"restaurants" json:"restaurants"`
}

type Store struct {
	ID                             string                         `bson:"_id,omitempty" json:"id"`
	Token                          string                         `bson:"token" json:"token"`
	Name                           string                         `bson:"name" json:"name"`
	MenuID                         string                         `bson:"menu_id,omitempty" json:"menu_id"`
	PosType                        string                         `bson:"pos_type" json:"pos_type"`
	City                           string                         `bson:"city" json:"city"`
	StorePhoneNumber               string                         `bson:"store_phone_number" json:"store_phone_number"`
	WhatsappChatId                 string                         `bson:"whatsapp_chat_id,omitempty" json:"whatsapp_chat_id"`
	SendWhatsappNotification       bool                           `bson:"send_whatsapp_notification" json:"send_whatsapp_notification"`
	WhatsappPaymentChatId          string                         `bson:"whatsapp_payment_chat_id,omitempty" json:"whatsapp_payment_chat_id"`
	Address                        StoreAddress                   `bson:"address" json:"address"`
	QRMenu                         StoreQRMenuConfig              `bson:"qr_menu" json:"qr_menu"`
	StarterApp                     StoreStarterAppConfig          `bson:"starter_app" json:"starter_app"`
	Glovo                          StoreGlovoConfig               `bson:"glovo" json:"glovo"`
	Wolt                           StoreWoltConfig                `bson:"wolt" json:"wolt"`
	KwaakaAdmin                    StoreKwaakaAdminConfig         `bson:"kwaaka_admin" json:"kwaaka_admin"`
	Chocofood                      StoreChocofoodConfig           `bson:"chocofood" json:"chocofood"`
	Express24                      StoreExpress24Config           `bson:"express24" json:"express24"`
	Deliveroo                      StoreDeliverooConfig           `bson:"deliveroo" json:"deliveroo"`
	RKeeper                        StoreRKeeperConfig             `bson:"rkeeper" json:"rkeeper"`
	RKeeper7XML                    StoreRKeeper7XMLConfig         `bson:"rkeeper7_xml" json:"rkeeper7_xml"`
	Yaros                          StoreYarosConfig               `bson:"yaros" json:"yaros"`
	Talabat                        StoreTalabatConfig             `bson:"talabat" json:"talabat"`
	ExternalConfig                 []StoreExternalConfig          `bson:"external" json:"external"` // json is external as in ~/pkg/store/dto/update_store.go UpdateStore
	MoySklad                       StoreMoySkladConfig            `bson:"moysklad" json:"moy_sklad"`
	IikoCloud                      StoreIikoConfig                `bson:"iiko_cloud" json:"iiko_cloud"`
	TillyPad                       TillyPadConfig                 `bson:"tillypad" json:"tillypad"`
	YTimes                         YTimesConfig                   `bson:"ytimes" json:"ytimes"`
	Paloma                         StorePalomaConfig              `bson:"paloma" json:"paloma"`
	Poster                         StorePosterConfig              `bson:"poster" json:"poster"`
	Jowi                           StoreJowiConfig                `bson:"jowi" json:"jowi"`
	Posist                         StorePosistConfig              `bson:"posist" json:"posist"`
	Telegram                       StoreTelegramConfig            `bson:"telegram" json:"telegram"`
	ExternalPosIntegrationSettings ExternalPosIntegrationSettings `bson:"external_pos_integration_settings" json:"external_pos_integration_settings"`
	Notification                   Notification                   `bson:"notification" json:"notification"`
	CallCenter                     StoreCallcenter                `bson:"callcenter" json:"call_center"`
	Delivery                       []StoreDelivery                `bson:"delivery" json:"delivery"`
	Menus                          StoreDSMenus                   `bson:"menus" json:"menus"`
	Settings                       Settings                       `bson:"settings" json:"settings"`
	IntegrationDate                time.Time                      `bson:"integration_date" json:"integration_date"`
	UpdatedAt                      time.Time                      `bson:"updated_at" json:"updated_at"`
	CreatedAt                      time.Time                      `bson:"created_at" json:"created_at"`
	Payments                       []Payment                      `bson:"payments" json:"payments"`
	RestaurantGroupID              string                         `bson:"restaurant_group_id" json:"restaurant_group_id"`
	BillParameter                  Parameters                     `bson:"bill_parameters" json:"bill_parameter"`
	TipsTypes                      []TipsType                     `bson:"tips_types" json:"tips_types"`
	Yandex                         YandexConfig                   `bson:"yandex" json:"yandex"`
	SendToPOS                      bool                           `bson:"send_to_pos" json:"send_to_pos"`
	IsDeleted                      bool                           `bson:"is_deleted" json:"is_deleted"`
	Users                          []string                       `bson:"users" json:"users"`
	StoreSchedule                  StoreSchedule                  `bson:"store_schedule"`
	AutoUpdatePermission           bool                           `bson:"auto_update_permission" json:"auto_update_permission"`
	DeferSubmission                DeferSubmission                `bson:"defer_submission" json:"defer_submission"`
	Kwaaka3PL                      Kwaaka3PL                      `bson:"kwaaka_3pl" json:"kwaaka_3pl"`
	LegalEntityId                  string                         `bson:"legal_entity_id" json:"legal_entity_id"`
	AccountManagerId               string                         `bson:"account_manager_id" json:"account_manager_id"`
	SalesManagerId                 string                         `bson:"sales_manager_id" json:"sales_manager_id"`
	Contacts                       []Contact                      `bson:"contacts" json:"contacts"`
	ExternalLinks                  []Link                         `bson:"external_links" json:"external_links"`
	PaymentSystems                 []PaymentSystem                `bson:"payment_systems" json:"payment_systems"`
	SocialMediaLinks               []StoreSocialMediaLinks        `bson:"social_media_links" json:"social_media_links"`
	RestaurantCharge               RestaurantCharge               `bson:"restaurant_charge" json:"restaurant_charge"`
	RestaurantTips                 RestaurantTips                 `bson:"restaurant_tips" json:"restaurant_tips"`
	WhatsappConfig                 WhatsappUltraMsgConfig         `bson:"whatsapp_config" json:"whatsapp_config"`
	CompensationCount              int                            `bson:"compensation_count,omitempty" json:"compensation_count,omitempty"`
	ValidationSettings             ValidationSettings             `bson:"validation_settings,omitempty" json:"validation_settings,omitempty"`
	AutoUpdateSettings             AutoUpdateSettings             `bson:"auto_update_settings" json:"auto_update_settings"`
	OrderAutoCloseSettings         OrderAutoCloseSettings         `bson:"order_auto_close_settings" json:"order_auto_close_settings"`
}

type StoreStarterAppConfig struct {
	ApiKey             string                     `bson:"api_key" json:"api_key"`
	ShopID             string                     `bson:"shop_id" json:"shop_id"`
	SendToPos          bool                       `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketPlace      bool                       `bson:"is_market_place" json:"is_market_place"`
	StoreID            []string                   `bson:"store_id" json:"store_id"`
	IgnoreStatusUpdate bool                       `bson:"ignore_status_update"`
	PaymentTypes       DeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	CookingTime        int32                      `bson:"cooking_time" json:"cooking_time"`
}

type RestaurantTips struct {
	IsRestaurantTipsOn bool `bson:"is_restaurant_tips_on" json:"is_restaurant_tips_on"`
}

type DeferSubmission struct {
	IsDeferSumbission bool `bson:"is_defer_submission" json:"is_defer_sumbission"`
	DefaultTime       int  `bson:"default_time" json:"default_time"`
	BusyTime          int  `bson:"busy_time" json:"busy_time"`
}

type WhatsappUltraMsgConfig struct {
	InstanceId  string `bson:"instance_id" json:"instance_id"`
	AuthToken   string `bson:"auth_token" json:"auth_token"`
	PhoneNumber string `bson:"phone_number" json:"phone_number"`
}

type RestaurantCharge struct {
	IsRestaurantChargeOn    bool    `bson:"is_restaurant_charge_on" json:"is_restaurant_charge_on"`
	MinRestaurantCharge     float64 `bson:"min_restaurant_charge" json:"min_restaurant_charge"`
	MaxRestaurantCharge     float64 `bson:"max_restaurant_charge" json:"max_restaurant_charge"`
	RestaurantChargePercent float64 `bson:"restaurant_charge_percent" json:"restaurant_charge_percent"`
}

type PaymentSystem struct {
	Name           string `json:"name" bson:"name"`
	IsActive       bool   `json:"is_active" bson:"is_active"`
	PaymentURL     string `json:"payment_url" bson:"payment_url"`
	PosPaymentType string `json:"pos_payment_type" bson:"pos_payment_type"`
}

type Link struct {
	Name      string `json:"name" bson:"name"`
	Url       string `json:"url" bson:"url"`
	ImageLink string `json:"image_link" bson:"image_link"`
}

type StoreSocialMediaLinks struct {
	Name string `bson:"name" json:"name"`
	URL  string `bson:"url" json:"url"`
	Logo string `bson:"logo" json:"logo"`
}

type Kwaaka3PL struct {
	Is3pl                  bool      `bson:"is_3pl" json:"is_3pl"`
	WoltDriveStoreID       string    `bson:"wolt_drive_store_id" json:"wolt_drive_store_id"`
	IndriveStoreID         string    `bson:"indrive_store_id" json:"indrive_store_id"`
	KwaakaChargeAbsolute   float64   `bson:"kwaaka_charge_absolute" json:"delivery_service_fee_absolut"`
	KwaakaChargePercentage float64   `bson:"kwaaka_charge_percentage" json:"delivery_service_fee_percentage"`
	IndriveAvailable       bool      `json:"indrive_available" bson:"indrive_available"`
	WoltDriveAvailable     bool      `json:"wolt_drive_available" bson:"wolt_drive_available"`
	YandexAvailable        bool      `json:"yandex_available" bson:"yandex_available"`
	IikoCouriersAvailable  bool      `json:"iiko_couriers_available" bson:"iiko_couriers_available"`
	Polygons               []Polygon `bson:"polygons" json:"polygons"`
	CPO                    float64   `bson:"cpo" json:"cpo"`
	IsDynamic              bool      `bson:"is_dynamic" json:"is_dynamic"`
	IsInstantCall          bool      `bson:"is_instant_call" json:"is_instant_call"`
	TaxiClass              string    `bson:"taxi_class" json:"taxi_class"`
	ChatID                 string    `bson:"chat_id" json:"chat_id"`
	DeliveryPosProductId   string    `bson:"delivery_pos_product_id" json:"delivery_pos_product_id"`
}

type ValidationSettings struct {
	ForbiddenUpsert bool `bson:"forbidden_upsert,omitempty" json:"forbidden_upsert,omitempty"`
}

type StoreDSMenus []StoreDSMenu

func (s Store) VerifyMenuOwnership(menuId string) bool {
	for _, menuDs := range s.Menus {
		if menuDs.ID == menuId {
			return true
		}
	}

	return false
}

func (s Store) GetMenuMarkupPercent(menuId string) int {
	for _, menuDs := range s.Menus {
		if menuDs.ID == menuId {
			return menuDs.MarkupPercent
		}
	}

	return 0
}

func (s StoreDSMenus) GetActiveMenu(name AggregatorName) StoreDSMenu {

	sort.Slice(s, func(i, j int) bool {
		return s[i].UpdatedAt.Before(s[j].UpdatedAt)
	})

	for i := len(s) - 1; i >= 0; i-- {
		if name.String() != s[i].Delivery {
			continue
		}
		return s[i]
	}

	return StoreDSMenu{}

}

type StoreSchedule struct {
	GlovoSchedule  AggregatorSchedule `bson:"glovo_schedule"`
	WoltSchedule   AggregatorSchedule `bson:"wolt_schedule"`
	DirectSchedule []DirectSchedule   `bson:"direct_schedule"`
}

func (s Store) GetAggregatorStoreIDs(name string) []string {
	switch name {
	case "glovo":
		if len(s.Glovo.StoreID) != 0 {
			return s.Glovo.StoreID
		}
	case "wolt":
		if len(s.Wolt.StoreID) != 0 {
			return s.Wolt.StoreID
		}
	case "yandex", "emenu":
		for _, config := range s.ExternalConfig {
			if config.Type == name {
				return config.StoreID
			}
		}
	case "chocofood":
		if len(s.Chocofood.StoreID) != 0 {
			return s.Chocofood.StoreID
		}
	case "moysklad":
		return []string{
			"moysklad_store_id",
		}
	case "express24", "express24_v2":
		if len(s.Express24.StoreID) != 0 {
			return s.Express24.StoreID
		}
	case "talabat":
		if len(s.Talabat.BranchID) != 0 {
			return s.Talabat.BranchID
		}
	case "starter_app":
		if s.StarterApp.ShopID != "" {
			return []string{s.StarterApp.ShopID}
		}
	}

	return nil
}

type ExternalPosIntegrationSettings struct {
	StopListIsOn    bool `bson:"stoplist_is_on" json:"stoplist_is_on"`
	PayOrderIsOn    bool `bson:"pay_order_is_on" json:"pay_order_is_on"`
	PrePayOrderIsOn bool `bson:"pre_pay_order_is_on" json:"pre_pay_order_is_on"`
}

type Parameters struct {
	IsActive       bool           `bson:"is_active" json:"is_active"`
	BillParameters BillParameters `bson:"parameters" json:"bill_parameters"`
}

type BillParameters struct {
	AddPaymentType     bool `bson:"add_payment_type" json:"add_payment_type"`
	AddOrderCode       bool `bson:"add_order_code" json:"add_order_code"`
	AddComments        bool `bson:"add_comments" json:"add_comments"`
	AddDelivery        bool `bson:"add_delivery" json:"add_delivery"`
	AddAddress         bool `bson:"add_address" json:"add_address"`
	AddQuantityPersons bool `bson:"add_quantity_persons" json:"add_quantity_persons"`
}

type IikoWoltStatus struct {
	Iiko string `bson:"iiko" json:"iiko"`
	Wolt string `bson:"wolt" json:"wolt"`
}

type Payment struct {
	Type            string `bson:"type" json:"type"`
	Service         string `bson:"service" json:"service"`
	Username        string `bson:"username" json:"username"`
	Password        string `bson:"password" json:"password"`
	MerchantService string `bson:"merchant_service" json:"merchant_service"`
}

func (s Store) GetStorePaymentService(service string) (Payment, error) {
	for _, pay := range s.Payments {
		if pay.Service == service {
			return pay, nil
		}
	}
	return Payment{}, errors.New("not found payment service")
}

type Settings struct {
	TimeZone               TimeZone                `bson:"timezone" json:"timezone"`
	Currency               string                  `bson:"currency" json:"currency"`
	LanguageCode           string                  `bson:"language_code" json:"language_code"`
	SendToPos              bool                    `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketplace          bool                    `bson:"is_marketplace" json:"is_marketplace"`
	PriceSource            string                  `bson:"price_source" json:"price_source"`
	OrderDestination       OrderDestination        `bson:"order_destination" json:"order_destination"`
	IsAutoUpdate           bool                    `bson:"is_auto_update" json:"is_auto_update"`
	IsDeleted              bool                    `bson:"is_deleted" json:"is_deleted"`
	CallCenter             StoreCallcenter         `bson:"callcenter" json:"callcenter"`
	Group                  MenuGroup               `bson:"group" json:"group"`
	CommentSetting         CommentSetting          `bson:"comment_setting" json:"comment_setting"`
	RetrySetting           RetrySetting            `bson:"retry_setting" json:"retry_setting"`
	HasVirtualStore        bool                    `bson:"has_virtual_store" json:"has_virtual_store"`
	IsVirtualStore         bool                    `bson:"is_virtual_store" json:"is_virtual_store"`
	StopListClosingActions []StopListClosingAction `bson:"stoplist_closing_actions" json:"stoplist_closing_actions"`
	StopListByBalance      bool                    `bson:"stop_list_by_balance" json:"stoplist_by_balance"`
	StatusList             []string                `bson:"status_list" json:"status_list"`
	ScheduledStatusChange  ScheduledStatusChange   `bson:"scheduled_status_change"`
	Email                  string                  `bson:"email" json:"email"`
	IgnoreUpdateStopList   bool                    `bson:"ignore_update_stoplist" json:"ignore_update_stoplist"` //если true игнорирует рестораны и не обновляет стоп листы для них
	UtensilsProductID      string                  `bson:"utensils_product_id" json:"utensils_product_id"`       //айди продукта “Приборы“ (NAVAT)
	SendUtensilsToPos      bool                    `bson:"send_utensils_to_pos" json:"send_utensils_to_pos"`
}

type StopListClosingAction struct {
	RestaurantID string `bson:"restaurant_id" json:"restaurant_id"`
	Opening      string `bson:"opening" json:"opening"` // format 15:04:05
	Closing      string `bson:"closing" json:"closing"` // format 15:04:05
	Status       string `bson:"status" json:"status"`
}

type RetrySetting struct {
	Message  string `bson:"message" json:"message"`
	IsActive bool   `bson:"is_active" json:"is_active"`
}

type TimeZone struct {
	TZ        string  `bson:"tz" json:"tz"`
	UTCOffset float64 `bson:"utc_offset" json:"utc_offset"`
}

type MenuGroup struct {
	MenuID  string `bson:"menu_id" json:"menu_id"`
	GroupID string `bson:"group_id" json:"group_id"`
}

type StoreManagement struct {
	StoreInfo       []StoreInfo `bson:"store_info" json:"store_info"`
	DeliveryService string      `bson:"delivery_service" json:"delivery_service"`
}

type StoreInfo struct {
	RestaurantId string `bson:"restaurant_id" json:"restaurant_id"`
	StoreID      string `bson:"store_id" json:"store_id"`
	StoreStatus  bool   `bson:"store_status" json:"store_status"`
}
type StoreManagementResponse struct {
	ErrMessage      string `bson:"err_message" json:"err_message"`
	Success         bool   `bson:"success" json:"success"`
	RestaurantId    string `bson:"restaurant_id" json:"restaurant_id"`
	StoreID         string `bson:"store_id" json:"store_id"`
	IsOpen          bool   `bson:"is_open" json:"is_open"`
	DeliveryService string `bson:"delivery_service" json:"delivery_service"`
}

type TipsType struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	IsActive bool   `json:"is_active"`
}

func (s Store) GetYandexConfig() (StoreExternalConfig, error) {
	for _, config := range s.ExternalConfig {
		if config.Type != "yandex" {
			continue
		}
		return config, nil
	}
	return StoreExternalConfig{}, errors.New("yandex not found")
}

type YandexConfig struct {
	StoreIDs      []string               `bson:"store_id" json:"store_ids"`
	PaymentTypes  map[string]PaymentType `bson:"payment_types,omitempty" json:"payment_types"`
	MenuUrl       string                 `bson:"menu_url" json:"menu_url"`
	SendToPos     bool                   `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketPlace bool                   `bson:"is_marketplace" json:"is_market_place"`
}

type UpdateStoreSchedule struct {
	RestaurantID    string
	StoreSchedule   AggregatorSchedule
	DeliveryService string
}

type ScheduledStatusChange struct {
	IsActive         bool     `bson:"is_active"`
	SwitchInterval   int      `bson:"switch_interval"`
	DeliveryServices []string `bson:"delivery_services"`
}

type OpenTimeDuration struct {
	DeliveryService        string
	ActualOpenTimeDuration time.Duration // sum of time differences between each timeslot
	TotalOpenTimeDuration  time.Duration // difference between first timeslot.opening_time and last timeslot.closing_time
}

type CallCenterRestaurant struct {
	ID           string `bson:"_id" json:"_id"`
	Name         string `bson:"name" json:"name"`
	IsIntegrated bool   `bson:"is_integrated" json:"is_integrated"`
}

type DirectRestaurant struct {
	ID                  string `bson:"_id" json:"_id"`
	Name                string `bson:"name" json:"name"`
	QRMenuIsMarketplace bool   `bson:"qr_menu_is_marketplace" json:"qr_menu"`
}

type Restaurants struct {
	Restaurant []Restaurant
}

type Restaurant struct {
	RestaurantID string `bson:"restaurant_id" json:"restaurant_id"`
	Name         string `bson:"name" json:"name"`
}

type AutoUpdateSettings struct {
	AutoUpdateOn bool         `bson:"auto_update_on" json:"auto_update_on"`
	Wolt         bool         `bson:"wolt" json:"wolt"`
	Glovo        bool         `bson:"glovo" json:"glovo"`
	Yandex       bool         `bson:"yandex" json:"yandex"`
	UpdateFields UpdateFields `bson:"update_fields" json:"update_fields"`
}

type UpdateFields struct {
	ProductName          bool `json:"product_name" bson:"product_name"`
	ProductPrice         bool `json:"product_price" bson:"product_price"`
	ProductDescription   bool `json:"product_description" bson:"product_description"`
	ProductImage         bool `json:"product_image" bson:"product_image"`
	AttributeGroupName   bool `json:"attribute_group_name" bson:"attribute_group_name"`
	AttributeGroupMinMax bool `json:"attribute_group_min_max" bson:"attribute_group_min_max"`
	AttributeName        bool `json:"attribute_name" bson:"attribute_name"`
	AttributePrice       bool `json:"attribute_price" bson:"attribute_price"`
}

type OrderAutoCloseSettings struct {
	OrderAutoClose     bool `bson:"order_auto_close" json:"order_auto_close"`
	OrderAutoCloseTime int  `bson:"order_auto_close_time" json:"order_auto_close_time"`
}
