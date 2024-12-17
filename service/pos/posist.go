package pos

import (
	"context"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/config"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"strings"

	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	"github.com/kwaaka-team/orders-core/pkg/posist"
	"github.com/kwaaka-team/orders-core/pkg/posist/clients"
	posistModels "github.com/kwaaka-team/orders-core/pkg/posist/clients/models"
	"github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/pkg/errors"
)

type posistPosService struct {
	*BasePosService
	posistCli   clients.Posist
	customerKey string
	tabId       string
}

func newPosistPosService(bps *BasePosService, baseUrl, authBasic, customerKey, tabId string) (*posistPosService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "posist pos constructor error")
	}

	posistClient, err := posist.New(&clients.Config{
		Protocol:    "http",
		BaseURL:     baseUrl,
		AuthBasic:   authBasic,
		CustomerKey: customerKey,
	})
	if err != nil {
		return nil, err
	}

	return &posistPosService{
		BasePosService: bps,
		posistCli:      posistClient,
		customerKey:    customerKey,
		tabId:          tabId,
	}, nil
}

func (p *posistPosService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (coreOrderModels.PosStatus, error) {
	switch posStatus {
	case "synced":
		return coreOrderModels.ACCEPTED, nil
	case "preparing":
		return coreOrderModels.COOKING_STARTED, nil
	case "out for delivery":
		return coreOrderModels.ON_WAY, nil
	case "payment received":
		return coreOrderModels.CLOSED, nil
	case "canceled":
		return coreOrderModels.CANCELLED_BY_POS_SYSTEM, nil
	}
	return 0, coreOrderModels.StatusIsNotExist
}

func (p *posistPosService) toItems(req coreOrderModels.Order) []posistModels.OrderItem {
	orderItems := make([]posistModels.OrderItem, 0, len(req.Products))

	for _, product := range req.Products {
		addons := make([]posistModels.AddOn, 0, len(product.Attributes))
		comboItems := make([]posistModels.ComboItem, 0, len(product.Attributes))

		for _, attribute := range product.Attributes {
			switch attribute.IsComboAttribute {
			case true:
				comboItems = append(comboItems, posistModels.ComboItem{
					Id:       attribute.ID,
					Quantity: attribute.Quantity * product.Quantity,
				})
			default:
				addons = append(addons, posistModels.AddOn{
					Id:       attribute.ID,
					Quantity: attribute.Quantity * product.Quantity,
				})
			}
		}

		orderItems = append(orderItems, posistModels.OrderItem{
			Id:            product.ID,
			Quantity:      product.Quantity,
			Rate:          int(product.Price.Value),
			AddOns:        addons,
			MapComboItems: comboItems,
		})
	}

	return orderItems
}

func (p *posistPosService) getCustomerInfo(req coreOrderModels.Order, store coreStoreModels.Store) posistModels.Customer {
	return posistModels.Customer{
		Firstname: req.Customer.Name,
		Mobile:    "+77771111111",
		AddType:   "home",
		Address1:  "N/A",
		Address2:  "N/A",
		City:      store.Address.City,
	}
}

func (p *posistPosService) constructPosOrder(req coreOrderModels.Order, store coreStoreModels.Store) (coreOrderModels.Order, posistModels.Order, error) {
	posOrderId := uuid.New().String()

	req = setPosOrderId(req, posOrderId)

	items := p.toItems(req)
	customerInfo := p.getCustomerInfo(req, store)

	order := posistModels.Order{
		Source: posistModels.Source{
			Name:    "Posist",
			Id:      "Syx3IXCI",
			OrderId: posOrderId,
		},
		Customer: customerInfo,
		Items:    items,
		TabType:  "delivery",
		Payments: posistModels.Payments{
			Type: "COD",
		},
	}

	return req, order, nil
}

func (p *posistPosService) CreateOrder(ctx context.Context, order coreOrderModels.Order, globalConfig config.Configuration,
	store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli store.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (coreOrderModels.Order, error) {
	var err error

	order, err = prepareAnOrder(ctx, order, store, menu, aggregatorMenu, menuClient)
	if err != nil {
		return order, err
	}

	utils.Beautify("prepared order", order)

	var posOrder posistModels.Order

	order, posOrder, err = p.constructPosOrder(order, store)
	if err != nil {
		return order, err
	}

	order, err = p.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	utils.Beautify("pos order body", posOrder)

	err = p.posistCli.CreateOrder(ctx, p.customerKey, posOrder)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (p *posistPosService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {
	var stopListItems coreMenuModels.StopListItems

	items, err := p.posistCli.GetStopList(ctx, p.customerKey)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}

	for _, item := range items {
		stopListItems = append(stopListItems, coreMenuModels.StopListItem{
			ProductID: item.Id,
		})
	}

	return stopListItems, nil
}

func (p *posistPosService) toEntities(ctx context.Context, posMenu posistModels.Menu, store coreStoreModels.Store) ([]coreMenuModels.Product, []coreMenuModels.Section, []coreMenuModels.MenuCollection, []coreMenuModels.AttributeGroup, []coreMenuModels.Attribute, error) {
	systemProducts := make([]coreMenuModels.Product, 0)
	sections := make([]coreMenuModels.Section, 0)
	collections := make([]coreMenuModels.MenuCollection, 0)
	attributeGroups := make([]coreMenuModels.AttributeGroup, 0, len(posMenu.Modifiers))
	attributes := make([]coreMenuModels.Attribute, 0)

	uniqueProducts := make(map[string]struct{})
	uniqueSections := make(map[string]struct{})
	uniqueCollections := make(map[string]struct{})
	uniqueAttributes := make(map[string]struct{})

	for _, collection := range posMenu.Categories {

		// collection
		if _, ok := uniqueCollections[collection.ID]; ok {
			continue
		} else {
			uniqueCollections[collection.ID] = struct{}{}

			collections = append(collections, coreMenuModels.MenuCollection{
				ExtID:           collection.ID,
				Name:            collection.Name,
				CollectionOrder: collection.Order,
			})
		}

		// category

		for _, category := range collection.SubCategories {
			if _, exist := uniqueSections[category.ID]; exist {
				continue
			} else {
				uniqueSections[category.ID] = struct{}{}

				sections = append(sections, coreMenuModels.Section{
					ExtID:        category.ID,
					Name:         category.Name,
					SectionOrder: category.Order,
				})
			}

			for _, item := range category.Entities {
				if !item.IsActive {
					continue
				}

				if _, ok := uniqueProducts[item.ID]; ok {
					continue
				} else {
					uniqueProducts[item.ID] = struct{}{}
				}

				imageUrl := item.ImageURL

				if len(item.AggregatorImage) > 0 {
					if item.AggregatorImage[0].Jpg != "" {
						imageUrl = item.AggregatorImage[0].Jpg
					}
				}

				systemProduct := coreMenuModels.Product{
					ExtID:     item.ID,
					PosID:     item.ID,
					ProductID: item.ID,
					Description: []coreMenuModels.LanguageDescription{
						{
							LanguageCode: store.Settings.LanguageCode,
							Value:        item.Description,
						},
					},
					Name: []coreMenuModels.LanguageDescription{
						{
							LanguageCode: store.Settings.LanguageCode,
							Value:        item.Name,
						},
					},
					Price: []coreMenuModels.Price{
						{
							Value:        item.Price,
							CurrencyCode: store.Settings.Currency,
						},
					},
					Section:          category.ID,
					ImageURLs:        []string{imageUrl},
					IsAvailable:      true,
					IsIncludedInMenu: item.IsActive,
					AttributesGroups: item.Modifiers,
				}

				for _, modifierGroup := range item.Modifiers {
					if strings.Contains(modifierGroup, "combo") {
						systemProduct.IsCombo = true
						break
					}
				}

				systemProducts = append(systemProducts, systemProduct)
			}
		}
	}

	for _, modifier := range posMenu.Modifiers {
		if !modifier.IsActive {
			continue
		}

		attributeIds := make([]string, 0, len(modifier.ConstituentItems))

		isCombo := false

		switch modifier.Type {
		case "Combo":
			isCombo = true
		}

		for _, consItem := range modifier.ConstituentItems {
			if !consItem.IsActive {
				continue
			}

			attributeIds = append(attributeIds, consItem.ID)

			if _, exist := uniqueAttributes[consItem.ID]; exist {
				continue
			}

			attributes = append(attributes, coreMenuModels.Attribute{
				ExtID:       consItem.ID,
				Name:        consItem.Name,
				Price:       consItem.Price,
				IsAvailable: consItem.IsActive,
				Description: []coreMenuModels.LanguageDescription{
					{
						Value:        consItem.Description,
						LanguageCode: store.Settings.LanguageCode,
					},
				},
				IsComboAttribute: isCombo,
			})
		}

		attributeGroups = append(attributeGroups, coreMenuModels.AttributeGroup{
			Min:        modifier.Min,
			Max:        modifier.Max,
			Name:       modifier.Name,
			ExtID:      modifier.ID,
			Attributes: attributeIds,
		})
	}

	return systemProducts, sections, collections, attributeGroups, attributes, nil
}

func (p *posistPosService) GetOrderStatus(ctx context.Context, order coreOrderModels.Order) (string, error) {
	orderStatus, err := p.posistCli.GetOrderStatus(ctx, order.PosOrderID)
	if err != nil {
		return "", err
	}

	return orderStatus.Status, nil
}

func (p *posistPosService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	posMenu, err := p.posistCli.GetMenu(ctx, p.customerKey, p.tabId)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	systemProducts, systemSections, systemCollections, systemAttributeGroups, systemAttributes, err := p.toEntities(ctx, posMenu, store)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	systemMenu := coreMenuModels.Menu{
		Name:             "posist pos menu",
		Products:         systemProducts,
		Sections:         systemSections,
		Collections:      systemCollections,
		Attributes:       systemAttributes,
		AttributesGroups: systemAttributeGroups,
	}

	return systemMenu, nil
}

func (p *posistPosService) CancelOrder(ctx context.Context, order coreOrderModels.Order, store coreStoreModels.Store) error {
	//TODO implement me
	panic("implement me")
}

func (s *posistPosService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (s *posistPosService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (s *posistPosService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
