package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/config"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type ApiToken interface {
	FindStores(ctx context.Context, query selector.ApiToken) ([]models.Store, error)
}

type ApiTokenManager struct {
	globalConfig       config.Configuration
	apiTokenRepository drivers.ApiTokensRepository
}

func NewApiTokenManager(
	globalConfig config.Configuration,
	ds drivers.Datastore) ApiToken {
	return &ApiTokenManager{
		globalConfig:       globalConfig,
		apiTokenRepository: ds.ApiTokensRepository(),
	}
}

func (s *ApiTokenManager) FindStores(ctx context.Context, query selector.ApiToken) ([]models.Store, error) {
	stores, err := s.apiTokenRepository.GetStores(ctx, query)
	if err != nil {
		return nil, err
	}

	return stores, nil
}
