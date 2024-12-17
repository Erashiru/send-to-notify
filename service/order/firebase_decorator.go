package order

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	firebase_client "github.com/kwaaka-team/orders-core/pkg/firebase"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type FirebaseServiceDecorator struct {
	service CreationService

	storeClient        storeClient.Client
	firebaseMsgService *firebase_client.MessageService
}

func NewFirebaseServiceDecorator(
	service CreationService,
	storeClient storeClient.Client,
	firebaseMsgService *firebase_client.MessageService,
) (*FirebaseServiceDecorator, error) {
	if service == nil {
		return nil, errors.New("order service is nil")
	}
	if storeClient == nil {
		return nil, errors.New("storeClient is nil")
	}
	if firebaseMsgService == nil {
		return nil, errors.New("firebaseMsgService is nil")
	}
	return &FirebaseServiceDecorator{
		service:            service,
		storeClient:        storeClient,
		firebaseMsgService: firebaseMsgService,
	}, nil
}

func (s *FirebaseServiceDecorator) CreateOrder(ctx context.Context, externalStoreID, deliveryService string, aggReq interface{}, storeSecret string) (models.Order, error) {
	order, err := s.service.CreateOrder(ctx, externalStoreID, deliveryService, aggReq, storeSecret)
	if err != nil {
		return order, err
	}
	if decoratorErr := s.sendNotification(ctx, order); decoratorErr != nil {
		log.Err(validator.ErrNotifyingInBrowser).Msg(decoratorErr.Error())
	}

	return order, err
}

func (s *FirebaseServiceDecorator) sendNotification(ctx context.Context, order models.Order) error {
	usersToNotify, err := s.storeClient.FindUserStores(ctx, storeModels.UserStore{
		StoreId:          order.RestaurantID,
		SendNotification: true,
	})
	if err != nil {
		return err
	}

	var tokens []string
	for _, user := range usersToNotify {
		tokens = append(tokens, user.FCMTokens...)
	}

	return s.firebaseMsgService.SendOrderNotification(ctx, tokens, order)
}
