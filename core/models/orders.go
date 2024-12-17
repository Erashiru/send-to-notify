package models

import (
	models3 "github.com/kwaaka-team/orders-core/service/kwaaka_3pl/models"
	"time"
)

type OrderStatus string

const (
	STATUS_NEW                           OrderStatus = "NEW"
	STATUS_WAIT_SENDING                  OrderStatus = "WAIT_SENDING"
	STATUS_PENDING                       OrderStatus = "PENDING"
	STATUS_PROCESSING                    OrderStatus = "PROCESSING"
	STATUS_ACCEPTED                      OrderStatus = "ACCEPTED"
	STATUS_SKIPPED                       OrderStatus = "SKIPPED"
	STATUS_COOKING_STARTED               OrderStatus = "COOKING_STARTED"
	STATUS_COOKING_COMPLETE              OrderStatus = "COOKING_COMPLETE"
	STATUS_READY_FOR_PICKUP              OrderStatus = "READY_FOR_PICKUP"
	STATUS_OUT_FOR_DELIVERY              OrderStatus = "OUT_FOR_DELIVERY"
	STATUS_PICKED_UP_BY_CUSTOMER         OrderStatus = "PICKED_UP_BY_CUSTOMER"
	STATUS_DELIVERED                     OrderStatus = "DELIVERED"
	STATUS_CLOSED                        OrderStatus = "CLOSED"
	STATUS_CANCELLED_BY_DELIVERY_SERVICE OrderStatus = "CANCELLED_BY_DELIVERY_SERVICE"
	STATUS_CANCELLED_BY_POS_SYSTEM       OrderStatus = "CANCELLED_BY_POS_SYSTEM"
	STATUS_FAILED                        OrderStatus = "FAILED"
	STATUS_CANCELLED                     OrderStatus = "CANCELLED"
)

func (os OrderStatus) ToString() string {
	return string(os)
}

const (
	PAYMENT_METHOD_DELAYED = "DELAYED"
	PAYMENT_METHOD_CASH    = "CASH"
	PAYMENT_METHOD_CARD    = "CARD"
)

const (
	POS_PAYMENT_CASH = "Cash"
	POS_PAYMENT_CARD = "Card"
)

const (
	ORDER_DELIVERY_CLIENT        = "DeliveryByClient"
	FOODBAND_DELIVERY_AGGREGATOR = "aggregator_delivery"
	FOODBAND_DELIVERY_RESTAURANT = "self_delivery"
	FOODBAND_CUSTOMER_PICKUP     = "customer_delivery"
)

const (
	PROMO_TYPE_FIXED      = "FIXED"
	PROMO_TYPE_PERCENTAGE = "PERCENTAGE"
)

const (
	ORDER_TYPE_INSTANT  = "INSTANT"
	ORDER_TYPE_PREORDER = "PREORDER"
)

type PreOrder struct {
	Time   TransactionTime `bson:"time" json:"time"`
	Status string          `bson:"status" json:"status"`
}

type Price struct {
	Value        float64 `bson:"value" json:"value"`
	CurrencyCode string  `bson:"currency_code" json:"currency_code"`
}

type Courier struct {
	Name        string `bson:"name" json:"name"`
	PhoneNumber string `bson:"phone_number" json:"phone_number"`
}

type CustomerInvoicingDetails struct {
	CompanyName    string `bson:"company_name" json:"company_name"`
	CompanyAddress string `bson:"company_address" json:"company_address"`
	TaxID          string `bson:"tax_id" json:"tax_id"`
}

type Customer struct {
	Name                string                   `bson:"name" json:"name"`
	PhoneNumber         string                   `bson:"phone_number" json:"phone_number"`
	PhoneNumberWithPlus string                   `bson:"phone_number_with_plus"`
	Hash                string                   `bson:"hash" json:"hash"`
	InvoicingDetails    CustomerInvoicingDetails `bson:"invoicing_details" json:"invoicing_details"`
	FCMToken            string                   `bson:"fcm_token" json:"fcm_token"`
	Email               string                   `bson:"email" json:"email"`
}

type ProductAttribute struct {
	ID               string `bson:"id" json:"id"`
	GroupID          string `bson:"group_id" json:"group_id"`
	Quantity         int    `bson:"quantity" json:"quantity"`
	Price            Price  `bson:"price" json:"price"`
	Name             string `bson:"name" json:"name"`
	IsComboAttribute bool   `bson:"is_combo_attribute" json:"is_combo_attribute"`
}

type CreationResult struct {
	Message          string    `bson:"message,omitempty" json:"message"`
	CorrelationId    string    `bson:"correlationId,omitempty" json:"correlation_id"`
	OrderInfo        OrderInfo `bson:"orderInfo,omitempty" json:"order_info"`
	ErrorDescription string    `bson:"errorDescription,omitempty" json:"error_description"`
}

type OrderInfo struct {
	ID             string `bson:"id" json:"id"`
	OrganizationID string `bson:"organizationId" json:"organization_id"`
	Timestamp      int64  `bson:"timestamp" json:"timestamp"`
	CreationStatus string `bson:"creationStatus" json:"creation_status"`
}
type OrderProduct struct {
	ID                   string             `bson:"id" json:"id"`
	PurchasedProductID   string             `bson:"purchased_product_id" json:"purchased_product_id"`
	Quantity             int                `bson:"quantity" json:"quantity"`
	Price                Price              `bson:"price" json:"price"`
	PriceWithoutDiscount Price              `bson:"price_without_discount" json:"price_without_discount"`
	Name                 string             `bson:"name" json:"name"`
	IsCombo              bool               `bson:"is_combo" json:"is_combo"`
	Attributes           []ProductAttribute `bson:"attributes" json:"attributes"`
	SizeId               string             `bson:"size_id"`
	Promos               []Promo            `bson:"promos" json:"promos"`
	SourceActionID       string             `bson:"source_action_id" json:"source_action_id"`
	ProgramID            string             `bson:"program_id" json:"program_id"`
	ImageURLs            []string           `bson:"image_urls" json:"image_urls"`
}

type DeliveryAddress struct {
	Label        string  `bson:"label" json:"label"`
	Latitude     float64 `bson:"latitude" json:"latitude"`
	Longitude    float64 `bson:"longitude" json:"longitude"`
	City         string  `bson:"city" json:"city"`
	Comment      string  `bson:"comment" json:"comment"`
	BuildingName string  `json:"building_name"`
	Street       string  `json:"street"`
	Flat         string  `json:"flat"`
	Porch        string  `json:"porch"`
	Floor        string  `json:"floor"`
	HouseNumber  string  `json:"house_number"`
	Entrance     string  `json:"entrance"`
	Intercom     string  `json:"intercom"`
}

type OrderStatusUpdate struct {
	Name string    `bson:"name" json:"name"`
	Time time.Time `bson:"time" json:"time"`
}

type CancelReason struct {
	Reason      string `bson:"reason" json:"reason"`
	Description string `bson:"description" json:"description"`
	ExtraData   string `bson:"extra_data" json:"extra_data"`
}

type CancelOrder struct {
	ID           string       `bson:"_id,omitempty" json:"id"`
	OrderID      string       `bson:"order_id" json:"order_id"`
	Comment      string       `bson:"comment" json:"comment"`
	CancelReason CancelReason `bson:"cancel_reason" json:"cancel_reason"`
}

type CancelOrderInPos struct {
	OrderID         string       `bson:"order_id" json:"order_id"`
	DeliveryService string       `bson:"delivery_sevice" json:"delivery_sevice"`
	CancelReason    CancelReason `bson:"cancel_reason" json:"cancel_reason"`
	PaymentStrategy string       `bson:"payment_strategy" json:"payment_strategy"`
}

type ViewStatus struct {
	Username string    `json:"by_username" bson:"by_username"`
	ReadTime time.Time `json:"read_time" bson:"read_time"`
}

type Order struct {
	ID                              string               `bson:"_id,omitempty" json:"_id,omitempty"`
	Type                            string               `bson:"type" json:"type"`
	Discriminator                   string               `bson:"discriminator" json:"discriminator,omitempty"`
	DeliveryService                 string               `bson:"delivery_service" json:"delivery_service,omitempty"`
	PosType                         string               `bson:"pos_type" json:"pos_type,omitempty"`
	RestaurantID                    string               `bson:"restaurant_id" json:"restaurant_id,omitempty"`
	RestaurantName                  string               `bson:"restaurant_name" json:"restaurant_name,omitempty"`
	Preorder                        PreOrder             `bson:"preorder" json:"preorder,omitempty"`
	OrderID                         string               `bson:"order_id" json:"order_id,omitempty"`
	StoreID                         string               `bson:"store_id" json:"store_id,omitempty"`
	ViewStatus                      ViewStatus           `bson:"view_status" json:"view_status"`
	HasAggregatorAndPartnerDiscount bool                 `bson:"has_aggregator_and_partner_discount" json:"has_aggregator_and_partner_discount"`
	OrderCode                       string               `bson:"order_code" json:"order_code,omitempty"`
	OrderCodePrefix                 string               `bson:"order_code_prefix" json:"order_code_prefix,omitempty"`
	PosOrderID                      string               `bson:"pos_order_id" json:"pos_order_id,omitempty"`
	PosGuid                         string               `bson:"pos_guid" json:"pos_guid,omitempty"`
	PickUpCode                      string               `bson:"pick_up_code" json:"pick_up_code,omitempty"`
	IsMarketplace                   bool                 `bson:"is_marketplace" json:"is_marketplace,omitempty"`
	EatsID                          string               `bson:"eats_id" json:"eats_id,omitempty"`
	IsDeferSubmission               bool                 `bson:"is_defer_submission" json:"is_defer_submission"`
	Status                          string               `bson:"status" json:"status,omitempty"`
	StatusesHistory                 []OrderStatusUpdate  `bson:"statuses_history" json:"statuses_history,omitempty"`
	OrderTime                       TransactionTime      `bson:"order_time" json:"order_time,omitempty"`
	EstimatedPickupTime             TransactionTime      `bson:"estimated_pickup_time" json:"estimated_pickup_time,omitempty"`
	UtcOffsetMinutes                string               `bson:"utc_offset_minutes" json:"utc_offset_minutes,omitempty"`
	PaymentMethod                   string               `bson:"payment_method" json:"payment_method,omitempty"`
	PaymentSystem                   string               `bson:"payment_system" json:"payment_system,omitempty"`
	Currency                        string               `bson:"currency" json:"currency,omitempty"`
	AllergyInfo                     string               `bson:"allergy_info" json:"allergy_info,omitempty"`
	SpecialRequirements             string               `bson:"special_requirements" json:"special_requirements,omitempty"`
	EstimatedTotalPrice             Price                `bson:"estimated_total_price" json:"estimated_total_price,omitempty"`
	DeliveryFee                     Price                `bson:"delivery_fee" json:"delivery_fee,omitempty"`
	MinimumBasketSurcharge          Price                `bson:"minimum_basket_surcharge" json:"minimum_basket_surcharge,omitempty"`
	CustomerCashPaymentAmount       Price                `bson:"customer_cash_payment_amount" json:"customer_cash_payment_amount,omitempty"`
	PartnerDiscountsProducts        Price                `bson:"partner_discounts_products" json:"partner_discounts_products,omitempty"`
	PartnerDiscountedProductsTotal  Price                `bson:"partner_discounted_products_total" json:"partner_discounted_products_total,omitempty"`
	TotalCustomerToPay              Price                `bson:"total_customer_to_pay" json:"total_customer_to_pay,omitempty"`
	RestaurantCharge                Price                `bson:"restaurant_charge" json:"restaurant_charge"`
	Courier                         Courier              `bson:"courier" json:"courier,omitempty"`
	Customer                        Customer             `bson:"customer" json:"customer,omitempty"`
	Persons                         int                  `bson:"persons" json:"persons,omitempty"`
	Products                        []OrderProduct       `bson:"products" json:"products,omitempty"`
	DeliveryAddress                 DeliveryAddress      `bson:"delivery_address" json:"delivery_address,omitempty"`
	BundledOrders                   []string             `bson:"bundled_orders" json:"bundled_orders,omitempty"`
	IsPickedUpByCustomer            bool                 `bson:"is_picked_up_by_customer" json:"is_picked_up_by_customer,omitempty"`
	CutleryRequested                bool                 `bson:"cutlery_requested" json:"cutlery_requested,omitempty"`
	CreationResult                  CreationResult       `bson:"creation_result" json:"creation_result,omitempty"`
	CancelReason                    CancelReason         `bson:"cancel_reason" json:"cancel_reason,omitempty"`
	PaymentStrategy                 string               `bson:"payment_strategy" json:"payment_strategy,omitempty"`
	LoyaltyCard                     string               `bson:"loyalty_card" json:"loyalty_card,omitempty"`
	CreatedAt                       Time                 `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt                       Time                 `bson:"updated_at" json:"updated_at,omitempty"`
	Promos                          []Promo              `bson:"promos" json:"promos,omitempty"`
	DiscountInfo                    PosDiscountsInfo     `bson:"discountsInfo" json:"discount_info,omitempty"`
	TableID                         string               `bson:"table_id" json:"table_id,omitempty"`
	PaymentID                       string               `bson:"payment_id" json:"payment_id,omitempty"`
	RestaurantHash                  string               `bson:"restaurant_hash" json:"restaurant_hash,omitempty"`
	OrganizationID                  string               `bson:"organization_id" json:"organization_id,omitempty"`
	RetryCount                      int                  `bson:"retry_count" json:"retry_count,omitempty"`
	IsRetry                         bool                 `bson:"is_retry" json:"is_retry,omitempty"`
	Errors                          []Error              `bson:"errors" json:"errors,omitempty"`
	IsPaidStatus                    bool                 `bson:"is_paid" json:"is_paid_status,omitempty"`
	IsParentOrder                   bool                 `bson:"is_parent_order" json:"is_parent_order,omitempty"`
	IsChildOrder                    bool                 `bson:"is_child_order" json:"is_child_order,omitempty"`
	ReadingTime                     []TransactionTime    `bson:"reading_time" json:"reading_time,omitempty"`
	VirtualStoreComment             string               `bson:"virtual_store_comment" json:"virtual_store_comment,omitempty"`
	ServiceFeeSum                   float64              `bson:"service_fee_sum" json:"service_fee_sum,omitempty"`
	HasServiceFee                   bool                 `bson:"has_service_fee" json:"has_service_fee,omitempty"`
	LogLinks                        LogLinks             `bson:"log_links" json:"log_links,omitempty"`
	LogMessages                     LogMessages          `bson:"log_messages" json:"log_messages,omitempty"`
	PosPaymentInfo                  PosPaymentInfo       `bson:"pos_payment_info" json:"pos_payment_info,omitempty"`
	CookingCompleteTime             time.Time            `bson:"cooking_complete_time" json:"cooking_complete_time"`
	CookingTime                     int32                `bson:"cooking_time" json:"cooking_time"`
	FailReason                      FailReason           `bson:"fail_reason" json:"fail_reason,omitempty"`
	Proposals                       []Proposal           `bson:"proposals,omitempty" json:"proposals"`
	IsInstantDelivery               bool                 `bson:"is_instant_delivery" json:"is_instant_delivery"`
	DeliveryOrderID                 string               `bson:"delivery_order_id" json:"delivery_order_id"`
	DeliveryDispatcher              string               `bson:"delivery_dispatcher" json:"delivery_dispatcher"`
	DispatcherDeliveryTime          int32                `bson:"dispatcher_delivery_time" json:"dispatcher_delivery_time"`
	DeliveryOrderPromiseID          string               `bson:"delivery_order_promise_id" json:"delivery_order_promise_id"`
	FullDeliveryPrice               float64              `bson:"full_delivery_price" json:"full_delivery_price"`
	ClientDeliveryPrice             float64              `bson:"client_delivery_price" json:"client_delivery_price"`
	RestaurantPayDeliveryPrice      float64              `bson:"restaurant_pay_delivery_price" json:"restaurant_pay_delivery_price"`
	KwaakaChargedDeliveryPrice      float64              `bson:"kwaaka_charged_delivery_price" json:"delivery_service_fee"`
	DeliveryDropOffScheduleTime     string               `bson:"delivery_drop_off_schedule_time" json:"delivery_drop_off_schedule_time"`
	RestaurantSelfDelivery          bool                 `bson:"restaurant_self_delivery" json:"restaurant_self_delivery"`
	IsActive                        bool                 `bson:"is_active" json:"is_active"`
	Review                          Review               `bson:"review,omitempty" json:"review,omitempty"`
	SendCourier                     bool                 `json:"send_courier" bson:"send_courier"`
	PromoCode                       string               `json:"promo_code" bson:"promo_code"` // promo code from qr menu
	OperatorID                      string               `json:"operator_id" bson:"operator_id"`
	OperatorName                    string               `json:"operator_name" bson:"operator_name"`
	History3plDeliveryInfo          []History3plDelivery `json:"history_3pl_delivery_info" bson:"history_3pl_delivery_info,omitempty"`
	DeliveryDispatcherPrice         float64              `bson:"delivery_dispatcher_price" json:"delivery_dispatcher_price,omitempty"`
	IsTestOrder                     bool                 `bson:"is_test_order" json:"is_test_order,omitempty"`
	PositionsOnStop                 []PositionsOnStop    `bson:"positions_on_stop" json:"positions_on_stop,omitempty"`
	// Todo temporary `Canceled3PlDeliveryInfo` field for kwaaka report analytics. Delete after a couple of months
	Canceled3PlDeliveryInfo []Cancelled3PLDelivery `bson:"canceled_3pl_delivery_info,omitempty" json:"canceled_3pl_delivery_info,omitempty"`
	IsCashPayment           bool                   `bson:"is_cash_payment" json:"is_cash_payment,omitempty"`
}

type FailReason struct {
	Code         string `bson:"code" json:"code"`
	Message      string `bson:"message" json:"message"`
	BusinessName string `bson:"business_name" json:"business_name"`
	Reason       string `bson:"reason" json:"reason"`
	Solution     string `bson:"solution" json:"solution"`
}

type PosPaymentInfo struct {
	PaymentTypeID            string `bson:"payment_type_id" json:"payment_type_id"`
	PaymentTypeKind          string `bson:"payment_type_kind" json:"iiko_payment_type_kind"`
	PromotionPaymentTypeID   string `bson:"promotion_payment_type_id"`
	OrderType                string `bson:"order_type" json:"order_type"`
	OrderTypeService         string `bson:"order_type_service" json:"order_type_service"`
	OrderTypeForVirtualStore string `bson:"order_type_for_virtual_store" json:"order_type_for_virtual_store"`
	IsProcessedExternally    *bool  `bson:"is_processed_externally,omitempty" json:"is_processed_externally"`
}

type LogLinks struct {
	LogStreamLink          string `bson:"log_stream_link" json:"log_stream_link"`
	LogStreamLinkByOrderId string `bson:"log_stream_link_by_order_id" json:"log_stream_link_by_order_id"`
}

type LogMessages struct {
	FromDelivery string `bson:"from_delivery" json:"from_delivery"`
	ToPos        string `bson:"to_pos" json:"to_pos"`
}

type Error struct {
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Code      int       `bson:"error_code" json:"error_code"`
	Message   string    `bson:"error_message" json:"error_message"`
}

type PosDiscountsInfo struct {
	Discounts []PosDiscount `bson:"discounts" json:"discounts"`
}

type PosDiscount struct {
	DiscountTypeId string `bson:"discountTypeId" json:"discount_type_id"`
	Type           string `bson:"type" json:"type"`
}

type Promo struct {
	Type           string `bson:"type" json:"type"`
	Discount       int    `bson:"discount" json:"discount"`
	Percent        int    `bson:"percent" json:"percent"`
	IikoDiscountId string `json:"iiko_discount_id" bson:"iiko_discount_id"`
}

type ShaurmaFoodOrdersInfo struct {
	RestaurantID    string              `bson:"restaurant_id" json:"restaurant_id,omitempty"`
	DeliveryService string              `bson:"delivery_service" json:"delivery_service,omitempty"`
	Products        []OrderProduct      `bson:"products" json:"products,omitempty"`
	StatusesHistory []OrderStatusUpdate `bson:"statuses_history" json:"statuses_history,omitempty"`
	OrderCode       string              `bson:"order_code" json:"order_code,omitempty"`
	PosOrderID      string              `bson:"pos_order_id" json:"pos_order_id,omitempty"`
	CreationResult  CreationResult      `bson:"creation_result" json:"creation_result,omitempty"`
	PosPaymentInfo  PosPaymentInfo      `bson:"pos_payment_info" json:"pos_payment_info,omitempty"`
	Errors          []Error             `bson:"errors" json:"errors,omitempty"`
}

type ShaurmaFoodOrdersInfoResponse struct {
	Orders          []ShaurmaFoodOrdersInfo `json:"orders"`
	PagesTotalCount int                     `json:"pages_total_count"`
}

type Review struct {
	ReviewContent string  `bson:"review,omitempty" json:"review,omitempty"`
	Rating        float32 `bson:"rating,omitempty" json:"rating,omitempty"`
}

type Proposal struct {
	Price               models3.Price `json:"price" bson:"price"`
	TimeEstimateMinutes int           `json:"time_estimate_minutes" bson:"time_estimate_minutes"`
	ProviderService     string        `json:"provider_service" bson:"provider_service"`
	Priority            int           `json:"priority" bson:"priority"`
}

type History3plDelivery struct {
	DeliveryOrderID            string          `bson:"delivery_order_id" json:"delivery_order_id"`
	DeliveryDispatcher         string          `bson:"delivery_dispatcher" json:"delivery_dispatcher"`
	DeliveryDispatcherPrice    float64         `bson:"delivery_dispatcher_price" json:"delivery_dispatcher_price"`
	FullDeliveryPrice          float64         `bson:"full_delivery_price" json:"full_delivery_price"`
	RestaurantPayDeliveryPrice float64         `bson:"restaurant_pay_delivery_price" json:"restaurant_pay_delivery_price"`
	KwaakaChargedDeliveryPrice float64         `bson:"kwaaka_charged_delivery_price" json:"delivery_service_fee"`
	DeliveryAddress            DeliveryAddress `bson:"delivery_address" json:"delivery_address"`
	Customer                   Customer        `bson:"customer" json:"customer"`
}

type OrderInfoForTelegramMsg struct {
	RestaurantName        string `json:"restaurant_name"`
	RestaurantAddress     string `json:"restaurant_address"`
	RestaurantPhoneNumber string `json:"restaurant_phone_number"`
	OrderId               string `json:"order_id"`
	Id3plOrder            string `json:"id_3pl_order"`
	CustomerName          string `json:"customer_name"`
	CustomerPhoneNumber   string `json:"customer_phone_number"`
	CustomerAddress       string `json:"customer_address"`
	DeliveryService       string `json:"delivery_service"`
}

type ActualizeSelfDeliveryStatusRequest struct {
	IikoStatus          string `json:"iiko_status"`
	CustomerPhoneNumber string `json:"customer_phone_number"`
	StoreID             string `json:"store_id"`
}

// Todo temporary
type Cancelled3PLDelivery struct {
	DeliveryOrderID            string  `bson:"delivery_order_id,omitempty" json:"delivery_order_id,omitempty"`
	DeliveryDispatcher         string  `bson:"delivery_dispatcher,omitempty" json:"delivery_dispatcher,omitempty"`
	DeliveryDispatcherPrice    float64 `bson:"delivery_dispatcher_price,omitempty" json:"delivery_dispatcher_price,omitempty"`
	FullDeliveryPrice          float64 `bson:"full_delivery_price,omitempty" json:"full_delivery_price,omitempty"`
	RestaurantPayDeliveryPrice float64 `bson:"restaurant_pay_delivery_price,omitempty" json:"restaurant_pay_delivery_price,omitempty"`
	KwaakaChargedDeliveryPrice float64 `bson:"kwaaka_charged_delivery_price,omitempty" json:"delivery_service_fee,omitempty"`
}

type DeliveryHistory struct {
	DeliveryOrderID            string  `bson:"delivery_order_id,omitempty" json:"delivery_order_id,omitempty"`
	DeliveryDispatcher         string  `bson:"delivery_dispatcher,omitempty" json:"delivery_dispatcher,omitempty"`
	DeliveryDispatcherPrice    float64 `bson:"delivery_dispatcher_price,omitempty" json:"delivery_dispatcher_price,omitempty"`
	FullDeliveryPrice          float64 `bson:"full_delivery_price,omitempty" json:"full_delivery_price,omitempty"`
	RestaurantPayDeliveryPrice float64 `bson:"restaurant_pay_delivery_price,omitempty" json:"restaurant_pay_delivery_price,omitempty"`
	KwaakaChargedDeliveryPrice float64 `bson:"kwaaka_charged_delivery_price,omitempty" json:"delivery_service_fee,omitempty"`
	Status                     string  `json:"status" bson:"status"`
}

type SaveToHistory struct {
	DeliveryAddress DeliveryAddress `json:"delivery_address"`
	Customer        Customer        `json:"customer"`
}

type PositionsOnStop struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}
