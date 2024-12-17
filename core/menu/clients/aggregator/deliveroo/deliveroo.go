package deliveroo

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	deliverooConf "github.com/kwaaka-team/orders-core/pkg/deliveroo"
	deliverooCli "github.com/kwaaka-team/orders-core/pkg/deliveroo/clients"
	deliverooModels "github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/dto"
	deliverooHttpCli "github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/http"
	pkgUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type mnm struct {
	cli      deliverooHttpCli.Client
	s3Bucket menu.S3_BUCKET
}

func (m mnm) ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error) {
	panic("implement me")
}

func (m mnm) UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error) {
	menuDeliveroo := toDeliverooMenu(menu)

	objectUrl, err := pkgUtils.UploadMenuToS3(store.ID, m.s3Bucket.KwaakaMenuFilesBucket, models.DELIVEROO.String(), m.s3Bucket.ShareMenuBaseUrl, menuDeliveroo, sv3)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{err.Error()},
		}, err
	}
	menuObj := deliverooModels.Menu{
		Name:     menu.Name,
		MenuData: menuDeliveroo,
		SiteIDs:  nil, //need to define
	}
	err = m.cli.UploadMenu(ctx, menuObj, extStoreId)
	if err != nil {
		log.Err(err).Msgf("can't publicate deliveroo menu %s", menu.ID)
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{err.Error()},
			MenuUrl:    objectUrl,
		}, err
	}

	return models.ExtTransaction{
		Status:     models.SUCCESS.String(),
		MenuID:     menuId,
		ExtStoreID: extStoreId,
		Details:    []string{},
		MenuUrl:    objectUrl,
	}, nil
}

func (m mnm) BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error) {
	stoplistProducts := make([]string, 0, len(products))

	for _, product := range products {
		stoplistProducts = append(stoplistProducts, product.ExtID)
	}

	var menuID string

	for _, menu := range store.Menus {
		if menu.Delivery == models.DELIVEROO.String() {
			menuID = menu.ID
		}
	}

	err := m.cli.UpdateUnavailabileItems(ctx, deliverooModels.UpdateUnavailabilitesRequest{
		SiteID:  storeID,
		BrandID: restaurantID,
		MenuID:  menuID,
		UpdateUnavailabilitesRequestBody: deliverooModels.UpdateUnavailabilitesRequestBody{
			UnavailableIDs: stoplistProducts,
		},
	})
	if err != nil {
		return "", err
	}
	return "", nil
}

func (m mnm) VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error) {
	return "", errors.New("method not implemented")
}

func (m mnm) BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error) {
	return "", errors.New("method not implemented")
}

func (m mnm) GetMenu(ctx context.Context, extStoreId string) (models.Menu, error) {
	return models.Menu{}, errors.New("method not implemented")
}

func (m mnm) ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error) {
	return models.MenuUploadTransaction{}, errors.New("method not implemented")
}

func NewManager(ctx context.Context, username, password, baseUrl string) (mnm, error) {

	cli, err := deliverooConf.NewDeliverooClient(&deliverooCli.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize Glovo client ")
		return mnm{}, err
	}
	return mnm{
		cli: *cli,
	}, nil
}
