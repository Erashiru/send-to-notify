package dto

import "github.com/kwaaka-team/orders-core/core/externalapi/models"

type GetOrdersResponse struct {
	Infos []OrderInfo `json:"result"`
}

type OrderInfo struct {
	RestaurantID string         `json:"restaurant_id"`
	Orders       []models.Order `json:"orders"`
	Total        int            `json:"total"`
}

type GetOrdersRequest struct {
	RestaurantIDs []string `json:"restaurant_ids"`
	StartDate     string   `json:"start_date"`
	EndDate       string   `json:"end_date"`
}
