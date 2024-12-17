package managers

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/menu/managers/validator"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/custom"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	entityChangesHistoryModels "github.com/kwaaka-team/orders-core/service/entity_changes_history/models"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

func (m *mnm) AttributesStopList(ctx context.Context, storeId string, items []models.ItemStopList, history entityChangesHistoryModels.EntityChangesHistoryRequest) error {
	store, err := m.storeRepo.Get(ctx, selector.Store{
		ID: storeId,
	})
	if err != nil {
		return err
	}

	aggregatorStopList, err := m.updatePosMenuAttributesStopList(ctx, store.MenuID, items, history)
	if err != nil {
		return err
	}

	result := models.StopListTransaction{
		StoreID:   store.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	errs := custom.Error{}
	for _, menu := range store.Menus {

		if !menu.IsActive {
			log.Trace().Msgf("menu %s not active", menu.ID)
			continue
		}

		attributes, err := m.updateAggregatorMenuAttributesStopList(ctx, menu.ID, aggregatorStopList, history)
		if err != nil {
			log.Trace().Err(err).Msgf("could not get/update menu %s", menu.ID)
			continue
		}

		if len(attributes) == 0 {
			log.Trace().Msgf("no products && attributes changed")
			continue
		}

		trx, err := m.updateAttributesInAggregator(ctx, storeModels.AggregatorName(menu.Delivery), store, attributes.Unique())
		if err != nil {
			log.Trace().Err(err).Msgf("could not update attribute store %s, menu %s", store.ID, menu.ID)
			errs.Append(err)
		}

		if trx != nil {
			result.Transactions = append(result.Transactions, trx...)
		}

	}

	if result.Transactions == nil {
		log.Info().Msgf("something wrong, transactions is null")
	}

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		// save to DB
		m.stm.Insert(context.Background(), []models.StopListTransaction{result})
	}()

	go func() {
		defer wg.Done()
		// send to telegram
		m.sendToNotify(context.Background(), store, result)
	}()
	wg.Wait()

	return errs.ErrorOrNil()
}

func (m *mnm) updateStopListAggregator(ctx context.Context, menuID string, stopLists []string, hasVirtualStore bool, restaurantID string, history entityChangesHistoryModels.EntityChangesHistoryRequest) (models.Products, models.Attributes, error) {
	menu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuID))
	if err != nil {
		return nil, nil, err
	}

	existStopLists := make(map[string]struct{}, len(stopLists))
	for _, id := range stopLists {
		existStopLists[id] = struct{}{}
	}

	products := make(models.Products, 0, len(menu.Products))

	for i := 0; i < len(menu.Products); i++ {

		id := menu.Products[i].ExtID

		if menu.Products[i].PosID != "" && menu.Products[i].PosID != menu.Products[i].ExtID {
			id = menu.Products[i].PosID
		}

		if hasVirtualStore {
			arr := strings.Split(id, "_")
			if len(arr) >= 2 {
				id = arr[1]
			}

			if arr[0] != restaurantID {
				continue
			}
		}

		// if product disabled by admin or by cron with section, then skip it
		if menu.Products[i].IsDisabled {
			log.Info().Msgf("product %s is disabled", id)
			continue
		}

		if menu.Products[i].IsDeleted {
			menu.Products[i].IsAvailable = false
			products = append(products, menu.Products[i])
			continue
		}

		// if product available is equal true add him to update slice
		if _, ok := existStopLists[id]; ok {
			// if balance not zero we do not update this product in aggregator
			menu.Products[i].IsAvailable = false
			products = append(products, menu.Products[i])
			continue
		}

		// if product available is equal false, but not exist in POS stop list, add him to update slice
		if !menu.Products[i].IsAvailable {
			menu.Products[i].IsAvailable = true
			products = append(products, menu.Products[i])
			continue
		}

		products = append(products, menu.Products[i])
	}

	attributes := make(models.Attributes, 0, len(menu.Attributes))
	for i := 0; i < len(menu.Attributes); i++ {

		id := menu.Attributes[i].ExtID

		if menu.Attributes[i].PosID != "" && menu.Attributes[i].PosID != menu.Attributes[i].ExtID {
			id = menu.Attributes[i].PosID
		}

		if hasVirtualStore {
			arr := strings.Split(id, "_")
			if len(arr) >= 2 {
				id = arr[1]
			}

			if arr[0] != restaurantID {
				continue
			}
		}

		// if product disabled by admin or by cron with section, then skip it
		if menu.Attributes[i].IsDisabled {
			log.Info().Msgf("attribute %s is disabled", menu.Attributes[i].ExtID)
			continue
		}

		if menu.Attributes[i].IsDeleted {
			menu.Attributes[i].IsAvailable = false
			attributes = append(attributes, menu.Attributes[i])
			continue
		}

		if _, ok := existStopLists[id]; ok {
			menu.Attributes[i].IsAvailable = false
			attributes = append(attributes, menu.Attributes[i])
			continue
		}

		if !menu.Attributes[i].IsAvailable {
			menu.Attributes[i].IsAvailable = true
			attributes = append(attributes, menu.Attributes[i])
			continue
		}

		attributes = append(attributes, menu.Attributes[i])
	}

	menu.StopLists = stopLists

	if err = m.menuRepo.Update(ctx, menu, entityChangesHistoryModels.EntityChangesHistory{
		CallFunction: "updateStopListAggregator",
		Author:       history.Author,
		TaskType:     history.TaskType,
	}); err != nil {
		log.Err(err).Msgf("could not update menu %s", menuID)
		return nil, nil, err
	}

	return products, attributes, nil
}

func (m *mnm) bulkUpdateAggregator(ctx context.Context,
	aggregatorName storeModels.AggregatorName,
	store storeModels.Store,
	products models.Products,
	attributes models.Attributes) ([]models.TransactionData, error) {

	aggregatorManager, err := m.getAggregatorManager(ctx, store, aggregatorName)
	if err != nil {
		log.Trace().Err(err).Msgf("could not get aggregator manager")
		return nil, err
	}

	extStoreIDs := store.GetAggregatorStoreIDs(aggregatorName.String())
	if len(extStoreIDs) == 0 {
		return nil, validator.ErrEmptyStores
	}

	transactions := make([]models.TransactionData, 0, len(extStoreIDs))
	for _, storeAggrID := range extStoreIDs {
		if len(products) != 0 {
			transaction := models.TransactionData{
				Delivery:   aggregatorName.String(),
				Products:   models.ToStopListProducts(products),
				Attributes: models.ToStopListAttributes(attributes),
			}

			rsp, err := aggregatorManager.BulkUpdate(ctx, store.ID, storeAggrID, products, attributes, store)
			if err != nil {
				log.Trace().Err(err).Msgf("could not bulk update store %s with external store_id from aggregator %s", store.ID, storeAggrID)
				transaction.Status = models.ERROR
				transaction.Message = err.Error()
			}

			transaction.ID = rsp
			transaction.StoreID = storeAggrID

			transactions = append(transactions, transaction)
		}

		if len(attributes) != 0 {
			transaction := models.TransactionData{
				Delivery:   aggregatorName.String(),
				Attributes: models.ToStopListAttributes(attributes),
			}

			response, err := aggregatorManager.BulkAttribute(ctx, storeAggrID, attributes)
			if err != nil {
				log.Trace().Err(err).Msgf("could not bulk attributes store %s with external store_id from aggregator %s", store.ID, storeAggrID)
				transaction.Status = models.ERROR
				transaction.Message = err.Error()
			}

			transaction.ID = response
			transaction.StoreID = storeAggrID

			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

func (m *mnm) sendToNotify(ctx context.Context, store storeModels.Store, trx models.StopListTransaction) {

	const stopListQueue = "stoplist-telegram"
	hasToBeSent := false
	stopLists, err := m.getStopListName(ctx, store.MenuID)
	if err != nil {
		log.Err(err).Msgf("could not get stop list in store_id %s from menu_id: %s", store.ID, store.MenuID)
	}

	for _, tr := range trx.Transactions {
		if tr.Status == models.ERROR {
			hasToBeSent = true
		}
	}

	if hasToBeSent {
		message := getMessage(ctx, store, trx, stopLists)

		if err = m.notifyCli.SendMessage(stopListQueue, message, "-1001830413167", store.Telegram.TelegramBotToken); err != nil {
			log.Err(err).Msgf("could not send message %s to sqs chat_id: %s", message, store.Telegram.StopListChatID)
		}
	}

}

func getMessage(ctx context.Context, store storeModels.Store, trx models.StopListTransaction, stopLists []string) string {

	var msg strings.Builder
	msg.WriteString("<b>[❌] Обновление стоп-листов</b>\n")

	msg.WriteString("<b>Ресторан: ")
	msg.WriteString(store.Name)
	msg.WriteString("</b>\n")

	msg.WriteString("<b>Город: ")
	msg.WriteString(store.Address.City)
	msg.WriteString("</b>\n")

	msg.WriteString("<b>Дата создания: ")
	msg.WriteString(time.Now().Format("2006.01.02 15:04:05"))
	msg.WriteString("</b>\n\n")

	var (
		products   string
		attributes string
	)

	for _, tr := range trx.Transactions {

		msg.WriteString("<b>Аггрегатор: ")
		msg.WriteString(tr.Delivery)
		msg.WriteString("</b>\n")

		msg.WriteString("<b>Аггрегатор ID: ")
		msg.WriteString(tr.StoreID)
		msg.WriteString("</b>\n")

		if tr.Status == models.ERROR {
			msg.WriteString("<b>Статус: ")
			msg.WriteString(tr.Status.String())
			msg.WriteString("</b>\n")
			msg.WriteString("<b>Ошибка: ")
			msg.WriteString(tr.Message)
			msg.WriteString("</b>\n")
		}

		for _, product := range tr.Products {
			sign := "🔴"
			if product.IsAvailable {
				sign = "\U0001F7E2"
			}
			products += fmt.Sprintf("<b>%v %s</b>\n", sign, strings.ToUpper(product.Name))
		}

		for _, attribute := range tr.Attributes {
			sign := "🔴"
			if attribute.IsAvailable {
				sign = "\U0001F7E2"
			}
			attributes += fmt.Sprintf("<b>%v %s</b>\n", sign, strings.ToUpper(attribute.AttributeName))
		}

		msg.WriteString("<b>Измененные позиции: \n\n")
		msg.WriteString(fmt.Sprintf("Продукты[%s]:\n", tr.Delivery))
		msg.WriteString(products)
		msg.WriteString(fmt.Sprintf("Аттрибуты[%s]:\n", tr.Delivery))
		msg.WriteString(attributes)
		msg.WriteString("</b>\n\n")
	}

	msg.WriteString("<b>Актуальный стоп-лист: \n\n")

	for _, s := range stopLists {
		msg.WriteString("<b>")
		msg.WriteString(s)
		msg.WriteString("</b>\n")
	}

	msg.WriteString("</b>\n")

	return msg.String()
}

func (m *mnm) updateProductsInAggregator(ctx context.Context,
	aggregatorName storeModels.AggregatorName,
	store storeModels.Store,
	products models.Products) ([]models.TransactionData, error) {

	aggregatorManager, err := m.getAggregatorManager(ctx, store, aggregatorName)
	if err != nil {
		log.Trace().Err(err).Msgf("could not get aggregator manager")
		return nil, err
	}

	extStoreIDs := store.GetAggregatorStoreIDs(aggregatorName.String())
	if len(extStoreIDs) == 0 {
		return nil, validator.ErrEmptyStores
	}

	transactions := make([]models.TransactionData, 0, len(extStoreIDs))
	errs := custom.Error{}
	for _, storeAggrID := range extStoreIDs {
		if len(products) != 0 {
			transaction := models.TransactionData{
				Delivery: aggregatorName.String(),
				Products: models.ToStopListProducts(products),
			}

			rsp, err := aggregatorManager.BulkUpdate(ctx, store.ID, storeAggrID, products, nil, store)
			if err != nil {
				log.Trace().Err(err).Msgf("could not update products store %s with external store_id from aggregator %s", store.ID, storeAggrID)
				transaction.Status = models.ERROR
				transaction.Message = err.Error()
				errs.Append(err)
			}

			transaction.ID = rsp
			transaction.StoreID = storeAggrID

			transactions = append(transactions, transaction)
		}
	}

	return transactions, errs.ErrorOrNil()
}

func (m *mnm) updateAttributesInAggregator(ctx context.Context,
	aggregatorName storeModels.AggregatorName,
	store storeModels.Store,
	attributes models.Attributes) ([]models.TransactionData, error) {

	aggregatorManager, err := m.getAggregatorManager(ctx, store, aggregatorName)
	if err != nil {
		log.Trace().Err(err).Msgf("could not get aggregator manager")
		return nil, err
	}

	extStoreIDs := store.GetAggregatorStoreIDs(aggregatorName.String())
	if len(extStoreIDs) == 0 {
		return nil, validator.ErrEmptyStores
	}

	transactions := make([]models.TransactionData, 0, len(extStoreIDs))
	errs := custom.Error{}
	for _, storeAggrID := range extStoreIDs {
		if len(attributes) != 0 {
			transaction := models.TransactionData{
				Delivery:   aggregatorName.String(),
				Attributes: models.ToStopListAttributes(attributes),
			}

			rsp, err := aggregatorManager.BulkAttribute(ctx, storeAggrID, attributes)
			if err != nil {
				log.Trace().Err(err).Msgf("could not update attributes store %s with external store_id from aggregator %s", store.ID, storeAggrID)
				transaction.Status = models.ERROR
				transaction.Message = err.Error()
				errs.Append(err)
			}

			transaction.ID = rsp
			transaction.StoreID = storeAggrID

			transactions = append(transactions, transaction)
		}
	}

	return transactions, errs.ErrorOrNil()
}

func (m *mnm) updatePosMenuAttributesStopList(ctx context.Context, menuID string, items []models.ItemStopList, history entityChangesHistoryModels.EntityChangesHistoryRequest) ([]models.ItemStopList, error) {
	menu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuID))
	if err != nil {
		return nil, err
	}

	stopListMap := make(map[string]struct{}, len(menu.StopLists))
	for _, id := range menu.StopLists {
		stopListMap[id] = struct{}{}
	}
	itemsMap := make(map[string]models.ItemStopList, len(items))
	for _, item := range items {
		itemsMap[item.ID] = item
	}

	aggregatorStopList := make([]models.ItemStopList, len(items))

	for i := 0; i < len(menu.Attributes); i++ {
		// if product disabled by admin or by cron with section, then skip it
		if menu.Attributes[i].IsDisabled {
			log.Info().Msgf("attribute %s is disabled", menu.Attributes[i].ExtID)
			continue
		}
		if menu.Attributes[i].IsDeleted {
			menu.Attributes[i].IsAvailable = false
			continue
		}

		if item, ok := itemsMap[menu.Attributes[i].ExtID]; ok {
			aggregatorStopList = append(aggregatorStopList, models.ItemStopList{
				ID:          menu.Attributes[i].ExtID,
				IsAvailable: item.IsAvailable,
				Price:       item.Price,
			})
			menu.Attributes[i].IsAvailable = item.IsAvailable
			menu.Attributes[i].Price = item.Price
			if item.IsAvailable {
				delete(stopListMap, menu.Attributes[i].ExtID)
				continue
			}
			stopListMap[menu.Attributes[i].ExtID] = struct{}{}
		}

	}

	stopList := make([]string, 0, len(stopListMap))
	for key := range stopListMap {
		stopList = append(stopList, key)
	}
	menu.StopLists = stopList

	if err = m.menuRepo.Update(ctx, menu, entityChangesHistoryModels.EntityChangesHistory{
		CallFunction: "updatePosMenuAttributesStopList",
		Author:       history.Author,
		TaskType:     history.TaskType,
	}); err != nil {
		log.Err(err).Msgf("could not update POS menu %s", menuID)
		return nil, err
	}

	return aggregatorStopList, nil
}

func (m *mnm) updateAggregatorMenuAttributesStopList(ctx context.Context, menuID string, items []models.ItemStopList, history entityChangesHistoryModels.EntityChangesHistoryRequest) (models.Attributes, error) {
	menu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuID))
	if err != nil {
		return nil, err
	}

	stopListMap := make(map[string]struct{}, len(menu.StopLists))
	for _, id := range menu.StopLists {
		stopListMap[id] = struct{}{}
	}
	itemsMap := make(map[string]models.ItemStopList, len(items))
	for _, item := range items {
		itemsMap[item.ID] = item
	}

	aggregatorStopList := make([]models.Attribute, 0, len(items))

	for i := 0; i < len(menu.Attributes); i++ {
		id := menu.Attributes[i].ExtID
		if menu.Attributes[i].PosID != "" && menu.Attributes[i].PosID != menu.Attributes[i].ExtID {
			id = menu.Attributes[i].PosID
		}

		// if product disabled by admin or by cron with section, then skip it
		if menu.Attributes[i].IsDisabled {
			log.Info().Msgf("attribute %s is disabled", id)
			continue
		}
		if menu.Attributes[i].IsDeleted {
			menu.Attributes[i].IsAvailable = false
			continue
		}

		if item, ok := itemsMap[id]; ok {
			menu.Attributes[i].IsAvailable = item.IsAvailable
			menu.Attributes[i].Price = item.Price
			aggregatorStopList = append(aggregatorStopList, menu.Attributes[i])
			if item.IsAvailable {
				delete(stopListMap, id)
				continue
			}
			stopListMap[id] = struct{}{}
		}
	}

	stopList := make([]string, 0, len(stopListMap))
	for key := range stopListMap {
		if key != "" {
			stopList = append(stopList, key)
		}
	}
	menu.StopLists = stopList

	if err = m.menuRepo.Update(ctx, menu, entityChangesHistoryModels.EntityChangesHistory{
		CallFunction: "updateAggregatorMenuAttributesStopList",
		Author:       history.Author,
		TaskType:     history.TaskType,
	}); err != nil {
		log.Err(err).Msgf("could not update POS menu %s", menuID)
		return nil, err
	}

	return aggregatorStopList, nil
}

// todo: dry
func (m *mnm) updateAggregatorMenuProducts(ctx context.Context, menuID string, items []models.Product, isAvailable bool, history entityChangesHistoryModels.EntityChangesHistoryRequest) (models.Products, error) {
	menu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuID))
	if err != nil {
		return nil, err
	}

	stopListMap := make(map[string]struct{}, len(menu.StopLists))
	for _, id := range menu.StopLists {
		stopListMap[id] = struct{}{}
	}
	itemsMap := make(map[string]models.Product, len(items))
	for _, item := range items {
		itemsMap[item.ExtID] = item
	}

	aggregatorStopList := make([]models.Product, 0, len(items))

	for i := 0; i < len(menu.Products); i++ {
		id := menu.Products[i].ExtID
		if menu.Products[i].ProductID != "" && menu.Products[i].ProductID != menu.Products[i].ExtID {
			id = menu.Products[i].ProductID
		}
		// if product disabled by admin or by cron with section, then skip it
		if isAvailable {
			menu.Products[i].IsDisabled = false
		} else {
			menu.Products[i].IsDisabled = true
		}

		if menu.Products[i].IsDeleted {
			menu.Products[i].IsAvailable = false
			continue
		}

		if item, ok := itemsMap[id]; ok {
			menu.Products[i].IsAvailable = item.IsAvailable
			menu.Products[i].Price = item.Price
			aggregatorStopList = append(aggregatorStopList, menu.Products[i])
			if item.IsAvailable {
				delete(stopListMap, id)
				continue
			}
			stopListMap[id] = struct{}{}
		}
	}

	stopList := make([]string, 0, len(stopListMap))
	for key := range stopListMap {
		if key != "" {
			stopList = append(stopList, key)
		}
	}
	menu.StopLists = stopList

	if err = m.menuRepo.Update(ctx, menu, entityChangesHistoryModels.EntityChangesHistory{
		CallFunction: "updateAggregatorMenuProducts",
		Author:       history.Author,
		TaskType:     history.TaskType,
	}); err != nil {
		log.Err(err).Msgf("could not update POS menu %s", menuID)
		return nil, err
	}

	return aggregatorStopList, nil
}
