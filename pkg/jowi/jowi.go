package jowi

import (
	"context"
	dto2 "github.com/kwaaka-team/orders-core/pkg/jowi/client/dto"
)

type Config struct {
	Protocol string
	BaseURL  string
	Insecure bool

	ApiKey    string
	ApiSecret string
	Sig       string
}

type Jowi interface {
	GetRestaurants(ctx context.Context) (dto2.ResponseListRestaurants, error)
	GetStopList(ctx context.Context, restaurantID string) (dto2.ResponseStopList, error)
	GetCourses(ctx context.Context, restaurantID string) (dto2.ResponseCourse, error)
	GetCourseCategories(ctx context.Context, restaurantID string) (dto2.ResponseCourseCategory, error)

	CreateOrder(ctx context.Context, order dto2.RequestCreateOrder) (dto2.ResponseOrder, error)
	GetOrder(ctx context.Context, restaurantID, orderID string) (dto2.ResponseOrder, error)
	CancelOrder(ctx context.Context, orderID string, cancelOrder dto2.RequestCancelOrder) (dto2.ResponseOrder, error)
}
