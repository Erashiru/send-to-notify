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
	rkeeper7_xml "github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml"
	rkeeperXMLConf "github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/create_order_request"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/create_order_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_menu_by_category"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_menu_modifiers_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_groups"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_schema_details"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_schemas"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/save_order_request"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

type rkeeper7XMLService struct {
	*BasePosService
	rkeeper7XMLCli rkeeperXMLConf.RKeeper7
}

var (
	seqNumberError         string = "В инстансе неверно заполнен атрибут seqNumber"
	seqNumberIncreaseError string = "SeqNumber должен быть увеличен"
	seqNumberError2        string = "В инстансе неверно заполнен атрибут SeqNumber"
)

func newRkeeper7XMLService(bps *BasePosService, baseUrl, username, password, ucsUsername, ucsPassword, token, licenseBaseUrl,
	anchor, objectId, stationId, stationCode, licenseInstanceGUID string, childItems,
	classificatorItemIdent int, classificatorPropMask, menuItemsPropMask, propFilter, cashier string) (*rkeeper7XMLService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "rkeeper7XMLService constructor error")
	}

	client, err := rkeeper7_xml.NewClient(&rkeeperXMLConf.Config{
		Protocol:               "http",
		BaseURL:                baseUrl,
		Username:               username,
		Password:               password,
		UCSUsername:            ucsUsername,
		UCSPassword:            ucsPassword,
		Token:                  token,
		LicenseInstanceGUID:    licenseInstanceGUID,
		LicenseBaseURL:         licenseBaseUrl,
		Anchor:                 anchor,
		ObjectID:               objectId,
		StationID:              stationId,
		StationCode:            stationCode,
		ChildItems:             childItems,
		ClassificatorItemIdent: classificatorItemIdent,
		ClassificatorPropMask:  classificatorPropMask,
		MenuItemsPropMask:      menuItemsPropMask,
		PropFilter:             propFilter,
		Cashier:                cashier,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize RKeeperXML Client.")
		return nil, err
	}

	_, err = client.SetLicense(context.Background())
	if err != nil {
		return nil, err
	}

	return &rkeeper7XMLService{
		BasePosService: bps,
		rkeeper7XMLCli: client,
	}, nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	menuItems, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetMenuByCategory(ctx)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	menuModifiers, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetMenuModifiers(ctx)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	menuModifierGroups, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetMenuModifierGroups(ctx)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	menuModifierSchemaDetails, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetMenuModifierSchemaDetails(ctx)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	menuModifierSchemas, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetMenuModifierSchemas(ctx)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	orderMenu, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetOrderMenu(ctx)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	mappingProducts := make(map[string]string)
	mappingAttributes := make(map[string]string)

	for _, dish := range orderMenu.Dishes.Item {
		if dish.Quantity != "" {
			quantity, err := strconv.Atoi(dish.Quantity)
			if err != nil {
				continue
			}

			if quantity <= 0 {
				continue
			}
		}

		mappingProducts[dish.Ident] = dish.Price
	}

	for _, modifier := range orderMenu.Modifiers.Item {
		mappingAttributes[modifier.Ident] = modifier.Price
	}

	return rkeeper7xmlSvc.menuFromClient(menuItems, menuModifiers, store.Settings, mappingProducts, mappingAttributes, menuModifierGroups, menuModifierSchemas, menuModifierSchemaDetails), nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {
	menuItems, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetMenuItems(ctx)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}

	menuModifiers, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetMenuModifiers(ctx)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}

	orderMenu, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetOrderMenu(ctx)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}

	availableItems := make(map[string]struct{}, len(orderMenu.Dishes.Item))
	availableModifiers := make(map[string]struct{}, len(orderMenu.Modifiers.Item))

	for _, item := range orderMenu.Dishes.Item {
		availableItems[item.Ident] = struct{}{}
	}

	for _, modifier := range orderMenu.Modifiers.Item {
		availableModifiers[modifier.Ident] = struct{}{}
	}

	stoplistItems := make(coreMenuModels.StopListItems, 0)

	for _, item := range menuItems.RK7Reference.Items.Item {
		if _, ok := availableItems[item.Ident]; !ok {
			stoplistItems = append(stoplistItems, coreMenuModels.StopListItem{
				ProductID: item.Ident,
			})
		}
	}

	for _, modifier := range menuModifiers.RK7Reference.Items.Item {
		if _, ok := availableModifiers[modifier.Ident]; !ok {
			stoplistItems = append(stoplistItems, coreMenuModels.StopListItem{
				ProductID: modifier.Ident,
			})
		}
	}

	return stoplistItems, nil
}

func (rkeeper7XMLService *rkeeper7XMLService) GetOrderStatus(ctx context.Context, order models.Order) (string, error) {
	response, err := rkeeper7XMLService.rkeeper7XMLCli.GetOrder(ctx, order.PosOrderID)
	if err != nil {
		return "", err
	}

	return response.CommandResult.Order.Finished, nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {
	switch posStatus {
	case "0":
		return models.ACCEPTED, nil
	case "1":
		return models.CLOSED, nil
	default:
		return 0, models.StatusIsNotExist
	}
}

func (rkeeper7xmlSvc *rkeeper7XMLService) setDeliveryTypeToOrder(ctx context.Context, order models.Order) error {
	return rkeeper7xmlSvc.rkeeper7XMLCli.SetDeliveryTypeToOrder(ctx, order.PosOrderID, order.PosPaymentInfo.OrderType)
}

func (rkeeper7xmlSvc *rkeeper7XMLService) sendOrder(ctx context.Context, posOrder create_order_request.RK7Query) (create_order_response.CreateOrderResponse, error) {

	response, err := rkeeper7xmlSvc.rkeeper7XMLCli.CreateOrder(ctx, posOrder.RK7CMD.Order.Table.Code, posOrder.RK7CMD.Order.Station.ID, posOrder.RK7CMD.Order.PersistentComment, posOrder.RK7CMD.Order.OrderType.Code)
	if err != nil {
		return create_order_response.CreateOrderResponse{}, err
	}

	return response, nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) constructPosOrder(order models.Order, store coreStoreModels.Store) (create_order_request.RK7Query, models.Order) {
	if order.TableID == "" {
		order.TableID = store.RKeeper7XML.DefaultTable
	}

	if store.RKeeper7XML.StationID == "" {
		store.RKeeper7XML.StationID = "1"
	}

	clientComment := ""
	if order.SpecialRequirements != "" || order.AllergyInfo != "" {
		clientComment += "\nКомментарий клиента: " + order.SpecialRequirements + order.AllergyInfo
	}

	request := create_order_request.RK7Query{
		RK7CMD: create_order_request.RK7CMD{
			Order: create_order_request.CreateOrderRequest{
				PersistentComment: "Код заказа: " + order.PickUpCode + clientComment,
				Table: create_order_request.Table{
					Code: order.TableID, // TableID
				},
				OrderType: create_order_request.OrderType{
					Code: store.RKeeper7XML.OrderTypeCode,
				},
				Station: create_order_request.Station{
					ID: store.RKeeper7XML.StationID,
				},
			},
		},
	}

	return request, order
}

func (rkeeper7xmlSvc *rkeeper7XMLService) saveOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, coreStoreModels.Store, error) {
	dishes := make(save_order_request.Dishes, 0, len(order.Products))

	for _, product := range order.Products {
		dish := save_order_request.Dish{
			ID:       product.ID,
			Price:    strconv.Itoa(int(product.Price.Value)),
			Quantity: "1",
		}

		for _, attribute := range product.Attributes {
			dish.Modi = append(dish.Modi, save_order_request.Modi{
				ID:    attribute.ID,
				Count: strconv.Itoa(attribute.Quantity),
				Price: strconv.Itoa(int(attribute.Price.Value)),
			})
		}

		for i := 0; i < product.Quantity; i++ {
			dishes = append(dishes, dish)
		}
	}

	_, err := rkeeper7xmlSvc.rkeeper7XMLCli.SaveOrder(ctx, order.PosOrderID, strconv.Itoa(store.RKeeper7XML.SeqNumber), order.PosPaymentInfo.PaymentTypeID, store.RKeeper7XML.PrepayReasonId, dishes, strconv.Itoa(int(order.EstimatedTotalPrice.Value)), store.RKeeper7XML.IsLifeTimeLicence)
	if err != nil {
		if rkeeper7xmlSvc.isWrongSeqNumber(err) && !store.RKeeper7XML.IsLifeTimeLicence {
			res, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetSeqNumber(ctx)
			if err != nil {
				return order, store, err
			}

			seqNumberForOrder, seqNumberForDB, err := rkeeper7xmlSvc.ConvertSeqNumber(res.LicenseInfo.LicenseInstance.SeqNumber)
			if err != nil {
				return order, store, fmt.Errorf("seq number converting error: %s", err)
			}

			_, err = rkeeper7xmlSvc.rkeeper7XMLCli.SaveOrder(ctx, order.PosOrderID, seqNumberForOrder, order.PosPaymentInfo.PaymentTypeID, store.RKeeper7XML.PrepayReasonId, dishes, strconv.Itoa(int(order.EstimatedTotalPrice.Value)), store.RKeeper7XML.IsLifeTimeLicence)
			if err != nil {
				return order, store, err
			}
			store.RKeeper7XML.SeqNumber = seqNumberForDB

			return order, store, nil
		}

		return order, store, err
	}

	if store.ExternalPosIntegrationSettings.PayOrderIsOn {
		log.Info().Msgf("start to send PayOrder for order: %s", order.PosGuid)
		_, err = rkeeper7xmlSvc.rkeeper7XMLCli.PayOrder(ctx, order.PosGuid, order.PosPaymentInfo.PaymentTypeID, strconv.Itoa(int(order.EstimatedTotalPrice.Value)), store.RKeeper7XML.StationCode)
		if err != nil {
			log.Error().Msgf("rkeeper7XMLService error: PayOrder error %s", err)
		}
	}
	return order, store, nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration,
	store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client,
	aggregatorMenu coreMenuModels.Menu, storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	defer func(store *coreStoreModels.Store) {
		if err := storeCli.Update(ctx, storeModels.UpdateStore{
			ID: &store.ID,
			RKeeper7XML: &storeModels.UpdateStoreRKeeper7XMLConfig{
				SeqNumber: &store.RKeeper7XML.SeqNumber,
			},
		}); err != nil {
			return
		}
	}(&store)

	var err error

	order, err = prepareAnOrder(ctx, order, store, menu, aggregatorMenu, menuClient)
	if err != nil {
		return order, err
	}

	posOrder, _ := rkeeper7xmlSvc.constructPosOrder(order, store)

	order, err = rkeeper7xmlSvc.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	errorSolutions, err := errSolution.GetAllErrorSolutions(ctx)
	if err != nil {
		log.Err(err).Msg("rkeeper7XMLService error: GetAllErrorSolutions")
		return models.Order{}, err
	}

	orderResponse, err := rkeeper7xmlSvc.sendOrder(ctx, posOrder)
	if err != nil {
		log.Err(err).Msgf("rkeeper7XMLService orderResponse.ErrorText: %s for order_id: %s", orderResponse.ErrorText, order.OrderID)
		failReason, _, setFailReasonErr := errSolution.SetFailReason(ctx, store, orderResponse.ErrorText, MatchingCodes(orderResponse.ErrorText, errorSolutions), "")
		if setFailReasonErr != nil {
			return order, setFailReasonErr
		}
		order.FailReason = failReason
		return order, err
	}

	order = setPosOrderId(order, orderResponse.VisitID)
	order.PosGuid = orderResponse.Guid
	order.CreationResult = models.CreationResult{
		Message: orderResponse.Status,
		OrderInfo: models.OrderInfo{
			ID:             orderResponse.VisitID,
			OrganizationID: store.RKeeper7XML.ObjectID,
			CreationStatus: orderResponse.Status,
		},
	}
	order.PosPaymentInfo.OrderType = store.RKeeper7XML.OrderTypeCode

	order, store, err = rkeeper7xmlSvc.saveOrder(ctx, order, store)
	if err != nil {
		failReason, _, setFailReasinErr := errSolution.SetFailReason(ctx, store, err.Error(), MatchingCodes(err.Error(), errorSolutions), "")
		if setFailReasinErr != nil {
			return order, setFailReasinErr
		}
		order.FailReason = failReason
		return order, err
	}
	store.RKeeper7XML.SeqNumber = store.RKeeper7XML.SeqNumber + 1

	return order, nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) menuFromClient(items get_menu_by_category.RK7QueryResult,
	modifiers get_menu_modifiers_response.RK7QueryResult,
	settings coreStoreModels.Settings,
	mappingProducts map[string]string,
	mappingAttributes map[string]string,
	modifierGroups get_modifier_groups.RK7QueryResult,
	modifierSchemas get_modifier_schemas.RK7QueryResult,
	modifierSchemaDetails get_modifier_schema_details.RK7QueryResult) coreMenuModels.Menu {

	relationShip := rkeeper7xmlSvc.buildRelationships(modifierSchemas, modifierSchemaDetails)

	menu := coreMenuModels.Menu{
		Name:             coreMenuModels.RKEEPER7XML.String(),
		ExtName:          coreMenuModels.MAIN.String(),
		Description:      "rkeeper pos menu",
		AttributesGroups: rkeeper7xmlSvc.modifierGroupsToModel(modifierGroups),
		Products:         rkeeper7xmlSvc.productsToModel(items, settings, mappingProducts, relationShip),
		Attributes:       rkeeper7xmlSvc.modifiersToModel(modifiers, settings, mappingAttributes),
		CreatedAt:        models.TimeNow(),
		UpdatedAt:        models.TimeNow(),
	}
	return menu
}

func (rkeeper7xmlSvc *rkeeper7XMLService) buildRelationships(modifierSchemas get_modifier_schemas.RK7QueryResult, modifierSchemaDetails get_modifier_schema_details.RK7QueryResult) map[string][]string {
	relationship := make(map[string][]string)

	for _, detail := range modifierSchemaDetails.RK7Reference.Items.Item {
		if val, ok := relationship[detail.ModiScheme]; ok {
			relationship[detail.ModiScheme] = append(val, detail.ModiGroup)
			continue
		}

		relationship[detail.ModiScheme] = []string{detail.ModiGroup}
	}

	return relationship
}

func (rkeeper7xmlSvc *rkeeper7XMLService) modifierGroupsToModel(modifierGroups get_modifier_groups.RK7QueryResult) coreMenuModels.AttributeGroups {
	var attributeGroups = make(coreMenuModels.AttributeGroups, 0, len(modifierGroups.RK7Reference.Items.Item))

	for _, modifierGroup := range modifierGroups.RK7Reference.Items.Item {
		ids := make([]string, 0, len(modifierGroup.Childs.Child))

		for _, child := range modifierGroup.Childs.Child {
			ids = append(ids, child.ChildIdent)
		}

		attributeGroups = append(attributeGroups, coreMenuModels.AttributeGroup{
			ExtID:      modifierGroup.Ident,
			Name:       modifierGroup.Name,
			Attributes: ids,
			Min:        0, // TODO: ?
			Max:        1, // TODO: ?
		})
	}

	return attributeGroups
}

func (rkeeper7xmlSvc *rkeeper7XMLService) getProductAttributeGroupsIDs(modiSchemeID string, relationShip map[string][]string) []string {
	if ids, ok := relationShip[modiSchemeID]; ok {
		return ids
	}

	return nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) productsToModel(req get_menu_by_category.RK7QueryResult, settings coreStoreModels.Settings, mappingProducts map[string]string, relationship map[string][]string) coreMenuModels.Products {
	results := make(coreMenuModels.Products, 0, len(mappingProducts))

	for _, product := range req.CommandResult[1].RK7Reference.Items.Item {
		attributeGroupsIDs := rkeeper7xmlSvc.getProductAttributeGroupsIDs(product.ModiScheme, relationship)

		res, err := rkeeper7xmlSvc.productToModel(product, settings, mappingProducts, attributeGroupsIDs)
		if err != nil {
			log.Err(err).Msgf("rkeeper7 xml cli err: get product %s", product.Ident)
			continue
		}

		results = append(results, res)
	}

	return results
}

func (rkeeper7xmlSvc *rkeeper7XMLService) productToModel(item get_menu_by_category.Item, setting coreStoreModels.Settings, mapping map[string]string, attributeGroupIDs []string) (coreMenuModels.Product, error) {
	res := coreMenuModels.Product{
		ProductID: item.Ident, // TODO: id?
		ExtID:     item.Ident, // TODO: id?
		ExtName:   item.Name,
		Name: []coreMenuModels.LanguageDescription{
			{
				Value: item.Name,
			},
		},
		ImageURLs: []string{item.Genphotolink}, // TODO: image?
		ProductsCreatedAt: coreMenuModels.ProductsCreatedAt{
			Value:     models.TimeNow(),
			Timezone:  setting.TimeZone.TZ,
			UTCOffset: setting.TimeZone.UTCOffset,
		},
		AttributesGroups: attributeGroupIDs,
		UpdatedAt:        models.TimeNow(),
		IsIncludedInMenu: true,
	}

	if val, ok := mapping[item.Ident]; !ok {
		res.IsAvailable = false
		res.Price = []coreMenuModels.Price{
			{
				Value:        0,
				CurrencyCode: setting.Currency,
			},
		}
	} else {
		price, err := strconv.Atoi(val)
		if err != nil {
			return coreMenuModels.Product{}, err
		}

		price = price / 100

		res.Price = []coreMenuModels.Price{
			{
				Value:        float64(price),
				CurrencyCode: setting.Currency,
			},
		}
		res.IsAvailable = true
	}

	return res, nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) modifiersToModel(menuModifiers get_menu_modifiers_response.RK7QueryResult, settings coreStoreModels.Settings, mappingAttributes map[string]string) coreMenuModels.Attributes {
	results := make(coreMenuModels.Attributes, 0, len(mappingAttributes))

	for _, item := range menuModifiers.RK7Reference.Items.Item {
		res, err := rkeeper7xmlSvc.modifierToModel(item, settings, mappingAttributes)
		if err != nil {
			log.Err(err).Msgf("rkeeper7 xml cli err: get attribute %s", item.Ident)
			continue
		}

		results = append(results, res)
	}

	return results
}

func (rkeeper7xmlSvc *rkeeper7XMLService) modifierToModel(item get_menu_modifiers_response.Item, setting coreStoreModels.Settings, mappingAttributes map[string]string) (coreMenuModels.Attribute, error) {
	res := coreMenuModels.Attribute{
		ExtID:     item.Ident, // TODO: id?
		ExtName:   item.Name,
		Name:      item.Name,
		UpdatedAt: models.TimeNow().Time,
	}

	if val, ok := mappingAttributes[item.Ident]; !ok {
		res.IsAvailable = false
		res.IncludedInMenu = false
	} else {
		price, err := strconv.Atoi(val)
		if err != nil {
			return coreMenuModels.Attribute{}, err
		}

		price = price / 100

		res.Price = float64(price)
		res.IsAvailable = true
		res.IncludedInMenu = true
	}

	return res, nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	return nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) isWrongSeqNumber(err error) bool {
	return strings.Contains(err.Error(), seqNumberError) || strings.Contains(err.Error(), seqNumberIncreaseError) || strings.Contains(err.Error(), seqNumberError2)
}

func (rkeeper7xmlSvc *rkeeper7XMLService) GetSeqNumber(ctx context.Context) (string, error) {
	res, err := rkeeper7xmlSvc.rkeeper7XMLCli.GetSeqNumber(ctx)
	if err != nil {
		return "", err
	}
	return res.LicenseInfo.LicenseInstance.SeqNumber, nil
}

func (rkeeper7xmlSvc *rkeeper7XMLService) ConvertSeqNumber(seqNumber string) (string, int, error) {
	seqNumberInt, err := strconv.Atoi(seqNumber)
	if err != nil {
		return "", 0, err
	}

	seqNumberInt = seqNumberInt + 1

	return strconv.Itoa(seqNumberInt), seqNumberInt, nil
}

func (r *rkeeper7XMLService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (r *rkeeper7XMLService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
