package pos

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/config"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/jowi"
	jowiClient "github.com/kwaaka-team/orders-core/pkg/jowi/client"
	jowiDto "github.com/kwaaka-team/orders-core/pkg/jowi/client/dto"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

type jowiService struct {
	*BasePosService
	jowiCli      jowi.Jowi
	restaurantID string
}

func newJowiService(bps *BasePosService, baseUrl, apiKey, apiSecret, restaurantID string) (*jowiService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "burgerKingService constructor error")
	}
	client, err := jowiClient.New(jowi.Config{
		Protocol:  "http",
		BaseURL:   baseUrl,
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Jowi Client.")
		return nil, err
	}
	return &jowiService{
		jowiCli:        client,
		BasePosService: bps,
		restaurantID:   restaurantID,
	}, nil
}

func (jowiSvc *jowiService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	courses, err := jowiSvc.jowiCli.GetCourses(ctx, jowiSvc.restaurantID)
	if err != nil {
		log.Trace().Err(err).Msgf("Get courses error")
		return coreMenuModels.Menu{}, err
	}

	sections, collections, err := jowiSvc.getCategories(ctx)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	stopList, err := jowiSvc.GetStopList(ctx)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	return jowiSvc.menuFromClient(courses, store.Settings, sections, collections, stopList), nil
}

func (jowiSvc *jowiService) getCategories(ctx context.Context) ([]coreMenuModels.Section, []coreMenuModels.MenuCollection, error) {
	courseCategories, err := jowiSvc.jowiCli.GetCourseCategories(ctx, jowiSvc.restaurantID)
	if err != nil {
		log.Trace().Err(err).Msgf("Get course categories error")
		return nil, nil, err
	}

	return jowiSvc.categoryFromClient(courseCategories)
}

func (jowiSvc *jowiService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {
	items, err := jowiSvc.jowiCli.GetStopList(ctx, jowiSvc.restaurantID)
	if err != nil {
		log.Trace().Err(err).Msgf("Get stop list error")
		return coreMenuModels.StopListItems{}, err
	}

	return jowiSvc.stopListFromClient(items)
}

func (jowiSvc *jowiService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {
	switch posStatus {
	case "WAIT_SENDING":
		return models.WAIT_SENDING, nil
	case "0":
		return models.NEW, nil
	case "1":
		return models.ACCEPTED, nil
	case "2":
		return models.CANCELLED_BY_POS_SYSTEM, nil
	case "3":
		return models.OUT_FOR_DELIVERY, nil
	case "4":
		return models.CLOSED, nil
	}
	return 0, models.StatusIsNotExist
}

func (jowiSvc *jowiService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration, store coreStoreModels.Store,
	menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	posOrder, _, err := jowiSvc.constructPosOrder(order, store)

	if err != nil {
		return order, err
	}

	order, err = jowiSvc.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	responseOrder, err := jowiSvc.sendOrder(ctx, posOrder)
	if err != nil {
		order.FailReason.Code = OTHER_FAIL_REASON_CODE
		order.FailReason.Message = responseOrder.ErrorResponse.Message
		return order, err
	}

	order = setPosOrderId(order, responseOrder.Order.Id)

	order.CreationResult = models.CreationResult{
		Message: responseOrder.Message,
		OrderInfo: models.OrderInfo{
			CreationStatus: strconv.Itoa(responseOrder.Order.Status),
			OrganizationID: responseOrder.Order.RestaurantId,
		},
		ErrorDescription: responseOrder.ErrorResponse.Message,
	}

	return order, nil
}

func (jowiSvc *jowiService) sendOrder(ctx context.Context, posOrder jowiDto.RequestCreateOrder) (jowiDto.ResponseOrder, error) {
	log.Info().Msgf("Jowi Request Body: %+v", posOrder)
	createResponse, err := jowiSvc.jowiCli.CreateOrder(ctx, posOrder)
	if err != nil {
		log.Err(err).Msg("Jowi create order error")
		return jowiDto.ResponseOrder{}, err
	}
	return createResponse, nil
}

func (jowiSvc *jowiService) constructPosOrder(order models.Order, store coreStoreModels.Store) (jowiDto.RequestCreateOrder, models.Order, error) {
	jowiOrder := jowiDto.RequestCreateOrder{
		RestaurantID: store.Token,
		Order: jowiDto.RequestCreateOrderBody{
			Address:     order.DeliveryAddress.Label,
			Phone:       order.Customer.PhoneNumber,
			Contact:     order.Customer.Name,
			AmountOrder: strconv.Itoa(int(order.EstimatedTotalPrice.Value)),
		},
	}

	if order.Persons != 0 {
		jowiOrder.Order.PeopleCount = order.Persons
	}

	var courses = make([]jowiDto.RequestCreateOrderCourse, 0)
	for _, product := range order.Products {
		courses = append(courses, jowiDto.RequestCreateOrderCourse{
			CourseId: product.ID,
			Count:    product.Quantity,
			Price:    int(product.Price.Value),
		})

		for _, attribute := range product.Attributes {
			courses = append(courses, jowiDto.RequestCreateOrderCourse{
				CourseId: attribute.ID,
				Count:    attribute.Quantity,
				Price:    int(attribute.Price.Value),
			})
		}
	}

	switch order.PaymentMethod {
	case "CASH":
		jowiOrder.Order.PaymentType = 0
		jowiOrder.Order.PaymentMethod = 0
	case "DELAYED":
		jowiOrder.Order.PaymentType = 1
		jowiOrder.Order.PaymentMethod = 1
	}

	var isMarketplace bool
	switch order.DeliveryService {
	case "glovo":
		isMarketplace = store.Glovo.IsMarketplace
	case "wolt":
		isMarketplace = store.Wolt.IsMarketplace
	case "chocofood":
		isMarketplace = store.Chocofood.IsMarketplace
	case "qr_menu":
		isMarketplace = store.QRMenu.IsMarketplace
	default:
		for _, deliveryService := range store.ExternalConfig {
			if deliveryService.Type == order.DeliveryService {
				isMarketplace = deliveryService.IsMarketplace
			}
		}
	}

	jowiOrder.Order.OrderType = 1
	if isMarketplace {
		jowiOrder.Order.OrderType = 0
	}

	jowiOrder.Order.Courses = courses

	return jowiOrder, order, nil
}

func (jowiSvc *jowiService) menuFromClient(req jowiDto.ResponseCourse, store coreStoreModels.Settings, sections []coreMenuModels.Section, collections []coreMenuModels.MenuCollection, stopList coreMenuModels.StopListItems) coreMenuModels.Menu {

	menu := coreMenuModels.Menu{
		Products: jowiSvc.getProducts(req, store),
	}

	menu.Sections = sections
	menu.Collections = collections
	menu.StopLists = stopList.Products()

	return menu
}

func (jowiSvc *jowiService) categoryFromClient(req jowiDto.ResponseCourseCategory) ([]coreMenuModels.Section, []coreMenuModels.MenuCollection, error) {
	var (
		sections    []coreMenuModels.Section
		collections []coreMenuModels.MenuCollection
	)

	var hasCollection bool
	for _, courseCategory := range req.CourseCategories {
		if courseCategory.ParentId != "" {
			hasCollection = true
		}
	}

	for _, courseCategory := range req.CourseCategories {
		if (hasCollection && courseCategory.ParentId != "") || !hasCollection {
			sections = append(sections, coreMenuModels.Section{
				ExtID:      courseCategory.Id,
				Name:       courseCategory.Title,
				Collection: courseCategory.ParentId,
			})
			continue
		}

		collections = append(collections, coreMenuModels.MenuCollection{
			ExtID: courseCategory.Id,
			Name:  courseCategory.Title,
		})
	}

	return sections, collections, nil
}

func (jowiSvc *jowiService) stopListFromClient(req jowiDto.ResponseStopList) (coreMenuModels.StopListItems, error) {
	var stopList coreMenuModels.StopListItems

	for _, course := range req.CourseCounts {
		count, err := strconv.ParseFloat(course.Count, 64)
		if err != nil {
			log.Err(err).Msgf("couldn't convert count string into float, %T", course.Count)
			return coreMenuModels.StopListItems{}, err
		}

		stopList = append(stopList, coreMenuModels.StopListItem{
			ProductID: course.Id,
			Balance:   count,
		})
	}

	return stopList, nil
}

func (jowiSvc *jowiService) getProducts(req jowiDto.ResponseCourse, settings coreStoreModels.Settings) coreMenuModels.Products {

	products := make(coreMenuModels.Products, 0, len(req.Courses))

	for _, product := range req.Courses {
		resProduct := jowiSvc.productToModel(product, settings)
		products = append(products, resProduct)
	}

	return products
}

func (jowiSvc *jowiService) productToModel(req jowiDto.Course, setting coreStoreModels.Settings) coreMenuModels.Product {

	price, err := strconv.ParseFloat(req.PriceForOnlineOrder, 64)
	if err != nil {
		price = 0
	}

	weight, err := strconv.ParseFloat(req.Weight, 64)
	if err != nil {
		weight = 0
	}

	product := coreMenuModels.Product{
		ExtID:            req.Id,
		ProductID:        req.Id,
		Section:          req.CourseCategoryId, // req.Section is product categories in iiko, but here we used linked groups
		ExtName:          req.Title,
		Code:             req.PackageCode,
		ImageURLs:        []string{req.ImageUrl},
		Weight:           weight,
		MeasureUnit:      req.UnitName,
		IsAvailable:      req.OnlineOrder,
		IsIncludedInMenu: req.OnlineOrder,
		ProductsCreatedAt: coreMenuModels.ProductsCreatedAt{
			Value:     models.TimeNow(),
			Timezone:  setting.TimeZone.TZ,
			UTCOffset: setting.TimeZone.UTCOffset,
		},
		Name: []coreMenuModels.LanguageDescription{
			{
				Value:        req.Title,
				LanguageCode: setting.LanguageCode,
			},
		},
		Price: []coreMenuModels.Price{
			{
				Value:        price,
				CurrencyCode: setting.Currency,
			},
		},
		UpdatedAt: models.TimeNow(),
	}

	return product
}

func (jowiSvc *jowiService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	_, err := jowiSvc.jowiCli.CancelOrder(ctx, order.PosOrderID, jowiDto.RequestCancelOrder{}) //TODO: unused function before
	if err != nil {
		return err
	}

	return nil
}

func (s *jowiService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (s *jowiService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (s *jowiService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
