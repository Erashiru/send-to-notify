// Dear programmer
// When I wrote this code, only god and
// I knew how it worked.
// Now, only god knows it!
//
// Therefore, if you are trying to optimize
// this routine and it fails (most surely),
// please increase this counter as a
// warning for the next person:
//
// totalHoursWastedHere = 667
//

package managers

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator"
	"github.com/kwaaka-team/orders-core/core/menu/managers/validator"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/pointer"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	entityChangesHistoryModels "github.com/kwaaka-team/orders-core/service/entity_changes_history/models"
	"time"

	"github.com/rs/zerolog/log"
)

// UpdateStopListStores allows to update product to/from stoplist
func (m *mnm) UpdateStopListStores(ctx context.Context, req models.UpdateStopListProduct, history entityChangesHistoryModels.EntityChangesHistoryRequest) ([]models.StopListTransaction, error) {

	// TODO add parent level for query times for product to stop list
	transactionResults := make([]models.StopListTransaction, 0, len(req.Data))

	// fixme: mb use async method?
	for _, store := range req.Data {

		res := models.StopListTransaction{
			StoreID: store.ID,
			Products: []models.StopListProduct{ // todo fix here
				{
					ProductID: req.ProductID,
				},
			},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		transactions, err := m.updateProductStore(ctx, req, store, history)
		if err != nil {
			log.Trace().Msg(err.Error())
			continue
		}

		res.Transactions = transactions
		transactionResults = append(transactionResults, res)
	}

	// fixme temporary method to save stop list transactions in DB
	// here we have to use another context, cause main context will be done or cancelled
	go m.stm.Insert(context.Background(), transactionResults)

	return transactionResults, nil
}

// updateProductStore utils to update one product to many stores in aggregators
func (m *mnm) updateProductStore(ctx context.Context, req models.UpdateStopListProduct, data models.UpdateStoreData, history entityChangesHistoryModels.EntityChangesHistoryRequest) ([]models.TransactionData, error) {

	product := models.Product{
		ExtID: req.ProductID,
	}

	product.IsAvailable = req.SetToStop

	transactions := make([]models.TransactionData, 0, len(data.Aggregators))

	store, err := m.storeRepo.Get(ctx, selector.EmptyStoreSearch().
		SetID(data.ID).SetIsActiveMenu(pointer.OfBool(true)))
	if err != nil {
		log.Trace().Msgf("could not find stores from %s", store.ID)
		return nil, err
	}

	for _, aggr := range data.Aggregators {

		if !aggr.IsActive {
			log.Trace().Msgf("aggregator %s not active", aggr.Name)
			transactions = append(transactions, models.TransactionData{
				ID:       data.ID,
				StoreID:  data.ID,
				Delivery: aggr.Name.String(),
				Status:   models.ERROR,
				Message:  fmt.Sprintf("aggregator %s not active", aggr.Name),
			})
			continue
		}

		transactionResults, err := m.updateAggregatorProduct(ctx, storeModels.AggregatorName(aggr.Name), store, product, history)
		if err != nil {
			log.Trace().Msgf("could not update from aggregator stores from %s with aggregator name %s", store.ID, aggr.Name)
			transactions = append(transactions, models.TransactionData{
				ID:       store.ID,
				StoreID:  store.ID,
				Delivery: aggr.Name.String(),
				Status:   models.ERROR,
				Message:  err.Error(),
			})
			continue
		}

		// cause different aggregators in one store (many to one)
		transactions = append(transactions, transactionResults...)

	}

	if len(transactions) == 0 {
		return nil, validator.ErrProductNotModifier
	}

	// update product in pos Menu
	if err = m.updateProductAvailable(ctx, store.MenuID, product, history); err != nil {
		log.Trace().Msgf("could not update pos menu %s in store id %s", store.MenuID, store.ID)
		// return nil, err
	}

	return transactions, nil
}

// updateAggregatorProduct update product in aggregators
func (m *mnm) updateAggregatorProduct(ctx context.Context, aggregatorName storeModels.AggregatorName,
	store storeModels.Store, product models.Product, history entityChangesHistoryModels.EntityChangesHistoryRequest) ([]models.TransactionData, error) {

	aggrManager, err := m.getAggregatorManager(ctx, store, aggregatorName)
	if err != nil {
		return nil, err
	}

	extStoreIDs := store.GetAggregatorStoreIDs(aggregatorName.String())
	if len(extStoreIDs) == 0 {
		return nil, validator.ErrEmptyStores
	}

	transactions := make([]models.TransactionData, 0, len(extStoreIDs))
	for _, storeAggrID := range extStoreIDs {

		transaction := models.TransactionData{
			ID:       store.ID,
			StoreID:  storeAggrID,
			Delivery: aggregatorName.String(),
		}

		rsp, err := aggrManager.ModifyProduct(ctx, storeAggrID, product)
		if err != nil {
			transaction.Status = models.ERROR
			transaction.Message = err.Error()
			transactions = append(transactions, transaction)
			continue
		}

		transaction.ID = rsp.ExtID
		transaction.Products = append(transaction.Products, models.StopListProduct{
			ExtID:       product.ExtID,
			IsAvailable: product.IsAvailable,
		})
		transaction.Status = models.SUCCESS

		transactions = append(transactions, transaction)
	}

	activeMenu := store.Menus.GetActiveMenu(aggregatorName)
	if activeMenu.ID == "" {
		return nil, fmt.Errorf("there is no active menu on %s menu", aggregatorName)
	}

	// update aggregators Menu
	if err = m.updateProductAvailable(ctx, activeMenu.ID, product, history); err != nil {
		log.Trace().Msgf("could not update aggregator menu %s in store id %s", activeMenu.ID, store.ID)
		// return nil, err
	}

	return transactions, nil
}

// getAggregatorManager implements new aggregator manager
func (m *mnm) getAggregatorManager(ctx context.Context, store storeModels.Store, aggregatorName storeModels.AggregatorName) (aggregator.Base, error) {

	aggrManager, err := aggregator.NewManager(ctx, m.globalConfig, aggregatorName, store, m.mspRepo, m.stRepo, m.storeRepo, m.menuRepo, m.restGroupMenuRepo, m.storeCli)
	if err != nil {
		log.Trace().Msgf("could not initialize aggregator manager with aggregator name %s from store %s", aggregatorName, store.ID)
		return nil, err
	}

	return aggrManager, nil
}

// updateProductAvailable allows to update POS menu or aggregator Menu
func (m *mnm) updateProductAvailable(ctx context.Context, menuId string, product models.Product, history entityChangesHistoryModels.EntityChangesHistoryRequest) error {

	menu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuId))
	if err != nil {
		return err
	}
	var IsProductFound bool

	for i := range menu.Products {
		if menu.Products[i].ExtID == product.ExtID {
			menu.Products[i].IsAvailable = product.IsAvailable
			IsProductFound = true
		}
	}
	if !IsProductFound {
		return fmt.Errorf("there is no %s in menuId %s", product.Name, menuId)
	}

	if err = m.menuRepo.Update(ctx, menu, entityChangesHistoryModels.EntityChangesHistory{
		CallFunction: "updateProductAvailable",
		Author:       history.Author,
		TaskType:     history.TaskType,
	}); err != nil {
		return err
	}

	return nil
}

func (m *mnm) UpdateMatchingProduct(ctx context.Context, req models.MatchingProducts) error {
	return m.menuRepo.UpdateProductForMatching(ctx, req)
}
