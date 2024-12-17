package aggregator

import (
	"context"
	expressModels "github.com/kwaaka-team/orders-core/core/express24/models"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models3 "github.com/kwaaka-team/orders-core/core/wolt/models"
	expressConf "github.com/kwaaka-team/orders-core/pkg/express24_v2"
	expressCli "github.com/kwaaka-team/orders-core/pkg/express24_v2/clients"
	"github.com/kwaaka-team/orders-core/pkg/express24_v2/clients/dto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

type express24ServiceV2 struct {
	deliveryServiceName models.Aggregator
	cli                 expressCli.Express24V2
}

func newExpress24v2Service(baseUrl string, store storeModels.Store) (*express24ServiceV2, error) {

	cli, err := expressConf.NewExpress24Client(&expressCli.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		Token:    store.Express24.Token,
	})
	if err != nil {
		return nil, errors.Wrap(constructorError, err.Error())
	}

	return &express24ServiceV2{
		models.EXPRESS24, cli,
	}, nil
}

func (s *express24ServiceV2) MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	return ""
}

func (s *express24ServiceV2) OpenStore(ctx context.Context, aggregatorStoreId string) error {
	return errors.New("method not implemented")
}

func (s *express24ServiceV2) GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error) {
	return false, errors.New("method not implemented")
}

func (s *express24ServiceV2) GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error) {
	return storeModels.AggregatorSchedule{}, errors.New("method not implemented")
}

func (s *express24ServiceV2) UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error {
	return nil
}

func (s *express24ServiceV2) IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error) {
	return store.Express24.IsMarketplace, nil
}

func (s *express24ServiceV2) SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	return nil, nil
}

func (s *express24ServiceV2) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	order, ok := req.(expressModels.Order)
	if !ok {
		return "", errors.New("casting error")
	}

	return order.Store.Branch.ExternalId, nil
}

func (s *express24ServiceV2) GetSystemCreateOrderRequestByAggregatorRequest(req interface{}, store storeModels.Store) (models.Order, error) {

	order, ok := req.(expressModels.Order)
	if !ok {
		return models.Order{}, errors.New("casting error")
	}

	return toOrderRequest(order, store), nil
}

func (s *express24ServiceV2) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error) {
	branch, err := strconv.Atoi(aggregatorStoreID)
	if err != nil {
		return "", err
	}

	err = s.cli.StopListByProducts(ctx, dto.StopListByProductsRequest{
		Products:  s.toProducts(products, isAvailable),
		BranchIDs: []int{branch},
	})

	if err != nil {
		log.Info().Msgf("express24_v2 stoplist ERROR MESSAGE: %s", err.Error())
		return "", err
	}

	return "", nil
}

func (s *express24ServiceV2) toProducts(req menuModels.Products, isAvailable bool) dto.Products {
	items := make([]dto.ProductItem, 0, len(req))

	for i := range req {
		items = append(items, dto.ProductItem{
			ExternalID:  req[i].ExtID,
			Quantity:    menuModels.BASEQUANTITY,
			IsAvailable: isAvailable,
		})
	}
	return dto.Products{
		Items:                   items,
		MakeAvailableOtherItems: false,
	}
}

func (s *express24ServiceV2) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error) {
	branch, err := strconv.Atoi(aggregatorStoreID)
	if err != nil {
		return "", err
	}

	err = s.cli.StopListByProducts(ctx, dto.StopListByProductsRequest{
		Products:  s.toProductsBulk(products),
		BranchIDs: []int{branch},
	},
	)

	if err != nil {
		log.Info().Msgf("express24_v2 stoplist ERROR MESSAGE: %s", err.Error())
		return "", err
	}

	return "", nil
}

func (s *express24ServiceV2) toProductsBulk(req menuModels.Products) dto.Products {
	items := make([]dto.ProductItem, 0, len(req))

	for i := range req {
		items = append(items, dto.ProductItem{
			ExternalID:  req[i].ExtID,
			Quantity:    menuModels.BASEQUANTITY,
			IsAvailable: req[i].IsAvailable,
		})
	}
	return dto.Products{
		Items:                   items,
		MakeAvailableOtherItems: false,
	}
}

func (s *express24ServiceV2) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error) {
	branch, err := strconv.Atoi(aggregatorStoreID)
	if err != nil {
		return "", err
	}

	err = s.cli.StopListByAttributes(ctx, dto.StopListByAttributesRequest{
		Modifiers: s.toAttributes(attributes),
		BranchIDs: []int{branch},
	})

	if err != nil {
		log.Info().Msgf("express24_v2 stoplist ERROR MESSAGE: %s", err.Error())
		return "", err
	}

	return "", nil
}

func (s *express24ServiceV2) toAttributes(req menuModels.Attributes) dto.Modifiers {
	items := make([]dto.AttributeItem, 0, len(req))

	for i := range req {
		items = append(items, dto.AttributeItem{
			ExternalID:  req[i].ExtID,
			IsAvailable: req[i].IsAvailable,
		})
	}
	return dto.Modifiers{
		Items:                   items,
		MakeAvailableOtherItems: false,
	}
}

func (s *express24ServiceV2) GetAggregatorOrder(ctx context.Context, orderID string) (models3.Order, error) {
	return models3.Order{}, nil
}

func (s *express24ServiceV2) SendOrderErrorNotification(ctx context.Context, req interface{}) error {
	return nil
}

func (s *express24ServiceV2) SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error {
	return nil
}
