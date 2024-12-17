package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/starterapp/clients/dto"
)

type Config struct {
	Protocol string
	BaseURL  string
	ApiKey   string
}

type StarterApp interface {
	CreateCategories(ctx context.Context, req []dto.CategoryRequest) (dto.CreateMenuResponse, error)
	UpdateCategories(ctx context.Context, req []dto.CategoryRequest) error
	CreateModifierGroups(ctx context.Context, req []dto.ModifierGroupRequest) (dto.CreateMenuResponse, error)
	UpdateModifierGroups(ctx context.Context, req []dto.ModifierGroupRequest) error
	CreateModifiers(ctx context.Context, req []dto.ModifiersRequest) (dto.CreateMenuResponse, error)
	UpdateModifiers(ctx context.Context, req []dto.ModifiersRequest) error
	CreateMeals(ctx context.Context, req []dto.MealRequest) (dto.CreateMenuResponse, error)
	UpdateMeals(ctx context.Context, req []dto.MealRequest) error
	CreateMealOffers(ctx context.Context, req []dto.MealOfferRequest, shopID int) (dto.CreateMenuResponse, error)
	UpdateMealOffers(ctx context.Context, req []dto.MealOfferRequest, shopID int) error
	CreateModifierOffers(ctx context.Context, req []dto.ModifierOfferRequest) (dto.CreateMenuResponse, error)
	UpdateModifierOffers(ctx context.Context, req []dto.ModifierOfferRequest) error
	SendOrderErrorNotification(ctx context.Context, req dto.SendOrderErrorNotificationRequest, orderID string) error
	ChangeOrderStatus(ctx context.Context, req dto.ChangeOrderStatusRequest, orderID string) error
}
