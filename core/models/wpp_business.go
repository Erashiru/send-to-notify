package models

type SendUserRequest struct {
	PhoneNum          string `json:"phone_num"`
	RestaurantGroupId string `json:"restaurant_group_id"`
}

type RedisCodeRestaurantGroup struct {
	Code              string `json:"code"`
	RestaurantGroupId string `json:"restaurant_group_id"`
}
