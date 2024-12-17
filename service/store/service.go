package store

import (
	"context"
	"errors"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	kwaakaAdminModels "github.com/kwaaka-team/orders-core/core/kwaaka_admin/models"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	selector2 "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"strconv"
	"strings"
)

type Service interface {
	IsSendToPos(store models.Store, deliveryService string) (bool, error)
	IsAutoAccept(store models.Store, deliveryService string) (bool, error)
	IsPostAutoAccept(store models.Store, deliveryService string) (bool, error)
	GetByID(ctx context.Context, storeID string) (models.Store, error)
	GetStoresByToken(ctx context.Context, token string) ([]models.Store, error)
	FindAllStores(ctx context.Context) ([]models.Store, error)
	GetStoresByStoreGroupID(ctx context.Context, storeGroupID string) ([]models.Store, error)
	GetByExternalIdAndDeliveryService(ctx context.Context, externalStoreID string, deliveryService string) (models.Store, error)
	IsMarketplace(store models.Store, deliveryService string) (bool, error)
	GetPaymentTypes(store models.Store, deliveryService string, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error)
	IsSecretValid(store models.Store, deliveryService, secret string) (bool, error)
	GetStoreExternalIds(store models.Store, deliveryService string) ([]string, error)
	FindStoresByDeliveryService(ctx context.Context, deliveryService string) ([]models.Store, error)
	UpdateStoreSchedule(ctx context.Context, req models.UpdateStoreSchedule) error
	FindStoresByPosType(ctx context.Context, posType string) ([]models.Store, error)
	GetOrderCodePrefix(ctx context.Context, store models.Store, deliveryService string) (string, error)
	IgnoreStatusUpdate(store models.Store, deliveryService string) (bool, error)
	GetByYarosRestaurantID(ctx context.Context, restaurantID string) (models.Store, error)
	GetRestaurantsByGroupId(ctx context.Context, pagination selector2.Pagination, restId string) ([]models.Store, error)
	AddMenuObjectToMenus(ctx context.Context, storeId, deliveryService, menuId string) error
	GetStoresInRestGroupByName(ctx context.Context, restaurantGroupId, name string, legalEntities []string) ([]models.Store, error)
	UpdateMenuId(ctx context.Context, storeId, menuId string) error
	CreatePolygon(ctx context.Context, request kwaakaAdminModels.PolygonRequest) error
	UpdatePolygon(ctx context.Context, request kwaakaAdminModels.PolygonRequest) error
	GetPolygonByRestaurantID(ctx context.Context, restaurantID string) (kwaakaAdminModels.GetPolygonResponse, error)
	GetTwoGisReviewLink(ctx context.Context, restaurantID string) (string, error)
	CreateStorePhoneEmail(ctx context.Context, restaurantID string, request kwaakaAdminModels.StorePhoneEmail) error
	UpdateKwaakaAdminBusyMode(ctx context.Context, req []dto.BusyModeRequest) error
	GetStoresByIIKOOrganizationId(ctx context.Context, organizationId string) ([]models.Store, error)
	GetStoresByWppPhoneNum(ctx context.Context, phoneNum string) ([]models.Store, error)
	SetActualRkeeper7xmlSeqNumber(ctx context.Context, storeID, seqNumber string) error
	UpdateStoreByFields(ctx context.Context, req storeModels.UpdateStore) error
	GetStoresBySelectorFilter(ctx context.Context, query selector2.Store) ([]models.Store, error)
}

type storeServiceByDeliveryService interface {
	IsSendToPos(store models.Store) (bool, error)
	IsAutoAccept(store models.Store) (bool, error)
	IsPostAutoAccept(store models.Store) (bool, error)
	IsMarketplace(store models.Store) (bool, error)
	GetPaymentTypes(store models.Store, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error)
	GetStoreExternalIds(store models.Store) ([]string, error)
	GetOrderCodePrefix(ctx context.Context, store models.Store) (string, error)
	IgnoreStatusUpdate(store models.Store) bool
	IsSecretValid(store models.Store, secret string) (bool, error)
	GetStoreSchedulePrefix() string
}

type ServiceImpl struct {
	storeRepository Repository
	glovo           *glovoStoreService
	wolt            *woltStoreService
	express24       *express24Service
	deliveroo       *deliverooService
	external        map[string]*externalStoreService
	talabat         *talabatStoreService
	qrMenu          *qrMenuStoreService
	kwaakaAdmin     *kwaakaAdminStoreService
	starterApp      *starterAppService
}

func NewService(storeRepository Repository) (*ServiceImpl, error) {
	glovo, err := newGlovoStoreService()
	if err != nil {
		return nil, err
	}

	wolt, err := newWoltStoreService()
	if err != nil {
		return nil, err
	}

	express24, err := newExpress24Service()
	if err != nil {
		return nil, err
	}

	deliveroo, err := newDeliverooService()
	if err != nil {
		return nil, err
	}

	talabat, err := newTalabatStoreService()
	if err != nil {
		return nil, err
	}

	qrMenu, err := newQrMenuStoreService()
	if err != nil {
		return nil, err
	}

	kwaakaAdmin, err := newKwaakaAdminStoreService()
	if err != nil {
		return nil, err
	}

	starterApp, err := newStarterAppService()
	if err != nil {
		return nil, err
	}

	return &ServiceImpl{
		storeRepository: storeRepository,
		glovo:           glovo,
		wolt:            wolt,
		express24:       express24,
		deliveroo:       deliveroo,
		external:        make(map[string]*externalStoreService),
		talabat:         talabat,
		qrMenu:          qrMenu,
		kwaakaAdmin:     kwaakaAdmin,
		starterApp:      starterApp,
	}, nil
}

func (s *ServiceImpl) GetStoresByIIKOOrganizationId(ctx context.Context, organizationId string) ([]models.Store, error) {
	return s.storeRepository.GetStoresByIIKOOrganizationId(ctx, organizationId)
}

func (s *ServiceImpl) GetStoresByWppPhoneNum(ctx context.Context, phoneNum string) ([]models.Store, error) {
	return s.storeRepository.GetStoresByWppPhoneNumber(ctx, phoneNum)
}

func (s *ServiceImpl) AddMenuObjectToMenus(ctx context.Context, storeId, deliveryService, menuId string) error {
	store, err := s.storeRepository.GetById(ctx, storeId)
	if err != nil {
		return err
	}

	store.Menus = append(store.Menus, models.StoreDSMenu{
		ID:        menuId,
		Name:      "generated " + deliveryService + " menu",
		IsActive:  false,
		IsDeleted: false,
		IsSync:    true,
		Delivery:  deliveryService,
	})

	if err = s.storeRepository.UpdateMenus(ctx, store.ID, store.Menus); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) getByDeliveryService(deliveryService string) (storeServiceByDeliveryService, error) {
	if deliveryService == "" {
		return nil, errors.New("delivery service undefined")
	}
	switch deliveryService {
	case "glovo":
		return s.glovo, nil
	case "wolt":
		return s.wolt, nil
	case "express24":
		return s.express24, nil
	case "deliveroo":
		return s.deliveroo, nil
	case "talabat":
		return s.talabat, nil
	case "qr_menu":
		return s.qrMenu, nil
	case "kwaaka_admin":
		return s.kwaakaAdmin, nil
	case "starter_app":
		return s.starterApp, nil

	default:
		if service, ok := s.external[deliveryService]; ok {
			return service, nil
		} else {
			stService, err := newExternalStoreService(deliveryService)
			if err != nil {
				return nil, err
			}
			s.external[deliveryService] = stService
			return stService, nil
		}
	}
}

func (s *ServiceImpl) UpdateMenuId(ctx context.Context, storeId, menuId string) error {
	return s.storeRepository.UpdateMenuId(ctx, storeId, menuId)
}

func (s *ServiceImpl) GetByID(ctx context.Context, storeID string) (models.Store, error) {
	store, err := s.storeRepository.GetById(ctx, storeID)
	if err != nil {
		return models.Store{}, err
	}
	return *store, nil
}

func (s *ServiceImpl) GetByExternalIdAndDeliveryService(ctx context.Context, externalStoreID string, deliveryService string) (models.Store, error) {
	if store, err := s.storeRepository.GetByExternalIdAndAggregator(ctx, externalStoreID, deliveryService); err != nil {
		return models.Store{}, err
	} else {
		return *store, nil
	}
}

func (s *ServiceImpl) IsSecretValid(store models.Store, deliveryService, secret string) (bool, error) {
	byDeliveryService, err := s.getByDeliveryService(deliveryService)
	if err != nil {
		return false, err
	}
	return byDeliveryService.IsSecretValid(store, secret)
}

func (s *ServiceImpl) IsSendToPos(store models.Store, deliveryService string) (bool, error) {
	byDeliveryService, err := s.getByDeliveryService(deliveryService)
	if err != nil {
		return false, err
	}
	return byDeliveryService.IsSendToPos(store)
}

func (s *ServiceImpl) IsAutoAccept(store models.Store, deliveryService string) (bool, error) {
	byDeliveryService, err := s.getByDeliveryService(deliveryService)
	if err != nil {
		return false, err
	}
	return byDeliveryService.IsAutoAccept(store)
}

func (s *ServiceImpl) IsPostAutoAccept(store models.Store, deliveryService string) (bool, error) {
	byDeliveryService, err := s.getByDeliveryService(deliveryService)
	if err != nil {
		return false, err
	}
	return byDeliveryService.IsPostAutoAccept(store)
}

func (s *ServiceImpl) IgnoreStatusUpdate(store models.Store, deliveryService string) (bool, error) {
	byDeliveryService, err := s.getByDeliveryService(deliveryService)
	if err != nil {
		return false, err
	}

	return byDeliveryService.IgnoreStatusUpdate(store), nil
}

func (s *ServiceImpl) IsMarketplace(store models.Store, deliveryService string) (bool, error) {
	byDeliveryService, err := s.getByDeliveryService(deliveryService)
	if err != nil {
		return false, err
	}
	return byDeliveryService.IsMarketplace(store)
}

func (s *ServiceImpl) GetPaymentTypes(store models.Store, deliveryService string, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error) {
	byDeliveryService, err := s.getByDeliveryService(deliveryService)
	if err != nil {
		return models.DeliveryServicePaymentType{}, err
	}
	return byDeliveryService.GetPaymentTypes(store, paymentInfo)
}

func (s *ServiceImpl) GetStoreExternalIds(store models.Store, deliveryService string) ([]string, error) {
	byDeliveryService, err := s.getByDeliveryService(deliveryService)
	if err != nil {
		return nil, err
	}
	return byDeliveryService.GetStoreExternalIds(store)
}

func (s *ServiceImpl) FindStoresByDeliveryService(ctx context.Context, deliveryService string) ([]models.Store, error) {
	if deliveryService == "" {
		return nil, errors.New("delivery service is empty")
	}
	return s.storeRepository.GetStoresByDeliveryService(ctx, deliveryService)
}

func (s *ServiceImpl) UpdateStoreSchedule(ctx context.Context, req models.UpdateStoreSchedule) error {
	byDeliveryService, err := s.getByDeliveryService(req.DeliveryService)
	if err != nil {
		return err
	}

	prefix := byDeliveryService.GetStoreSchedulePrefix()

	return s.storeRepository.UpdateStoreSchedule(ctx, req.RestaurantID, req.StoreSchedule, prefix)
}

func (s *ServiceImpl) GetOrderCodePrefix(ctx context.Context, store models.Store, deliveryService string) (string, error) {
	byDeliveryService, err := s.getByDeliveryService(deliveryService)
	if err != nil {
		return "", err
	}

	return byDeliveryService.GetOrderCodePrefix(ctx, store)
}

func (s *ServiceImpl) GetStoresByToken(ctx context.Context, token string) ([]models.Store, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	return s.storeRepository.GetStoresByToken(ctx, token)
}

func (s *ServiceImpl) GetStoresByStoreGroupID(ctx context.Context, storeGroupID string) ([]models.Store, error) {
	if storeGroupID == "" {
		return nil, errors.New("token is empty")
	}

	return s.storeRepository.FindStoresByStoreGroupID(ctx, storeGroupID)
}

func (s *ServiceImpl) FindAllStores(ctx context.Context) ([]models.Store, error) {
	return s.storeRepository.FindAllStores(ctx)
}

func (s *ServiceImpl) FindStoresByPosType(ctx context.Context, posType string) ([]models.Store, error) {
	if posType == "" {
		return nil, errors.New("pos type is empty")
	}
	return s.storeRepository.FindStoresByPosType(ctx, posType)
}

func (s *ServiceImpl) FindStoresByTimeZone(ctx context.Context, timeZone string) ([]models.Store, error) {
	if timeZone == "" {
		return nil, errors.New("pos type is empty")
	}
	return s.storeRepository.FindStoresByTimeZone(ctx, timeZone)
}

func (s *ServiceImpl) GetByYarosRestaurantID(ctx context.Context, restaurantID string) (models.Store, error) {
	return s.storeRepository.GetByYarosRestaurantID(ctx, restaurantID)
}

func (s *ServiceImpl) GetRestaurantsByGroupId(ctx context.Context, pagination selector2.Pagination, restId string) ([]models.Store, error) {
	return s.storeRepository.GetRestaurantsByGroupId(ctx, pagination, restId)
}

func (s *ServiceImpl) GetStoresInRestGroupByName(ctx context.Context, restaurantGroupId, name string, legalEntities []string) ([]models.Store, error) {
	var stores []models.Store
	var err error
	if name == "" {
		stores, err = s.storeRepository.GetRestaurantsByGroupId(ctx, selector2.Pagination{}, restaurantGroupId)
		if err != nil {
			return nil, err
		}
	} else {
		stores, err = s.storeRepository.FindStoreInRestGroupByName(ctx, name, restaurantGroupId)
		if err != nil {
			return nil, err
		}
	}

	res := []models.Store{}
	if len(legalEntities) == 0 {
		res = stores
		return res, nil
	}
	for _, store := range stores {
		if s.FindStoresByLegalEntityId(store, legalEntities) {
			res = append(res, store)
		}
	}

	return res, nil
}

func (s *ServiceImpl) FindStoresByLegalEntityId(store models.Store, legalEntities []string) bool {
	for _, entity := range legalEntities {
		if entity == store.LegalEntityId {
			return true
		}
	}

	return false
}

func (s *ServiceImpl) CreatePolygon(ctx context.Context, request kwaakaAdminModels.PolygonRequest) error {
	polygon, err := s.storeRepository.GetPolygonByRestaurantID(ctx, request.RestaurantID)
	if err != nil {
		return err
	}
	if polygon.Coordinates != nil {
		return errors.New("polygon already exists")
	}

	return s.storeRepository.CreatePolygon(ctx, request.RestaurantID, models.Geometry{
		PercentageModifier: request.Percentage,
		Coordinates:        request.Coordinates,
		Type:               "polygon",
	})
}

func (s *ServiceImpl) UpdatePolygon(ctx context.Context, request kwaakaAdminModels.PolygonRequest) error {
	return s.storeRepository.UpdatePolygon(ctx, request.RestaurantID, models.Geometry{
		PercentageModifier: request.Percentage,
		Coordinates:        request.Coordinates,
	})
}

func (s *ServiceImpl) GetPolygonByRestaurantID(ctx context.Context, restaurantID string) (kwaakaAdminModels.GetPolygonResponse, error) {

	polygon, err := s.storeRepository.GetPolygonByRestaurantID(ctx, restaurantID)
	if err != nil {
		return kwaakaAdminModels.GetPolygonResponse{}, err
	}

	result := kwaakaAdminModels.GetPolygonResponse{
		Percentage:  polygon.PercentageModifier,
		Coordinates: polygon.Coordinates,
	}

	return result, nil
}

func (s *ServiceImpl) GetTwoGisReviewLink(ctx context.Context, restaurantID string) (string, error) {
	store, err := s.storeRepository.GetById(ctx, restaurantID)
	if err != nil {
		return "", err
	}

	for _, url := range store.ExternalLinks {
		if url.Name == "2gis" {
			twoGisLink := url.Url
			linkSubstr := "?m="
			reviewLinkSubstr := "/tab/reviews/addreview"

			i := strings.Index(twoGisLink, linkSubstr)
			if i == -1 {
				return twoGisLink + reviewLinkSubstr, nil
			} else {
				return twoGisLink[:i] + reviewLinkSubstr + twoGisLink[i:], nil
			}
		}
	}

	return "", errors.New("2gis external link name not found")
}

func (s *ServiceImpl) CreateStorePhoneEmail(ctx context.Context, restaurantID string, request kwaakaAdminModels.StorePhoneEmail) error {
	return s.storeRepository.CreateStorePhoneEmail(ctx, restaurantID, request)
}

func (s *ServiceImpl) UpdateKwaakaAdminBusyMode(ctx context.Context, req []dto.BusyModeRequest) error {
	for _, rq := range req {
		store, err := s.GetByID(ctx, rq.RestaurantID)
		if err != nil {
			return err
		}
		if store.QRMenu.SendToPos && rq.BusyMode {
			for storeID := range store.QRMenu.StoreID {
				err = s.storeRepository.UpdateKwaakaAdminBusyMode(ctx, store.QRMenu.StoreID[storeID], rq.BusyMode, rq.BusyModeMinute)
				if err != nil {
					return errors.New("error: update direct busy mode status in direct store")
				}
			}
		}
	}
	return nil
}

func (s *ServiceImpl) SetActualRkeeper7xmlSeqNumber(ctx context.Context, storeID, seqNumber string) error {
	seqNum, err := strconv.Atoi(seqNumber)
	if err != nil {
		return err
	}

	seqNum = seqNum + 1

	if err := s.UpdateStoreByFields(ctx, storeModels.UpdateStore{
		ID: &storeID,
		RKeeper7XML: &storeModels.UpdateStoreRKeeper7XMLConfig{
			SeqNumber: &seqNum,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) UpdateStoreByFields(ctx context.Context, req storeModels.UpdateStore) error {
	return s.storeRepository.UpdateStoreByFields(ctx, req)
}

func (s *ServiceImpl) GetStoresBySelectorFilter(ctx context.Context, query selector2.Store) ([]models.Store, error) {
	return s.storeRepository.GetStoresBySelectorFilter(ctx, query)
}
