package pos

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/config"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/pkg/ytimes"
	"github.com/kwaaka-team/orders-core/pkg/ytimes/clients"
	ytimesModels "github.com/kwaaka-team/orders-core/pkg/ytimes/clients/models"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"time"
)

type ytimesPosService struct {
	*BasePosService
	ytimesCli clients.Client
	store     coreStoreModels.Store
}

func newYtimesPosService(bps *BasePosService, baseUrl string, store coreStoreModels.Store) (*ytimesPosService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "ytimes pos constructor error")
	}

	ytimesClient := ytimes.New(clients.Config{
		BaseUrl: baseUrl,
		Token:   store.YTimes.AuthToken,
	})

	return &ytimesPosService{
		BasePosService: bps,
		ytimesCli:      ytimesClient,
		store:          store,
	}, nil
}

func (ks *ytimesPosService) toItems(req models.Order) []ytimesModels.OrderItemList {
	orderItemList := make([]ytimesModels.OrderItemList, 0, len(req.Products))

	for i := 0; i < len(req.Products); i++ {
		posProduct := ytimesModels.OrderItemList{
			PriceWithDiscount: req.Products[i].Price.Value,
			Quantity:          req.Products[i].Quantity,
			SupplementList:    make(map[string]int),
		}

		if req.Products[i].SizeId == "" {
			posProduct.GoodsItemGuid = &req.Products[i].ID
		} else {
			posProduct.MenuItemGuid = &req.Products[i].ID
			posProduct.MenuTypeGuid = &req.Products[i].SizeId

			for _, attribute := range req.Products[i].Attributes {
				posProduct.SupplementList[attribute.ID] = attribute.Quantity
			}
		}

		orderItemList = append(orderItemList, posProduct)
	}

	return orderItemList
}

func (ks *ytimesPosService) constructPosOrder(req models.Order, store coreStoreModels.Store) (models.Order, ytimesModels.Order, error) {
	comment, _, _ := ks.constructOrderComments(req, store)

	items := ks.toItems(req)

	posOrder := ytimesModels.Order{
		Type:             "DELIVERY",
		ShopGuid:         store.YTimes.PointId,
		PrintFiscalCheck: true,
		ItemList:         items,
		Comment:          comment,
		PaidValue:        req.EstimatedTotalPrice.Value - req.PartnerDiscountsProducts.Value,
	}

	return req, posOrder, nil
}

func (ks *ytimesPosService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration, store coreStoreModels.Store,
	menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	var err error

	order, err = prepareAnOrder(ctx, order, store, menu, aggregatorMenu, menuClient)
	if err != nil {
		return order, err
	}

	var posOrder ytimesModels.Order

	order, posOrder, err = ks.constructPosOrder(order, store)
	if err != nil {
		return order, err
	}

	order, err = ks.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	response, err := ks.ytimesCli.CreateOrder(ctx, posOrder)
	if err != nil {
		return order, err
	}

	var posOrderId string

	if len(response.Rows) > 0 {
		posOrderId = response.Rows[0].Guid
	}

	order = setPosOrderId(order, posOrderId)

	return order, nil
}

func (ks *ytimesPosService) IsAliveStatus(ctx context.Context, store coreStoreModels.Store) (bool, error) {
	return true, nil
}

func (ks *ytimesPosService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {
	switch posStatus {
	case "ACCEPTED":
		return models.ACCEPTED, nil
	case "CLOSED":
		return models.CLOSED, nil
	case "CANCELLED":
		return models.CANCELLED_BY_POS_SYSTEM, nil
	}

	return 0, fmt.Errorf("undefined pos status: %s", posStatus)
}

func (ks *ytimesPosService) checkForClosed(now, orderTime time.Time, minute int) bool {
	if now.After(orderTime.Add(time.Duration(minute) * time.Minute)) {
		return true
	}

	return false
}

func (ks *ytimesPosService) GetOrderStatus(ctx context.Context, order models.Order) (string, error) {
	var (
		now                  = time.Now().UTC()
		orderTime            = order.OrderTime.Value.Time
		adjustedPickupMinute int
	)

	switch order.DeliveryService {
	case models.GLOVO.String():
		adjustedPickupMinute = ks.store.Glovo.AdjustedPickupMinutes
	case models.WOLT.String():
		adjustedPickupMinute = ks.store.Wolt.AdjustedPickupMinutes
	case models.EXPRESS24.String():
		adjustedPickupMinute = ks.store.Express24.AdjustedPickupMinutes
	case models.QRMENU.String():
		adjustedPickupMinute = ks.store.QRMenu.AdjustedPickupMinutes
	case models.YANDEX.String(), models.EMENU.String():
		for _, externalCfg := range ks.store.ExternalConfig {
			if externalCfg.Type == order.DeliveryService {
				adjustedPickupMinute = externalCfg.AdjustedPickupMinutes
				break
			}
		}
	}

	if adjustedPickupMinute > 0 && ks.checkForClosed(now, orderTime, adjustedPickupMinute) {
		return "CLOSED", nil
	}

	return "ACCEPTED", nil
}

func (ks *ytimesPosService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {
	return coreMenuModels.StopListItems{}, ErrUnsupportedMethod
}

func (ks *ytimesPosService) getSections(rows []ytimesModels.MenuRow) []coreMenuModels.Section {
	sections := make([]coreMenuModels.Section, 0, len(rows))

	for _, row := range rows {
		sections = append(sections, coreMenuModels.Section{
			ExtID:        row.Guid,
			Name:         row.Name,
			SectionOrder: row.Priority,
			ImageUrl:     row.ImageLink,
		})
	}

	return sections
}

func (ks *ytimesPosService) itemToSystemProducts(item ytimesModels.ItemList, sectionId string) []coreMenuModels.Product {
	systemProducts := make([]coreMenuModels.Product, 0, len(item.TypeList))

	productName := item.Name

	systemProduct := coreMenuModels.Product{
		ExtID:     item.Guid,
		ProductID: item.Guid,
		ImageURLs: []string{item.ImageLink},
		Description: []coreMenuModels.LanguageDescription{
			{
				Value: item.Description,
			},
		},
		Section:          sectionId,
		IsAvailable:      true,
		IsIncludedInMenu: true,
		IsSync:           true,
	}

	for key := range item.SupplementCategoryToFreeCount {
		systemProduct.AttributesGroups = append(systemProduct.AttributesGroups, key)
	}

	if len(item.TypeList) > 1 {
		for i := 0; i < len(item.TypeList); i++ {
			systemProduct.ExtID = item.Guid + item.TypeList[i].Guid
			systemProduct.SizeID = item.TypeList[i].Guid
			systemProduct.Name = []coreMenuModels.LanguageDescription{
				{
					Value: item.TypeList[i].Name + " " + productName,
				},
			}
			systemProduct.Price = []coreMenuModels.Price{
				{
					Value: item.TypeList[i].Price,
				},
			}
			systemProducts = append(systemProducts, systemProduct)
		}
	} else if len(item.TypeList) == 1 {
		systemProduct.Name = []coreMenuModels.LanguageDescription{
			{
				Value: item.TypeList[0].Name + " " + productName,
			},
		}
		systemProduct.SizeID = item.TypeList[0].Guid
		systemProduct.Price = []coreMenuModels.Price{
			{
				Value: item.TypeList[0].Price,
			},
		}
		systemProducts = append(systemProducts, systemProduct)
	}

	// TODO: default supplements docs is empty

	return systemProducts
}

func (ks *ytimesPosService) goodsItemToSystemProduct(goodsItem ytimesModels.GoodsList, sectionId string) coreMenuModels.Product {
	return coreMenuModels.Product{
		ExtID:     goodsItem.Guid,
		ProductID: goodsItem.Guid,
		Name: []coreMenuModels.LanguageDescription{
			{
				Value: goodsItem.Name,
			},
		},
		ImageURLs: []string{goodsItem.ImageLink},
		Description: []coreMenuModels.LanguageDescription{
			{
				Value: goodsItem.Description,
			},
		},
		Price: []coreMenuModels.Price{
			{
				Value: float64(goodsItem.Price),
			},
		},
		Section:          sectionId,
		IsAvailable:      true,
		IsIncludedInMenu: true,
		IsSync:           true,
	}
}

func (ks *ytimesPosService) toSystemProducts(rows []ytimesModels.MenuRow) []coreMenuModels.Product {
	var (
		systemProducts = make([]coreMenuModels.Product, 0, 4)
	)

	for _, row := range rows {
		for _, categoryList := range row.CategoryList {
			for _, item := range categoryList.ItemList {
				systemProducts = append(systemProducts, ks.itemToSystemProducts(item, categoryList.Guid)...)
			}

			for _, good := range categoryList.GoodsList {
				systemProduct := ks.goodsItemToSystemProduct(good, categoryList.Guid)
				systemProducts = append(systemProducts, systemProduct)
			}
		}

		for _, item := range row.ItemList {
			systemProducts = append(systemProducts, ks.itemToSystemProducts(item, row.Guid)...)
		}

		for _, goodsItem := range row.GoodsList {
			systemProduct := ks.goodsItemToSystemProduct(goodsItem, row.Guid)
			systemProducts = append(systemProducts, systemProduct)
		}
	}

	return systemProducts
}

func (ks *ytimesPosService) toSystemAttributesAndAttributeGroups(rows []ytimesModels.SupplementRow) ([]coreMenuModels.AttributeGroup, []coreMenuModels.Attribute) {
	var (
		attributeGroups    = make([]coreMenuModels.AttributeGroup, 0, len(rows))
		attributes         = make([]coreMenuModels.Attribute, 0, 4)
		existingAttributes = make(map[string]bool)
	)

	for _, row := range rows {
		attributeIds := make([]string, 0, len(row.ItemList))

		for _, attribute := range row.ItemList {
			if !existingAttributes[attribute.Guid] {
				attributes = append(attributes, coreMenuModels.Attribute{
					ExtID: attribute.Guid,
					Name:  attribute.Name,
					Price: float64(attribute.DefaultPrice),
				})

				existingAttributes[attribute.Guid] = true
			}

			attributeIds = append(attributeIds, attribute.Guid)
		}

		attributeGroups = append(attributeGroups, coreMenuModels.AttributeGroup{
			ExtID:          row.Guid,
			Name:           row.Name,
			Min:            0,
			Max:            row.MaxSelectedCount,
			MultiSelection: row.AllowSeveralItem,
			Attributes:     attributeIds,
		})

	}

	return attributeGroups, attributes
}

func (ks *ytimesPosService) toSystemMenu(posMenu ytimesModels.Menu, posModifierGroupList ytimesModels.SupplementList) coreMenuModels.Menu {
	systemMenu := coreMenuModels.Menu{
		Sections: ks.getSections(posMenu.Rows),
	}

	systemProducts := ks.toSystemProducts(posMenu.Rows)

	systemAttributeGroups, systemAttributes := ks.toSystemAttributesAndAttributeGroups(posModifierGroupList.Rows)

	systemMenu.Products = systemProducts
	systemMenu.Attributes = systemAttributes
	systemMenu.AttributesGroups = systemAttributeGroups

	return systemMenu
}

func (ks *ytimesPosService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	responseMenu, err := ks.ytimesCli.GetMenu(ctx, store.YTimes.PointId)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	responseSupplementList, err := ks.ytimesCli.GetSupplementList(ctx, store.YTimes.PointId)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	systemMenu := ks.toSystemMenu(responseMenu, responseSupplementList)

	return systemMenu, nil
}

func (ks *ytimesPosService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	return nil
}

func (s *ytimesPosService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (s *ytimesPosService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (s *ytimesPosService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
