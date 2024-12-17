package pos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/custom"
	selector2 "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	rkeeperClient "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite"
	rkeeperConf "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients"
	rkeeperDto "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/dto"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type rkeeperService struct {
	*BasePosService
	rkeeperCli      rkeeperConf.RKeeper
	rKeeperObjectID int
}

func newRkeeperService(bps *BasePosService, rKeeperObjectID int, storeToken, apiKey, baseUrl string) (*rkeeperService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "burgerKingService constructor error")
	}

	if storeToken != "" {
		apiKey = storeToken
	}

	client, err := rkeeperClient.NewRKeeperClient(&rkeeperConf.Config{
		Protocol: "http",
		ApiKey:   apiKey,
		BaseURL:  baseUrl,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize RKeeper Client.")
		return nil, err
	}

	return &rkeeperService{
		rKeeperObjectID: rKeeperObjectID,
		rkeeperCli:      client,
		BasePosService:  bps}, nil
}

func (rkeeperSvc *rkeeperService) GetOrderStatus(ctx context.Context, order models.Order) (string, error) {
	orderTask, err := rkeeperSvc.rkeeperCli.GetOrder(ctx, order.PosOrderID, rkeeperSvc.rKeeperObjectID)
	if err != nil {
		return "", err
	}

	var status string

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		log.Info().Msgf("retry get order status attempt %d", i)

		response, err := rkeeperSvc.rkeeperCli.GetOrderTask(ctx, orderTask.ResponseCommon.TaskGUID)
		if err != nil {
			continue
		}

		if response.TaskResponse.OrderResponseBody.Status.Value == "" {
			continue
		}

		status = response.TaskResponse.OrderResponseBody.Status.Value
		return status, nil
	}

	return "", errors.New("unknown status")
}

func (rkeeperSvc *rkeeperService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	if err := rkeeperSvc.UpdateMenu(ctx, store); err != nil {
		return coreMenuModels.Menu{}, err
	}

	menu, err := rkeeperSvc.rkeeperCli.GetMenu(ctx, store.RKeeper.ObjectId)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	// exist products in DB
	products, posProducts, err := rkeeperSvc.existProducts(ctx, systemMenuInDb.Products)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	stopList, err := rkeeperSvc.rkeeperCli.GetStopList(ctx, store.RKeeper.ObjectId)
	if err != nil {
		return coreMenuModels.Menu{}, fmt.Errorf("couldn't get stoplist %w", err)
	}

	return rkeeperSvc.menuFromClient(menu.TaskResponse.Menu, products, store.Settings, posProducts, stopList), nil
}

func (rkeeperSvc *rkeeperService) UpdateMenu(ctx context.Context, store coreStoreModels.Store) error {
	_, err := rkeeperSvc.rkeeperCli.UpdateMenu(ctx, store.RKeeper.ObjectId)
	if err != nil {
		return fmt.Errorf("couldn't update menu %w", err)
	}

	log.Info().Msgf("success update menu for rkeeper object id %d", store.RKeeper.ObjectId)

	return nil
}

func (rkeeperSvc *rkeeperService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {

	switch posStatus {
	case "WAIT_SENDING":
		return models.WAIT_SENDING, nil
	case "Canceled":
		return models.CANCELLED_BY_POS_SYSTEM, nil
	case "Created":
		return models.ACCEPTED, nil
	case "Cooking":
		return models.COOKING_STARTED, nil
	case "Ready":
		return models.COOKING_COMPLETE, nil
	case "Completed", "Complited", "IssuedOut":
		return models.CLOSED, nil
	default:
		return 0, models.StatusIsNotExist
	}
}

func (rkeeperSvc *rkeeperService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration,
	store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	order, err := prepareAnOrder(ctx, order, store, menu, aggregatorMenu, menuClient)
	if err != nil {
		return order, err
	}

	posOrderBody, err := rkeeperSvc.constructPosOrder(order, store)

	if err != nil {
		return order, err
	}

	posOrder := rkeeperDto.CreateOrderRequest{
		Params: rkeeperDto.CreateOrderParam{
			Async: rkeeperDto.Sync{
				ObjectID: store.RKeeper.ObjectId,
				Timeout:  120,
			},
			Order: posOrderBody,
		},
	}

	order, err = rkeeperSvc.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	errorSolutions, err := errSolution.GetAllErrorSolutions(ctx)
	if err != nil {
		log.Err(err).Msg("rkeeper7XMLService error: GetAllErrorSolutions")
		return models.Order{}, err
	}

	syncResponse, err := rkeeperSvc.sendOrder(ctx, posOrder)
	if err != nil {
		errMessage := syncResponse.Error.WsError.Desc + syncResponse.Error.AgentError.Desc + err.Error()
		failReason, _, failReasonErr := errSolution.SetFailReason(ctx, store, err.Error(), MatchingCodes(errMessage, errorSolutions), "")
		if failReasonErr != nil {
			return order, failReasonErr
		}
		order.FailReason = failReason

		return order, errors.New(syncResponse.Error.WsError.Desc + syncResponse.Error.AgentError.Desc)
	}

	posOrderID, err := rkeeperSvc.CreateOrderTask(ctx, syncResponse.ResponseCommon.TaskGUID)
	if err != nil {
		return order, err
	}

	var payOrder rkeeperDto.SyncResponse
	if posOrderID != "" && store.ExternalPosIntegrationSettings.PayOrderIsOn {
		payOrder, err = rkeeperSvc.payOrder(ctx, store.RKeeper.ObjectId, int(order.EstimatedTotalPrice.Value), posOrderID, order.PosPaymentInfo.PaymentTypeID)
		if err != nil {
			return models.Order{}, err
		}
	}

	slog.Info("rkeeperService send order result ", "order_id", order.OrderID, "syncResponse", syncResponse)

	order = setPosOrderId(order, posOrderID)

	storeGroup, err := storeCli.FindStoreGroup(ctx, selector2.StoreGroup{
		ID: store.RestaurantGroupID,
	})
	if err != nil {
		return models.Order{}, err
	}

	if order.PosOrderID != "" || payOrder.Error.WsError.Desc != "" || payOrder.Error.WsError.Code != "" || payOrder.Error.AgentError.Desc != "" || payOrder.Error.AgentError.Code != "" {
		if order.RetryCount <= storeGroup.RetryCount {
			slog.Info("rkeeperService send sqs message to queue name:order-retry", "orderGuid", payOrder.TaskResponse.Order.OrderGuid)
			if retryErr := notifyQueue.SendSQSMessage(ctx, "order-retry", payOrder.TaskResponse.Order.OrderGuid); retryErr != nil {
				slog.Error("rkeeperService: notifyQueue.SendSQSMessage", "error", retryErr.Error())
			}
		}
	}

	order.CreationResult = models.CreationResult{
		Message: syncResponse.TaskResponse.Status,
		OrderInfo: models.OrderInfo{
			ID:             syncResponse.ResponseCommon.TaskGUID,
			OrganizationID: strconv.Itoa(syncResponse.ResponseCommon.ObjectID),
			CreationStatus: syncResponse.TaskResponse.Status,
		},
		ErrorDescription: syncResponse.Error.WsError.Desc + " " + syncResponse.Error.AgentError.Desc,
	}

	if syncResponse.Error.WsError.Desc != "" || syncResponse.Error.AgentError.Desc != "" {
		order.Status = "FAILED"
		return order, err
	}

	return order, nil
}

func (rkeeperSvc *rkeeperService) payOrder(ctx context.Context, objectID, amount int, orderId, currency string) (rkeeperDto.SyncResponse, error) {
	var (
		payOrder rkeeperDto.SyncResponse
		err      error
	)

	for i := 0; i < 4; i++ {

		payOrder, err = rkeeperSvc.rkeeperCli.PayOrder(ctx, objectID, amount, orderId, currency)
		if err != nil {
			continue
		}
		if payOrder.Error.WsError.Code != "" {
			continue
		}
		timeout := 60
		// for loop every 10 seconds until timeout passed
		for j := 0; j <= timeout; j += 10 {
			time.Sleep(10 * time.Second)
			task, err := rkeeperSvc.rkeeperCli.CreateOrderTask(ctx, payOrder.ResponseCommon.TaskGUID)
			if err != nil || task.Error.AgentError.Code != "" || task.Error.WsError.Code != "" {
				log.Info().Msgf("pay order error: %+v", task)
				continue
			}
			log.Info().Msgf("pay order success: %+v", payOrder)
			break
		}

	}

	if err != nil {
		return payOrder, err
	}

	if payOrder.Error.WsError.Code != "" || payOrder.Error.AgentError.Code != "" {
		return payOrder, fmt.Errorf("pay order error: %+v", payOrder.Error)
	}

	return payOrder, nil
}

func (rkeeperSvc *rkeeperService) sendOrder(ctx context.Context, posOrder rkeeperDto.CreateOrderRequest) (rkeeperDto.SyncResponse, error) {
	var errs custom.Error

	log.Info().Msgf("rkeeper Request Body: %+v", posOrder)

	createResponse, err := rkeeperSvc.rkeeperCli.CreateOrder(ctx, posOrder.Params.Async.ObjectID, posOrder.Params.Order)
	if err != nil {
		log.Err(err).Msg("rkeeper create order error")
		errs.Append(err, validator.ErrIgnoringPos)
		return rkeeperDto.SyncResponse{}, errs
	}

	return createResponse, nil
}

func (rkeeperSvc *rkeeperService) fillTaker(isPickUpByCustomer bool) string {
	if isPickUpByCustomer {
		return "customer"
	}

	return "courier"
}

func (rkeeperSvc *rkeeperService) fillExpiditionType(isPickUpByCustomer bool) string {
	if isPickUpByCustomer {
		return "pickup"
	}

	return "delivery"
}

func (rkeeperSvc *rkeeperService) fillPayment(order models.Order, store coreStoreModels.Store, rkeeperOrder rkeeperDto.Order) rkeeperDto.Order {
	if !store.ExternalPosIntegrationSettings.PrePayOrderIsOn {
		rkeeperOrder.Price = &rkeeperDto.Price{
			Total: int(order.EstimatedTotalPrice.Value),
		}
	} else {
		rkeeperOrder.PrePayments = []rkeeperDto.PrePayment{
			{
				Amount:   int(order.EstimatedTotalPrice.Value),
				Currency: order.PosPaymentInfo.PaymentTypeID,
			},
		}
	}

	return rkeeperOrder
}

func (rkeeperSvc *rkeeperService) fillOrderTypeCode(orderType string, rkeeperOrder rkeeperDto.Order) rkeeperDto.Order {
	if orderType != "" {
		orderType, err := strconv.Atoi(orderType)
		if err != nil {
			log.Err(err).Msgf("order type is incorrect to convert to int")
		} else {
			rkeeperOrder.OrderTypeCode = orderType
		}
	}

	return rkeeperOrder
}

func (rkeeperSvc *rkeeperService) fillPaymentType(paymentMethod string, rkeeperOrder rkeeperDto.Order) rkeeperDto.Order {
	//switch paymentMethod {
	//case "CASH":
	//	rkeeperOrder.Payment.Type = "cash"
	//case "DELAYED":
	//	rkeeperOrder.Payment.Type = "online"
	//}

	return rkeeperOrder
}

func (rkeeperSvc *rkeeperService) fillProducts(order models.Order, rkeeperOrder rkeeperDto.Order) rkeeperDto.Order {
	var products = make([]rkeeperDto.CreateOrderProduct, 0)
	for _, product := range order.Products {
		var ingredients = make([]rkeeperDto.CreateOrderIngredient, 0, len(product.Attributes))

		for _, attribute := range product.Attributes {
			ingredients = append(ingredients, rkeeperDto.CreateOrderIngredient{
				Id:       attribute.ID,
				Name:     attribute.Name,
				Quantity: attribute.Quantity,
			})
		}

		products = append(products, rkeeperDto.CreateOrderProduct{
			Id:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			Ingredients: ingredients,
		})
	}

	rkeeperOrder.Products = products

	return rkeeperOrder
}

func (rkeeperSvc *rkeeperService) constructPosOrder(order models.Order, store coreStoreModels.Store) (rkeeperDto.Order, error) {
	orderComment, _, _ := rkeeperSvc.constructOrderComments(order, store)

	rkeeperOrder := rkeeperDto.Order{
		OriginalOrderId: order.OrderID,
		Customer: &rkeeperDto.PersonInfo{
			Name:  order.Customer.Name,
			Phone: order.Customer.PhoneNumber,
		},
		ExpeditionType:  rkeeperSvc.fillExpiditionType(order.IsPickedUpByCustomer),
		Comment:         strings.ToTitle(order.DeliveryService) + "\n" + orderComment,
		PersonsQuantity: order.Persons,
	}

	rkeeperOrder = rkeeperSvc.fillPickupDelivery(order, rkeeperOrder.ExpeditionType, store, rkeeperOrder)

	rkeeperOrder = rkeeperSvc.fillPayment(order, store, rkeeperOrder)

	rkeeperOrder = rkeeperSvc.fillOrderTypeCode(order.PosPaymentInfo.OrderType, rkeeperOrder)

	rkeeperOrder = rkeeperSvc.fillPaymentType(order.PaymentMethod, rkeeperOrder)

	rkeeperOrder = rkeeperSvc.fillProducts(order, rkeeperOrder)

	return rkeeperOrder, nil
}

func (rkeeperSvc *rkeeperService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {

	resp, err := rkeeperSvc.rkeeperCli.GetStopList(ctx, rkeeperSvc.rKeeperObjectID)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}

	stopLists := make(coreMenuModels.StopListItems, 0, len(resp.TaskResponse.StopList.Dishes))

	for _, dish := range resp.TaskResponse.StopList.Dishes {
		stopLists = append(stopLists, coreMenuModels.StopListItem{
			ProductID: dish.ID,
		})
	}

	return stopLists, nil
}

func (rkeeperSvc *rkeeperService) menuFromClient(req rkeeperDto.Menu, productsExist map[string]string, settings coreStoreModels.Settings, posProducts map[string]coreMenuModels.Product, stopList rkeeperDto.StopListResponse) coreMenuModels.Menu {

	menu := coreMenuModels.Menu{
		Name:             coreMenuModels.RKEEPER.String(),
		ExtName:          coreMenuModels.MAIN.String(),
		Description:      "rkeeper pos menu",
		AttributesGroups: rkeeperSvc.attributeGroupsToModel(req),
		Products:         rkeeperSvc.productsToModel(req, productsExist, settings),
		Groups:           rkeeperSvc.groupsToModel(req.Categories),
		CreatedAt:        models.TimeNow(),
		UpdatedAt:        models.TimeNow(),
	}

	menu.Attributes = rkeeperSvc.attributesToModel(req.Ingredients, menu.AttributesGroups)
	sections, collections, supercollections := rkeeperSvc.getCollections(req)
	menu.Sections = sections
	menu.Collections = collections
	menu.SuperCollections = supercollections

	menu.Products = append(menu.Products, rkeeperSvc.existProductsByStopList(posProducts, stopList)...)
	menu.Products = menu.Products.Unique()
	return menu
}

func (rkeeperSvc *rkeeperService) existProductsByStopList(posProducts map[string]coreMenuModels.Product, stopList rkeeperDto.StopListResponse) coreMenuModels.Products {
	products := make(coreMenuModels.Products, 0, len(stopList.TaskResponse.StopList.Dishes))

	for _, dish := range stopList.TaskResponse.StopList.Dishes {
		if val, ok := posProducts[dish.ID]; ok {
			val.IsAvailable = false
			products = append(products, val)
		}
	}

	return products
}

func (rkeeperSvc *rkeeperService) attributesToModel(ingredients rkeeperDto.Ingredients, attributeGroups coreMenuModels.AttributeGroups) coreMenuModels.Attributes {

	attrGroupsExist := make(map[string]coreMenuModels.AttributeGroup, len(attributeGroups))

	for _, attrGroup := range attributeGroups {
		for _, ingredient := range attrGroup.Attributes {
			attrGroupsExist[ingredient] = attrGroup
		}
	}

	res := make(coreMenuModels.Attributes, 0, len(ingredients))

	for _, attribute := range ingredients {
		attr, err := rkeeperSvc.attributeToModel(attribute)
		if err != nil {
			log.Err(err).Msg("rkeeper cli err: get attribute")
			continue
		}

		if attrGroup, ok := attrGroupsExist[attribute.ID]; ok {
			attr.ParentAttributeGroup = attrGroup.ExtID
			attr.HasAttributeGroup = true
			attr.AttributeGroupName = attrGroup.Name
			attr.AttributeGroupMin = attrGroup.Min
			attr.AttributeGroupMax = attrGroup.Max
		}

		res = append(res, attr)
	}

	return res
}

func (rkeeperSvc *rkeeperService) attributeToModel(attributes rkeeperDto.Ingredient) (coreMenuModels.Attribute, error) {

	res := coreMenuModels.Attribute{
		ExtID:       attributes.ID,
		Name:        attributes.Name,
		ExtName:     attributes.Description,
		IsAvailable: true,
		UpdatedAt:   time.Now(),
	}

	price, err := strconv.ParseFloat(attributes.Price, 64)
	if err != nil {
		return coreMenuModels.Attribute{}, err
	}

	res.Price = price

	return res, nil
}

func (rkeeperSvc *rkeeperService) attributeGroupsToModel(req rkeeperDto.Menu) coreMenuModels.AttributeGroups {

	groupsInSchemes := make(map[string]rkeeperDto.IngredientsSchemeGroup, len(req.IngredientsSchemes))
	for _, scheme := range req.IngredientsSchemes {
		if scheme.IngredientsGroups != nil {
			for _, attributeGroups := range scheme.IngredientsGroups {
				groupsInSchemes[attributeGroups.ID] = attributeGroups
			}
		}
	}

	res := make(coreMenuModels.AttributeGroups, 0, len(req.IngredientsGroups))
	for _, ingredientGroup := range req.IngredientsGroups {

		attributeGroup := rkeeperSvc.attributeGroupToModel(ingredientGroup)

		if data, ok := groupsInSchemes[attributeGroup.ExtID]; ok {
			attributeGroup.Min = data.MinCount
			attributeGroup.Max = data.MaxCount
		}

		res = append(res, attributeGroup)
	}
	return res
}

func (rkeeperSvc *rkeeperService) attributeGroupToModel(req rkeeperDto.IngredientsGroup) coreMenuModels.AttributeGroup {
	return coreMenuModels.AttributeGroup{
		ExtID:      req.ID,
		Name:       req.Name,
		Attributes: req.Ingredients,
	}
}

func (rkeeperSvc *rkeeperService) groupsToModel(req rkeeperDto.Categories) coreMenuModels.Groups {

	res := make(coreMenuModels.Groups, 0, len(req))
	for _, group := range req {
		res = append(res, rkeeperSvc.groupToModel(group))
	}
	return res
}

func (rkeeperSvc *rkeeperService) groupToModel(req rkeeperDto.Category) coreMenuModels.Group {
	return coreMenuModels.Group{
		ID:          req.ID,
		Name:        req.Name,
		ParentGroup: req.ParentId,
	}
}

func (rkeeperSvc *rkeeperService) getSections(req rkeeperDto.Products) (coreMenuModels.Sections, map[string]struct{}) {

	sectionExist := make(map[string]struct{}, len(req))
	sections := make(coreMenuModels.Sections, 0, len(req))

	for _, product := range req {
		if product.SchemeId == "" {
			continue
		}
		sectionExist[product.SchemeId] = struct{}{}
		sections = append(sections, coreMenuModels.Section{
			ExtID: product.ID,
		})

	}

	return sections, sectionExist
}

func (rkeeperSvc *rkeeperService) getCollections(req rkeeperDto.Menu) (coreMenuModels.Sections, coreMenuModels.MenuCollections, coreMenuModels.MenuSuperCollections) {

	superCollections := make(coreMenuModels.MenuSuperCollections, 0, len(req.Categories))
	collections := make(coreMenuModels.MenuCollections, 0, len(req.Categories))

	sections, sectionExist := rkeeperSvc.getSections(req.Products)

	categoryExist := make(map[string]rkeeperDto.Category, len(rkeeperDto.Categories{}))
	for _, category := range req.Categories {
		categoryExist[category.ID] = category
	}

	for _, category := range req.Categories {

		// third level is section
		if _, ok := sectionExist[category.ID]; ok {

			if category.ParentId != "" {

				// second level is collection
				if parent, ok := categoryExist[category.ParentId]; ok {

					collections = append(collections, coreMenuModels.MenuCollection{
						ExtID: parent.ID,
						Name:  parent.Name,
					})

					if parent.ParentId != "" {
						if proParent, ok := categoryExist[parent.ParentId]; ok {
							superCollections = append(superCollections, coreMenuModels.MenuSuperCollection{
								ExtID: proParent.ID,
								Name:  proParent.Name,
							})
						}
					}

				}
			}

			if i, ok := sections.GetIndex(category.ID); ok {
				sections[i] = coreMenuModels.Section{
					ExtID:      category.ID,
					Name:       category.Name,
					Collection: category.ParentId,
				}
			}
		}

		if category.ParentId == "" {
			superCollections = append(superCollections, coreMenuModels.MenuSuperCollection{
				ExtID: category.ID,
				Name:  category.Name,
			})
			continue
		}

	}

	return sections, collections, superCollections
}

func (rkeeperSvc *rkeeperService) productsToModel(req rkeeperDto.Menu, productExist map[string]string, settings coreStoreModels.Settings) coreMenuModels.Products {

	results := make(coreMenuModels.Products, 0, len(req.Products))

	for _, product := range req.Products {

		res, err := rkeeperSvc.productToModel(req, product, productExist, settings)
		if err != nil {
			log.Err(err).Msgf("rkeeper cli err: get product %s", product.ID)
			continue
		}

		results = append(results, res)
	}

	return results
}

func (rkeeperSvc *rkeeperService) productToModel(
	req rkeeperDto.Menu,
	product rkeeperDto.Product,
	productExist map[string]string,
	setting coreStoreModels.Settings) (coreMenuModels.Product, error) {

	res := coreMenuModels.Product{
		ProductID:     product.ID,
		ExtID:         uuid.NewString(),
		ParentGroupID: product.CategoryID,
		Section:       product.CategoryID,
		IsAvailable:   true,
		ExtName:       product.Name,
		Name: []coreMenuModels.LanguageDescription{
			{
				Value: product.Name,
			},
		},
		Description: []coreMenuModels.LanguageDescription{
			{
				Value: product.Description,
			},
		},
		ImageURLs:   product.ImageUrls,
		MeasureUnit: product.Measure.Unit,
		ProductsCreatedAt: coreMenuModels.ProductsCreatedAt{
			Value:     models.TimeNow(),
			Timezone:  setting.TimeZone.TZ,
			UTCOffset: setting.TimeZone.UTCOffset,
		},
		UpdatedAt: models.TimeNow(),
	}

	price, err := strconv.ParseFloat(product.Price, 64)
	if err != nil {
		return coreMenuModels.Product{}, err
	}

	res.Price = []coreMenuModels.Price{
		{
			Value: price,
		},
	}

	if product.Measure.Value == "" {
		res.Weight = 0
	}

	measureValue, err := strconv.ParseFloat(product.Measure.Value, 64)
	if err != nil {
		log.Err(err).Msgf("product %s measure value %s", product.ID, product.Measure)
		measureValue = 0
	}

	res.Weight = measureValue

	// TODO: check this point
	// checking ext_id has in db
	key := strings.TrimSpace(res.ProductID + res.ParentGroupID)
	extID, ok := productExist[key]
	if ok {
		res.ExtID = extID
	}

	if product.SchemeId != "" {
		attributeGroups, attributes := rkeeperSvc.collectAttributes(product.SchemeId, req)
		res.AttributesGroups = attributeGroups
		res.Attributes = attributes
	}

	return res, nil
}

func (rkeeperSvc *rkeeperService) existProducts(ctx context.Context, products []coreMenuModels.Product) (map[string]string, map[string]coreMenuModels.Product, error) {
	// add to hash map
	productExist := make(map[string]string, len(products))
	posProducts := make(map[string]coreMenuModels.Product, len(products))
	for _, product := range products {
		key := strings.TrimSpace(product.ProductID + product.ParentGroupID)

		// cause has cases if product_id && parent_id same, size_id different
		productExist[key] = product.ExtID
		posProducts[product.ProductID] = product
	}

	return productExist, posProducts, nil
}

func (rkeeperSvc *rkeeperService) collectAttributes(schemeId string, menu rkeeperDto.Menu) ([]string, []string) {

	existAttributeGroup := make(map[string]struct{}, len(menu.IngredientsSchemes))

	attributeGroups := make([]string, 0, len(menu.IngredientsGroups))

	for _, scheme := range menu.IngredientsSchemes {
		if schemeId == scheme.ID {
			if len(scheme.IngredientsGroups) != 0 {
				for _, ingredientGroup := range scheme.IngredientsGroups {
					existAttributeGroup[ingredientGroup.ID] = struct{}{}
					attributeGroups = append(attributeGroups, ingredientGroup.ID)
				}
			}
		}
	}

	attributes := make([]string, 0, len(menu.Ingredients))

	for _, attributeGroup := range menu.IngredientsGroups {
		if _, ok := existAttributeGroup[attributeGroup.ID]; ok {
			attributes = append(attributes, attributeGroup.Ingredients...)
		}
	}

	return attributeGroups, attributes
}

func (rkeeperSvc *rkeeperService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	_, err := rkeeperSvc.rkeeperCli.CancelOrder(ctx, store.RKeeper.ObjectId, order.PosOrderID)
	if err != nil {
		return err
	}

	return nil
}

func (rkeeperSvc *rkeeperService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (rkeeperSvc *rkeeperService) fillPickupDelivery(order models.Order, expeditionType string, store coreStoreModels.Store, rkeeperOrder rkeeperDto.Order) rkeeperDto.Order {
	expTimeStr := order.EstimatedPickupTime.Value.Time.Add(time.Duration(store.Settings.TimeZone.UTCOffset) * time.Hour).Format("2006-01-02T15:04:05") // :TODO
	expTime, _ := time.Parse("2006-01-02T15:04:05", expTimeStr)
	if expeditionType == "pickup" {
		rkeeperOrder.Pickup = rkeeperDto.PickUp{
			Courier: &rkeeperDto.PersonInfo{
				Name:  order.Courier.Name,
				Phone: order.Courier.PhoneNumber,
			},
			ExpectedTime: expTime,
			Taker:        rkeeperSvc.fillTaker(order.IsPickedUpByCustomer),
		}

		return rkeeperOrder
	}

	rkeeperOrder.Delivery = rkeeperDto.Delivery{
		ExpectedTime: expTime,
		Address: &rkeeperDto.DeliveryAddress{
			FullAddress: order.DeliveryAddress.Label,
		},
	}
	return rkeeperOrder

}

func (r *rkeeperService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (r *rkeeperService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}

func (rkeeperSvc *rkeeperService) CreateOrderTask(ctx context.Context, taskGUID string) (string, error) {

	var (
		err         error
		posOrderId  string
		posResponse rkeeperDto.CreateOrderTaskResponse
	)

	timeout := 210
	for i := 0; i <= timeout; i += 20 {
		time.Sleep(20 * time.Second)
		log.Info().Msgf("retry count attempt %d", i)
		posResponse, err = rkeeperSvc.rkeeperCli.CreateOrderTask(ctx, taskGUID)
		if err != nil {
			log.Info().Msgf("rkeeperSvc.rkeeperCli.CreateOrderTask error: %s", err.Error())
			continue
		}

		if posResponse.TaskResponse.Order.OrderGuid == "" {
			continue
		}

		posOrderId = posResponse.TaskResponse.Order.OrderGuid
		break
	}
	slog.Info("rkeeperSvc", "posOrderID", posOrderId)
	if posOrderId == "" && (posResponse.Error.WsError.Desc != "" || posResponse.Error.WsError.Code != "" ||
		posResponse.Error.AgentError.Desc != "" || posResponse.Error.AgentError.Code != "") {
		log.Info().Msg("pos order id is empty and posResponse.Error is not null")
		return "", fmt.Errorf("create order task error: %+v", posResponse.Error)
	} else if posOrderId == "" && posResponse.Error.WsError.Desc == "" && posResponse.Error.WsError.Code == "" &&
		posResponse.Error.AgentError.Desc == "" && posResponse.Error.AgentError.Code == "" {
		return "", errors.New("wsa timeout")
	}

	return posOrderId, nil
}

func (rkeeperSvc *rkeeperService) ValidateCreateOrderTaskError(ctx context.Context, errSolution error_solutions.Service, syncResponse rkeeperDto.SyncResponse, order models.Order, err error) (models.Order, error) {

	log.Info().Msgf("ValidateCreateOrderTaskError for err: %s", err.Error())

	errorSolutions, err := errSolution.GetAllErrorSolutions(ctx)
	if err != nil {
		log.Err(err).Msg("rkeeperSvc error: GetAllErrorSolutions")
		return models.Order{}, err
	}

	errMessage := syncResponse.Error.WsError.Desc + syncResponse.Error.AgentError.Desc + err.Error()

	switch err {
	case errors.New("wsa timeout"):
		failReason, _, setFailReasonErr := errSolution.SetFailReason(ctx, coreStoreModels.Store{}, errMessage, "", WSA_TIMEOUT)
		if setFailReasonErr != nil {
			return models.Order{}, setFailReasonErr
		}
		order.FailReason = failReason
		return order, errors.New(errMessage)
	default:
		failReason, _, setFailReasonErr := errSolution.SetFailReason(ctx, coreStoreModels.Store{}, errMessage, MatchingCodes(errMessage, errorSolutions), "")
		if setFailReasonErr != nil {
			return models.Order{}, setFailReasonErr
		}
		order.FailReason = failReason
		return order, errors.New(errMessage)
	}
}
