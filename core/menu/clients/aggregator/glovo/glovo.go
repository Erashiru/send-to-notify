package glovo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/custom"
	"github.com/kwaaka-team/orders-core/core/menu/models/pointer"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	pkgUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	"time"

	glovoConf "github.com/kwaaka-team/orders-core/pkg/glovo"
	glovoCli "github.com/kwaaka-team/orders-core/pkg/glovo/clients"
	glovoModels "github.com/kwaaka-team/orders-core/pkg/glovo/clients/dto"
	"github.com/rs/zerolog/log"
)

type Glovo interface {
	ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error)
	UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error)
	BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error)
	VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error)
	BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error)
	GetMenu(ctx context.Context, extStoreId string) (models.Menu, error)
	ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error)
}

type mnm struct {
	cli      glovoCli.Glovo
	s3Bucket menu.S3_BUCKET
}

func NewManager(ctx context.Context, cfg menu.Configuration) (Glovo, error) {

	cli, err := glovoConf.NewGlovoClient(&glovoCli.Config{
		Protocol: "http",
		BaseURL:  cfg.GlovoConfiguration.BaseURL,
		ApiKey:   cfg.GlovoConfiguration.Token,
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize Glovo client ")
		return nil, err
	}
	return &mnm{
		cli:      cli,
		s3Bucket: cfg.S3_BUCKET,
	}, nil
}

func (m mnm) BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error) {
	res, err := m.cli.BulkUpdate(ctx, storeID, glovoModels.BulkUpdateRequest{
		Attributes: toAttributes(storeID, attributes),
	})

	if err != nil {
		return "", err
	}

	return res, nil
}

func (m mnm) GetMenu(ctx context.Context, extStoreId string) (models.Menu, error) {
	panic("implement me")
}

func (m mnm) ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error) {

	req := glovoModels.ProductModifyRequest{
		ID:          product.ExtID,
		StoreID:     storeID,
		IsAvailable: pointer.OfBool(product.IsAvailable),
	}

	// sometimes we don't need to change price of products
	if len(product.Price) != 0 {
		if product.Price[0].Value != 0 {
			req.Price = pointer.OfFloat64(product.Price[0].Value)
		}
	}

	resp, err := m.cli.ModifyProduct(ctx, req)
	if err != nil {
		return models.ProductModifyResponse{}, err
	}

	// Todo unmarshalling response
	return productModifierFromClient(resp), nil
}

func (m mnm) BulkUpdate(ctx context.Context, restaurantId, storeID string, reqProducts models.Products, reqAttributes models.Attributes, store storeModels.Store) (string, error) {

	res, err := m.cli.BulkUpdate(ctx, storeID, glovoModels.BulkUpdateRequest{
		Products:   toProducts(storeID, reqProducts),
		Attributes: toAttributes(storeID, reqAttributes),
	})

	if err != nil {
		return "", err
	}

	return res, nil
}

func (m mnm) UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error) {
	menuGlovo := ToGlovoMenu(menu, offers)

	objectUrl, err := pkgUtils.UploadMenuToS3(store.ID, m.s3Bucket.KwaakaMenuFilesBucket, models.GLOVO.String(), m.s3Bucket.ShareMenuBaseUrl, menuGlovo, sv3)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{err.Error()},
		}, err
	}

	res, err := m.cli.UploadMenu(ctx, glovoModels.UploadMenuRequest{
		MenuURL: objectUrl,
		StoreId: extStoreId,
	})
	if err != nil {
		log.Err(err).Msgf("can't publicate glovo menu")
		return models.ExtTransaction{
			ID:         res.TransactionID,
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{err.Error()},
			MenuUrl:    objectUrl,
		}, err
	}

	sCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	status, err := m.VerifyMenu(sCtx, models.ExtTransaction{
		ID:         res.TransactionID,
		ExtStoreID: extStoreId,
	})
	if err != nil {
		log.Err(err).Msgf("glovo upload verify menu err: menu_id %s in store_id %s", menuId, store.ID)
	}

	res.Status = glovoModels.Status(status)

	return uploadMenuToTransaction(res, menuId, extStoreId, objectUrl), err
}

// VerifyMenu check menu upload
func (m mnm) VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error) {

	resp, err := m.cli.VerifyMenu(ctx, transaction.ExtStoreID, transaction.ID)
	if err != nil {
		return models.NOT_PROCESSED, err
	}

	if !validStatus(resp.Status, []glovoModels.Status{
		glovoModels.PROCESSING,
		glovoModels.SUCCESS,
	}) {
		return models.NOT_PROCESSED, nil
	}

	return models.Status(resp.Status), nil
}

// ValidStatus used to paused and completed installments
func validStatus(status glovoModels.Status, validStatuses []glovoModels.Status) bool {

	for _, valid := range validStatuses {
		if valid.String() == status.String() {
			return true
		}
	}
	return false
}

func (m mnm) ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error) {
	menuGlovo := ToGlovoMenu(request.Menu, request.OffersBK)

	res, err := m.cli.ValidateMenu(ctx, menuGlovo)
	if err != nil {
		return request.MenuUploadTransaction, err
	}

	if !res.Valid {
		var errs custom.Error
		for i := 0; i < len(request.MenuUploadTransaction.ExtTransactions); i++ {
			request.MenuUploadTransaction.ExtTransactions[i].Details = res.Errors
		}

		for _, r := range res.Errors {
			errs.Append(fmt.Errorf(r))
		}
		return request.MenuUploadTransaction, errs
	}
	return request.MenuUploadTransaction, nil
}
