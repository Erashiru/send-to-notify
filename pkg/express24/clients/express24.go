package clients

import (
	"context"
	dto2 "github.com/kwaaka-team/orders-core/pkg/express24/clients/dto"
)

type Config struct {
	Protocol           string
	BaseURL            string
	Insecure           bool
	Username, Password string
}
type Express24 interface {
	GetBranches(ctx context.Context) (dto2.GetBranchesResponse, error)
	UpdateBranches(ctx context.Context, req dto2.UpdateBranchesRequest) (dto2.UpdateBranchesResponse, error)
	UpdateOffers(ctx context.Context, req dto2.UpdateOffersRequest) (dto2.UpdateOffersResponse, error)
	UpdateProducts(ctx context.Context, req dto2.UpdateProductsRequest) (dto2.UpdateProductsResponse, error)
}
