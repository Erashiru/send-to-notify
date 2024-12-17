package stoplist

import (
	"context"
	"fmt"
	"sync"

	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func (s *ServiceImpl) ActualizeStopListByStoreID(ctx context.Context, storeID string) error {
	store, err := s.storeService.GetByID(ctx, storeID)
	if err != nil {
		return err
	}

	return s.actualizeStopList(ctx, store)
}

func (s *ServiceImpl) ActualizeStopListByToken(ctx context.Context, token string) error {
	stores, err := s.storeService.GetStoresByToken(ctx, token)
	if err != nil {
		return err
	}
	fmt.Println("get stores by token")
	mockstore := storeModels.Store{
		ID:   "01",
		Name: "mock",
		Telegram: storeModels.StoreTelegramConfig{
			TelegramBotToken: "7241041762:AAHGzpJExb_KUKdbv1fhubXi95ZnMd3gcU0",
		},
	}
	stoplistprod := menuModels.StopListProduct{
		ExtID:       "01",
		IsAvailable: true,
		Price:       1,
	}
	stoplistprods := menuModels.StopListProducts{
		stoplistprod,
	}
	attributes := menuModels.StopListAttributes{
		menuModels.StopListAttribute{
			AttributeID:   "01",
			AttributeName: "mock",
			IsAvailable:   true,
			Price:         1,
		},
	}
	transactionData := menuModels.TransactionData{
		StoreID:    "01",
		Delivery:   "Delivery",
		Products:   stoplistprods,
		Attributes: attributes,
	}
	s.sendToNotify(ctx, mockstore, transactionData)
	fmt.Println("send to notify")
	for i := range stores {
		store := stores[i]
		if err = s.actualizeStopList(ctx, store); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) ActualizeStopListByPosTypes(ctx context.Context, posTypes []string) error {
	if len(posTypes) == 0 {
		return nil
	}

	stores := make([]storeModels.Store, 0)
	listOfStores := make([]storeModels.Store, 0)
	for i := range posTypes {
		posType := posTypes[i]
		storesByPos, err := s.storeService.FindStoresByPosType(ctx, posType)
		if err != nil {
			return err
		}
		listOfStores = append(listOfStores, storesByPos...)
	}

	if len(listOfStores) == 0 {
		return nil
	}

	for _, st := range listOfStores {
		if !st.Settings.IgnoreUpdateStopList {
			stores = append(stores, st)
		}
	}

	var g errgroup.Group

	for i := range stores {
		store := stores[i]

		g.Go(func() error {
			return s.actualizeStopList(ctx, store)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) ActualizeStopListByPosType(ctx context.Context, posType string) error {
	stores, err := s.storeService.FindStoresByPosType(ctx, posType)
	if err != nil {
		return err
	}

	jobs := make(chan storeModels.Store)

	var wg sync.WaitGroup
	for i := 0; i < s.concurrencyLevel; i++ {
		wg.Add(1)
		go s.startWorker(ctx, jobs, &wg)
	}

	for i := range stores {
		store := stores[i]
		log.Info().Msgf("start processing. Total stores = %d, current = %d", len(stores), i)
		jobs <- store
	}

	close(jobs)
	wg.Wait()

	return nil
}

func (s *ServiceImpl) ActualizeStoplistbyYarosStoreID(ctx context.Context, storeID string) error {
	store, err := s.storeService.GetByYarosRestaurantID(ctx, storeID)
	if err != nil {
		return err
	}

	return s.actualizeStopList(ctx, store)
}

func (s *ServiceImpl) startWorker(ctx context.Context, jobs chan storeModels.Store, wg *sync.WaitGroup) {
	defer wg.Done()

	for store := range jobs {
		if err := s.actualizeStopList(ctx, store); err != nil {
			log.Err(err).Msgf("error while actualize stop list by pos type %s, store_id = %s", store.PosType, store.ID)
		}
	}
}

func (s *ServiceImpl) actualizeStopList(ctx context.Context, store storeModels.Store) error {
	log.Info().Msgf("actualize stop list for store_id = %s, store_name = %s", store.ID, store.Name)
	posService, err := s.posFactory.GetPosService(models.Pos(store.PosType), store)
	if err != nil {
		return err
	}

	stopListItems, err := posService.GetStopList(ctx)
	if err != nil {
		return err
	}

	posMenu, err := s.menuService.GetMenuById(ctx, store.MenuID)
	if err != nil {
		return err
	}

	log.Info().Msgf("pos menu %s items for stoplist %v", store.MenuID, stopListItems.Products())

	stopListItems, err = posService.SortStoplistItemsByIsIgnored(ctx, posMenu, stopListItems)
	if err != nil {
		return err
	}

	log.Info().Msgf("stoplist after sorting %v", stopListItems.Products())

	isByBalance := posService.IsStopListByBalance(ctx, store)
	balanceLimit := float64(posService.GetBalanceLimit(ctx, store))

	newStopList, err := s.actualizeStopListPosMenu(ctx, store, stopListItems, isByBalance, balanceLimit)
	if err != nil {
		return err
	}

	for _, menu := range store.Menus {
		if err = s.actualizeStopListAggregatorMenu(ctx, store, menu, newStopList, isByBalance, balanceLimit, stopListItems); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) actualizeStopListPosMenu(ctx context.Context, store storeModels.Store, stopListItems menuModels.StopListItems, isByBalance bool, balanceLimit float64) (menuModels.StopListItems, error) {

	products, attributes, err := s.getItemsForUpdatePosMenu(ctx, store.MenuID, stopListItems, isByBalance, balanceLimit)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 && len(attributes) == 0 {
		return nil, fmt.Errorf("nothing to update in pos-menu-stoplist, store_id = %s, store_name = %s", store.ID, store.Name)
	}

	menuID := store.MenuID
	if err = s.menuService.UpdateStopList(ctx, menuID, stopListItems.Products()); err != nil {
		return nil, err
	}

	if len(products) != 0 {
		if err = s.updateStopListByProductIDsInDatabase(ctx, menuID, products); err != nil {
			return nil, err
		}
	}

	if len(attributes) != 0 {
		if err = s.updateStopListByAttributesInDatabase(ctx, menuID, attributes); err != nil {
			return nil, err
		}
	}

	stopListForAggMenu := s.getStopListItemsForAggregatorMenu(stopListItems, products, attributes)

	return stopListForAggMenu, nil
}

func (s *ServiceImpl) getStopListItemsForAggregatorMenu(stopListItems []menuModels.StopListItem, products []menuModels.Product, attributes []menuModels.Attribute) []menuModels.StopListItem {
	if len(stopListItems) == 0 {
		return stopListItems
	}

	items := make(map[string][]string)
	for _, product := range products {
		i := product.ProductID
		if extIDs, ok := items[i]; ok {
			extIDs = append(extIDs, product.ExtID)
			items[i] = extIDs
		} else {
			extIDs = make([]string, 0)
			extIDs = append(extIDs, product.ExtID)
			items[i] = extIDs
		}
	}

	for _, attribute := range attributes {
		items[attribute.ExtID] = []string{attribute.ExtID}
	}

	stopListItemsForAggregatorMenu := make([]menuModels.StopListItem, 0, len(stopListItems))
	for i := range stopListItems {

		stopListItem := stopListItems[i]
		productID := stopListItem.ProductID

		extIDs, ok := items[productID]
		if !ok {
			continue
		}

		for i := range extIDs {
			extID := extIDs[i]
			newStopListItem := menuModels.StopListItem{
				ProductID: extID,
				Balance:   stopListItem.Balance,
			}
			stopListItemsForAggregatorMenu = append(stopListItemsForAggregatorMenu, newStopListItem)
		}

	}
	return stopListItemsForAggregatorMenu
}

func (s *ServiceImpl) actualizeStopListAggregatorMenu(ctx context.Context, store storeModels.Store, menu storeModels.StoreDSMenu, newStopListItems menuModels.StopListItems, isByBalance bool, balanceLimit float64, posStopListItems menuModels.StopListItems) error {
	if !menu.IsActive {
		return nil
	}

	products, attributes, err := s.getItemsForUpdateAggregatorMenu(ctx, store.ID, menu.ID, newStopListItems, isByBalance, balanceLimit)
	if err != nil {
		return err
	}
	log.Info().Msgf("list of stopList products: %+v, attributes: %+v", products, attributes)

	menuID := menu.ID
	if err = s.menuService.UpdateStopList(ctx, menuID, newStopListItems.Products()); err != nil {
		return err
	}

	if err = s.updateStopListByProductIDInAggregator(ctx, store, products, menu.Delivery, posStopListItems); err != nil {
		return err
	}
	if err = s.updateStopListByProductIDsInDatabase(ctx, menuID, products); err != nil {
		return err
	}

	if err = s.updateStopListByAttributesInAggregator(ctx, store, attributes, menu.Delivery, posStopListItems); err != nil {
		return err
	}
	if err = s.updateStopListByAttributesInDatabase(ctx, menuID, attributes); err != nil {
		return err
	}

	// for yandex
	if err := s.sendStopListUpdateNotification(ctx, store, menu.Delivery); err != nil {
		//TODO: refactor
		log.Err(err).Msgf("failed to send stop list update to yandex")
		//return err
	}

	return nil
}

func (s *ServiceImpl) getItemsForUpdatePosMenu(ctx context.Context, menuID string, stopListItems []menuModels.StopListItem,
	isByBalance bool, balanceLimit float64) ([]menuModels.Product, []menuModels.Attribute, error) {
	posMenu, err := s.menuService.FindById(ctx, menuID)
	if err != nil {
		return nil, nil, err
	}

	extractor := posMenuIDExtractor{}
	slc, err := newStopListMenuComparator(posMenu, stopListItems, isByBalance, balanceLimit, extractor)
	if err != nil {
		return nil, nil, err
	}

	products, attributes := slc.process()
	return products, attributes, nil

}

func (s *ServiceImpl) getItemsForUpdateAggregatorMenu(ctx context.Context, storeID string, menuID string, stopListItems []menuModels.StopListItem,
	isByBalance bool, balanceLimit float64) ([]menuModels.Product, []menuModels.Attribute, error) {
	menu, err := s.menuService.FindById(ctx, menuID)
	if err != nil {
		return nil, nil, err
	}

	extractor := aggregatorMenuIDExtractor{false, storeID}
	slc, err := newStopListMenuComparator(menu, stopListItems, isByBalance, balanceLimit, extractor)
	if err != nil {
		return nil, nil, err
	}

	products, attributes := slc.process()
	return s.updateProductsFilter(menu, products), s.updateAttributeFilter(menu, attributes), nil
}

func (s *ServiceImpl) updateStopListByAttributesInAggregator(ctx context.Context, store storeModels.Store, attributes []menuModels.Attribute,
	deliveryService string, posStopListItems menuModels.StopListItems) error {

	if len(attributes) == 0 {
		return nil
	}

	externalStoreIDs, err := s.storeService.GetStoreExternalIds(store, deliveryService)
	if err != nil {
		return err
	}

	aggregatorService, err := s.aggregatorFactory.GetAggregator(deliveryService, store)
	if err != nil {
		return err
	}

	sla := s.toStopListAttributes(attributes)
	transaction := menuModels.StopListTransaction{
		StoreID:          store.ID,
		Attributes:       sla,
		PosStopListItems: posStopListItems,
	}

	for _, storeID := range externalStoreIDs {
		transactionData := menuModels.TransactionData{
			StoreID:    storeID,
			Delivery:   deliveryService,
			Attributes: sla,
		}
		trID, err := aggregatorService.UpdateStopListByAttributesBulk(ctx, storeID, attributes)
		if err != nil {
			transactionData.Status = menuModels.ERROR
			transactionData.Message = err.Error()
			s.sendToNotify(ctx, store, transactionData)
		} else {
			transactionData.Status = menuModels.SUCCESS
		}

		transactionData.ID = trID
		transaction.Transactions = append(transaction.Transactions, transactionData)
	}

	if err = s.repo.InsertStopListTransaction(ctx, transaction); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) toStopListAttributes(attributes []menuModels.Attribute) menuModels.StopListAttributes {
	res := make(menuModels.StopListAttributes, 0, len(attributes))

	for _, attribute := range attributes {
		res = append(res, menuModels.StopListAttribute{
			AttributeID:   attribute.ExtID,
			AttributeName: attribute.Name,
			IsAvailable:   attribute.IsAvailable,
			Price:         attribute.Price,
		})
	}

	return res
}

func (s *ServiceImpl) updateStopListByAttributesInDatabase(ctx context.Context, menuID string, attributes []menuModels.Attribute) error {
	if err := s.updateStopListAvailableStatusByAttributeIDsInDatabase(ctx, menuID, attributes); err != nil {
		return err
	}
	if err := s.updateDisabledStatusByAttributeIDsInDatabase(ctx, menuID, attributes); err != nil {
		return err
	}
	return nil
}

func (s *ServiceImpl) updateStopListAvailableStatusByAttributeIDsInDatabase(ctx context.Context, menuID string, attributes []menuModels.Attribute) error {
	var (
		attributeIdsWithAvailabilityFalse = make([]string, 0, len(attributes))
		attributeIdsWithAvailabilityTrue  = make([]string, 0, len(attributes))
	)

	for _, attribute := range attributes {
		if attribute.IsAvailable {
			attributeIdsWithAvailabilityTrue = append(attributeIdsWithAvailabilityTrue, attribute.ExtID)
		} else {
			attributeIdsWithAvailabilityFalse = append(attributeIdsWithAvailabilityFalse, attribute.ExtID)
		}
	}

	if len(attributeIdsWithAvailabilityTrue) != 0 {
		if err := s.menuService.UpdateAttributesAvailabilityStatus(ctx, menuID, attributeIdsWithAvailabilityTrue, true); err != nil {
			return err
		}
	}

	if len(attributeIdsWithAvailabilityFalse) != 0 {
		if err := s.menuService.UpdateAttributesAvailabilityStatus(ctx, menuID, attributeIdsWithAvailabilityFalse, false); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) updateProductsFilter(menu *menuModels.Menu, products menuModels.Products) menuModels.Products {

	result := make(menuModels.Products, 0, len(products))

	oldProductMap := make(map[string]bool)
	for _, product := range menu.Products {
		oldProductMap[product.ExtID] = product.IsAvailable
	}

	for _, product := range products {
		if oldAvailability, found := oldProductMap[product.ExtID]; found {
			if oldAvailability != product.IsAvailable {
				result = append(result, product)
			}
		}
	}
	//	return result
	return products
}

func (s *ServiceImpl) updateAttributeFilter(menu *menuModels.Menu, attributes menuModels.Attributes) menuModels.Attributes {

	result := make(menuModels.Attributes, 0, len(attributes))

	oldAttributeMap := make(map[string]bool)
	for _, attribute := range menu.Attributes {
		oldAttributeMap[attribute.ExtID] = attribute.IsAvailable
	}

	for _, attribute := range attributes {
		if oldAvailability, found := oldAttributeMap[attribute.ExtID]; found {
			if oldAvailability != attribute.IsAvailable {
				result = append(result, attribute)
			}
		}
	}
	//return result
	return attributes
}

func (s *ServiceImpl) sendStopListUpdateNotification(ctx context.Context, store storeModels.Store, deliveryService string) error {
	externalStoreIDs, err := s.storeService.GetStoreExternalIds(store, deliveryService)
	if err != nil {
		return err
	}

	aggregatorService, err := s.aggregatorFactory.GetAggregator(deliveryService, store)
	if err != nil {
		return err
	}

	for _, storeID := range externalStoreIDs {
		if err := aggregatorService.SendStopListUpdateNotification(ctx, storeID); err != nil {
			return err
		}
	}

	return nil
}
