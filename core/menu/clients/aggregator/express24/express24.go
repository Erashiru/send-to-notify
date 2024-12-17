package express24

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/errors"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/utils"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	expressConf "github.com/kwaaka-team/orders-core/pkg/express24"
	expressCli "github.com/kwaaka-team/orders-core/pkg/express24/clients"
	"github.com/kwaaka-team/orders-core/pkg/express24/clients/dto"
	"github.com/rs/zerolog/log"
	"strconv"
)

type Express24 interface {
	BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error)
	ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error)
	UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error)
	VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error)
	BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error)
	GetMenu(ctx context.Context, extStoreId string) (models.Menu, error)
	ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error)
}

type mnm struct {
	cli      expressCli.Express24
	s3Bucket menu.S3_BUCKET
}

func NewManager(ctx context.Context, cfg menu.Configuration) (Express24, error) {

	cli, err := expressConf.NewExpress24Client(&expressCli.Config{
		Protocol: "http",
		BaseURL:  cfg.Express24Configuration.BaseURL,
		Username: cfg.Express24Configuration.Username,
		Password: cfg.Express24Configuration.Password,
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize Express24 client ")
		return nil, err
	}
	return &mnm{
		cli:      cli,
		s3Bucket: cfg.S3_BUCKET,
	}, nil
}

func (m mnm) BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error) {
	branch, err := strconv.Atoi(storeID)
	if err != nil {
		return "", err
	}

	utils.Beautify("stoplist update request send", branch)

	resp, err := m.cli.UpdateProducts(ctx, dto.UpdateProductsRequest{
		Data: dto.UpdateProductData{
			Branches: []int{branch},
			Products: toProducts(products),
		},
	})

	if err != nil {
		return "", err
	}

	if resp.Failed != nil {
		for _, fail := range resp.Failed {
			log.Info().Msgf("PRODUCT ID: %s, ERROR MESSAGE: %s", fail.ExternalId, fail.Message)
		}
	}

	return "", nil
}

func (m mnm) ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error) {
	return models.MenuUploadTransaction{}, errors.ErrNotImplemented
}

func (m mnm) ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error) {
	return models.ProductModifyResponse{}, errors.ErrNotImplemented
}

func (m mnm) UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error) {
	return models.ExtTransaction{}, errors.ErrNotImplemented
}

func (m mnm) VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error) {
	return "", errors.ErrNotImplemented
}

func (m mnm) BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error) {
	return "", errors.ErrNotImplemented
}

func (m mnm) GetMenu(ctx context.Context, extStoreId string) (models.Menu, error) {
	return models.Menu{}, errors.ErrNotImplemented

}
