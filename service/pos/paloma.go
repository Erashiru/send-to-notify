package pos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	"github.com/kwaaka-team/orders-core/pkg/paloma"
	palomaConf "github.com/kwaaka-team/orders-core/pkg/paloma/clients"
	palomaModels "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type palomaService struct {
	*BasePosService
	palomaCli           palomaConf.Paloma
	pointID             string
	isStopListByBalance bool
}

func newPalomaService(bps *BasePosService, baseUrl, apikey, class, pointID string, isStopListByBalance bool) (*palomaService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "burgerKingService constructor error")
	}

	client, err := paloma.New(&palomaConf.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		ApiKey:   apikey,
		Class:    class,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Paloma Client.")
		return nil, err
	}

	return &palomaService{
		palomaCli:           client,
		BasePosService:      bps,
		pointID:             pointID,
		isStopListByBalance: isStopListByBalance,
	}, nil
}

func (palomaSvc *palomaService) IsStopListByBalance(ctx context.Context, store coreStoreModels.Store) bool {
	return store.Paloma.StopListByBalance
}

func (palomaSvc *palomaService) GetBalanceLimit(ctx context.Context, store coreStoreModels.Store) int {
	return store.Paloma.StopListBalanceLimit
}

func (palomaSvc *palomaService) GetOrderStatus(ctx context.Context, order models.Order) (string, error) {
	orderStatus, err := palomaSvc.palomaCli.GetOrderStatus(ctx, order.PosOrderID)
	if err != nil {
		return "", err
	}

	return orderStatus.Status, nil
}

func (palomaSvc *palomaService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	menu, err := palomaSvc.palomaCli.GetMenu(ctx, store.Paloma.PointID)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	products, err := palomaSvc.existProducts(ctx, systemMenuInDb.Products)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	return palomaSvc.menuFromClient(menu, store.Settings, products), nil
}

func (palomaSvc *palomaService) convertToItemsWithBalance(ctx context.Context, menu palomaModels.Menu) coreMenuModels.StopListItems {
	stopListItems := make(coreMenuModels.StopListItems, 0, 4)

	for _, itemGroup := range menu.ItemGroups {
		for _, item := range itemGroup.Items {
			stopListItems = append(stopListItems, coreMenuModels.StopListItem{
				ProductID: strconv.Itoa(item.ObjectId),
				Balance:   item.Quantity,
			})
		}
	}

	return stopListItems
}

func (palomaSvc *palomaService) getStopListByBalance(ctx context.Context) (coreMenuModels.StopListItems, error) {
	menu, err := palomaSvc.palomaCli.GetMenu(ctx, palomaSvc.pointID)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}

	return palomaSvc.convertToItemsWithBalance(ctx, menu), nil
}

func (palomaSvc *palomaService) getStopListWithoutBalance(ctx context.Context) (coreMenuModels.StopListItems, error) {
	resp, err := palomaSvc.palomaCli.GetStopList(ctx, palomaSvc.pointID)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}

	stopLists := make(coreMenuModels.StopListItems, 0, len(resp.StopListItems))

	for _, item := range resp.StopListItems {
		stopLists = append(stopLists, coreMenuModels.StopListItem{
			ProductID: strconv.Itoa(item.ObjectId),
		})
	}

	return stopLists, nil
}

func (palomaSvc *palomaService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {
	if palomaSvc.isStopListByBalance {
		return palomaSvc.getStopListByBalance(ctx)
	}

	return palomaSvc.getStopListWithoutBalance(ctx)
}

func (palomaSvc *palomaService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {
	switch posStatus {
	case "new":
		return models.ACCEPTED, nil
	case "cooking":
		return models.COOKING_STARTED, nil
	case "on_way":
		return models.ON_WAY, nil
	case "completed":
		return models.CLOSED, nil
	case "canceled":
		return models.CANCELLED_BY_POS_SYSTEM, nil
	}
	return 0, models.StatusIsNotExist
}

func (palomaSvc *palomaService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration,
	store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	order, err := prepareAnOrder(ctx, order, store, menu, aggregatorMenu, menuClient)
	if err != nil {
		return order, err
	}

	palomaOrder, err := palomaSvc.constructPosOrder(order, store)
	if err != nil {
		return order, err
	}

	order, err = palomaSvc.SetPosRequestBodyToOrder(order, palomaOrder)
	if err != nil {
		return order, err
	}

	orderResponse, err := palomaSvc.sendOrder(ctx, palomaOrder, store)
	if err != nil {
		order.FailReason.Code = OTHER_FAIL_REASON_CODE
		order.FailReason.Message = orderResponse.Status
		return order, err
	}

	order = setPosOrderId(order, order.ID)

	order.CreationResult = models.CreationResult{
		Message: orderResponse.Status,
		OrderInfo: models.OrderInfo{
			ID:             orderResponse.OrderId,
			OrganizationID: store.Paloma.PointID,
			CreationStatus: orderResponse.Status,
		},
	}

	return order, nil
}

func (palomaSvc *palomaService) fillItems(req models.Order) ([]palomaModels.OrderItem, error) {
	var items = make([]palomaModels.OrderItem, 0, len(req.Products))

	for _, product := range req.Products {

		productID, err := strconv.Atoi(product.ID)
		if err != nil {
			return nil, err
		}

		item := palomaModels.OrderItem{
			ObjectId: productID,
			Name:     product.Name,
			Price:    int(product.Price.Value),
			Count:    product.Quantity,
		}

		var itemPriceForCombo int

		for _, attribute := range product.Attributes {

			attributeID, err := strconv.Atoi(attribute.ID)
			if err != nil {
				return nil, err
			}

			switch product.IsCombo {
			case true:
				item.ComplexItems = append(item.ComplexItems, palomaModels.OrderComplexItem{
					ObjectId: attributeID,
					Name:     attribute.Name,
					Count:    attribute.Quantity * product.Quantity,
					Price:    int(attribute.Price.Value),
				})
				itemPriceForCombo += int(attribute.Price.Value) * attribute.Quantity
			default:
				item.Modifications = append(item.Modifications, palomaModels.OrderModification{
					ObjectId: attributeID,
					Name:     attribute.Name,
					Count:    attribute.Quantity * product.Quantity,
					Price:    int(attribute.Price.Value),
				})
			}
		}

		if product.IsCombo {
			item.Price = item.Price + itemPriceForCombo
		}

		items = append(items, item)
	}

	return items, nil
}

func (palomaSvc *palomaService) setDeliveryType(req models.Order, palomaOrder palomaModels.Order) palomaModels.Order {
	switch req.IsPickedUpByCustomer {
	case false:
		palomaOrder.DeliveryType = palomaDeliveryTypeCourier
	default:
		palomaOrder.DeliveryType = palomaDeliveryTypeCustomer
	}

	if req.SendCourier {
		palomaOrder.DeliveryType = palomaDeliveryTypeCourier
	}

	return palomaOrder
}

func (palomaSvc *palomaService) setPaymentType(req models.Order, palomaOrder palomaModels.Order) palomaModels.Order {
	switch req.PaymentMethod {
	case "CASH":
		palomaOrder.IsCash = true
		palomaOrder.IsPayed = false
	case "DELAYED":
		palomaOrder.IsCash = false
		palomaOrder.IsPayed = true
	default:
		palomaOrder.IsCash = false
		palomaOrder.IsPayed = true
	}

	return palomaOrder
}

func (palomaSvc *palomaService) setNameAndPhoneNumber(req models.Order, palomaOrder palomaModels.Order) palomaModels.Order {
	if req.IsMarketplace {
		switch req.DeliveryService {
		case models.GLOVO.String():
			palomaOrder.Name = "Glovo"
			palomaOrder.Phone = "+77777777771"
		case models.WOLT.String():
			palomaOrder.Name = "Wolt"
			palomaOrder.Phone = "+77777777772"
		case models.YANDEX.String():
			palomaOrder.Name = "Yandex"
			palomaOrder.Phone = "+77777777773"
		case models.EMENU.String():
			palomaOrder.Name = "Emenu"
			palomaOrder.Phone = "+77777777774"
		case models.EXPRESS24.String():
			palomaOrder.Name = "Express24"
			palomaOrder.Phone = "+77777777775"
		}
	}

	return palomaOrder
}

func (palomaSvc *palomaService) constructOrderComments(req models.Order, store coreStoreModels.Store) string {
	var (
		addressName        = "Адрес"
		orderCodeName      = "Код заказа"
		paymentCashName    = "Наличный"
		paymentDelayedName = "Безналичный"
		commentName        = "Комментарий"
		deliveryName       = "Доставка"
		allergyName        = "Аллергия"
		paymentTypeName    = "Тип оплаты"
		pickUpToName       = "Приготовить к"
		quantityPerson     = "Количество людей"
	)

	commentSettings := store.Settings.CommentSetting

	if commentSettings.HasCommentSetting {
		addressName = commentSettings.AddressName
		orderCodeName = commentSettings.OrderCodeName
		paymentCashName = commentSettings.CashPaymentName
		paymentDelayedName = commentSettings.DelayedPaymentName
		commentName = commentSettings.CommentName
		deliveryName = commentSettings.DeliveryName
		allergyName = commentSettings.Allergy
		paymentTypeName = commentSettings.PaymentTypeName
		pickUpToName = commentSettings.PickUpToName
		quantityPerson = commentSettings.QuantityPerson
	}

	comment := fmt.Sprintf("%s: %s\n", orderCodeName, req.PickUpCode)

	deliveryService := req.DeliveryService
	if deliveryService == models.QRMENU.String() {
		deliveryService = "Kwaaka"
	}

	comment += fmt.Sprintf("%s: %s\n", deliveryName, deliveryService)

	switch req.PaymentMethod {
	case models.PAYMENT_METHOD_CASH:
		comment += fmt.Sprintf("%s: %s\n", paymentTypeName, paymentCashName)
	case models.PAYMENT_METHOD_DELAYED:
		comment += fmt.Sprintf("%s: %s\n", paymentTypeName, paymentDelayedName)
	}

	if !req.IsMarketplace && req.DeliveryAddress.Label != "" {
		comment += fmt.Sprintf("%s: %s\n", addressName, req.DeliveryAddress.Label)
	}

	if req.AllergyInfo != "" {
		comment += fmt.Sprintf("%s: %s\n", allergyName, req.AllergyInfo)
	}

	if req.SpecialRequirements != "" {
		comment += fmt.Sprintf("%s: %s\n", commentName, req.SpecialRequirements)
	}

	if req.Persons != 0 {
		comment += fmt.Sprintf("%s: %d\n", quantityPerson, req.Persons)
	}

	if !req.EstimatedPickupTime.Value.IsZero() {
		comment += fmt.Sprintf("%s: %s\n", pickUpToName, req.EstimatedPickupTime.Value.Time.
			Add(time.Duration(store.Settings.TimeZone.UTCOffset)*time.Hour).
			Format("15:04:05"))
	}

	return comment
}

func (palomaSvc *palomaService) constructPosOrder(req models.Order, store coreStoreModels.Store) (palomaModels.Order, error) {
	order := palomaModels.Order{
		OrderId:        req.ID,
		Date:           req.EstimatedPickupTime.Value.Time.Add(time.Duration(store.Settings.TimeZone.UTCOffset) * time.Hour).String(),
		Name:           req.Customer.Name,
		Phone:          req.Customer.PhoneNumber,
		Address:        req.DeliveryAddress.Label,
		CoordinateLong: strconv.Itoa(int(req.DeliveryAddress.Longitude)),
		CoordinateLat:  strconv.Itoa(int(req.DeliveryAddress.Latitude)),
		Comment:        palomaSvc.constructOrderComments(req, store),
		PersonAmount:   req.Persons,
		TotalPrice:     int(req.EstimatedTotalPrice.Value),
		DiscountAmount: 0,
	}

	order = palomaSvc.setNameAndPhoneNumber(req, order)

	order = palomaSvc.setPaymentType(req, order)

	order = palomaSvc.setDeliveryType(req, order)

	items, err := palomaSvc.fillItems(req)
	if err != nil {
		return palomaModels.Order{}, err
	}

	order.OrderItems = items

	return order, nil
}

func (palomaSvc *palomaService) sendOrder(ctx context.Context, order any, store coreStoreModels.Store) (palomaModels.OrderResponse, error) {
	//var errs custom.Error
	posOrder, ok := order.(palomaModels.Order)

	if !ok {
		return palomaModels.OrderResponse{}, validator.ErrCastingPos
	}

	utils.Beautify("paloma request body", posOrder)

	createResponse, err := palomaSvc.palomaCli.CreateOrder(ctx, store.Paloma.PointID, posOrder)
	if err != nil {
		log.Err(err).Msg("paloma create order error in palomaService")
		//errs.Append(err, validator.ErrIgnoringPos)
		return palomaModels.OrderResponse{}, err
	}

	return createResponse, nil
}

func (palomaSvc *palomaService) menuFromClient(req palomaModels.Menu, settings coreStoreModels.Settings, productsExist map[string]coreMenuModels.Product) coreMenuModels.Menu {

	menu := coreMenuModels.Menu{
		Name:        coreMenuModels.PALOMA.String(),
		ExtName:     coreMenuModels.MAIN.String(),
		Description: "paloma pos menu",
		CreatedAt:   models.TimeNow(),
		UpdatedAt:   models.TimeNow(),
	}

	products, groups, attributeGroups, attributes := palomaSvc.toEntities(req, settings, productsExist)

	menu.Products = products
	menu.Groups = groups
	menu.AttributesGroups = attributeGroups
	menu.Attributes = attributes

	return menu
}

func (palomaSvc *palomaService) toEntities(req palomaModels.Menu, settings coreStoreModels.Settings, productsExist map[string]coreMenuModels.Product) ([]coreMenuModels.Product, []coreMenuModels.Group, []coreMenuModels.AttributeGroup, []coreMenuModels.Attribute) {
	products := make([]coreMenuModels.Product, 0, 10)
	groups := make([]coreMenuModels.Group, 0, len(req.ItemGroups))
	attributeGroups := make([]coreMenuModels.AttributeGroup, 0, 4)
	attributes := make([]coreMenuModels.Attribute, 0, 4)

	var mapAttributeGroup = make(map[string]struct{})
	var mapAttribute = make(map[string]struct{})

	for _, itemGroup := range req.ItemGroups {
		for _, item := range itemGroup.Items {
			product, attributeGroups_, attributes_ := palomaSvc.toEntity(item, settings, productsExist)

			// unique attribute groups
			for _, attributeGroup := range attributeGroups_ {
				if _, ok := mapAttributeGroup[attributeGroup.ExtID]; ok {
					continue
				}

				mapAttributeGroup[attributeGroup.ExtID] = struct{}{}
				attributeGroups = append(attributeGroups, attributeGroup)
			}

			// unique attributes
			for _, attribute := range attributes_ {
				if _, ok := mapAttribute[attribute.ExtID]; ok {
					continue
				}

				mapAttribute[attribute.ExtID] = struct{}{}
				attributes = append(attributes, attribute)
			}

			product.ParentGroupID = strconv.Itoa(itemGroup.ObjectId)
			products = append(products, product)
		}

		groups = append(groups, coreMenuModels.Group{
			ID:     strconv.Itoa(itemGroup.ObjectId),
			Name:   itemGroup.Name,
			Images: []string{itemGroup.Image},
			InMenu: true,
		})
	}

	return products, groups, attributeGroups, attributes
}

func (palomaSvc *palomaService) toEntity(req palomaModels.Item, settings coreStoreModels.Settings, productsExist map[string]coreMenuModels.Product) (coreMenuModels.Product, []coreMenuModels.AttributeGroup, []coreMenuModels.Attribute) {
	extID := uuid.New().String()

	posProduct, ok := productsExist[strconv.Itoa(req.ObjectId)]
	if ok {
		extID = posProduct.ExtID
	}

	product := coreMenuModels.Product{
		ExtID:     extID,
		ProductID: strconv.Itoa(req.ObjectId),
		Name: []coreMenuModels.LanguageDescription{
			{
				Value:        req.Name,
				LanguageCode: settings.LanguageCode,
			},
		},
		Description: []coreMenuModels.LanguageDescription{
			{
				Value:        req.Description,
				LanguageCode: settings.LanguageCode,
			},
		},
		Price: []coreMenuModels.Price{
			{
				Value:        req.Price,
				CurrencyCode: settings.Currency,
			},
		},
		ImageURLs:   []string{req.Image},
		IsAvailable: true,
	}

	if req.IUseInMenu == palomaIntTrue {
		product.IsIncludedInMenu = true
	}

	if req.MarkDeleted == palomaIntTrue {
		product.IsDeleted = true
	}

	for _, defaultAttribute := range posProduct.MenuDefaultAttributes {
		if defaultAttribute.ByAdmin {
			product.MenuDefaultAttributes = append(product.MenuDefaultAttributes, defaultAttribute)
		}
	}

	if len(req.ComplexGroups) > 0 {
		product.IsCombo = true
	}

	var attributeGroups = make([]coreMenuModels.AttributeGroup, 0, len(req.ModifierGroups))
	var attributes = make([]coreMenuModels.Attribute, 0, 4)

	for _, complexGroup := range req.ComplexGroups {
		attributeGroup := coreMenuModels.AttributeGroup{
			ExtID: strconv.Itoa(complexGroup.ObjectId),
			Name:  complexGroup.Name,
			Min:   complexGroup.MinCount,
			Max:   complexGroup.MaxCount,
		}

		product.AttributesGroups = append(product.AttributesGroups, strconv.Itoa(complexGroup.ObjectId))

		for _, modifier := range complexGroup.ComplexItems {
			attribute := coreMenuModels.Attribute{
				ExtID:                strconv.Itoa(modifier.ObjectId),
				Name:                 modifier.Name,
				Price:                modifier.Price,
				IsAvailable:          true,
				ParentAttributeGroup: strconv.Itoa(complexGroup.ObjectId),
			}

			if modifier.IUseInMenu == "1" {
				attribute.IncludedInMenu = true
			}

			if modifier.MarkDeleted == "1" {
				attribute.IsDeleted = true
			}

			attributes = append(attributes, attribute)
			attributeGroup.Attributes = append(attributeGroup.Attributes, strconv.Itoa(modifier.ObjectId))
		}

		attributeGroups = append(attributeGroups, attributeGroup)
	}

	for _, modifierGroup := range req.ModifierGroups {
		attributeGroup := coreMenuModels.AttributeGroup{
			ExtID: strconv.Itoa(modifierGroup.ObjectId),
			Name:  modifierGroup.Name,
		}

		product.AttributesGroups = append(product.AttributesGroups, strconv.Itoa(modifierGroup.ObjectId))

		for _, modifier := range modifierGroup.Modifiers {
			attribute := coreMenuModels.Attribute{
				ExtID:                strconv.Itoa(modifier.ObjectId),
				Name:                 modifier.Name,
				Price:                modifier.Price,
				IsAvailable:          true,
				ParentAttributeGroup: strconv.Itoa(modifierGroup.ObjectId),
			}

			if modifier.IUseInMenu == palomaIntTrue {
				attribute.IncludedInMenu = true
			}

			if modifier.MarkDeleted == palomaIntTrue {
				attribute.IsDeleted = true
			}

			attributes = append(attributes, attribute)
			attributeGroup.Attributes = append(attributeGroup.Attributes, strconv.Itoa(modifier.ObjectId))
		}

		attributeGroups = append(attributeGroups, attributeGroup)
	}

	return product, attributeGroups, attributes
}

func (palomaSvc *palomaService) existProducts(ctx context.Context, products []coreMenuModels.Product) (map[string]coreMenuModels.Product, error) {
	// add to hash map
	productExist := make(map[string]coreMenuModels.Product, len(products))
	for _, product := range products {
		// cause has cases if product_id && parent_id same, size_id different
		productExist[product.ProductID] = product
	}

	return productExist, nil
}

func (palomaSvc *palomaService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	_, err := palomaSvc.palomaCli.CancelOrder(ctx, order.PosOrderID) //TODO: unused function before
	if err != nil {
		return err
	}

	return nil
}

func (s *palomaService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (s *palomaService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (s *palomaService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
