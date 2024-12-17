package pos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/config"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/custom"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/domain/logger"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	menuUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	posterClient "github.com/kwaaka-team/orders-core/pkg/poster"
	posterConf "github.com/kwaaka-team/orders-core/pkg/poster/clients"
	posterModels "github.com/kwaaka-team/orders-core/pkg/poster/clients/models"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	models2 "github.com/kwaaka-team/orders-core/service/pos/models/poster"
	"github.com/kwaaka-team/orders-core/service/pos/repository"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type PosterService struct {
	*BasePosService
	posterCli         posterConf.Poster
	storeAuthRepo     repository.PosterStoreAuthsRepository
	storeCli          storeClient.Client
	menuCli           menuCore.Client
	applicationID     string
	applicationSecret string
	redirectURI       string
	logger            *zap.SugaredLogger
}

func NewPosterService(bps *BasePosService, baseURL, token string, storeAuthRepo repository.PosterStoreAuthsRepository,
	storeCli storeClient.Client, menuCli menuCore.Client,
	applicationID, applicationSecret, redirectURI string,
	logger *zap.SugaredLogger) (*PosterService, error) {

	posterCli, err := posterClient.NewClient(&posterConf.Config{
		Protocol: "http",
		BaseURL:  baseURL,
		Token:    token,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Poster Client.")
		return nil, err
	}

	return &PosterService{
		BasePosService:    bps,
		posterCli:         posterCli,
		storeAuthRepo:     storeAuthRepo,
		storeCli:          storeCli,
		menuCli:           menuCli,
		applicationSecret: applicationSecret,
		applicationID:     applicationID,
		redirectURI:       redirectURI,
		logger:            logger,
	}, nil
}

func (s *PosterService) GetOrderStatus(ctx context.Context, order models.Order) (string, error) {
	if order.CookingCompleteTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
		return "", nil
	}
	if order.CookingCompleteTime.After(time.Now().UTC()) {
		return "", nil
	}
	if s.hasStatusInStatusesHistory("COOKING_COMPLETE", order.StatusesHistory) {
		return "", nil
	}
	log.Info().Msgf("update order status in aggregator for id=%s, order_id=%s, order_code=%s, delivery=%s", order.ID, order.OrderID, order.OrderCode, order.DeliveryService)

	return "ready", nil
}

func (s *PosterService) hasStatusInStatusesHistory(status string, history []models.OrderStatusUpdate) bool {
	for _, s := range history {
		if s.Name == status {
			return true
		}
	}
	return false
}

func (s *PosterService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {
	switch posStatus {
	case "closed":
		return models.ACCEPTED, nil
	case "ready":
		return models.CLOSED, nil
	}
	return 0, models.StatusIsNotExist
}

func (s *PosterService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration,
	store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	var err error

	order, err = s.prepareAnOrder(ctx, order, store, menu, aggregatorMenu, menuClient)
	if err != nil {
		return order, err
	}

	posOrder, _, err := s.constructPosOrder(ctx, order, store)
	if err != nil {
		return order, validator.ErrCastingPos
	}

	posOrder, ok := posOrder.(posterModels.CreateOrderRequest)
	if !ok {
		return order, validator.ErrCastingPos
	}

	order, err = s.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	errorSolutions, err := errSolution.GetAllErrorSolutions(ctx)
	if err != nil {
		log.Err(err).Msg("PosterService error: GetAllErrorSolutions")
		return models.Order{}, err
	}

	responseOrder, err := s.sendOrder(ctx, posOrder)
	if err != nil {
		log.Err(err).Msgf("PosterService orderResponse.ErrorText: %s for order_id: %s", responseOrder.ErrorResponse.Message, order.OrderID)
		failReason, _, failReasonErr := errSolution.SetFailReason(ctx, store, err.Error(), MatchingCodes(responseOrder.ErrorResponse.Message, errorSolutions), "")
		if failReasonErr != nil {
			return order, failReasonErr
		}
		order.FailReason = failReason
		return order, err
	}

	order = setPosOrderId(order, store.Poster.AccountNumberString+strconv.Itoa(responseOrder.Response.IncomingOrderID))

	order.CreationResult = models.CreationResult{
		Message: responseOrder.Message,
		OrderInfo: models.OrderInfo{
			CreationStatus: strconv.Itoa(responseOrder.Response.Status),
			OrganizationID: strconv.Itoa(responseOrder.Response.SpotId),
		},
		ErrorDescription: responseOrder.ErrorResponse.Message,
	}
	log.Info().Msgf("sending msg to queue, posOrderID: %s", order.PosOrderID)

	return order, nil
}

func (p *PosterService) constructPhoneNumber(order models.Order) string {
	phoneNumber := order.Customer.PhoneNumber
	if phoneNumber == "" {
		phoneNumber = "+77771111111"
	}

	return phoneNumber
}

func (p *PosterService) constructPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	layout := "2006-01-02 15:04:05"
	orderComment, _, _ := p.BasePosService.constructOrderComments(order, store)
	utcOffset := time.Duration(int64(store.Settings.TimeZone.UTCOffset)) * time.Hour
	deliveryTime := order.EstimatedPickupTime.Value.Add(utcOffset)

	posterOrder := posterModels.CreateOrderRequest{
		SpotID: store.Poster.SpotId,
		Phone:  p.constructPhoneNumber(order),
		ClientAddress: posterModels.CreateOrderAddressRequest{
			Address1: order.DeliveryAddress.Label,
		},
		Comment:      orderComment,
		DeliveryTime: deliveryTime.Format(layout),
	}

	if order.DeliveryAddress.Latitude != 0 && order.DeliveryAddress.Longitude != 0 {
		posterOrder.ClientAddress.Longitude = fmt.Sprintf("%v", order.DeliveryAddress.Longitude)
		posterOrder.ClientAddress.Latitude = fmt.Sprintf("%v", order.DeliveryAddress.Latitude)
	}

	if order.PaymentMethod == "DELAYED" {
		posterOrder.Payment = &posterModels.CreateOrderPaymentRequest{
			Type:     1,
			Sum:      int(order.EstimatedTotalPrice.Value * 100),
			Currency: store.Settings.Currency,
		}
	}

	if store.Poster.IgnorePaymentType {
		posterOrder.Payment = &posterModels.CreateOrderPaymentRequest{
			Type:     int(store.Poster.PaymentType),
			Sum:      p.getSumByPaymentType(store.Poster.PaymentType, order.EstimatedTotalPrice.Value),
			Currency: store.Settings.Currency,
		}
	}

	switch order.IsPickedUpByCustomer {
	case true:
		posterOrder.ServiceMode = 2
	case false:
		posterOrder.ServiceMode = 3
		posterOrder.DeliveryPrice = int(order.DeliveryFee.Value * 100)
	}

	var posterProducts = make([]posterModels.CreateOrderProductRequest, 0)

	for _, product := range order.Products {
		productID, err := strconv.Atoi(product.ID)
		if err != nil {
			return nil, models.Order{}, err
		}
		attributesPrice := 0.0
		productAttributes := make([]posterModels.CreateOrderModificationRequest, 0, len(product.Attributes))
		for _, attribute := range product.Attributes {
			attributeId, err := strconv.Atoi(attribute.ID)
			if err != nil {
				return nil, models.Order{}, err
			}
			posterAttribute := posterModels.CreateOrderModificationRequest{
				M: attributeId,
				A: attribute.Quantity,
			}
			productAttributes = append(productAttributes, posterAttribute)
			attributesPrice = attributesPrice + attribute.Price.Value*float64(attribute.Quantity)
		}
		posterProduct := posterModels.CreateOrderProductRequest{
			ProductID:     productID,
			ModificatorID: 0,
			Modifications: productAttributes,
			Count:         strconv.Itoa(product.Quantity),
			Price:         int(product.Price.Value+attributesPrice) * 100,
		}
		posterProducts = append(posterProducts, posterProduct)
	}

	posterOrder.Products = posterProducts

	return posterOrder, order, nil
}

func (p *PosterService) getSumByPaymentType(paymentType int32, sum float64) int {
	if paymentType == 0 {
		return 0
	}

	return int(sum * 100)
}

func (p *PosterService) sendOrder(ctx context.Context, order any) (posterModels.CreateOrderResponse, error) {
	var errs custom.Error

	posOrder, ok := order.(posterModels.CreateOrderRequest)

	if !ok {
		return posterModels.CreateOrderResponse{}, validator.ErrCastingPos
	}

	utils.Beautify("Poster Request Body", posOrder)

	createResponse, err := p.posterCli.CreateOrder(ctx, posOrder)
	if err != nil {
		log.Err(err).Msg("poster error")
		errs.Append(err, validator.ErrIgnoringPos)
		return posterModels.CreateOrderResponse{}, errs
	}
	return createResponse, nil
}

func (p *PosterService) GetStopList(ctx context.Context) (result coreMenuModels.StopListItems, err error) {

	stocks, err := p.posterCli.GetStopList(ctx)
	if err != nil {
		return result, err
	}

	for _, stock := range stocks.StopLists {
		balance, err := strconv.ParseFloat(stock.IngredientLeft, 64)
		if err != nil {
			return result, err
		}
		if balance <= 0 {
			result = append(result, coreMenuModels.StopListItem{
				ProductID: stock.IngredientId,
				Balance:   balance,
			})
		}
	}

	log.Info().Msgf("MC[INFO:] [GetStopList] stock poster:  %+v \n  result %+v ", stocks, result)

	return result, nil
}

func (p *PosterService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	menu, err := p.posterCli.GetProducts(ctx)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	products, err := p.existProducts(ctx, systemMenuInDb.Products)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	return p.menuFromClient(menu, store, products)
}

func (p *PosterService) existProducts(ctx context.Context, products []coreMenuModels.Product) (map[string]coreMenuModels.Product, error) {
	// add to hash map
	productExist := make(map[string]coreMenuModels.Product, len(products))
	for _, product := range products {
		// cause has cases if product_id && parent_id same, size_id different
		productExist[product.ProductID] = product
	}

	return productExist, nil
}

func (p *PosterService) menuFromClient(posterProducts posterModels.GetProductsResponse, store coreStoreModels.Store, productsExist map[string]coreMenuModels.Product) (coreMenuModels.Menu, error) {
	products, attributeGroups, attributes, err := p.toEntities(posterProducts.Response, store, productsExist)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	menu := coreMenuModels.Menu{
		Name:             coreMenuModels.POSTER.String(),
		ExtName:          coreMenuModels.MAIN.String(),
		CreatedAt:        models.TimeNow(),
		UpdatedAt:        models.TimeNow(),
		Products:         products,
		AttributesGroups: attributeGroups,
		Attributes:       attributes,
	}

	return menu, nil
}

func (p *PosterService) toEntities(posterProducts []posterModels.GetProductsResponseBody, store coreStoreModels.Store, productsExist map[string]coreMenuModels.Product) ([]coreMenuModels.Product, []coreMenuModels.AttributeGroup, []coreMenuModels.Attribute, error) {
	var (
		products                = make([]coreMenuModels.Product, 0, len(posterProducts))
		attributeGroups         = make([]coreMenuModels.AttributeGroup, 0, 4)
		attributes              = make([]coreMenuModels.Attribute, 0, 4)
		uniqueAttributeGroups   = make(map[string]struct{})
		uniqueAttributes        = make(map[string]struct{})
		uniqueAvailableProducts = make(map[string]bool)
	)

	for _, item := range posterProducts {
		productAvailable := p.fromSpotsVisibleToAvailable(store.Poster.SpotId, item.Spots)
		product := coreMenuModels.Product{
			ExtID:        item.ProductId,
			ProductID:    item.ProductId,
			PosID:        item.ProductId,
			IngredientID: item.IngredientId,
			Section:      item.MenuCategoryId,
			Name: []coreMenuModels.LanguageDescription{
				{
					Value:        item.ProductName,
					LanguageCode: store.Settings.LanguageCode,
				},
			},
			IsAvailable:      productAvailable,
			ImageURLs:        []string{item.Photo},
			IsIncludedInMenu: p.fromHiddenToBool(item.Hidden),
		}
		uniqueAvailableProducts[item.ProductId] = productAvailable

		if item.Price.Field1 != "" {
			price, err := strconv.Atoi(item.Price.Field1)
			if err != nil {
				log.Err(err).Msgf("price for product %s is not exist", item.ProductName)
				continue
			}
			product.Price = []coreMenuModels.Price{
				{
					Value:        float64(price) / 100,
					CurrencyCode: store.Settings.Currency,
				},
			}
		}

		//teh-karta - have group-modifications
		for _, modifierGroup := range item.GroupModifications {
			log.Info().Msgf("GroupModifications  %v", modifierGroup.DishModificationGroupId)
			attributeGroupID := strconv.Itoa(modifierGroup.DishModificationGroupId)
			_, ok := uniqueAttributeGroups[attributeGroupID]
			if !ok {
				attributeGroup := coreMenuModels.AttributeGroup{
					PosID: attributeGroupID,
					ExtID: attributeGroupID,
					Name:  modifierGroup.Name,
				}

				attributeGroup.Min = modifierGroup.NumMin
				attributeGroup.Max = modifierGroup.NumMax

				for _, modifier := range modifierGroup.Modifications {
					attributeID := strconv.Itoa(modifier.DishModificationId)
					ingredientID := strconv.Itoa(modifier.IngredientId)
					_, exist := uniqueAttributes[attributeID]
					if !exist {
						attributeAvailable := true
						available, existsInProducts := uniqueAvailableProducts[ingredientID]
						if existsInProducts {
							attributeAvailable = available
						}
						attribute := coreMenuModels.Attribute{
							ExtID:        ingredientID,
							PosID:        ingredientID,
							IngredientID: attributeID,
							Name:         modifier.Name,
							IsAvailable:  attributeAvailable,
						}
						attribute.Price = float64(modifier.Price)
						attributes = append(attributes, attribute)
						uniqueAttributes[attributeID] = struct{}{}
					}

					//attributeIDs = append(attributeIDs, id)
					attributeGroup.Attributes = append(attributeGroup.Attributes, ingredientID)
				}

				attributeGroups = append(attributeGroups, attributeGroup)
				uniqueAttributeGroups[attributeGroupID] = struct{}{}
			}

			product.AttributesGroups = append(product.AttributesGroups, attributeGroupID)
		}

		if posProduct, posProductExists := productsExist[product.ExtID]; posProductExists {
			if posProduct.MenuDefaultAttributes != nil {
				for _, defaultAttribute := range posProduct.MenuDefaultAttributes {
					if defaultAttribute.ByAdmin {
						product.MenuDefaultAttributes = append(product.MenuDefaultAttributes, defaultAttribute)
					}
				}
			}
			product.CookingTime = posProduct.CookingTime
		}
		products = append(products, product)

		//size product; len(product) create
		for _, modifier := range item.Modifications {
			log.Info().Msgf("Modifications len Sizes len mods %v  mod_id% v", len(item.Modifications), modifier.ModificatorID)

			var sizeProduct = product

			sizeProduct.SizeID = modifier.ModificatorID
			sizeProduct.ExtID = product.ExtID + "_" + modifier.ModificatorID

			price, err := strconv.ParseFloat(modifier.ModificatorSelfprice, 64)
			if err != nil {
				return nil, nil, nil, err
			}

			sizeProduct.Price = []coreMenuModels.Price{{
				Value:        price / 100,
				CurrencyCode: store.Settings.Currency,
			}}

			if posSizeProduct, posSizeProductExists := productsExist[sizeProduct.ExtID]; posSizeProductExists {
				if posSizeProduct.MenuDefaultAttributes != nil {
					for _, defaultAttribute := range posSizeProduct.MenuDefaultAttributes {
						if defaultAttribute.ByAdmin {
							sizeProduct.MenuDefaultAttributes = append(sizeProduct.MenuDefaultAttributes, defaultAttribute)
						}
					}
				}
			}

			products = append(products, sizeProduct)
		}

	}

	log.Info().Msgf("len AG %+v len attr %v ", len(attributeGroups), len(attributes))

	return products, attributeGroups, attributes, nil
}

func (p *PosterService) fromHiddenToBool(hidden string) bool {
	return hidden != "1"
}

func (p *PosterService) fromSpotsVisibleToAvailable(storeSpotID string, posterProductSpots []posterModels.GetProductsResponseStop) bool {
	isAvailable := false
	for _, productSpot := range posterProductSpots {
		if productSpot.SpotId != storeSpotID {
			continue
		}
		isAvailable = productSpot.Visible != "0"
	}
	return isAvailable
}

func (man *PosterService) StoreAuth(ctx context.Context, code, account string) error {
	postUrl := fmt.Sprintf("https://%s.joinposter.com/api/v2/auth/access_token", account)
	data := url.Values{
		"application_id":     {man.applicationID},
		"application_secret": {man.applicationSecret},
		"grant_type":         {"authorization_code"},
		"redirect_uri":       {man.redirectURI},
		"code":               {code},
	}

	resp, err := http.PostForm(postUrl, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resErr models2.PosterStoreAuthError
	if err = json.Unmarshal(body, &resErr); err != nil {
		return err
	}
	if resErr.Code > 200 {
		man.logger.Error(logger.LoggerInfo{
			System:   "poster response error",
			Response: string(body),
		})
		return errors.New(resErr.ErrorMessage)
	}

	var res models2.PosterStoreAuth

	if err = json.Unmarshal(body, &res); err != nil {
		return err
	}

	_, err = man.storeAuthRepo.InsertCustomer(ctx, res)
	if err != nil {
		return err
	}

	return nil
}

func (man *PosterService) GetStoplistByBalanceItems(ctx context.Context, wh models2.WHEvent) (map[string]models2.RestaurantStoplistItems, error) {

	stores, err := man.getStoresByAccountNumber(ctx, wh.AccountNumber)
	if err != nil {
		return nil, err
	}

	data, err := models2.CastData(wh.Data)
	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "poster response error",
			Response: err,
		})
		return nil, err
	}

	if data.Type != 2 && data.Type != 4 { //tovar || teh-karta
		return nil, errors.New("item type is invalid")
	}

	man.logger.Info(logger.LoggerInfo{
		System:   "poster info",
		Response: fmt.Sprintf("stoplist item is tovar || teh karta %+v \n  accountNumber %v \n", data, wh.AccountNumber),
	})

	res := make(map[string]models2.RestaurantStoplistItems, len(stores))

	for _, store := range stores {
		// for poster stoplist realized only by products balance
		if !store.Poster.StopListByBalance {
			man.logger.Error(logger.LoggerInfo{
				System:   "poster response error",
				Response: errors.New("stop list by balance false"),
			})
			return nil, errors.New("stop list by balance false")
		}
		isAvailable := false
		if data.ValueAbsolute > 0 {
			isAvailable = true
		}

		storeItems, err := man.getStoreStoplistItems(ctx, store, strconv.Itoa(data.ElementId), isAvailable)
		if err != nil {
			return nil, err
		}
		res[store.ID] = storeItems
	}

	return res, nil
}

func (man *PosterService) getAttributeID(posMenu coreMenuModels.Menu, posterIngredientID string) (coreMenuModels.Attribute, error) {
	for _, attribute := range posMenu.Attributes {
		if attribute.ExtID == posterIngredientID {
			return attribute, nil
		}
	}

	return coreMenuModels.Attribute{}, errors.Errorf("not found attribute in posMenu by id: %s", posterIngredientID)
}

func (man *PosterService) getProductID(posMenu coreMenuModels.Menu, posterProductID string) (coreMenuModels.Product, error) {
	for _, product := range posMenu.Products {
		if product.ProductID == posterProductID {
			return product, nil
		}
	}
	return coreMenuModels.Product{}, errors.Errorf("not found product in posMenu by id: %s", posterProductID)
}

func (man *PosterService) GetStoplistItems(ctx context.Context, wh models2.WHEvent) (map[string]models2.RestaurantStoplistItems, error) {
	stores, err := man.getStoresByAccountNumber(ctx, wh.AccountNumber)
	if err != nil {
		return nil, err
	}
	if len(stores) == 0 {
		return nil, errors.New("not found stores")
	}

	if !man.isStopListByAdmin(stores) {
		return nil, nil
	}

	posterCli, err := posterClient.NewClient(&posterConf.Config{
		Protocol: "http",
		BaseURL:  "https://joinposter.com",
		Token:    stores[0].Poster.Token,
	})

	if err != nil {
		log.Trace().Err(err).Msg("cant initialize Poster Client")
	}

	posterProduct, err := posterCli.GetProduct(ctx, wh.ObjectID)
	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "poster response error",
			Response: err,
		})
		return nil, err
	}

	res, err := man.getSpotsStoplistItems(ctx, wh.ObjectID, posterProduct.Spots, stores)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (man *PosterService) isStopListByAdmin(stores []coreStoreModels.Store) bool {
	return !(stores[0].Poster.IgnoreStopListFromAdmin)
}

func (man *PosterService) getSpotsStoplistItems(ctx context.Context, itemID string, spots []posterModels.GetProductsResponseStop, stores []coreStoreModels.Store) (map[string]models2.RestaurantStoplistItems, error) {
	res := make(map[string]models2.RestaurantStoplistItems, len(stores))

	for _, spot := range spots {
		for _, store := range stores {
			if store.Poster.SpotId != spot.SpotId {
				continue
			}

			storeItems, err := man.getStoreStoplistItems(ctx, store, itemID, spot.Visible != "0")
			if err != nil {
				return nil, err
			}
			res[store.ID] = storeItems
		}
	}

	return res, nil
}

func (man *PosterService) getStoreStoplistItems(ctx context.Context, rst coreStoreModels.Store, elementID string, isAvailable bool) (models2.RestaurantStoplistItems, error) {
	posMenu, err := man.menuCli.GetMenuByID(ctx, rst.MenuID)
	if err != nil {
		man.logger.Info(logger.LoggerInfo{
			System:   "poster info",
			Response: fmt.Sprintf("get pos menu err %s, id %s", err.Error(), rst.MenuID),
		})
		return models2.RestaurantStoplistItems{}, err
	}
	attribute, err := man.getAttributeID(posMenu, elementID)
	if err != nil {
		man.logger.Info(logger.LoggerInfo{
			System:   "poster info",
			Response: err,
		})
	}
	resAttributeID := ""
	if err == nil && attribute.IsAvailable != isAvailable {
		resAttributeID = attribute.ExtID
	}

	product, err := man.getProductID(posMenu, elementID)
	if err != nil {
		man.logger.Info(logger.LoggerInfo{
			System:   "poster info",
			Response: err,
		})
	}
	resProductID := ""
	if err == nil && product.IsAvailable != isAvailable {
		resProductID = product.ProductID
	}

	return models2.RestaurantStoplistItems{
		ProductID:   resProductID,
		AttributeID: resAttributeID,
		IsAvailable: isAvailable,
	}, nil
}

func (man *PosterService) getStoresByAccountNumber(ctx context.Context, accountNumber string) ([]coreStoreModels.Store, error) {
	storesDB, err := man.storeCli.FindStores(
		ctx, dto.StoreSelector{PosterAccountNumber: accountNumber},
	)
	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "poster response error",
			Response: err,
		})
		return nil, err
	}
	if len(storesDB) == 0 {
		man.logger.Error(logger.LoggerInfo{
			System:   "poster response error",
			Response: "not found stores",
		})
		return nil, err
	}

	return storesDB, nil
}

func (man *PosterService) prepareAnOrder(ctx context.Context, order models.Order, store coreStoreModels.Store, menu coreMenuModels.Menu, aggregatorMenu coreMenuModels.Menu, menuCli menuCore.Client) (models.Order, error) {

	var (
		serviceFee float64
		err        error
	)

	order = setCustomerPhoneNumber(ctx, order, store)

	order, promosMap, giftMap, promoWithPercentMap, err := getPromosMap(ctx, order, menuCli)
	if err != nil {
		return order, err
	}

	order, serviceFee, err = man.fullFillProducts(
		order, store,
		menuUtils.ProductsMap(menu),
		menuUtils.AtributesMap(menu),
		menuUtils.AtributeGroupsMap(menu),
		aggregatorMenu,
	)
	if err != nil {
		return order, err
	}

	if serviceFee != 0 {
		order.HasServiceFee = true
		order.ServiceFeeSum = serviceFee
		order.EstimatedTotalPrice.Value = order.EstimatedTotalPrice.Value - serviceFee
	}

	order = applyOrderDiscount(ctx, order, promosMap, giftMap, promoWithPercentMap)

	return order, nil
}

func (man *PosterService) fullFillProducts(
	req models.Order, store coreStoreModels.Store,
	productsMap map[string]coreMenuModels.Product, attributesMap map[string]coreMenuModels.Attribute,
	attributeGroupsMap map[string]coreMenuModels.AttributeGroup, aggregatorMenu coreMenuModels.Menu,
) (models.Order, float64, error) {

	var (
		serviceFee    float64
		orderProducts = make([]models.OrderProduct, 0, len(req.Products))
	)

	aggregatorProducts, aggregatorAttributes, _ := activeMenuPositions(aggregatorMenu)

	var cookingTime int32
	for _, product := range req.Products {
		var (
			orderProductAttributes = make([]models.ProductAttribute, 0, len(product.Attributes))
			modifiersPrice         float64
		)

		productPosID, exist := aggregatorProducts[product.ID]
		switch exist {
		case true:
			product.ID = productPosID
		default:
			log.Error().Msgf("product with ID %s %s not matched", product.ID, product.Name)
		}

		menuProduct, ok := productsMap[product.ID]
		if !ok {
			log.Info().Msgf("PRODUCT NOT FOUND IN POS MENU, ID %s, NAME %s", product.ID, product.Name)
			req.FailReason.Code = PRODUCT_MISSED_CODE
			req.FailReason.Message = PRODUCT_MISSED + product.ID
			return req, 0, errors.Wrap(errs.ErrProductNotFound, fmt.Sprintf("PRODUCT NOT FOUND IN POS MENU, ID %s, NAME %s", product.ID, product.Name))
		}

		var countCommonAttributes int

		if menuProduct.MenuDefaultAttributes != nil && len(menuProduct.MenuDefaultAttributes) > 0 {
			countCommonAttributes = len(menuProduct.DefaultAttributes)

			for _, defaultAttribute := range menuProduct.MenuDefaultAttributes {
				if defaultAttribute.DefaultAmount == 0 {
					defaultAttribute.DefaultAmount = 1
				}

				attribute, ok := attributesMap[defaultAttribute.ExtID]
				if !ok {
					log.Info().Msgf("attribute with ID %s %s not matched", defaultAttribute.ExtID, attribute.Name)
				}

				product.Attributes = append(product.Attributes, models.ProductAttribute{
					ID:       attribute.IngredientID,
					Name:     defaultAttribute.Name,
					Quantity: defaultAttribute.DefaultAmount,
					Price: models.Price{
						Value:        float64(defaultAttribute.Price),
						CurrencyCode: store.Settings.Currency,
					},
				})
			}
		}

		for index, attribute := range product.Attributes {

			if attribute.ID == models.ServiceFee {
				serviceFee += attribute.Price.Value * float64(product.Quantity) * float64(attribute.Quantity)

				if index < len(product.Attributes)-countCommonAttributes {
					modifiersPrice += attribute.Price.Value * float64(product.Quantity) * float64(attribute.Quantity)
				}
				continue
			}

			attributePosID, attributeExistInAggregatorMenu := aggregatorAttributes[attribute.ID]
			switch attributeExistInAggregatorMenu {
			case true:
				attribute.ID = attributePosID
			default:
				log.Info().Msgf("attribute with ID %s %s not matched", attribute.ID, attribute.Name)
			}

			menuAttribute, attributeExistInPosMenu := attributesMap[attribute.ID]
			if attributeExistInPosMenu {
				attributeGroupID := menuUtils.FindAttributeGroupID(product.ID, attribute.ID, productsMap, attributeGroupsMap, attributesMap, store.IikoCloud.IsExternalMenu, menuAttribute.ParentAttributeGroup)
				attribute.GroupID = attributeGroupID
				attribute.ID = menuAttribute.IngredientID
				if store.Settings.PriceSource == models.POSPriceSource {
					attribute.Price.Value = menuAttribute.Price
				}

				if index < len(product.Attributes)-countCommonAttributes {
					modifiersPrice += attribute.Price.Value * float64(product.Quantity) * float64(attribute.Quantity)
				}

				orderProductAttributes = append(orderProductAttributes, attribute)
				continue
			}

			menuAttributeProduct, productExistInPosMenu := productsMap[attribute.ID]
			if !productExistInPosMenu {
				req.FailReason.Code = ATTRIBUTE_MISSED_CODE
				req.FailReason.Message = ATTRIBUTE_MISSED + attribute.ID
				return req, 0, fmt.Errorf("ATTRIBUTE NOT FOUND IN POS MENU, ID %s, NAME %s", attribute.ID, attribute.Name)
			}

			orderProduct := models.OrderProduct{
				ID:       menuAttributeProduct.ProductID,
				Name:     attribute.Name,
				Quantity: attribute.Quantity,
				IsCombo:  menuAttributeProduct.IsCombo,
				Price: models.Price{
					Value:        attribute.Price.Value,
					CurrencyCode: store.Settings.Currency,
				},
			}

			if store.Settings.PriceSource == models.POSPriceSource {
				orderProduct.Price = models.Price{
					Value:        menuAttributeProduct.Price[0].Value,
					CurrencyCode: store.Settings.Currency,
				}
			}

			if index < len(product.Attributes)-countCommonAttributes {
				modifiersPrice += orderProduct.Price.Value * float64(product.Quantity) * float64(attribute.Quantity)
			}

			orderProducts = append(orderProducts, orderProduct)

			if store.PosType == models.Poster.String() {
				if menuProduct.CookingTime > cookingTime {
					cookingTime = menuProduct.CookingTime
				}
			}
		}

		switch store.Settings.PriceSource {
		case models.POSPriceSource:
			product.Price.Value = menuProduct.Price[0].Value
		default:
			if req.DeliveryService == models.YANDEX.String() {
				product.Price.Value = product.Price.Value - modifiersPrice/float64(product.Quantity)
			}
		}

		if menuProduct.ProductID != "" {
			product.ID = menuProduct.ProductID
		}

		product.IsCombo = menuProduct.IsCombo
		product.Attributes = orderProductAttributes
		orderProducts = append(orderProducts, product)

		if store.PosType == models.Poster.String() {
			if menuProduct.CookingTime > cookingTime {
				cookingTime = menuProduct.CookingTime
			}
		}
	}

	req.Products = orderProducts
	if store.PosType == models.Poster.String() {
		if cookingTime == 0 {
			cookingTime = store.Poster.CookingTime
		}
		req.CookingCompleteTime = time.Now().UTC().Add(time.Minute * time.Duration(cookingTime))
	}

	return req, serviceFee, nil
}

func (s *PosterService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	return nil
}

func (s *PosterService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (s *PosterService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (s *PosterService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
