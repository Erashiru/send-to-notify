package promo_code

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	models3 "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	models2 "github.com/kwaaka-team/orders-core/core/qrmenu/models"
	models4 "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/promo_code/dto"
	"github.com/kwaaka-team/orders-core/service/promo_code/repository"
	"github.com/kwaaka-team/orders-core/service/promo_code/user_repository"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math"
	"sync"
	"time"
)

type Service interface {
	Create(ctx context.Context, promoRequest models.PromoCode) error
	Update(ctx context.Context, updatePromoCodeRequest models.UpdatePromoCode) error
	GetByID(ctx context.Context, promoCodeID string) (models.PromoCode, error)
	GetPromoCodesByRestaurantId(ctx context.Context, restaurantId string, pagination selector.Pagination) ([]models.PromoCode, error)
	ValidatePromoCodeForUser(ctx context.Context, userPromoCode models.ValidateUserPromoCode) (bool, string, float64, int, []models2.Product, error)
	GetPromoCodeByCodeAndRestaurantId(ctx context.Context, promoCodeValue string, restaurantId string) (models.PromoCode, error)
	AddUserPromoCodeUsageTimeToDB(ctx context.Context, userId string, promoCodeValue string, restaurantId string) error
	GetAvailablePromoCodeByCode(ctx context.Context, promoCodeValue string) (models.PromoCode, error)
}

type ServiceImpl struct {
	repository        repository.Repository
	userPromoCodeRepo user_repository.Repository
	store             store.Service
	menuService       *menuServicePkg.Service
	logger            *zap.SugaredLogger
}

func NewPromoCodeService(logger *zap.SugaredLogger, repo repository.Repository, userRepo user_repository.Repository, store store.Service, menu *menuServicePkg.Service) (*ServiceImpl, error) {
	return &ServiceImpl{
		logger:            logger,
		repository:        repo,
		userPromoCodeRepo: userRepo,
		store:             store,
		menuService:       menu,
	}, nil
}

func (s *ServiceImpl) Create(ctx context.Context, promoRequest models.PromoCode) error {

	s.logger.Infof("service: create promo code for request: %v", promoRequest)

	if err := s.validatePromoCodeProducts(ctx, promoRequest.RestaurantIDs, promoRequest.Product); err != nil {
		return err
	}

	err := s.repository.CreatePromo(ctx, promoRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) validatePromoCodeProducts(ctx context.Context, restaurantIds []string, products []models2.Product) error {

	s.logger.Infof("service: validate promo code products exists in restaurants")

	for i := range restaurantIds {

		errorCh := s.validateMenusProduct(ctx, restaurantIds[i], products)
		for err := range errorCh {
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ServiceImpl) validateMenusProduct(ctx context.Context, restaurantId string, promoProductIds []models2.Product) chan error {

	errorCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		mapDbProducts := make(map[string]models3.Product)

		restaurant, err := s.store.GetByID(ctx, restaurantId)
		if err != nil {
			errorCh <- err
		}

		qrMenu := restaurant.Menus.GetActiveMenu(models4.QRMENU)

		dBProducts, _, err := s.menuService.ListProductsByMenuId(ctx, qrMenu.ID)
		if err != nil {
			errorCh <- err
		}

		for i := range dBProducts {
			product := dBProducts[i]
			mapDbProducts[product.ExtID] = product
		}

		for i := range promoProductIds {
			product := promoProductIds[i]
			if mapDbProducts[product.ProductID].ExtID == "" {
				errorCh <- errors.New("product with id: " + product.ProductID + ", name: " + product.Name + "does not exist in restaurant: " + restaurant.Name)
			}
		}

	}()

	go func() {
		wg.Wait()
		close(errorCh)
	}()

	return errorCh
}

func (s *ServiceImpl) Update(ctx context.Context, updatePromoCodeRequest models.UpdatePromoCode) error {

	s.logger.Infof("service: update promo code for request: %v", &updatePromoCodeRequest)

	return s.repository.UpdatePromo(ctx, updatePromoCodeRequest)
}

func (s *ServiceImpl) GetPromoCodesByRestaurantId(ctx context.Context, restaurantId string, pagination selector.Pagination) ([]models.PromoCode, error) {

	s.logger.Infof("service: get promo codes by restaurant id: %v", restaurantId)

	return s.repository.GetPromoCodesByRestaurantId(ctx, restaurantId, pagination)
}

func (s *ServiceImpl) GetByID(ctx context.Context, promoCodeID string) (models.PromoCode, error) {

	s.logger.Infof("service: get promo code by id: %v", promoCodeID)

	return s.repository.GetPromoCodeByID(ctx, promoCodeID)
}

func (s *ServiceImpl) GetAvailablePromoCodeByCode(ctx context.Context, promoCodeValue string) (models.PromoCode, error) {

	s.logger.Infof("service: get available promo code by code: %s", promoCodeValue)

	return s.repository.GetAvailablePromoCodeByCode(ctx, promoCodeValue)
}

func (s *ServiceImpl) ValidatePromoCodeForUser(ctx context.Context, userPromoCode models.ValidateUserPromoCode) (bool, string, float64, int, []models2.Product, error) {

	s.logger.Infof("service: start to validate promo code: %v for user id:%v", userPromoCode.PromoCode, userPromoCode.UserId)

	promo, err := s.repository.GetPromoByCodeAndRestaurantID(ctx, userPromoCode.PromoCode, userPromoCode.RestaurantID)
	if err != nil {
		s.logger.Infof("promo code: %v does not exist", userPromoCode.PromoCode)
		if errors.Is(err, dto.ErrPromoCodeNotFound) {
			return false, fmt.Sprintf("Промокод: %v не существует", userPromoCode.PromoCode), 0, 0, nil, nil
		}
		return false, "", 0, 0, nil, err
	}

	usageTime, err := s.userPromoCodeRepo.GetUsageCountForUser(ctx, userPromoCode.UserId, userPromoCode.PromoCode, userPromoCode.RestaurantID)
	if err != nil {
		return false, "", 0, 0, nil, err
	}

	comment, ok := s.checkPromoCode(promo, userPromoCode, usageTime)
	if !ok {
		return false, comment, 0, 0, nil, err
	}

	exist, totalSum, saleSum, products := s.checkPromoCodeType(promo, userPromoCode)
	if !exist {
		return false, "", 0, 0, nil, nil
	}

	return true, comment, totalSum, saleSum, products, nil
}

func (s *ServiceImpl) AddUserPromoCodeUsageTimeToDB(ctx context.Context, userId string, promoCodeValue string, restaurantId string) error {

	usageTime, err := s.userPromoCodeRepo.GetUsageCountForUser(ctx, userId, promoCodeValue, restaurantId)
	if err != nil {
		return err
	}

	promoCode, err := s.repository.GetPromoByCodeAndRestaurantID(ctx, promoCodeValue, restaurantId)
	if err != nil {
		return err
	}

	switch usageTime {
	case 0:
		if err := s.userPromoCodeRepo.CreateUserUsePromoCodeTime(ctx, userId, promoCodeValue, promoCode.RestaurantIDs); err != nil {
			return err
		}
	default:
		if err := s.userPromoCodeRepo.UpdateUsageTimeForUser(ctx, userId, promoCodeValue, restaurantId, usageTime); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) checkPromoCode(promoFromDB models.PromoCode, promoFromCart models.ValidateUserPromoCode, usagedTime int) (string, bool) {

	if !promoFromDB.Available {
		return fmt.Sprintf("Промокод %v сейчас не активный", promoFromCart.PromoCode), false
	}

	if !s.contains(promoFromDB.DeliveryType, promoFromCart.DeliveryType) {
		return fmt.Sprintf("Промокод %v действителен, когда тип доставки: %v", promoFromCart.PromoCode, promoFromDB.DeliveryType), false
	}

	if !(time.Now().After(promoFromDB.ValidFrom) && time.Now().Before(promoFromDB.ValidUntil)) {
		return fmt.Sprintf("Промокод %v действителен с %v до %v", promoFromCart.PromoCode, promoFromDB.ValidFrom, promoFromDB.ValidUntil), false
	}

	if promoFromDB.MinimumOrderPrice > promoFromCart.TotalSum {
		return fmt.Sprintf("Чтобы использовать промокод %v, дополните корзину до %v тенге", promoFromCart.PromoCode, promoFromDB.MinimumOrderPrice), false
	}

	if promoFromDB.UsageTime <= usagedTime {
		return fmt.Sprintf("Промокод %v можно использовать %v раз, вы исчерпали свой лимит", promoFromCart.PromoCode, promoFromDB.UsageTime), false
	}

	return fmt.Sprintf("Промокод %v сейчас активный", promoFromCart.PromoCode), true
}

func (s *ServiceImpl) checkPromoCodeType(promoFromDB models.PromoCode, promoFromCart models.ValidateUserPromoCode) (bool, float64, int, []models2.Product) {

	switch promoFromDB.PromoCodeCategory {
	case "gift":
		return true, float64(promoFromCart.TotalSum), 0, promoFromDB.Product
	case "sale":
		if promoFromDB.ForAllProduct || len(promoFromDB.Product) == 0 || s.containsProduct(promoFromCart.CartProducts, promoFromDB.Product) {
			switch promoFromDB.SaleType {
			case "percentage":
				return true, math.Ceil(float64(promoFromCart.TotalSum) - float64(promoFromCart.TotalSum*promoFromDB.Sale)/100), (promoFromCart.TotalSum * promoFromDB.Sale) / 100, []models2.Product{}
			case "currency":
				return true, float64(promoFromCart.TotalSum - promoFromDB.Sale), promoFromDB.Sale, []models2.Product{}
			}
		}
	}
	return false, 0, 0, []models2.Product{}
}

func (s *ServiceImpl) contains(array []string, e string) bool {
	for _, a := range array {
		if a == e {
			return true
		}
	}
	return false
}

func (s *ServiceImpl) containsProduct(cart []models2.Product, promoFromDB []models2.Product) bool {
	for i := range cart {
		for j := range promoFromDB {
			if cart[i].ProductID == promoFromDB[j].ProductID {
				return true
			}
		}
	}
	return false
}

func (s *ServiceImpl) GetPromoCodeByCodeAndRestaurantId(ctx context.Context, promoCodeValue string, restaurantId string) (models.PromoCode, error) {

	return s.repository.GetPromoByCodeAndRestaurantID(ctx, promoCodeValue, restaurantId)
}
