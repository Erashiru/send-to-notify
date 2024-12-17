package pos

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"strconv"
)

type Service interface {
	MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error)
	CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration,
		store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
		storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error)
	IsAliveStatus(ctx context.Context, store coreStoreModels.Store) (bool, error) // method for checking POS system is off or on
	GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error)
	GetOrderStatus(ctx context.Context, order models.Order) (string, error)
	GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error)
	AwakeTerminal(ctx context.Context, store coreStoreModels.Store) error
	IsStopListByBalance(ctx context.Context, store coreStoreModels.Store) bool
	GetBalanceLimit(ctx context.Context, store coreStoreModels.Store) int
	CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error
	GetSeqNumber(ctx context.Context) (string, error)
	SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error)
	CloseOrder(ctx context.Context, posOrderId string) error
}

type Factory interface {
	GetPosService(posType models.Pos, store coreStoreModels.Store) (Service, error)
}

type FactoryImpl struct {
	baseService    *BasePosService
	notifyQueue    notifyQueue.SQSInterface
	retryQueueName string

	iikoBaseURL                 string
	iikoTransportToFrontTimeout int
	syrveBaseURL                string
	posterBaseURL               string

	palomaBaseUrl string
	palomaClass   string

	jowiBaseUrl   string
	jowiApiKey    string
	jowiApiSecret string

	burgerKingAddress string
	bkOfferRepository drivers.BKOfferRepository

	rkeeperBaseUrl string
	rkeeperApiKey  string

	rkeeper7XMLLisenceUrl string

	yarosBaseUrl    string
	yarosInfoSystem string

	tillypadBaseUrl string
	ytimesBaseUrl   string
	ytimesToken     string

	posistBaseUrl string
}

func NewFactory(
	anotherBillRepository AnotherBillRepository,
	notifyQueue notifyQueue.SQSInterface,
	retryQueueName string,
	iikoBaseURL, iikoTransportToFrontTimeoutStr, posterBaseURL,
	palomaBaseUrl, palomaClass,
	jowiBaseUrl, jowiApiKey, jowiApiSecret,
	rkeeperBaseUrl, rkeeperApiKey,
	burgerKingAddress string, bkOfferRepository drivers.BKOfferRepository,
	rkeeper7XMLLisenceUrl,
	syrveBaseURL,
	yarosBaseUrl, yarosInfoSystem, tillypadBaseUrl, ytimesBaseUrl, ytimesToken string, posistBaseUrl string,
) (*FactoryImpl, error) {
	var err error
	transportToFrontTimeout := 0
	if iikoTransportToFrontTimeoutStr != "" {
		transportToFrontTimeout, err = strconv.Atoi(iikoTransportToFrontTimeoutStr)
		if err != nil {
			transportToFrontTimeout = 180
		}
	}

	anotherBillParams, err := anotherBillRepository.GetStoreIDs()
	if err != nil {
		return nil, err
	}
	items := make(map[string][]string)
	for _, billParam := range anotherBillParams {
		items[billParam.RestaurantID] = append(items[billParam.RestaurantID], billParam.Deliveries...)
	}

	bps := &BasePosService{
		anotherBillStoreIDs: items,
	}

	return &FactoryImpl{
		baseService:    bps,
		notifyQueue:    notifyQueue,
		retryQueueName: retryQueueName,

		iikoBaseURL:                 iikoBaseURL,
		iikoTransportToFrontTimeout: transportToFrontTimeout,

		posterBaseURL: posterBaseURL,

		palomaBaseUrl: palomaBaseUrl,
		palomaClass:   palomaClass,

		jowiBaseUrl:   jowiBaseUrl,
		jowiApiKey:    jowiApiKey,
		jowiApiSecret: jowiApiSecret,

		burgerKingAddress: burgerKingAddress,
		bkOfferRepository: bkOfferRepository,

		rkeeperBaseUrl: rkeeperBaseUrl,
		rkeeperApiKey:  rkeeperApiKey,

		rkeeper7XMLLisenceUrl: rkeeper7XMLLisenceUrl,

		syrveBaseURL: syrveBaseURL,

		yarosBaseUrl:    yarosBaseUrl,
		yarosInfoSystem: yarosInfoSystem,

		tillypadBaseUrl: tillypadBaseUrl,
		posistBaseUrl:   posistBaseUrl,

		ytimesBaseUrl: ytimesBaseUrl,
	}, nil
}

func (f *FactoryImpl) GetPosService(posType models.Pos, store coreStoreModels.Store) (Service, error) {
	switch posType {
	case models.BurgerKing:
		return f.getBKService()
	case models.FoodBand:
		return f.getFoodbandService(store)
	case models.IIKO:
		return f.getIikoService(store)
	case models.Syrve:
		return f.getSyrveService(store)
	case models.JOWI:
		return f.getJowiService(store)
	case models.Paloma:
		return f.getPalomaService(store)
	case models.Poster:
		return f.getPosterService(store)
	case models.RKeeper:
		return f.getRKeeperService(store)
	case models.RKeeper7XML:
		return f.getRKeeper7XML(store)
	case models.Yaros:
		return f.getYarosService(store)
	case models.CTMax:
		return f.getCTMaxService()
	case models.Kwaaka:
		return f.getKwaakaPosService()
	case models.TillyPad:
		return f.getTillypadPosService(store)
	case models.Ytimes:
		return f.getYtimesPosService(store)
	case models.Posist:
		return f.getPosistPosService(store)
	}

	return nil, errors.New("pos " + posType.String() + " is not found")
}

func (f *FactoryImpl) getPosistPosService(store coreStoreModels.Store) (*posistPosService, error) {
	svc, err := newPosistPosService(f.baseService, f.posistBaseUrl, store.Posist.AuthBasic, store.Posist.CustomerKey, store.Posist.TabId)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (f *FactoryImpl) getIikoService(store coreStoreModels.Store) (*iikoService, error) {
	s, err := newIikoService(f.baseService, f.iikoBaseURL, store.IikoCloud.OrganizationID, store.IikoCloud.TerminalID, store.IikoCloud.Key, f.iikoTransportToFrontTimeout, store.IikoCloud.CustomDomain)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *FactoryImpl) getSyrveService(store coreStoreModels.Store) (*iikoService, error) {
	s, err := newIikoService(f.baseService, f.syrveBaseURL, store.IikoCloud.OrganizationID, store.IikoCloud.TerminalID, store.IikoCloud.Key, f.iikoTransportToFrontTimeout, store.IikoCloud.CustomDomain)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *FactoryImpl) getPosterService(store coreStoreModels.Store) (*PosterService, error) {
	s, err := NewPosterService(f.baseService, f.posterBaseURL, store.Poster.Token, nil, nil, nil, "", "", "", nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *FactoryImpl) getFoodbandService(store coreStoreModels.Store) (*foodbandService, error) {
	posIntegrationCfg, err := getExternalCfgByName(store.ExternalConfig, models.FoodBand.String())
	if err != nil {
		return nil, err
	}

	posIntegrationCli, err := newFoodbandService(f.baseService, posIntegrationCfg.WebhookConfig.OrderCreate, posIntegrationCfg.WebhookConfig.OrderCancel, posIntegrationCfg.ClientSecret, posIntegrationCfg.StoreID[0], posIntegrationCfg.WebhookConfig.RetryMaxCount)
	if err != nil {
		return nil, err
	}

	return posIntegrationCli, nil
}

func (f *FactoryImpl) getYarosService(store coreStoreModels.Store) (*yarosService, error) {
	s, err := newYarosService(f.baseService, f.notifyQueue, f.retryQueueName, store.Yaros.StoreId, f.yarosInfoSystem, f.yarosBaseUrl, store.Yaros.Username, store.Yaros.Password)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *FactoryImpl) getCTMaxService() (*ctMaxService, error) {
	s, err := newCTMaxService(f.baseService)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (f *FactoryImpl) getPalomaService(store coreStoreModels.Store) (*palomaService, error) {
	s, err := newPalomaService(f.baseService, f.palomaBaseUrl, store.Paloma.ApiKey, f.palomaClass, store.Paloma.PointID, store.Paloma.StopListByBalance)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *FactoryImpl) getJowiService(store coreStoreModels.Store) (*jowiService, error) {
	s, err := newJowiService(f.baseService, f.jowiBaseUrl, f.jowiApiKey, f.jowiApiSecret, store.Jowi.RestaurantID)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *FactoryImpl) getRKeeperService(store coreStoreModels.Store) (*rkeeperService, error) {
	s, err := newRkeeperService(f.baseService, store.RKeeper.ObjectId, store.RKeeper.ApiKey, f.rkeeperApiKey, f.rkeeperBaseUrl)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *FactoryImpl) getBKService() (*burgerKingService, error) {
	s, err := newBurgerKingService(f.baseService, f.burgerKingAddress, f.bkOfferRepository)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *FactoryImpl) getTillypadPosService(store coreStoreModels.Store) (*tillypadPosService, error) {
	s, err := newTillypadPosService(f.baseService, f.tillypadBaseUrl, store.TillyPad.PointId, store.TillyPad.ClientId, store.TillyPad.ClientSecret, store.TillyPad.PathPrefix)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (f *FactoryImpl) getYtimesPosService(store coreStoreModels.Store) (*ytimesPosService, error) {
	s, err := newYtimesPosService(f.baseService, f.ytimesBaseUrl, store)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (f *FactoryImpl) getRKeeper7XML(store coreStoreModels.Store) (*rkeeper7XMLService, error) {
	s, err := newRkeeper7XMLService(f.baseService,
		store.RKeeper7XML.Domain,
		store.RKeeper7XML.Username,
		store.RKeeper7XML.Password,
		store.RKeeper7XML.UCSUsername,
		store.RKeeper7XML.UCSPassword,
		store.RKeeper7XML.Token,
		f.rkeeper7XMLLisenceUrl,
		store.RKeeper7XML.Anchor,
		store.RKeeper7XML.ObjectID,
		store.RKeeper7XML.StationID,
		store.RKeeper7XML.StationCode,
		store.RKeeper7XML.LicenseInstanceGUID,
		store.RKeeper7XML.ChildItems,
		store.RKeeper7XML.ClassificatorItemIdent,
		store.RKeeper7XML.ClassificatorPropMask,
		store.RKeeper7XML.MenuItemsPropMask,
		store.RKeeper7XML.PropFilter,
		store.RKeeper7XML.Cashier,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *FactoryImpl) getKwaakaPosService() (*kwaakaPosService, error) {
	s, err := newKwaakaPosService(f.baseService)
	if err != nil {
		return nil, err
	}
	return s, nil
}
