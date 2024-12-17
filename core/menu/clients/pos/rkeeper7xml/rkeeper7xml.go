package rkeeper_xml

import (
	"context"
	"errors"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/base"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	rkeeperXMLCli "github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml"
	rkeeperXMLConf "github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients"
	rkeeperXMLModels "github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/trade_group_details_response"
	"strconv"
)

type manager struct {
	cfg      menu.Configuration
	cli      rkeeperXMLConf.RKeeper7
	menuRepo drivers.MenuRepository
}

func NewManager(cfg menu.Configuration, menuRepo drivers.MenuRepository, store storeModels.Store) (base.Manager, error) {
	cli, err := rkeeperXMLCli.NewClient(&rkeeperXMLConf.Config{
		Protocol:               "http",
		BaseURL:                store.RKeeper7XML.Domain,
		Username:               store.RKeeper7XML.Username,
		Password:               store.RKeeper7XML.Password,
		UCSUsername:            store.RKeeper7XML.UCSUsername,
		UCSPassword:            store.RKeeper7XML.UCSPassword,
		Token:                  store.RKeeper7XML.Token,
		Anchor:                 store.RKeeper7XML.Anchor,
		ObjectID:               store.RKeeper7XML.ObjectID,
		LicenseBaseURL:         cfg.RKeeper7XMLConfiguration.LicenseBaseURL,
		StationID:              store.RKeeper7XML.StationID,
		StationCode:            store.RKeeper7XML.StationCode,
		LicenseInstanceGUID:    store.RKeeper7XML.LicenseInstanceGUID,
		ChildItems:             store.RKeeper7XML.ChildItems,
		ClassificatorItemIdent: store.RKeeper7XML.ClassificatorItemIdent,
		ClassificatorPropMask:  store.RKeeper7XML.ClassificatorPropMask,
		MenuItemsPropMask:      store.RKeeper7XML.MenuItemsPropMask,
		PropFilter:             store.RKeeper7XML.PropFilter,
		Cashier:                store.RKeeper7XML.Cashier,
	})
	if err != nil {
		return nil, err
	}

	return &manager{
		cli:      cli,
		cfg:      cfg,
		menuRepo: menuRepo,
	}, nil
}

func (man manager) deleteNonDeliveryItems(menuItems rkeeperXMLModels.MenuRK7QueryResult, tradeGroupEntities trade_group_details_response.RK7QueryResult) rkeeperXMLModels.MenuRK7QueryResult {
	tradeGroupItemsMap := make(map[string]struct{})

	for _, item := range tradeGroupEntities.RK7Reference.Items.Item {
		//if item.Flag != "gdfAdd" || item.TradeObject != "toDish" || item.RefCollName != "MENUITEMS" {
		//	continue
		//}

		tradeGroupItemsMap[item.ObjectSifr] = struct{}{}
	}

	items := make([]rkeeperXMLModels.Item, 0, len(tradeGroupEntities.RK7Reference.Items.Item))

	for _, item := range menuItems.RK7Reference.Items.Item {
		if _, ok := tradeGroupItemsMap[item.Ident]; !ok {
			continue
		}

		items = append(items, item)
	}

	menuItems.RK7Reference.Items.Item = items

	return menuItems
}

func (man manager) GetAggMenu(ctx context.Context, store storeModels.Store) ([]models.Menu, error) {
	return nil, errors.New("method not implemented")
}

func (man manager) GetMenu(ctx context.Context, store storeModels.Store) (models.Menu, error) {
	menuItems, err := man.cli.GetMenuItems(ctx)
	if err != nil {
		return models.Menu{}, err
	}

	tradeGroupEntities, err := man.cli.GetItemsByTradeGroup(ctx, store.RKeeper7XML.TradeGroupId)
	if err != nil {
		return models.Menu{}, err
	}

	menuModifiers, err := man.cli.GetMenuModifiers(ctx)
	if err != nil {
		return models.Menu{}, err
	}

	menuModifierGroups, err := man.cli.GetMenuModifierGroups(ctx)
	if err != nil {
		return models.Menu{}, err
	}

	menuModifierSchemaDetails, err := man.cli.GetMenuModifierSchemaDetails(ctx)
	if err != nil {
		return models.Menu{}, err
	}

	menuModifierSchemas, err := man.cli.GetMenuModifierSchemas(ctx)
	if err != nil {
		return models.Menu{}, err
	}

	menuItems = man.deleteNonDeliveryItems(menuItems, tradeGroupEntities)

	orderMenu, err := man.cli.GetOrderMenu(ctx)
	if err != nil {
		return models.Menu{}, err
	}

	mappingProducts := make(map[string]string)
	mappingAttributes := make(map[string]string)

	deliveryItems, err := man.cli.GetDeliveryPriceByPriceTypeId(ctx, store.RKeeper7XML.PriceTypeId)
	if err != nil {
		return models.Menu{}, err
	}

	entityPriceMap := map[string]string{}

	for _, deliveryItem := range deliveryItems.RK7Reference.Items.Item {
		if deliveryItem.Species != "psModifier" && deliveryItem.Species != "psDish" {
			continue
		}

		entityPriceMap[deliveryItem.ObjectID] = deliveryItem.Value
	}

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

		price := dish.Price

		if val, ok := entityPriceMap[dish.Ident]; ok {
			price = val
		}

		mappingProducts[dish.Ident] = price
	}

	for _, modifier := range orderMenu.Modifiers.Item {
		price := modifier.Price

		if val, ok := entityPriceMap[modifier.Ident]; ok {
			price = val
		}

		mappingAttributes[modifier.Ident] = price
	}

	posMenu, err := man.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(store.MenuID))
	if err != nil {
		return models.Menu{}, err
	}

	existProducts := make(map[string]models.Product)

	for _, product := range posMenu.Products {
		existProducts[product.ProductID] = product
	}

	return man.menuFromClient(menuItems, menuModifiers, store.Settings, mappingProducts, mappingAttributes, menuModifierGroups, menuModifierSchemas, menuModifierSchemaDetails, existProducts, entityPriceMap), nil
}
