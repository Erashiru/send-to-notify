package aggregator

import (
	"context"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/kwaaka-team/orders-core/core/config"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models3 "github.com/kwaaka-team/orders-core/core/wolt/models"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/pkg/errors"
)

var constructorError = errors.New("base pos service is nil")
var deliveryServiceNotFoundError = errors.New("delivery service not found")

type Aggregator interface {
	MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string
	UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error
	GetSystemCreateOrderRequestByAggregatorRequest(req interface{}, store storeModels.Store) (models.Order, error)
	UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error)
	UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error)
	UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error)
	GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error)
	GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error)
	OpenStore(ctx context.Context, aggregatorStoreId string) error
	IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error)
	SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error)
	GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error)
	GetAggregatorOrder(ctx context.Context, orderID string) (models3.Order, error)
	SendOrderErrorNotification(ctx context.Context, req interface{}) error
	SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error
}

type Factory interface {
	GetAggregator(aggName string, store storeModels.Store) (Aggregator, error)
}

type FactoryImpl struct {
	glovo       *glovoService
	yandex      *yandexService
	qrMenu      *qrMenuService
	kwaakaAdmin *kwaakaAdminService

	woltBaseUrl              string
	talabatMiddlewareBaseUrl string
	talabatMenuBaseUrl       string
	express24BaseUrl         string
	starterAppBaseUrl        string
	menuService              *menu.Service
}

func NewFactory(
	woltBaseUrl string,
	glovoBaseURL, glovoToken string,
	talabatMiddlewareBaseUrl, talabatMenuBaseUrl string,
	express24BaseUrl, starterAppBaseUrl string, menuService *menu.Service, cognito *cognitoidentityprovider.CognitoIdentityProvider, cfg config.Configuration,
) (*FactoryImpl, error) {

	glovo, err := newGlovoService(glovoBaseURL, glovoToken)
	if err != nil {
		return nil, err
	}
	yandex, err := newYandexService(cfg.YandexConfiguration)
	if err != nil {
		return nil, err
	}

	qrMenu, err := newQrMenuService()
	if err != nil {
		return nil, err
	}

	kwaakaAdminService, err := newKwaakaAdminService(cognito)
	if err != nil {
		return nil, err
	}

	return &FactoryImpl{
		woltBaseUrl:              woltBaseUrl,
		glovo:                    glovo,
		yandex:                   yandex,
		express24BaseUrl:         express24BaseUrl,
		talabatMiddlewareBaseUrl: talabatMiddlewareBaseUrl,
		talabatMenuBaseUrl:       talabatMenuBaseUrl,
		qrMenu:                   qrMenu,
		kwaakaAdmin:              kwaakaAdminService,
		menuService:              menuService,
		starterAppBaseUrl:        starterAppBaseUrl,
	}, nil
}

func (f *FactoryImpl) GetAggregator(aggName string, store storeModels.Store) (Aggregator, error) {
	switch aggName {
	case models.GLOVO.String():
		return f.glovo, nil
	case models.WOLT.String():
		return f.getAggregatorServiceForWolt(store, f.menuService)
	case models.YANDEX.String():
		return f.yandex, nil
	case models.EMENU.String():
		return f.getAggregatorServiceForEmenu(store)
	case models.EXPRESS24.String():
		return f.getAggregatorForExpress24V2(store)
	case models.TALABAT.String():
		return f.getAggregatorServiceForTalabat(store)
	case models.QRMENU.String():
		return f.qrMenu, nil
	case models.KWAAKA_ADMIN.String():
		return f.kwaakaAdmin, nil
	case models.STARTERAPP.String():
		return f.getAggregatorForStarterApp(store, f.menuService)
	}

	return nil, errors.Wrapf(deliveryServiceNotFoundError, "delivery service %s not found", aggName)
}

func (f *FactoryImpl) getAggregatorForExpress24(store storeModels.Store) (Aggregator, error) {
	return newExpress24Service(f.express24BaseUrl, store)
}

func (f *FactoryImpl) getAggregatorForExpress24V2(store storeModels.Store) (Aggregator, error) {
	return newExpress24v2Service(f.express24BaseUrl, store)
}

func (f *FactoryImpl) getAggregatorForStarterApp(store storeModels.Store, menuService *menu.Service) (Aggregator, error) {
	return newStarterAppService(f.starterAppBaseUrl, store, menuService)
}

func (f *FactoryImpl) getAggregatorServiceForEmenu(store storeModels.Store) (Aggregator, error) {
	return newEmenuService(store)
}

func (f *FactoryImpl) getAggregatorServiceForWolt(store storeModels.Store, menuService *menu.Service) (Aggregator, error) {
	return newWoltService(f.woltBaseUrl, store, menuService)
}

func (f *FactoryImpl) getAggregatorServiceForTalabat(store storeModels.Store) (Aggregator, error) {
	return newTalabatService(f.talabatMiddlewareBaseUrl, f.talabatMenuBaseUrl, store)
}
