package managers

import (
	"context"
	"github.com/google/uuid"
	generalConfig "github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/core/externalapi/database/drivers"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/core/externalapi/utils"
	storeCoreModel "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/store"
	storeCoreDto "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/rs/zerolog/log"
)

type AuthClient interface {
	FindByIDAndSecret(ctx context.Context, clientID string, clientSecret string) (models.AuthClient, error)
	FindByID(ctx context.Context, clientID string) (models.AuthClient, error)
	GenerateToken(ctx context.Context, req models.Credentials) (string, error)
	GetListID(ctx context.Context) ([]models.AuthClient, error)
	SetCredential(ctx context.Context, req models.SetCredsToStore) error
}

type AuthClientManager struct {
	ds             drivers.DataStore
	authClientRepo drivers.AuthClientRepository
	appSecret      string
	storeCli       store.Client
	emenuCfg       generalConfig.EmenuConfiguration
}

func NewAuthClientManager(ds drivers.DataStore, storeCli store.Client, appSecret string, emenuCfg generalConfig.EmenuConfiguration) AuthClient {
	return &AuthClientManager{
		ds:             ds,
		authClientRepo: ds.AuthClientRepository(),
		storeCli:       storeCli,
		appSecret:      appSecret,
		emenuCfg:       emenuCfg,
	}
}

func (manager AuthClientManager) FindByID(ctx context.Context, clientID string) (models.AuthClient, error) {
	authClient, err := manager.authClientRepo.FindByID(ctx, clientID)

	if err != nil {
		log.Trace().Err(err).Msgf("Auth client not found with given credentials")
		return models.AuthClient{}, err
	}

	return authClient, nil
}

func (manager AuthClientManager) FindByIDAndSecret(ctx context.Context, clientID string, clientSecret string) (models.AuthClient, error) {
	encodedData, err := utils.EncodeSecret(clientSecret, manager.appSecret)

	if err != nil {
		log.Trace().Err(err).Msgf("Cant encode secret")
		return models.AuthClient{}, err
	}

	log.Info().Msgf("Token: %v", encodedData)

	authClient, err := manager.authClientRepo.FindByIDAndSecret(ctx, clientID, encodedData)

	if err != nil {
		log.Trace().Err(err).Msgf("Auth client not found with given credentials.")
		return models.AuthClient{}, err
	}

	return authClient, nil
}

func (manager AuthClientManager) GenerateToken(ctx context.Context, req models.Credentials) (string, error) {

	encodedData, err := utils.EncodeSecret(req.AuthenticateData.ClientSecret, manager.appSecret)
	if err != nil {
		log.Trace().Err(err).Msgf("Cant encode secret")
		return "", err
	}

	if err := manager.authClientRepo.AuthClientExist(ctx, req.ClientID, encodedData); err != nil {
		return "", err
	}

	req.ClientSecret = encodedData

	res, err := manager.authClientRepo.CreateAuthClient(ctx, models.ToModelAuthClient(req))
	if err != nil {
		return "", err
	}

	return res, nil
}

func (manager AuthClientManager) GetListID(ctx context.Context) ([]models.AuthClient, error) {
	res, err := manager.authClientRepo.GetListID(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (manager AuthClientManager) SetCredential(ctx context.Context, req models.SetCredsToStore) error {

	res, err := manager.authClientRepo.GetAuthClientByID(ctx, req.AuthId)
	if err != nil {
		return err
	}

	store, err := manager.storeCli.FindStore(ctx, storeCoreDto.StoreSelector{ID: req.RestID})
	if err != nil {
		return err
	}

	externalUpdate := manager.setExternalConfigsFromStore(store.ExternalConfig)

	update := storeCoreDto.UpdateStoreExternalConfig{
		StoreID:       []string{uuid.New().String()},
		Type:          &res.Service,
		ClientSecret:  &res.ClientSecret,
		SendToPos:     &req.SendToPos,
		IsMarketplace: &req.IsMarketplace,
		MenuUrl:       &req.MenuUrl,
		PaymentTypes: &storeCoreDto.UpdateDeliveryServicePaymentType{
			CASH: &storeCoreDto.UpdateIIKOPaymentType{
				IikoPaymentTypeID:        &req.PaymentTypes.CASH.PaymentTypeID,
				IikoPaymentTypeKind:      &req.PaymentTypes.CASH.PaymentTypeKind,
				OrderType:                &req.PaymentTypes.CASH.OrderType,
				PromotionPaymentTypeID:   &req.PaymentTypes.CASH.PromotionPaymentTypeID,
				OrderTypeService:         &req.PaymentTypes.CASH.OrderTypeService,
				OrderTypeForVirtualStore: &req.PaymentTypes.CASH.OrderTypeForVirtualStore,
				IsProcessedExternally:    req.PaymentTypes.CASH.IsProcessedExternally,
			},
			DELAYED: &storeCoreDto.UpdateIIKOPaymentType{
				IikoPaymentTypeID:        &req.PaymentTypes.DELAYED.PaymentTypeID,
				IikoPaymentTypeKind:      &req.PaymentTypes.DELAYED.PaymentTypeKind,
				OrderType:                &req.PaymentTypes.DELAYED.OrderType,
				PromotionPaymentTypeID:   &req.PaymentTypes.DELAYED.PromotionPaymentTypeID,
				OrderTypeService:         &req.PaymentTypes.DELAYED.OrderTypeService,
				OrderTypeForVirtualStore: &req.PaymentTypes.DELAYED.OrderTypeForVirtualStore,
				IsProcessedExternally:    req.PaymentTypes.DELAYED.IsProcessedExternally,
			},
		},
	}

	if res.Service == models.EMENU {
		update.AuthToken = &manager.emenuCfg.EmenuAuthToken
		update.WebhookURL = &manager.emenuCfg.EmenuWebhookURL
		update.WebhookProductStoplist = &manager.emenuCfg.EmenuWebhookProductStoplist
		update.WebhookAttributeStoplist = &manager.emenuCfg.EmenuWebhookAttributeStoplist
	}

	externalUpdate = append(externalUpdate, update)

	if err := manager.storeCli.Update(ctx, storeCoreDto.UpdateStore{
		ID:       &req.RestID,
		External: externalUpdate,
	}); err != nil {
		return err
	}

	return nil
}

func (manager AuthClientManager) setExternalConfigsFromStore(storeExtConf []storeCoreModel.StoreExternalConfig) []storeCoreDto.UpdateStoreExternalConfig {
	switch storeExtConf {
	case nil:
		return []storeCoreDto.UpdateStoreExternalConfig{}
	default:
		var externalUpdate []storeCoreDto.UpdateStoreExternalConfig

		for i := 0; i < len(storeExtConf); i++ {
			update := storeCoreDto.UpdateStoreExternalConfig{
				StoreID:                  storeExtConf[i].StoreID,
				Type:                     &storeExtConf[i].Type,
				ClientSecret:             &storeExtConf[i].ClientSecret,
				SendToPos:                &storeExtConf[i].SendToPos,
				IsMarketplace:            &storeExtConf[i].IsMarketplace,
				MenuUrl:                  &storeExtConf[i].MenuUrl,
				AuthToken:                &storeExtConf[i].AuthToken,
				WebhookURL:               &storeExtConf[i].WebhookURL,
				WebhookProductStoplist:   &storeExtConf[i].WebhookProductStoplist,
				WebhookAttributeStoplist: &storeExtConf[i].WebhookAttributeStoplist,
				PaymentTypes: &storeCoreDto.UpdateDeliveryServicePaymentType{
					CASH: &storeCoreDto.UpdateIIKOPaymentType{
						IikoPaymentTypeID:   &storeExtConf[i].PaymentTypes.CASH.PaymentTypeID,
						IikoPaymentTypeKind: &storeExtConf[i].PaymentTypes.CASH.PaymentTypeKind,
					},
					DELAYED: &storeCoreDto.UpdateIIKOPaymentType{
						IikoPaymentTypeID:   &storeExtConf[i].PaymentTypes.DELAYED.PaymentTypeID,
						IikoPaymentTypeKind: &storeExtConf[i].PaymentTypes.DELAYED.PaymentTypeKind,
					},
				},
			}
			externalUpdate = append(externalUpdate, update)
		}
		return externalUpdate
	}
}
