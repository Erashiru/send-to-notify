package external

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	externalConf "github.com/kwaaka-team/orders-core/pkg/externalapi"
	externalCli "github.com/kwaaka-team/orders-core/pkg/externalapi/clients"
	"github.com/rs/zerolog/log"
)

type External interface {
	ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error)
	UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error)
	BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error)
	VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error)
	BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error)
	GetMenu(ctx context.Context, extStoreId string) (models.Menu, error)
	ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error)
}

type mnm struct {
	cli                      externalCli.Client
	webhookURL               string
	webhookProductStoplist   string
	webhookAttributeStoplist string
	storeID                  string
	globalConfig             menu.Configuration
}

func NewManager(ctx context.Context, cfg menu.Configuration, store storeModels.Store, delivery string) (External, error) {
	// Initialize new External client
	var (
		authToken                string
		webhookURL               string
		webhookProductStoplist   string
		webhookAttributeStoplist string
		storeID                  string
	)

	for _, config := range store.ExternalConfig {
		if config.Type == delivery && config.WebhookURL != "" {
			authToken = config.AuthToken
			webhookURL = config.WebhookURL
			webhookProductStoplist = config.WebhookProductStoplist
			webhookAttributeStoplist = config.WebhookAttributeStoplist

			if len(config.StoreID) == 0 {
				return nil, fmt.Errorf("store id is empty for %s service", delivery)
			}

			storeID = config.StoreID[0]
			break
		}
	}

	if webhookURL == "" {
		log.Trace().Err(ErrNoWebhookSubscription).Msgf("%s has not webhook subscription for %s delivery", store.Name, delivery)
		return nil, ErrNoWebhookSubscription
	}

	cli, err := externalConf.NewWebhookClient(&externalCli.Config{
		Protocol:  "http",
		AuthToken: authToken,
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize External client")
		return nil, err
	}

	if webhookProductStoplist == "" {
		webhookProductStoplist = webhookURL
	}

	if webhookAttributeStoplist == "" {
		webhookAttributeStoplist = webhookURL
	}

	return &mnm{
		cli:                      cli,
		storeID:                  storeID,
		globalConfig:             cfg,
		webhookURL:               webhookURL,
		webhookProductStoplist:   webhookProductStoplist,
		webhookAttributeStoplist: webhookAttributeStoplist,
	}, nil
}

func (m mnm) ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error) {
	return models.ProductModifyResponse{}, ErrNotImplemented
}

func (m mnm) UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error) {
	return models.ExtTransaction{}, ErrNotImplemented
}

func (m mnm) VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error) {
	return "", ErrNotImplemented
}

func (m mnm) GetMenu(ctx context.Context, extStoreId string) (models.Menu, error) {
	return models.Menu{}, ErrNotImplemented
}

func (m mnm) ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error) {
	return models.MenuUploadTransaction{}, ErrNotImplemented
}
