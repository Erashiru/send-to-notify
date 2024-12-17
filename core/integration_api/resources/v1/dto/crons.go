package dto

type UpdateOrderStatusCronRequest struct {
	PosTypes []string `json:"pos_types"`
}

type UpdateStopListCronRequest struct {
	PosTypes []string `json:"pos_types"`
}

type UpdateStopListBySectionCronRequest struct {
	WoltSectionIDs    []string `json:"wolt_ids"`
	GlovoSectionIDs   []string `json:"glovo_ids"`
	YandexSectionIDs  []string `json:"yandex_ids"`
	RestaurantGroupID string   `json:"restaurant_group_id"`
	IsAvailable       *bool    `json:"is_available"`
}

type SendDeferSumbissionRequest struct {
	RestaurantGroupId string `json:"restaurant_group_id"`
}
