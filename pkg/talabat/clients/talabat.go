package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/talabat/models"
)

type Config struct {
	Protocol string
	BaseURL  string
	Username string
	Password string
}

type TalabatMenu interface {
	CreateNewMenu(ctx context.Context, req models.CreateNewMenuRequest) error
	GetRequestStatus(ctx context.Context, requestID string) (models.GetRequestStatusResponse, error)
	Close(ctx context.Context)
	UpdateItemsAvailability(ctx context.Context, req models.UpdateItemsAvailabilityRequest) error
}

type TalabatMW interface {
	Close(ctx context.Context)
	AcceptOrder(ctx context.Context, req models.AcceptOrderRequest) error
	RejectOrder(ctx context.Context, req models.RejectOrderRequest) error
	OrderPickedUp(ctx context.Context, req models.OrderPickedUpRequest) error
	MarkOrderPrepared(ctx context.Context, orderToken string) error
	SubmitCatalog(ctx context.Context, req models.SubmitCatalogRequest) (models.SubmitCatalogResponse, error)
}
