package stoplist

import (
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/pkg/errors"
)

func NewStopListServiceCron(
	storeService store.Service,
	storeGroupService storegroup.Service,
	menuService *menu.Service,
	aggregatorFactory aggregator.Factory,
	posFactory pos.Factory,
	repo Repository,
	woltCfg config.WoltConfiguration,
	concurrencyLevel int, sqsCli que.SQSInterface) (*ServiceImpl, error) {

	cronStopListS, err := newCronStopList(menuService)
	if err != nil {
		return nil, err
	}

	s, err := newStopListService(
		storeService, storeGroupService, menuService,
		aggregatorFactory, posFactory, repo,
		concurrencyLevel, cronStopListS, woltCfg, sqsCli)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func NewStopListServicePosWebhook(
	storeService store.Service,
	storeGroupService storegroup.Service,
	menuService *menu.Service,
	aggregatorFactory aggregator.Factory,
	posFactory pos.Factory,
	repo Repository,
	woltCfg config.WoltConfiguration,
	concurrencyLevel int, sqsCli que.SQSInterface) (*ServiceImpl, error) {

	webhookStopListS, err := newWebhookStopList()
	if err != nil {
		return nil, err
	}

	s, err := newStopListService(
		storeService, storeGroupService, menuService,
		aggregatorFactory, posFactory, repo,
		concurrencyLevel, webhookStopListS, woltCfg, sqsCli)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func newStopListService(
	storeService store.Service,
	storeGroupService storegroup.Service,
	menuService *menu.Service,
	aggregatorFactory aggregator.Factory,
	posFactory pos.Factory,
	repo Repository,
	concurrencyLevel int,
	stopListIsDisabledBehavior stopListType,
	woltCfg config.WoltConfiguration,
	sqsCli que.SQSInterface,
) (*ServiceImpl, error) {
	if storeService == nil {
		return nil, errors.New("store service is nil")
	}
	if storeGroupService == nil {
		return nil, errors.New("store group service is nil")
	}
	if menuService == nil {
		return nil, errors.New("menu service is nil")
	}
	if aggregatorFactory == nil {
		return nil, errors.New("aggregator factory is nil")
	}
	if posFactory == nil {
		return nil, errors.New("pos factory is nil")
	}
	if concurrencyLevel <= 0 {
		return nil, errors.New("concurrencyLevel is invalid")
	}
	if stopListIsDisabledBehavior == nil {
		return nil, errors.New("stopListType is nil")
	}
	if sqsCli == nil {
		return nil, errors.New("sqsCli is nil")
	}
	return &ServiceImpl{
		storeService:      storeService,
		storeGroupService: storeGroupService,
		menuService:       menuService,
		aggregatorFactory: aggregatorFactory,
		posFactory:        posFactory,
		repo:              repo,
		concurrencyLevel:  concurrencyLevel,
		woltCfg:           woltCfg,
		stopListType:      stopListIsDisabledBehavior,
		notifyCli:         sqsCli,
	}, nil
}

func NewStopListServiceValidate(
	storeService store.Service,
	storeGroupService storegroup.Service,
	menuService *menu.Service,
	aggregatorFactory aggregator.Factory,
	posFactory pos.Factory,
	repo Repository,
	concurrencyLevel int,
	woltCfg config.WoltConfiguration,
	sqsCli que.SQSInterface,
) (*ServiceImpl, error) {

	validateStopListS := newValidateStopList()

	s, err := newStopListService(
		storeService, storeGroupService, menuService,
		aggregatorFactory, posFactory, repo,
		concurrencyLevel, validateStopListS, woltCfg, sqsCli)
	if err != nil {
		return nil, err
	}

	return s, nil
}
