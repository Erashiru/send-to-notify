package pos

import (
	"context"
	goErrors "errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/custom"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	"github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	yarosClient "github.com/kwaaka-team/orders-core/pkg/yaros"
	yarosConf "github.com/kwaaka-team/orders-core/pkg/yaros/clients"
	yarosModels "github.com/kwaaka-team/orders-core/pkg/yaros/models"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"regexp"
	"strconv"

	"github.com/kwaaka-team/orders-core/core/config"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
)

var (
	ErrRetry = goErrors.New("order is sent to retry")
)

type yarosService struct {
	*BasePosService
	notifyQueue     que.SQSInterface
	retryQueueName  string
	yarosCli        yarosConf.Yaros
	yarosStoreID    string
	yarosInfoSystem string
}

func newYarosService(bps *BasePosService, notifyQueue que.SQSInterface, retryQueueName, yarosStoreID, yarosInfoSystem, baseUrl,
	username, password string) (*yarosService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "Yaros Service constructor error")
	}
	client, err := yarosClient.NewClient(&yarosConf.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Yaros Client.")
		return nil, err
	}

	return &yarosService{
		BasePosService:  bps,
		notifyQueue:     notifyQueue,
		retryQueueName:  retryQueueName,
		yarosCli:        client,
		yarosStoreID:    yarosStoreID,
		yarosInfoSystem: yarosInfoSystem,
	}, nil
}

func (s *yarosService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {
	switch posStatus {
	case "rejected":
		return models.CANCELLED_BY_POS_SYSTEM, nil
	case "accepted":
		return models.ACCEPTED, nil
	case "created":
		return models.ACCEPTED, nil
	case "processing":
		return models.COOKING_STARTED, nil
	case "sent":
		return models.ON_WAY, nil
	case "modified":
		return models.NEW, nil
	case "completed":
		return models.COOKING_COMPLETE, nil
	default:
		return 0, models.StatusIsNotExist
	}
}

func isPhoneNumber(number string) bool {
	pattern := `^\+?[1-9]\d{1,14}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(number)
}

func (s *yarosService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration,
	store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {

	var err error
	phone := order.Customer.PhoneNumber

	order, err = prepareAnOrder(ctx, order, store, menu, aggregatorMenu, menuClient)
	if err != nil {
		return order, err
	}
	if isPhoneNumber(phone) {
		order.Customer.PhoneNumber = phone
	}

	posOrder, order, err := s.constructPosOrder(order, s.yarosInfoSystem, store)
	if err != nil {
		return order, validator.ErrCastingPos
	}

	order, err = s.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	responseOrder, err := s.sendOrder(ctx, posOrder, store)
	if err != nil {
		log.Err(err).Msgf("couldn't create order in YAROS pos")
		log.Info().Msgf("yaros error case: run RETRY")
		if err = s.orderRetry(ctx, order); err != nil {
			log.Err(err).Msgf("(CreateOrder - yarosService) error")
			return order, err
		}
		return order, ErrRetry
	}

	order.Status = responseOrder.Status
	order.CreationResult = models.CreationResult{
		Message: responseOrder.Message,
		OrderInfo: models.OrderInfo{
			CreationStatus: responseOrder.Status,
		},
		ErrorDescription: responseOrder.Message,
	}

	if len(responseOrder.Orders) > 0 {
		order = setPosOrderId(order, responseOrder.Orders[0].Id)
	}

	return order, nil
}

func (s *yarosService) orderRetry(ctx context.Context, order models.Order) error {
	if err := s.notifyQueue.SendSQSMessage(ctx, s.retryQueueName, fmt.Sprintf("yaros_%s", order.OrderID)); err != nil {
		log.Trace().Err(err).Msgf("SendSQSMessage creation-timeout error: %s", order.OrderID)
		return err
	}
	return nil
}

func (s *yarosService) sendOrder(ctx context.Context, order any, store coreStoreModels.Store) (yarosModels.OrderResponse, error) {
	var errs custom.Error

	posOrder, ok := order.(yarosModels.OrderRequest)
	if !ok {
		return yarosModels.OrderResponse{}, validator.ErrCastingPos
	}

	log.Info().Msgf("Yaros Request Body %+v", posOrder)

	createResponse, err := s.yarosCli.CreateOrder(ctx, store.Yaros.StoreId, posOrder)
	if err != nil {
		log.Err(err).Msg("yaros error")
		errs.Append(err, validator.ErrIgnoringPos)
		return yarosModels.OrderResponse{}, errs
	}
	return createResponse, nil
}

func (s *yarosService) constructPosOrder(order models.Order, infosystem string, store coreStoreModels.Store) (yarosModels.OrderRequest, models.Order, error) {
	timestamp := order.OrderTime.Value.Unix()

	var payMethod string
	switch order.PaymentMethod {
	case "DELAYED":
		payMethod = "cashless"
	case "CASH":
		payMethod = "cash"
	default:
		order.FailReason.Code = PAYMENT_TYPE_MISSED_CODE
		order.FailReason.Message = PAYMENT_TYPE_MISSED
		return yarosModels.OrderRequest{}, models.Order{}, nil
	}

	var yarosProducts = make([]yarosModels.OrderItem, 0, len(order.Products))

	for _, product := range order.Products {
		yarosProduct := yarosModels.OrderItem{
			ProductId: product.ID,
			Quantity:  strconv.Itoa(product.Quantity),
			Price:     strconv.Itoa(int(product.Price.Value)),
			Amount:    strconv.Itoa(product.Quantity * int(product.Price.Value)),
		}
		yarosProducts = append(yarosProducts, yarosProduct)
	}
	yarosOrder := yarosModels.OrderRequest{
		Orders: []yarosModels.PosOrder{
			{
				Id:              order.OrderID,
				Type:            "delivery",
				InfoSystem:      infosystem,
				DeliveryService: order.DeliveryService,
				OrderCode:       order.OrderCode,
				PickUpCode:      order.PickUpCode,
				Date:            strconv.Itoa(int(timestamp)),
				Change:          strconv.Itoa(int(order.MinimumBasketSurcharge.Value)),
				Total:           strconv.Itoa(int(order.EstimatedTotalPrice.Value)),
				Status:          "created", //надо узнать
				User: yarosModels.OrderUser{
					Name:  order.Customer.Name,
					Phone: order.Customer.PhoneNumber,
				},
				Address:   order.DeliveryAddress.Label,
				Comment:   s.constructComment(order.AllergyInfo, order.CutleryRequested),
				PayMethod: payMethod,
				Items:     yarosProducts,
			},
		},
	}

	if store.Yaros.Department != "" && len(yarosOrder.Orders) > 0 {
		yarosOrder.Orders[0].Department = store.Yaros.Department
	}

	return yarosOrder, order, nil
}

func (s *yarosService) constructComment(allergyInfo string, cutleryRequested bool) string {
	if !cutleryRequested {
		return allergyInfo
	}

	return "Положить приборы"
}

func (s *yarosService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {
	items, err := s.yarosCli.GetStopList(ctx, s.yarosStoreID)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}
	return s.toStopListItem(items.StopListItems), err
}

func (s *yarosService) toStopListItem(items []yarosModels.StopListItem) coreMenuModels.StopListItems {

	var stoplistItems coreMenuModels.StopListItems

	for _, item := range items {
		if item.Quantity > 0 {
			continue
		}
		stoplistItem := coreMenuModels.StopListItem{
			ProductID: item.Id,
			Balance:   float64(item.Quantity),
		}
		stoplistItems = append(stoplistItems, stoplistItem)
	}

	return stoplistItems
}

func (s *yarosService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	items, err := s.yarosCli.GetItems(ctx, s.yarosStoreID)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}
	categories, err := s.yarosCli.GetCategories(ctx, s.yarosStoreID)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}
	return menuFromClient(items, categories)
}

func menuFromClient(items yarosModels.GetItemsResponse, categories yarosModels.GetCategoriesResponse) (coreMenuModels.Menu, error) {
	otherCollection := coreMenuModels.MenuCollection{
		ExtID: uuid.New().String(),
		Name:  "Other",
	}

	otherSection := coreMenuModels.Section{
		ExtID:      uuid.New().String(),
		Collection: otherCollection.ExtID,
		Name:       "Other",
	}

	sections, collections := toSectionCollection(categories.Categories)

	products, otherSectionUsed, err := toProduct(items.Items, sections, otherSection)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}
	if otherSectionUsed {
		sections = append(sections, otherSection)
		collections = append(collections, otherCollection)
	}

	return coreMenuModels.Menu{
		Collections: collections,
		Name:        models.Yaros.String(),
		CreatedAt:   models.TimeNow(),
		UpdatedAt:   models.TimeNow(),
		Sections:    sections,
		Products:    products,
	}, nil
}

func toProduct(items []yarosModels.Item, sections coreMenuModels.Sections, otherSection coreMenuModels.Section) (coreMenuModels.Products, bool, error) {
	products := make(coreMenuModels.Products, 0, len(items))

	var otherSectionUsed bool

	sectionsMap := make(map[string]coreMenuModels.Section)
	for _, section := range sections {
		sectionsMap[section.ExtID] = section
	}

	for _, item := range items {
		if item.ImageUrl == "" {
			item.ImageUrl = "https://kwaaka-menu-files.s3.eu-west-1.amazonaws.com/images/default_image_for_product/7b05c838-2f84-423c-bc61-af2ab55d50c3.jpg"
		}
		product := coreMenuModels.Product{
			ExtID:            item.Id,
			PosID:            item.Id,
			Section:          item.CategoryId,
			IsIncludedInMenu: true,
			Name: []coreMenuModels.LanguageDescription{
				{
					Value:        item.Title,
					LanguageCode: "ru",
				},
			},
			ImageURLs: []string{item.ImageUrl},
			Description: []coreMenuModels.LanguageDescription{
				{
					Value:        item.Description,
					LanguageCode: "ru",
				},
			},
			MeasureUnit: item.Measure,
		}
		if item.Price == "" {
			item.Price = "0"
		}
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return coreMenuModels.Products{}, false, err
		}

		product.Price = []coreMenuModels.Price{
			{
				Value:        price,
				CurrencyCode: "", // what currency?
			},
		}
		if _, ok := sectionsMap[item.CategoryId]; !ok {
			product.Section = otherSection.ExtID
			otherSectionUsed = true
		}
		products = append(products, product)
	}
	return products, otherSectionUsed, nil
}

func toSectionCollection(categories []yarosModels.Category) (coreMenuModels.Sections, coreMenuModels.MenuCollections) {
	collections := make([]coreMenuModels.MenuCollection, 0, 4)
	sections := make([]coreMenuModels.Section, 0, 4)

	for _, category := range categories {
		if category.ParentId == "" {
			collections = append(collections, coreMenuModels.MenuCollection{
				ExtID:           category.Id,
				Name:            category.Title,
				CollectionOrder: category.SortPriority,
			})
			continue
		}
		sections = append(sections, coreMenuModels.Section{
			Collection:   category.ParentId,
			Name:         category.Title,
			ExtID:        category.Id,
			SectionOrder: category.SortPriority,
		})
	}
	return sections, collections
}

func (s *yarosService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	return nil
}

func (s *yarosService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (s *yarosService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (s *yarosService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
