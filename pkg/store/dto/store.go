package dto

import (
	models2 "github.com/kwaaka-team/orders-core/core/storecore/models"
)

func FromStoreManagementModel(req models2.StoreManagementResponse) StoreManagementResponse {
	return StoreManagementResponse{
		ErrMessage:      req.ErrMessage,
		Success:         req.Success,
		RestaurantId:    req.RestaurantId,
		StoreID:         req.StoreID,
		IsOpen:          req.IsOpen,
		DeliveryService: req.DeliveryService,
	}
}

func FromStoreManagementModels(req []models2.StoreManagementResponse) []StoreManagementResponse {
	var stores []StoreManagementResponse

	for _, store := range req {
		stores = append(stores, FromStoreManagementModel(store))
	}

	return stores
}

type OrderDestination string

type CreateStoreRequest struct {
	Usernames         []string                  `json:"usernames"`
	Name              string                    `json:"name"`
	Address           CreateStoreAddress        `json:"address"`
	Currency          string                    `json:"currency"`
	LanguageCode      string                    `json:"language_code"`
	StoreGroupId      string                    `json:"restaurant_group_id"`
	StoreQRMenuConfig StoreQRMenuConfig         `json:"qrmenu"`
	LegalEntityId     string                    `json:"legal_entity_id"`
	AccountManagerId  string                    `json:"account_manager_id"`
	SalesManagerId    string                    `json:"sales_manager_id"`
	Contacts          []CreateContacts          `json:"contacts"`
	Links             []CreateExternalLinks     `json:"links"`
	Telegram          CreateStoreTelegramConfig `json:"telegram"`
}

type CreateStoreTelegramConfig struct {
	CancelChatID string `json:"cancel_chat_id"`
}

type CreateExternalLinks struct {
	Name      string `json:"name"`
	Url       string `json:"url"`
	ImageLink string `json:"image_link"`
}

type CreateContacts struct {
	FullName string `json:"full_name"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Comment  string `json:"comment"`
}

type CreateStoreAddress struct {
	City        CreateStoreCity `json:"city"`
	Street      string          `json:"street"`
	Coordinates Coordinates     `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CreateStoreCity struct {
	Name     string              `json:"name"`
	Timezone CreateStoreTimezone `json:"timezone"`
}

type CreateStoreTimezone struct {
	Tz        string  `json:"tz"`
	UtcOffset float64 `json:"utc_offset"`
}

type StoreManagementRequest struct {
	StoreInfos      []StoreInfo `json:"store_infos"`
	DeliveryService string      `json:"delivery_service"`
}

type StoreManagementResponse struct {
	ErrMessage      string `json:"err_message"`
	Success         bool   `json:"success"`
	RestaurantId    string `json:"restaurant_id"`
	StoreID         string `json:"store_id"`
	IsOpen          bool   `json:"is_open"`
	DeliveryService string `json:"delivery_service"`
}
type StoreInfo struct {
	RestaurantId string `json:"restaurant_id"`
	StoreId      string `json:"store_id"`
	StoreStatus  bool   `json:"store_status"`
}

type DirectSchedule struct {
	DayOfWeek int                    `json:"day_of_week" bson:"day_of_week"`
	TimeSlots DirectScheduleTimeSlot `json:"time_slots" bson:"time_slots"`
}

type DirectScheduleTimeSlot struct {
	Opening string `json:"opening" bson:"opening"`
	Closing string `json:"closing" bson:"closing"`
}

func (s StoreManagementRequest) ToModel() models2.StoreManagement {
	storeInfo := make([]models2.StoreInfo, 0, len(s.StoreInfos))
	for _, store := range s.StoreInfos {
		storeInfo = append(storeInfo, models2.StoreInfo{
			RestaurantId: store.RestaurantId,
			StoreID:      store.StoreId,
			StoreStatus:  store.StoreStatus,
		})

	}

	return models2.StoreManagement{
		StoreInfo:       storeInfo,
		DeliveryService: s.DeliveryService,
	}
}
