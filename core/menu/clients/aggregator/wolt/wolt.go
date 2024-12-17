package wolt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/errors"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/custom"
	"github.com/kwaaka-team/orders-core/core/menu/models/pointer"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	pkgUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	woltCli "github.com/kwaaka-team/orders-core/pkg/wolt"
	wolt "github.com/kwaaka-team/orders-core/pkg/wolt/clients"
	woltModels "github.com/kwaaka-team/orders-core/pkg/wolt/clients/dto"
	"github.com/rs/zerolog/log"
)

const queueName = "wolt-discount-run"

type Wolt interface {
	ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error)
	UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error)
	BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error)
	VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error)
	BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error)
	GetMenu(ctx context.Context, extStoreId string) (models.Menu, error)
	ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error)
}

type mnm struct {
	cli      wolt.Wolt
	s3Bucket menu.S3_BUCKET
	sqsCli   notifyQueue.SQSInterface
}

type Config struct {
	Username, Password string
	ApiKey             string
}

func NewManager(ctx context.Context, cfg menu.Configuration) (Wolt, error) {

	cli, err := woltCli.NewWoltClient(&wolt.Config{
		Protocol: "http",
		BaseURL:  cfg.WoltConfiguration.BaseURL,
		ApiKey:   cfg.WoltConfiguration.ApiKey,
		Username: cfg.WoltConfiguration.Username,
		Password: cfg.WoltConfiguration.Password,
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize wolt client.")
		return nil, err
	}

	sqsCli := notifyQueue.NewSQS(sqs.NewFromConfig(cfg.AwsConfig))

	return &mnm{
		cli:      cli,
		s3Bucket: cfg.S3_BUCKET,
		sqsCli:   sqsCli,
	}, nil
}

func (m mnm) GetMenu(ctx context.Context, woltStoreId string) (models.Menu, error) {
	woltMenu, err := m.cli.GetMenu(ctx, woltStoreId)
	if err != nil {
		return models.Menu{}, nil
	}

	return fromWoltToSystemMenu(woltMenu), nil
}

func (m mnm) BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error) {
	updateAttributes := make([]woltModels.UpdateAttribute, 0, len(attributes))

	for _, attribute := range attributes {
		req := woltModels.UpdateAttribute{
			ExtID:       attribute.ExtID,
			IsAvailable: pointer.OfBool(attribute.IsAvailable),
		}

		updateAttributes = append(updateAttributes, req)
	}

	tx, err := m.cli.BulkAttribute(ctx, storeID, woltModels.UpdateAttributes{
		Attribute: updateAttributes,
	})

	if err != nil {
		return "", err
	}

	return tx, nil
}

func (m mnm) ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error) {

	updateProducts := make([]woltModels.UpdateProduct, 0, 1)

	req := woltModels.UpdateProduct{
		ExtID:       product.ExtID,
		IsAvailable: pointer.OfBool(product.IsAvailable),
	}

	updateProducts = append(updateProducts, req)

	tx, err := m.cli.BulkUpdate(ctx, storeID, woltModels.UpdateProducts{
		Product: updateProducts,
	})

	if err != nil {
		return models.ProductModifyResponse{}, err
	}

	return models.ProductModifyResponse{
		ExtID:       product.ExtID,
		IsAvailable: product.IsAvailable,
		Msg:         tx,
	}, nil
}

func (m mnm) BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error) {
	updateProducts := make([]woltModels.UpdateProduct, 0, len(products))

	for _, product := range products {
		req := woltModels.UpdateProduct{
			ExtID:       product.ExtID,
			IsAvailable: pointer.OfBool(product.IsAvailable),
		}

		updateProducts = append(updateProducts, req)
	}

	tx, err := m.cli.BulkUpdate(ctx, storeID, woltModels.UpdateProducts{
		Product: updateProducts,
	})

	if err != nil {
		return "", err
	}

	return tx, nil
}

func (m mnm) VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error) {

	log.Info().Msg("verify menu in wolt not found")

	if transaction.Status != "" {
		return models.Status(transaction.Status), nil
	}

	return models.ERROR, nil
}

func (m mnm) UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error) {

	if menu.HasWoltPromo && userRole != "Admins" {
		return models.ExtTransaction{}, errors.ErrNoPermissionForPublishWoltMenu
	}

	menuWolt := toWoltMenu(store, menu)

	objectUrl, err := pkgUtils.UploadMenuToS3(store.ID, m.s3Bucket.KwaakaMenuFilesBucket, models.WOLT.String(), m.s3Bucket.ShareMenuBaseUrl, menuWolt, sv3)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{err.Error()},
		}, err
	}

	err = m.cli.UploadMenu(ctx, menuWolt, extStoreId)
	if err != nil {
		log.Err(err).Msgf("can't publicate wolt menu")
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{err.Error()},
			MenuUrl:    objectUrl,
		}, err
	}

	if err := m.runDiscountNotification(ctx, store, menu); err != nil {
		log.Err(err).Msgf("send notification for run discount error, store id: %s, menu id: %s", store.ID, menu.ID)
	}

	return models.ExtTransaction{
		Status:     models.SUCCESS.String(),
		MenuID:     menuId,
		ExtStoreID: extStoreId,
		Details:    []string{},
		MenuUrl:    objectUrl,
	}, nil
}

func (m mnm) ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error) {
	if request.MenuUploadTransaction.Status == models.ERROR.String() {
		var errs custom.Error
		for _, extTrans := range request.MenuUploadTransaction.ExtTransactions {
			for _, detail := range extTrans.Details {
				errs.Append(fmt.Errorf(detail))
			}
		}
		return request.MenuUploadTransaction, errs
	}
	return request.MenuUploadTransaction, nil
}

func (m mnm) runDiscountNotification(ctx context.Context, store storeModels.Store, menu models.Menu) error {
	if !m.hasDiscount(menu) {
		log.Info().Msgf("doesn't have discount, store id: %s, menu id: %s", store.ID, menu.ID)
		return nil
	}

	queueURL, err := m.sqsCli.GetQueueURL(ctx, queueName)
	if err != nil {
		return err
	}

	req := struct {
		RestaurantID string `json:"restaurant_id"`
		MenuID       string `json:"menu_id"`
	}{
		RestaurantID: store.ID,
		MenuID:       menu.ID,
	}

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := m.sqsCli.SendSQSMessage(ctx, queueURL, string(b)); err != nil {
		return err
	}

	log.Info().Msgf("send notification for run discount success, store id: %s, menu id: %s", store.ID, menu.ID)

	return nil
}

func (m mnm) hasDiscount(menu models.Menu) bool {
	for _, product := range menu.Products {
		if product.DiscountPrice.IsActive {
			return true
		}
	}
	return false
}
