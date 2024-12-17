package talabat

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	pkgUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	talabatClient "github.com/kwaaka-team/orders-core/pkg/talabat"
	talabatConfig "github.com/kwaaka-team/orders-core/pkg/talabat/clients"
	talabatModels "github.com/kwaaka-team/orders-core/pkg/talabat/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type mnm struct {
	cli       talabatConfig.TalabatMenu
	mwCli     talabatConfig.TalabatMW
	storeRepo drivers.StoreRepository
	menuRepo  drivers.MenuRepository
	s3Bucket  menu.S3_BUCKET
}

func NewManager(ctx context.Context, cfg menu.Configuration, store storeModels.Store, storeRepo drivers.StoreRepository, menuRepo drivers.MenuRepository) (*mnm, error) {

	switch store.Talabat.IsNewMenu {
	case true:
		mwCli, err := newMwManager(ctx, cfg.TalabatConfiguration.MiddlewareBaseURL, store.Talabat.Username, store.Talabat.Password)
		if err != nil {
			log.Trace().Err(err).Msg("can't initialize talabat mw client")
			return nil, err
		}
		return &mnm{
			mwCli:     mwCli,
			storeRepo: storeRepo,
			menuRepo:  menuRepo,
			s3Bucket:  cfg.S3_BUCKET,
		}, nil
	default:
		cli, err := newMenuManager(ctx, cfg.TalabatConfiguration.MenuBaseURL, store.Talabat.Username, store.Talabat.Password)
		if err != nil {
			log.Trace().Err(err).Msg("can't initialize talabat menu client")
			return nil, err
		}
		return &mnm{
			cli:       cli,
			storeRepo: storeRepo,
			menuRepo:  menuRepo,
			s3Bucket:  cfg.S3_BUCKET,
		}, nil
	}
}

func newMenuManager(ctx context.Context, baseUrl, username, password string) (talabatConfig.TalabatMenu, error) {
	cli, err := talabatClient.NewMenuClient(&talabatConfig.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		Username: username,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	return cli, nil
}

func newMwManager(ctx context.Context, baseUrl, username, password string) (talabatConfig.TalabatMW, error) {
	cli, err := talabatClient.NewMiddlewareClient(&talabatConfig.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		Username: username,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	return cli, nil
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
	requestID := uuid.New().String()
	scheduledOn, err := time.Now().Add(2 * time.Minute).UTC().MarshalText()
	if err != nil {
		return "", err
	}

	err = m.cli.UpdateItemsAvailability(ctx, talabatModels.UpdateItemsAvailabilityRequest{
		RequestID:    requestID,
		RestaurantID: store.Talabat.RestaurantID,
		ScheduledOn:  string(scheduledOn),
		Availability: toAvailabilities(storeID, products, attributes),
	})

	if err != nil {
		return "", err
	}

	return requestID, nil
}

func (m mnm) VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error) {
	return "", errors.New("not implemented")
}

func (m mnm) UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error) {
	var requestID, objectUrl string
	var err error

	switch store.Talabat.IsNewMenu {
	case true:
		requestID, objectUrl, err = m.uploadNewMenu(ctx, store, sv3, menu)
	default:
		requestID, objectUrl, err = m.uploadMenu(ctx, store, sv3)
	}

	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{err.Error()},
			MenuUrl:    objectUrl,
			ID:         requestID,
		}, errors.Wrapf(err, "publicate talabat menu error")
	}

	return models.ExtTransaction{
		Status:     models.PROCESSING.String(),
		MenuID:     menuId,
		ExtStoreID: extStoreId,
		Details:    []string{},
		MenuUrl:    objectUrl,
		ID:         requestID,
	}, nil
}

func (m mnm) uploadNewMenu(ctx context.Context, store storeModels.Store, sv3 *s3.S3, menu models.Menu) (string, string, error) {
	talabatNewMenu, err := m.constructTalabatNewMenu(ctx, menu)
	if err != nil {
		return "", "", errors.Wrapf(err, "construct talabat new menu error")
	}

	talabatCatalog := talabatModels.SubmitCatalogRequest{
		ChainCode: store.Talabat.ChainID,
		Vendors: []string{
			store.Talabat.VendorID,
		},
		Catalog:     talabatNewMenu,
		CallbackUrl: "",
	}

	objectUrl, err := pkgUtils.UploadMenuToS3(store.ID, m.s3Bucket.KwaakaMenuFilesBucket, models.TALABAT.String(), m.s3Bucket.ShareMenuBaseUrl, talabatCatalog, sv3)
	if err != nil {
		return "", "", errors.Wrapf(err, "talabat new menu upload to s3 error")
	}

	res, err := m.mwCli.SubmitCatalog(ctx, talabatCatalog)

	if err != nil {
		return "", objectUrl, errors.Wrapf(err, "publicate talabat new menu error")
	}

	return res.CatalogImportId, objectUrl, nil
}

func (m mnm) uploadMenu(ctx context.Context, store storeModels.Store, sv3 *s3.S3) (string, string, error) {
	talabatMenu, err := m.constructTalabatMenu(ctx, store)
	if err != nil {
		return "", "", errors.Wrapf(err, "construct talabat old menu error")
	}

	objectUrl, err := pkgUtils.UploadMenuToS3(store.ID, m.s3Bucket.KwaakaMenuFilesBucket, models.TALABAT.String(), m.s3Bucket.ShareMenuBaseUrl, talabatMenu, sv3)
	if err != nil {
		return "", "", errors.Wrapf(err, "talabat menu old upload to s3 error")
	}

	requestID := uuid.New().String()

	err = m.cli.CreateNewMenu(ctx, talabatModels.CreateNewMenuRequest{
		Menu:         talabatMenu,
		RestaurantID: store.Talabat.RestaurantID,
		RequestID:    requestID,
	})

	if err != nil {
		return requestID, objectUrl, errors.Wrapf(err, "publicate talabat old menu error")
	}
	return requestID, objectUrl, nil
}

func (m mnm) ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error) {
	return models.MenuUploadTransaction{}, nil
}
