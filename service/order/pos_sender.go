package order

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	modelsMenu "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	MenuClient "github.com/kwaaka-team/orders-core/pkg/menu"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/stoplist"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type PosSender interface {
	SendPosRequest(ctx context.Context, order models.Order, store storeModels.Store, deliveryService string) (models.Order, error)
}

type PosSenderImpl struct {
	globalConfig    config.Configuration
	menuClient      MenuClient.Client
	storeClient     storeClient.Client
	menuService     *menu.Service
	storeService    store.Service
	repository      Repository
	ds              drivers.DataStore
	posFactory      pos.Factory
	errSolution     error_solutions.Service
	stopListService stoplist.Service
	telegramService TelegramService
	notifyQueue     notifyQueue.SQSInterface
}

func NewPosSender(
	globalConfig config.Configuration,
	menuClient MenuClient.Client,
	storeClient storeClient.Client,
	menuService *menu.Service,
	storeService store.Service,
	repository Repository,
	ds drivers.DataStore,
	posFactory pos.Factory,
	errSolution error_solutions.Service,
	stopListService stoplist.Service,
	telegramService TelegramService,
	notifyQueue notifyQueue.SQSInterface) (*PosSenderImpl, error) {
	posSender := &PosSenderImpl{
		globalConfig:    globalConfig,
		menuClient:      menuClient,
		storeClient:     storeClient,
		menuService:     menuService,
		storeService:    storeService,
		repository:      repository,
		ds:              ds,
		posFactory:      posFactory,
		errSolution:     errSolution,
		stopListService: stopListService,
		telegramService: telegramService,
		notifyQueue:     notifyQueue,
	}
	if err := posSender.validate(); err != nil {
		return nil, err
	}
	return posSender, nil
}

func (ps *PosSenderImpl) validate() error {
	if &ps.globalConfig == nil {
		return errors.Wrap(errConstructor, "global config is nil")
	}
	if ps.menuClient == nil {
		return errors.Wrap(errConstructor, "menu client is nil")
	}
	if ps.storeClient == nil {
		return errors.Wrap(errConstructor, "store client is nil")
	}
	if ps.menuService == nil {
		return errors.Wrap(errConstructor, "menu service is nil")
	}
	if ps.storeService == nil {
		return errors.Wrap(errConstructor, "store factory is nil")
	}
	if ps.repository == nil {
		return errors.Wrap(errConstructor, "repository is nil")
	}
	if ps.ds == nil {
		return errors.Wrap(errConstructor, "dataStore is nil")
	}
	if ps.posFactory == nil {
		return errors.Wrap(errConstructor, "pos factory is nil")
	}
	return nil
}

func (ps *PosSenderImpl) SendPosRequest(ctx context.Context, order models.Order, store storeModels.Store, deliveryService string) (models.Order, error) {

	aggMenu, err := ps.menuService.GetAggregatorMenuIfExists(ctx, store, deliveryService)
	if err != nil {
		log.Err(err).Msgf("Menu with ID=%v not found", store.MenuID)
		return ps.failOrder(ctx, order, err)
	}

	posMenu, err := ps.menuService.FindById(ctx, store.MenuID)
	if err != nil {
		log.Err(err).Msgf("Menu with ID=%v not found", store.MenuID)
		return ps.failOrder(ctx, order, err)
	}

	order, err = ps.createOrderInPos(ctx, order, store, *posMenu, aggMenu)
	if err == nil {
		return order, nil
	} else if errors.Is(err, pos.ErrRetry) {
		return order, err
	}

	order.Status = "ACCEPTED"

	if errors.Is(err, validator.ErrIgnoringPos) {
		order, errSuccess := ps.successOrder(ctx, order)
		if errSuccess != nil {
			return order, errSuccess
		}
		return order, nil
	}

	return ps.failOrder(ctx, order, err)
}

func (ps *PosSenderImpl) createOrderInPos(ctx context.Context, order models.Order,
	store storeModels.Store, posMenu, aggregatorMenu modelsMenu.Menu) (models.Order, error) {

	posService, err := ps.posFactory.GetPosService(models.Pos(order.PosType), store)
	if err != nil {
		return order, err
	}

	order, err = posService.CreateOrder(ctx, order, ps.globalConfig, store, posMenu, ps.menuClient, aggregatorMenu, ps.storeClient, ps.errSolution, ps.notifyQueue)
	if err != nil {
		if err1 := ps.SetErrorSolutionAndAddProductToStopList(ctx, store, err.Error(), order); err1 != nil {
			log.Err(err1).Msgf("SetErrorSolutionAndAddProductToStopList for store: %s", store.ID)
		}
		log.Trace().Err(err).Msg("cant send order to POS")
		return order, err
	}

	if err = ps.validatePosOrderId(order); err != nil {
		log.Trace().Err(err).Msg("")
		return order, err
	}

	return order, nil
}

func (ps *PosSenderImpl) SetErrorSolutionAndAddProductToStopList(ctx context.Context, st storeModels.Store, errorMessage string, order models.Order) error {

	errorSolutions, err := ps.errSolution.GetAllErrorSolutions(ctx)
	if err != nil {
		return err
	}

	errorSolutionByCode, addToStopListStatus, err := ps.errSolution.GetErrorSolutionByCode(ctx, st, pos.MatchingCodes(errorMessage, errorSolutions))
	if err != nil {
		return err
	}
	var product modelsMenu.Product
	if addToStopListStatus {
		productID := pos.GetProductIDFromRegexp(errorMessage, errorSolutionByCode)
		productErrorCodes := map[string]bool{
			"21": true, "4": true, "1": true,
		}
		attributeErrorCodes := map[string]bool{
			"25": true, "5": true, "7": true, "21": true, "27": true, "28": true,
		}
		if len(productID) > 0 {
			var err error
			switch {
			case productErrorCodes[errorSolutionByCode.Code]:
				err = ps.stopListService.UpdateStopListByPosProductID(ctx, false, st.ID, productID)
				if err == nil {
					log.Info().Msgf("successfully put product with id: %s to stop with error solution code: %s for store_id : %s", productID, errorSolutionByCode.Code, st.ID)
				}
				for _, orderProduct := range order.Products {
					if orderProduct.ID == productID {
						product.ExtID = orderProduct.ID
						product.Name = append(product.Name, modelsMenu.LanguageDescription{Value: orderProduct.Name})
					}
				}

			case attributeErrorCodes[errorSolutionByCode.Code]:
				err = ps.stopListService.UpdateStopListByAttributeID(ctx, false, st.ID, productID)
				if err == nil {
					log.Info().Msgf("successfully put attribute with id: %s to stop with error solution code: %s for store_id : %s", productID, errorSolutionByCode.Code, st.ID)
				}
				for _, orderProduct := range order.Products {
					if len(orderProduct.Attributes) > 0 {
						for _, orderAttribute := range orderProduct.Attributes {
							if orderAttribute.ID == productID {
								product.ExtID = orderAttribute.ID
								product.Name = append(product.Name, modelsMenu.LanguageDescription{Value: orderAttribute.Name})
							}
						}
					}
				}
			default:
				return fmt.Errorf("unsupported error code to update stoplist bu pos product/attribute id : %s", errorSolutionByCode.Code)
			}
			if err != nil {
				return err
			}
			log.Info().Msgf("send stoplist status true, for product/attribute id: %s, store_id:%s, error solution code:%s", productID, st.ID, errorSolutionByCode.Code)

			if errorSolutionByCode.SendToTelegram {
				if err := ps.telegramService.SendMessageToQueue(telegram.PutProductToStopListWithErrSolution, models.Order{}, st, order.FailReason.BusinessName, "", "", product); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (ps *PosSenderImpl) validatePosOrderId(order models.Order) error {
	if order.PosOrderID == "" {
		return errors.New("pos order id is empty")
	}

	return nil
}

func (ps *PosSenderImpl) successOrder(ctx context.Context, req models.Order) (models.Order, error) {

	log.Info().Msgf("Order created successfully, status: %v", req.Status)

	req.StatusesHistory = append(req.StatusesHistory, models.OrderStatusUpdate{
		Name: req.Status,
		Time: models.TimeNow().Time,
	})

	// ignoring a possible error during the order update
	if err := ps.repository.UpdateOrder(ctx, req); err != nil {
		return req, err
	}

	return req, nil
}

func (ps *PosSenderImpl) failOrder(ctx context.Context, req models.Order, err error) (models.Order, error) {

	log.Trace().Err(err).Msgf("fail  order: %s, err: %s, validator: %v", req.OrderID, err.Error(), validator.ErrFailed)

	req.Status = string(models.STATUS_FAILED)

	req.StatusesHistory = append(req.StatusesHistory, models.OrderStatusUpdate{
		Name: req.Status,
		Time: models.TimeNow().Time,
	})

	errSolutions, err2 := ps.errSolution.GetAllErrorSolutions(ctx)
	if err2 != nil {
		log.Err(err).Msgf("PosSenderImpl error: GetAllErrorSolutions")
		return models.Order{}, err
	}

	if err != nil {
		failReason, _, failReasonErr := ps.errSolution.SetFailReason(ctx, storeModels.Store{}, err.Error(), pos.MatchingCodes(err.Error(), errSolutions), "")
		if failReasonErr != nil {
			return models.Order{}, failReasonErr
		}
		req.FailReason = failReason
	}

	if updateErr := ps.repository.UpdateOrder(ctx, req); updateErr != nil {
		return req, errors.Wrap(errWithNotification, validator.ErrFailed.Error())
	}

	log.Trace().Err(err).Msg("Error while saving order")

	return req, errors.Wrap(errWithNotification, err.Error())
}
