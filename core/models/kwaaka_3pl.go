package models

import (
	"time"
)

type SetKwaaka3plDispatcherRequest struct {
	OrderID                string          `json:"order_id,omitempty"`
	Dispatcher             string          `json:"dispatcher,omitempty"`
	DeliveryOrderPromiseID string          `json:"delivery_order_promise_id,omitempty"`
	FullDeliveryPrice      float64         `json:"full_delivery_price,omitempty"`
	DeliveryAddress        DeliveryAddress `json:"delivery_address"`
	Customer               Customer        `json:"customer"`
}

type GetDeliveryInfoResp struct {
	DeliveryID    string                      `json:"delivery_id"`
	Statuses      []GetDeliveryStatus         `json:"statuses"`
	DeliveryOrder GetDeliveryOrderTrackingUrl `json:"delivery_order"`
	CookingTime   int                         `json:"cooking_time"`
	BusyMode      bool                        `json:"busy_mode"`
	CancelState   string                      `json:"cancel_state"`
}

type GetDeliveryStatus struct {
	Status      string    `json:"status" bson:"status"`
	Description string    `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

type GetDeliveryOrderTrackingUrl struct {
	TrackingUrl  string  `json:"tracking_url" bson:"tracking_url"`
	Longitude    float64 `json:"longitude" bson:"longitude"`
	Latitude     float64 `json:"latitude" bson:"latitude"`
	CourierPhone string  `json:"phone_number" bson:"phone_number"`
}

type DeliveryService struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type GetOrdersByCustomerPhoneRequest struct {
	CustomerPhone string `json:"customer_phone"`
	RestaurantID  string `json:"restaurant_id"`
	Pagination
}

type GetOrdersByCustomerPhoneResponse struct {
	Orders               []OrderByCustomerPhone `json:"orders"`
	CustomerOrderHistory CustomerOrderHistory   `json:"customer_order_history"`
}
type CustomerOrderHistory struct {
	Phone       string    `json:"phone"`
	Name        string    `json:"name"`
	AverageBill float64   `json:"average_bill"`
	Amount      int       `json:"amount"`
	LastOrder   time.Time `json:"last_order"`
}

type Pagination struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}
type OrderByCustomerPhone struct {
	ID                  string              `json:"_id"`
	DeliveryService     string              `json:"delivery_service"`
	RestaurantID        string              `json:"restaurant_id"`
	RestaurantName      string              `json:"restaurant_name"`
	OrderID             string              `json:"order_id"`
	OrderCode           string              `json:"order_code"`
	Status              string              `json:"status"`
	StatusesHistory     []OrderStatusUpdate `json:"statuses_history"`
	OrderTime           TransactionTime     `json:"order_time"`
	Customer            Customer            `json:"customer"`
	Products            []OrderProduct      `json:"products"`
	EstimatedTotalPrice Price               `json:"estimated_total_price"`
}

type Delivery3plOrder struct {
	Id                 string                      `bson:"_id,omitempty" json:"_id,omitempty"`
	DeliveryService    string                      `bson:"delivery_service" json:"delivery_service"`
	DeliveryExternalID string                      `bson:"delivery_external_id" json:"delivery_external_id"`
	CreatedAt          time.Time                   `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time                   `bson:"updated_at" json:"updated_at"`
	Courier            GetDeliveryOrderTrackingUrl `bson:"courier" json:"courier"`
	Status             string                      `bson:"status" json:"status"`
	StatusHistory      []GetDeliveryStatus         `bson:"status_history" json:"status_history"`
	VersionID          int                         `bson:"version_id" json:"version_id"`
	CancelState        string                      `bson:"cancel_state" json:"cancel_state"`
	ChangesHistory     ChangesHistory              `bson:"changes_history" json:"changes_history"`
}

type ChangesHistory struct {
	Username  string    `bson:"username" json:"username"`
	Action    string    `bson:"action" json:"action"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Todo Temporary
type DeliveryOrderPQ struct {
	Id                 string            `json:"id"`
	DeliveryService    DeliveryServicePQ `json:"delivery_service"`
	DeliveryExternalId string            `json:"delivery_external_id"`
	VersionId          int               `json:"version_id"`
	CreatedAt          time.Time         `json:"created_at"`
	Courier            CourierPQ         `json:"courier"`
	Address            Address           `json:"address"`
}

// Todo Temporary
type CourierPQ struct {
	Id          string `json:"-"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	TrackingUrl string `json:"tracking_url"`
}

// Todo Temporary
type CourierLocations struct {
	TrackingUrl string `json:"tracking_url"`
}

// Todo Temporary
type DeliveryServicePQ struct {
	ID        string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Todo Temporary
type Address struct {
	ID                string    `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID        string    `json:"customer_id" bson:"customer_id"`
	Latitude          float64   `json:"latitude" bson:"latitude"`
	Longitude         float64   `json:"longitude" bson:"longitude"`
	Label             string    `json:"label" bson:"label"`
	City              string    `json:"city" bson:"city"`
	Street            string    `json:"street" bson:"street"`
	LocationType      string    `json:"location_type" bson:"location_type"`
	Building          string    `json:"building" bson:"building"`
	Floor             string    `json:"floor" bson:"floor"`
	Apartment         string    `json:"apartment" bson:"apartment"`
	Entrance          string    `json:"entrance" bson:"entrance"`
	DoorBellInfo      string    `json:"door_bell_info" bson:"doorBellInfo"`
	Comment           string    `json:"comment" bson:"comment"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
	RestaurantGroupID string    `json:"restaurant_group_id" bson:"restaurant_group_id"`
}

// Todo Temporary
type DeliveryStatusHistory struct {
	ID          string    `json:"-" db:"id"`
	Status      string    `json:"status" db:"status"`
	DeliveryId  string    `json:"delivery_id" db:"delivery_id"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func ToModel(deliveryStatuses []DeliveryStatusHistory) []GetDeliveryStatus {
	var res = make([]GetDeliveryStatus, 0, len(deliveryStatuses))

	for _, deliveryStatus := range deliveryStatuses {
		res = append(res, GetDeliveryStatus{
			Status:      deliveryStatus.Status,
			Description: deliveryStatus.Description,
			CreatedAt:   deliveryStatus.CreatedAt,
		})
	}
	return res
}
