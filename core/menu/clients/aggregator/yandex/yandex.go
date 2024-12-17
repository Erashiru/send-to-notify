package yandex

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	yandexClient "github.com/kwaaka-team/orders-core/pkg/yandex"
	yandexConfig "github.com/kwaaka-team/orders-core/pkg/yandex/clients"
	yandexModels "github.com/kwaaka-team/orders-core/pkg/yandex/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type mnm struct {
	cli       yandexConfig.Yandex
	storeRepo drivers.StoreRepository
	menuRepo  drivers.MenuRepository
}

func NewManager(ctx context.Context, cfg menu.Configuration, store storeModels.Store, storeRepo drivers.StoreRepository, menuRepo drivers.MenuRepository) (*mnm, error) {
	cli, err := yandexClient.NewClient(&yandexConfig.Config{
		Protocol:     "http",
		BaseURL:      cfg.YandexConfiguration.BaseURL,
		ClientID:     cfg.YandexConfiguration.ClientID,
		ClientSecret: cfg.YandexConfiguration.ClientSecret,
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize yandex client")
		return nil, err
	}

	return &mnm{
		cli:       cli,
		storeRepo: storeRepo,
		menuRepo:  menuRepo,
	}, nil
}

func (m mnm) GetMenu(ctx context.Context, extStoreId string) (models.Menu, error) {
	return models.Menu{}, errors.New("not implemented")
}

func (m mnm) BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error) {
	return "", nil
}

func (m mnm) ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error) {
	return models.ProductModifyResponse{}, errors.New("not implemented")
}

func (m mnm) BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error) {
	return "", errors.New("not implemented")
}

func (m mnm) VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error) {
	return "", errors.New("not implemented")
}

func (m mnm) UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error) {
	err := m.cli.MenuImportInitiation(ctx, yandexModels.MenuInitiationRequest{
		RestaurantID:  extStoreId,
		OperationType: "menu",
	})

	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{err.Error()},
		}, err
	}

	return models.ExtTransaction{
		Status:     models.SUCCESS.String(),
		MenuID:     menuId,
		ExtStoreID: extStoreId,
	}, nil
}

func (m mnm) ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error) {
	return models.MenuUploadTransaction{}, nil
}
