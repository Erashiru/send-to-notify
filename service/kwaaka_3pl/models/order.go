package models

import (
	"time"
)

type Address struct {
	Label        string  `json:"label"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Comment      string  `json:"comment"`
	City         string  `json:"city"`
	BuildingName string  `json:"building_name"`
	Street       string  `json:"street"`
	Porch        string  `json:"porch"`
	Floor        string  `json:"floor"`
	Flat         string  `json:"flat"`
}

type Item struct {
	Name     string  `json:"name"`
	ID       string  `json:"id"`
	Weight   string  `json:"weight"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type StoreInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type CustomerInfo struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type CreateDeliveryRequest struct {
	ID                string       `json:"id"`
	FullDeliveryPrice float64      `json:"full_delivery_price"`
	Provider          string       `json:"provider"`
	Items             []Item       `json:"items"`
	PickUpTime        time.Time    `json:"pick_up_time"`
	DeliveryAddress   Address      `json:"delivery_address"`
	StoreAddress      Address      `json:"store_address"`
	CustomerInfo      CustomerInfo `json:"customer"`
	StoreInfo         StoreInfo    `json:"store_info"`
	PickUpCode        string       `json:"pick_up_code"`
	Comment           string       `json:"comment"`
	Currency          string       `json:"currency"`
	ExternalStoreID   string       `json:"external_store_id"`
	TaxiClass         string       `json:"taxi_class"`
}

type ListProvidersRequest struct {
	Address                   OrderAddress  `json:"address"`
	MinPreparationTimeMinutes int           `json:"min_preparation_time_minutes"` //default:30 min:0 max:60
	ScheduledDropoffTime      time.Time     `json:"scheduled_dropoff_time"`
	ItemsSettings             ItemsSettings `json:"items_settings"`
	RestaurantCoordinates     Coordinates   `json:"restaurant_coordinates"`
	KwaakaChargePercentage    float64       `json:"delivery_service_fee_percentage"`
	KwaakaChargeAbsolut       float64       `json:"delivery_service_fee_absolut"`
	IndriveAvailable          bool          `json:"indrive_available"`
	WoltAvailable             bool          `json:"wolt_drive_available"`
	YandexAvailable           bool          `json:"yandex_available"`
}

type OrderAddress struct {
	Street      string      `json:"street"`
	City        string      `json:"city"`
	PostCode    string      `json:"post_code"`
	Coordinates Coordinates `json:"coordinates"`
	Language    string      `json:"language"`
}

type ItemsSettings struct {
	Quantity int  `json:"quantity"`
	Size     Size `json:"size"`
	Weight   int  `json:"weight"`
}

type Size struct {
	Height float64 `json:"height"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type GetPromise struct {
	Price               Price  `json:"price"`
	TimeEstimateMinutes int    `json:"time_estimate_minutes"`
	ProviderService     string `json:"provider_service"`
	Priority            int    `json:"priority"`
}

type ProviderResponse struct {
	Provider *GetPromise `json:"provider,omitempty"`
	Error    string      `json:"error,omitempty"`
}

type Price struct {
	Amount          int     `json:"amount" bson:"amount"`
	Currency        string  `json:"currency" bson:"currency"`
	KwaakaChargeSum float64 `json:"delivery_service_fee" bson:"kwaaka_charge_sum"`
}

type Delivery struct {
	Dispatcher                 string  `json:"dispatcher" bson:"dispatcher"`
	DeliveryTime               int32   `json:"delivery_time" bson:"delivery_time"`
	ClientDeliveryPrice        float64 `json:"client_delivery_price" bson:"client_delivery_price"`
	FullDeliveryPrice          float64 `json:"full_delivery_price" bson:"full_delivery_price"`
	DropOffScheduleTime        string  `json:"drop_off_schedule_time" bson:"drop_off_schedule_time"`
	KwaakaChargedDeliveryPrice float64 `json:"delivery_service_fee" bson:"kwaaka_charged_delivery_price"`
	Priority                   int     `json:"priority"`
}

const (
	WoltDelivery      = "wolt"
	YandexDelivery    = "yandex"
	IndriveDelivery   = "indrive"
	InProcessing      = "IN_PROCESSING"
	PerformerLookup   = "PERFORMER_LOOKUP"
	ComingToPickup    = "COMING_TO_PICKUP"
	PickedUp          = "PICKED_UP"
	Delivered         = "DELIVERED"
	PayWaiting        = "PAY_WAITING"
	Returning         = "RETURNING"
	Returned          = "RETURNED"
	Failed            = "FAILED"
	Cancelled         = "CANCELLED"
	OrderCreated      = "ORDER_CREATED"
	New               = "NEW"
	CancelUnavailable = "unavailable"

	SetDeliveryDispatcher = "set 3pl dispatcher -> create delivery"
	CancelDelivery        = "cancel delivery"
)

type GetDeliveryDispatcherPricesRequest struct {
	DeliveryIDs []string `json:"delivery_ids"`
}
type GetDeliveryDispatcherPricesResponse struct {
	DeliveryPrices []DeliveryPrice `json:"delivery_prices"`
}

type DeliveryPrice struct {
	DeliveryID         string  `json:"delivery_id"`
	DeliveryExternalID string  `json:"delivery_external_id"`
	Price              float64 `json:"price"`
}
