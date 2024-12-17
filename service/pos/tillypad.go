package pos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/config"
	yandexDeliveryProtocolModels "github.com/kwaaka-team/orders-core/core/externalapi/models"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/pkg/tillypad/yandexDeliveryProtocolTillypad"
	"github.com/kwaaka-team/orders-core/pkg/tillypad/yandexDeliveryProtocolTillypad/clients"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type tillypadPosService struct {
	*BasePosService
	tillypadCli clients.Client
	pointId     string
}

func newTillypadPosService(bps *BasePosService, baseUrl, pointId, clientId, clientSecret, pathPrefix string) (*tillypadPosService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "tillypad pos constructor error")
	}

	tillypadClient, err := yandexDeliveryProtocolTillypad.NewTillypadClient(clients.Config{
		BaseURL:      baseUrl,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		PathPrefix:   pathPrefix,
	})
	if err != nil {
		return nil, err
	}

	return &tillypadPosService{
		BasePosService: bps,
		tillypadCli:    tillypadClient,
		pointId:        pointId,
	}, nil
}

func (ks *tillypadPosService) toItems(req models.Order) []yandexDeliveryProtocolModels.OrderItem {
	items := make([]yandexDeliveryProtocolModels.OrderItem, 0, len(req.Products))

	for _, product := range req.Products {
		modifications := make([]yandexDeliveryProtocolModels.OrderModification, 0, len(product.Attributes))

		for _, attribute := range product.Attributes {
			modifications = append(modifications, yandexDeliveryProtocolModels.OrderModification{
				Id:       attribute.ID,
				Name:     attribute.Name,
				Quantity: attribute.Quantity,
				Price:    int(attribute.Price.Value),
			})
		}

		items = append(items, yandexDeliveryProtocolModels.OrderItem{
			Id:            product.ID,
			Name:          product.Name,
			Quantity:      float64(product.Quantity),
			Price:         int(product.Price.Value),
			Modifications: modifications,
		})
	}

	return items
}

func (ks *tillypadPosService) toPayments(req models.Order) (models.Order, string, int, error) {
	itemCost := int(req.EstimatedTotalPrice.Value) - int(req.PartnerDiscountsProducts.Value)

	var paymentTypeKind string

	switch req.PosPaymentInfo.PaymentTypeKind {
	case "Cash":
		paymentTypeKind = "CASH"
	case "Card":
		paymentTypeKind = "CARD"
	default:
		req.FailReason.Code = PAYMENT_TYPE_MISSED_CODE
		req.FailReason.Message = PAYMENT_TYPE_MISSED
		return req, "", 0, errors.New("payment type is missed")
	}

	return req, paymentTypeKind, itemCost, nil
}

func (ks *tillypadPosService) getDeliveryInfo(req models.Order, store coreStoreModels.Store) yandexDeliveryProtocolModels.DeliveryInfo {
	completeBeforeDate := req.EstimatedPickupTime.Value.Time

	if completeBeforeDate.IsZero() {
		completeBeforeDate = time.Now().UTC().Add(time.Hour)
	}

	// Convert UTC time to stores local time
	completeBeforeDate = completeBeforeDate.Add(time.Duration(store.Settings.TimeZone.UTCOffset) * time.Hour)

	completeBefore := completeBeforeDate.Format(time.RFC3339Nano)

	var deliveryAddress *yandexDeliveryProtocolModels.MarketPlaceDeliveryAddress
	switch {
	case req.DeliveryAddress.Label == "":
		deliveryAddress = &yandexDeliveryProtocolModels.MarketPlaceDeliveryAddress{
			Full:      "N/A",
			Latitude:  "0.0",
			Longitude: "0.0",
		}
	default:
		deliveryAddress = &yandexDeliveryProtocolModels.MarketPlaceDeliveryAddress{
			Full:      req.DeliveryAddress.Label,
			Latitude:  strconv.FormatFloat(req.DeliveryAddress.Latitude, 'f', -1, 64),
			Longitude: strconv.FormatFloat(req.DeliveryAddress.Longitude, 'f', -1, 64),
		}
	}

	clientName := req.Customer.Name
	if clientName == "" {
		clientName = "N/A"
	}

	customerPhoneNumber := setCustomerPhoneNumberForDeliveryInfo(req.RestaurantID, req.DeliveryService, req.Customer.PhoneNumber)

	return yandexDeliveryProtocolModels.DeliveryInfo{
		ClientName:              clientName,
		PhoneNumber:             customerPhoneNumber,
		MarketPlaceDeliveryDate: completeBefore,
		DeliveryAddress:         deliveryAddress,
	}
}

func (ks *tillypadPosService) constructPosOrder(req models.Order, store coreStoreModels.Store) (models.Order, yandexDeliveryProtocolModels.Order, error) {
	posOrderId := uuid.New().String()

	orderComments, _, _ := ks.BasePosService.constructOrderComments(req, store)

	req = setPosOrderId(req, posOrderId)

	items := ks.toItems(req)

	req, paymentType, itemCost, err := ks.toPayments(req)
	if err != nil {
		return req, yandexDeliveryProtocolModels.Order{}, err
	}

	deliveryInfo := ks.getDeliveryInfo(req, store)

	codePrefix := ""
	switch req.DeliveryService {
	case models.WOLT.String():
		codePrefix = store.Wolt.OrderCodePrefix
	case models.GLOVO.String():
		codePrefix = store.Glovo.OrderCodePrefix
	}

	orderDate := time.Now().Format("060102") // YYMMDD format
	orderCode := fmt.Sprintf("%s_%s_%s_%s", codePrefix, strings.ToUpper(req.DeliveryService), orderDate, req.OrderCode)

	return req, yandexDeliveryProtocolModels.Order{
		Platform:      "YE",
		Discriminator: "marketplace",
		EatsId:        orderCode,
		RestaurantId:  ks.pointId,
		DeliveryInfo:  deliveryInfo,
		Comment:       orderComments,
		PaymentInfo: yandexDeliveryProtocolModels.PaymentInfo{
			PaymentType: paymentType,
			ItemsCost:   itemCost,
		},
		Items: items,
	}, nil
}

func (ks *tillypadPosService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration,
	store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	var err error

	order, err = prepareAnOrder(ctx, order, store, menu, aggregatorMenu, menuClient)
	if err != nil {
		return order, err
	}

	utils.Beautify("prepared order", order)

	var posOrder yandexDeliveryProtocolModels.Order

	order, posOrder, err = ks.constructPosOrder(order, store)
	if err != nil {
		return order, err
	}

	order, err = ks.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	utils.Beautify("pos order body", posOrder)

	response, err := ks.tillypadCli.CreateOrder(ctx, posOrder)
	if err != nil {
		return order, err
	}

	order = setPosOrderId(order, response.OrderId)

	return order, nil
}

func (ks *tillypadPosService) IsAliveStatus(ctx context.Context, store coreStoreModels.Store) (bool, error) {
	return true, nil
}

func (ks *tillypadPosService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {
	switch posStatus {
	case "NEW":
		return models.NEW, nil
	case "ACCEPTED_BY_RESTAURANT":
		return models.ACCEPTED, nil
	case "COOKING":
		return models.COOKING_STARTED, nil
	case "READY":
		return models.COOKING_COMPLETE, nil
	case "TAKEN_BY_COURIER":
		return models.CLOSED, nil
	case "DELIVERED":
		return models.DELIVERED, nil
	case "CANCELLED":
		return models.CANCELLED_BY_POS_SYSTEM, nil
	}

	return 0, fmt.Errorf("undefined pos status: %s", posStatus)
}

func (ks *tillypadPosService) GetOrderStatus(ctx context.Context, order models.Order) (string, error) {
	response, err := ks.tillypadCli.GetOrderStatus(ctx, order.PosOrderID)
	if err != nil {
		return "", err
	}

	return response.Status, nil
}

func (ks *tillypadPosService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {
	var stopListItems coreMenuModels.StopListItems

	response, err := ks.tillypadCli.GetAvailability(ctx, ks.pointId)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}

	for _, item := range response.Items {
		stopListItems = append(stopListItems, coreMenuModels.StopListItem{
			ProductID: item.ItemId,
			Balance:   float64(item.Stock),
		})
	}

	for _, modifier := range response.Modifiers {
		stopListItems = append(stopListItems, coreMenuModels.StopListItem{
			ProductID: modifier.ModifierId,
			Balance:   float64(modifier.Stock),
		})
	}

	return stopListItems, nil
}

func (ks *tillypadPosService) getSections(tillypadCategories []yandexDeliveryProtocolModels.Category) []coreMenuModels.Section {
	var (
		sections = make([]coreMenuModels.Section, 0, len(tillypadCategories))
	)

	for _, category := range tillypadCategories {
		if category.ParentId == "" {
			var imageUrl string

			if len(category.Images) != 0 {
				imageUrl = category.Images[0].Url
			}

			sections = append(sections, coreMenuModels.Section{
				ExtID:        category.Id,
				Name:         category.Name,
				SectionOrder: category.SortOrder,
				ImageUrl:     imageUrl,
			})
		}
	}

	return sections
}

func (ks *tillypadPosService) fillProduct(tillypadItem yandexDeliveryProtocolModels.Item) coreMenuModels.Product {
	images := make([]string, 0, len(tillypadItem.Images))

	for _, itemImage := range tillypadItem.Images {
		images = append(images, itemImage.Url)
	}

	return coreMenuModels.Product{
		ExtID:     tillypadItem.Id,
		ProductID: tillypadItem.Id,
		Section:   tillypadItem.CategoryId,
		Name: []coreMenuModels.LanguageDescription{
			{
				Value: tillypadItem.Name,
			},
		},
		Description: []coreMenuModels.LanguageDescription{
			{
				Value: tillypadItem.Description,
			},
		},
		Price: []coreMenuModels.Price{
			{
				Value: float64(tillypadItem.Price),
			},
		},
		Weight:           float64(tillypadItem.Measure),
		MeasureUnit:      tillypadItem.MeasureUnit,
		IsIncludedInMenu: true,
		ImageURLs:        images,
	}
}

func (ks *tillypadPosService) fillAttributeGroup(modifierGroup yandexDeliveryProtocolModels.ModifierGroup) coreMenuModels.AttributeGroup {
	return coreMenuModels.AttributeGroup{
		ExtID: modifierGroup.Id,
		Name:  modifierGroup.Name,
		Min:   modifierGroup.MinSelectedModifiers,
		Max:   modifierGroup.MaxSelectedModifiers,
	}
}

func (ks *tillypadPosService) addIdToMap(entity map[string]bool, id string) {
	entity[id] = true
}

func (ks *tillypadPosService) fillAttribute(modifier yandexDeliveryProtocolModels.Modifier) coreMenuModels.Attribute {
	return coreMenuModels.Attribute{
		ExtID: modifier.Id,
		PosID: modifier.Id,
		Name:  modifier.Name,
		Price: float64(modifier.Price),
		Min:   modifier.MinAmount,
		Max:   modifier.MaxAmount,
	}
}

func (ks *tillypadPosService) checkIfExists(entity map[string]bool, id string) bool {
	return entity[id]
}

func (ks *tillypadPosService) addAttributeGroupIdsToProduct(product coreMenuModels.Product, ids []string) coreMenuModels.Product {
	product.AttributesGroups = ids
	return product
}

func (ks *tillypadPosService) addAttributeIdsToAttributeGroup(attributeGroup coreMenuModels.AttributeGroup, ids []string) coreMenuModels.AttributeGroup {
	attributeGroup.Attributes = ids
	return attributeGroup
}

func (ks *tillypadPosService) getEntities(tillypadItems []yandexDeliveryProtocolModels.Item) ([]coreMenuModels.Product, []coreMenuModels.AttributeGroup, []coreMenuModels.Attribute) {
	var (
		systemProducts        = make([]coreMenuModels.Product, 0, len(tillypadItems))
		systemAttributeGroups = make([]coreMenuModels.AttributeGroup, 0, 4)
		systemAttributes      = make([]coreMenuModels.Attribute, 0, 4)
	)

	existingAttributeGroups := make(map[string]bool)
	existingAttributes := make(map[string]bool)

	for _, item := range tillypadItems {
		attributeGroupIds := make([]string, 0, len(item.ModifierGroups))

		for _, modifierGroup := range item.ModifierGroups {
			attributeGroupIds = append(attributeGroupIds, modifierGroup.Id)

			// if not exist
			if ks.checkIfExists(existingAttributeGroups, modifierGroup.Id) {

				attributeIds := make([]string, 0, len(modifierGroup.Modifiers))

				for _, modifier := range modifierGroup.Modifiers {
					attributeIds = append(attributeIds, modifier.Id)

					// if not exist
					if ks.checkIfExists(existingAttributes, modifier.Id) {
						ks.addIdToMap(existingAttributes, modifier.Id)
						systemAttribute := ks.fillAttribute(modifier)
						systemAttributes = append(systemAttributes, systemAttribute)
					}
				}

				ks.addIdToMap(existingAttributeGroups, modifierGroup.Id)

				systemAttributeGroup := ks.fillAttributeGroup(modifierGroup)
				systemAttributeGroup = ks.addAttributeIdsToAttributeGroup(systemAttributeGroup, attributeIds)
				systemAttributeGroups = append(systemAttributeGroups, systemAttributeGroup)
			}

		}

		systemProduct := ks.fillProduct(item)
		systemProduct = ks.addAttributeGroupIdsToProduct(systemProduct, attributeGroupIds)
		systemProducts = append(systemProducts, systemProduct)
	}

	return systemProducts, systemAttributeGroups, systemAttributes
}

func (ks *tillypadPosService) setProductsAvailability(systemProducts []coreMenuModels.Product, existingProductStopList map[string]bool) []coreMenuModels.Product {
	for i := range systemProducts {
		if !existingProductStopList[systemProducts[i].ProductID] {
			systemProducts[i].IsAvailable = true
		}
	}

	return systemProducts
}

func (ks *tillypadPosService) setAttributesAvailability(systemAttributes []coreMenuModels.Attribute, existingAttributeStopList map[string]bool) []coreMenuModels.Attribute {
	for i := range systemAttributes {
		if !existingAttributeStopList[systemAttributes[i].ExtID] {
			systemAttributes[i].IsAvailable = true
		}
	}

	return systemAttributes
}

func (ks *tillypadPosService) toSystemMenu(tillypadMenu yandexDeliveryProtocolModels.Menu, existingProductStopList map[string]bool, existingAttributeStopList map[string]bool) (coreMenuModels.Menu, error) {
	systemMenu := coreMenuModels.Menu{
		Sections: ks.getSections(tillypadMenu.Categories),
	}

	systemProducts, systemAttributesGroups, systemAttributes := ks.getEntities(tillypadMenu.Items)

	systemProducts = ks.setProductsAvailability(systemProducts, existingProductStopList)
	systemAttributes = ks.setAttributesAvailability(systemAttributes, existingAttributeStopList)

	systemMenu.Products = systemProducts
	systemMenu.Attributes = systemAttributes
	systemMenu.AttributesGroups = systemAttributesGroups

	return systemMenu, nil
}

func (ks *tillypadPosService) getStopListMapping(tillypadStopList yandexDeliveryProtocolModels.StopListResponse) (map[string]bool, map[string]bool) {
	var (
		existingProductStopList   = make(map[string]bool)
		existingAttributeStopList = make(map[string]bool)
	)

	for _, item := range tillypadStopList.Items {
		existingProductStopList[item.ItemId] = true
	}

	for _, modifier := range tillypadStopList.Modifiers {
		existingAttributeStopList[modifier.ModifierId] = true
	}

	return existingProductStopList, existingAttributeStopList
}

func (ks *tillypadPosService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	responseStopList, err := ks.tillypadCli.GetAvailability(ctx, store.TillyPad.PointId)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	existingProductStopList, existingAttributeStopList := ks.getStopListMapping(responseStopList)

	responseMenu, err := ks.tillypadCli.GetMenu(ctx, store.TillyPad.PointId)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	return ks.toSystemMenu(responseMenu, existingProductStopList, existingAttributeStopList)
}

func (ks *tillypadPosService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	return nil
}

func (ks *tillypadPosService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (ks *tillypadPosService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (ks *tillypadPosService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}

func setCustomerPhoneNumberForDeliveryInfo(storeID, deliveryService, customerPhoneNumber string) string {
	defNum := "+77771111111"
	if isYamiYami(storeID) && deliveryService == models.WOLT.String() {
		if len(customerPhoneNumber) > 0 && !strings.HasPrefix(customerPhoneNumber, "+7") {
			return defNum
		}
		return customerPhoneNumber
	}
	return defNum
}

func isYamiYami(storeId string) bool {
	yamiYamiMap := map[string]struct{}{
		"6634b5e93f89a1e89520d3f6": {},
		"6634b60bccef61bec9214907": {},
	}
	if _, ok := yamiYamiMap[storeId]; !ok {
		return false
	}
	return true
}
