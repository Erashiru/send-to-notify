package managers

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/integration_api/models"
	"github.com/kwaaka-team/orders-core/core/integration_api/repository"
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
	menuCli "github.com/kwaaka-team/orders-core/pkg/menu"
	orderCli "github.com/kwaaka-team/orders-core/pkg/order"
	orderModels "github.com/kwaaka-team/orders-core/pkg/order/dto"
	storeCli "github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/pkg/errors"
	"time"
)

var externalStatuses = map[string]struct{}{
	"ACCEPTED":         {},
	"COOKING_STARTED":  {},
	"COOKING_COMPLETE": {},
	"CLOSED":           {},
}

type ExternalPosIntegrationClient interface {
	UpdateOrderStatus(ctx context.Context, token, restaurantID, orderID, externalStatus, reason string) error
	UpdateStopList(ctx context.Context, token string, req models.StopListRequest, author string) error
	GetOrders(ctx context.Context, token, restaurantId, status, dateFrom, dateTo string) (models.GetOrdersResponse, error)
}

type ExternalPosIntegrationManager struct {
	orderCli                             orderCli.Client
	menuCli                              menuCli.Client
	storeCli                             storeCli.Client
	externalPosIntegrationAuthRepository repository.ExternalPosIntegrationAuthRepository
}

func NewExternalPosIntegrationManager(orderCli orderCli.Client, menuCli menuCli.Client, storeCli storeCli.Client, externalPosIntegrationAuthRepository repository.ExternalPosIntegrationAuthRepository) ExternalPosIntegrationClient {
	return &ExternalPosIntegrationManager{
		orderCli:                             orderCli,
		menuCli:                              menuCli,
		storeCli:                             storeCli,
		externalPosIntegrationAuthRepository: externalPosIntegrationAuthRepository,
	}
}

func (manager *ExternalPosIntegrationManager) validateRequest(restaurantId, orderId, externalStatus *string) error {
	if restaurantId != nil {
		if *restaurantId == "" {
			return fmt.Errorf("restaurant_id is missing")
		}
	}

	if orderId != nil {
		if *orderId == "" {
			return fmt.Errorf("order_id is missing")
		}
	}

	if externalStatus != nil {
		if *externalStatus == "" {
			return fmt.Errorf("status is missing")
		}
	}

	return nil
}

func (manager *ExternalPosIntegrationManager) UpdateOrderStatus(ctx context.Context, token, restaurantID, orderID, externalStatus, reason string) error {
	if err := manager.validateRequest(&restaurantID, &orderID, &externalStatus); err != nil {
		return err
	}

	// check if the token exists and if it belongs to the restaurant
	if err := manager.validateToken(ctx, token, restaurantID); err != nil {
		return err
	}

	store, err := manager.storeCli.FindStore(ctx, storeModels.StoreSelector{
		ID: restaurantID,
	})
	if err != nil {
		return fmt.Errorf("restaurant %w", err)
	}

	if !manager.validateExternalStatus(externalStatus) {
		return fmt.Errorf("%s status does not belong to the system", externalStatus)
	}

	if err = manager.orderCli.UpdateOrderStatus(ctx, orderID, store.PosType, externalStatus, reason); err != nil {
		return err
	}

	return nil
}

func (manager *ExternalPosIntegrationManager) UpdateStopList(ctx context.Context, token string, req models.StopListRequest, author string) error {
	if err := manager.validateRequest(&req.RestaurantID, nil, nil); err != nil {
		return err
	}

	// check if the token exists and if it belongs to the restaurant
	if err := manager.validateToken(ctx, token, req.RestaurantID); err != nil {
		return err
	}

	if err := manager.menuCli.PosIntegrationUpdateStopList(ctx, req.RestaurantID, req, author); err != nil {
		return err
	}

	return nil
}

func (manager *ExternalPosIntegrationManager) toOrderProductsResponse(order coreOrderModels.Order) []models.GetOrderBodyProduct {
	products := make([]models.GetOrderBodyProduct, 0, len(order.Products))

	for _, product := range order.Products {
		products = append(products, models.GetOrderBodyProduct{
			ID:                   product.ID,
			Price:                int(product.Price.Value),
			PriceWithoutDiscount: int(product.PriceWithoutDiscount.Value),
			Name:                 product.Name,
			Quantity:             product.Quantity,
			Modifiers:            manager.toOrderModifiersResponse(product.Attributes),
		})
	}

	return products
}

func (manager *ExternalPosIntegrationManager) toOrderModifiersResponse(attributes []coreOrderModels.ProductAttribute) []models.GetOrderBodyModifier {
	modifiers := make([]models.GetOrderBodyModifier, 0, len(attributes))

	for _, attribute := range attributes {
		modifiers = append(modifiers, models.GetOrderBodyModifier{
			ID:       attribute.ID,
			Name:     attribute.Name,
			Price:    int(attribute.Price.Value),
			Quantity: attribute.Quantity,
		})
	}

	return modifiers
}

func (manager *ExternalPosIntegrationManager) validateExternalStatus(externalStatus string) bool {
	if _, ok := externalStatuses[externalStatus]; ok {
		return true
	}

	return false
}

func (manager *ExternalPosIntegrationManager) validateToken(ctx context.Context, token, restaurantID string) error {
	authInfo, err := manager.externalPosIntegrationAuthRepository.GetAuthInfo(ctx, token)
	if err != nil {
		return errors.Wrap(err, "authorization token out of system")
	}

	var isExist bool
	for _, id := range authInfo.Restaurants {
		if id == restaurantID {
			isExist = true
			break
		}
	}

	if !isExist {
		return fmt.Errorf("token does not have access to the restaurant")
	}

	return nil
}

func (manager *ExternalPosIntegrationManager) toOrdersResponse(restaurantID string, orders []coreOrderModels.Order, status string) models.GetOrdersResponse {
	result := models.GetOrdersResponse{
		RestaurantID: restaurantID,
		Orders:       make([]models.GetOrderBody, 0, len(orders)),
	}

	for _, order := range orders {
		if status != "" && order.Status != status {
			continue
		}

		paymentType := "Card"
		if order.PaymentMethod == "CASH" {
			paymentType = "Cash"
		}

		orderType := "marketplace"

		if !order.IsMarketplace {
			orderType = "self-delivery"
		}

		if order.IsPickedUpByCustomer {
			orderType = "pickup"
		}

		result.Orders = append(result.Orders, models.GetOrderBody{
			ID:                  order.PosOrderID,
			Status:              order.Status,
			OrderCode:           order.OrderCode,
			PickUpCode:          order.PickUpCode,
			OrderTime:           order.OrderTime.Value.Format(time.RFC3339),
			EstimatedPickUpTime: order.EstimatedPickupTime.Value.Format(time.RFC3339),
			PeopleCount:         order.Persons,
			Address: models.GetOrderBodyAddress{
				Label: order.DeliveryAddress.Label,
			},
			Customer: models.GetOrderBodyCustomer{
				Name:  order.Customer.Name,
				Phone: order.Customer.PhoneNumber,
			},
			Comment:         order.AllergyInfo,
			OrderType:       orderType,
			DeliveryService: order.DeliveryService,
			PaymentInfo: models.GetOrderBodyPaymentInfo{
				Sum:  int(order.EstimatedTotalPrice.Value),
				Type: paymentType,
			},
			Products: manager.toOrderProductsResponse(order),
		})
	}

	return result
}

func (manager *ExternalPosIntegrationManager) GetOrders(ctx context.Context, token, restaurantId, status, dateFrom, dateTo string) (models.GetOrdersResponse, error) {
	if err := manager.validateRequest(&restaurantId, nil, nil); err != nil {
		return models.GetOrdersResponse{}, err
	}

	// check if the token exists and if it belongs to the restaurant
	if err := manager.validateToken(ctx, token, restaurantId); err != nil {
		return models.GetOrdersResponse{}, err
	}

	if !manager.validateExternalStatus(status) && status != "NEW" && status != "" {
		return models.GetOrdersResponse{}, fmt.Errorf("%s status does not belong to the system", status)
	}

	// check if the restaurant exists in the system
	_, err := manager.storeCli.FindStore(ctx, storeModels.StoreSelector{
		ID: restaurantId,
	})
	if err != nil {
		return models.GetOrdersResponse{}, errors.Wrap(err, "restaurant doesn't exist in system")
	}

	// getting orders in the last hour
	orders, err := manager.orderCli.GetActiveOrders(ctx, orderModels.ActiveOrderSelector{
		StoreID: restaurantId,
	})
	if err != nil {
		return models.GetOrdersResponse{}, err
	}

	return manager.toOrdersResponse(restaurantId, orders, status), nil
}
