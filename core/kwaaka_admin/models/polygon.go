package models

import "github.com/kwaaka-team/orders-core/core/storecore/models"

type PolygonRequest struct {
	RestaurantID string               `json:"restaurant_id,omitempty"`
	Percentage   int                  `json:"percentage_modifier,omitempty"`
	Coordinates  []models.Coordinates `json:"coordinates,omitempty" bson:"coordinates"`
}

type GetPolygonResponse struct {
	Percentage  int
	Coordinates []models.Coordinates
}
