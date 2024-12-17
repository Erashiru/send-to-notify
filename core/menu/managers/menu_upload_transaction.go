package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/managers/validator"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
)

type MenuUploadTransactionManager interface {
	Create(ctx context.Context, req models.MenuUploadTransaction) (string, error)
	Get(ctx context.Context, query selector.MenuUploadTransaction) (models.MenuUploadTransaction, error)
	List(ctx context.Context, query selector.MenuUploadTransaction) ([]models.MenuUploadTransaction, int64, error)
	Update(ctx context.Context, req models.MenuUploadTransaction) error
	Delete(ctx context.Context, id string) error
}

type mtm struct {
	globalConfig menu.Configuration
	mutRepo      drivers.MenuUploadTransactionRepository
	muValidator  validator.MenuUploadTransaction
}

func NewMenuUploadTransactionManager(
	globalConfig menu.Configuration,
	mutRepo drivers.MenuUploadTransactionRepository,
	muValidator validator.MenuUploadTransaction) MenuUploadTransactionManager {
	return &mtm{
		globalConfig: globalConfig,
		mutRepo:      mutRepo,
		muValidator:  muValidator,
	}
}

func (m *mtm) Create(ctx context.Context, req models.MenuUploadTransaction) (string, error) {
	res, err := m.mutRepo.Insert(ctx, req)
	if err != nil {
		return "", err
	}

	// if req.ExtTransactions.HasProcessingStatus() {
	// 	go m.verifyMenu(context.Background(), res, req)
	// }

	return res, nil
}

func (m *mtm) Get(ctx context.Context, query selector.MenuUploadTransaction) (models.MenuUploadTransaction, error) {
	return m.mutRepo.Get(ctx, query)
}

func (m *mtm) List(ctx context.Context, query selector.MenuUploadTransaction) ([]models.MenuUploadTransaction, int64, error) {
	return m.mutRepo.List(ctx, query)
}

func (m *mtm) Update(ctx context.Context, req models.MenuUploadTransaction) error {

	if err := m.mutRepo.Update(ctx, req); err != nil {
		return err
	}

	return nil
}

func (m *mtm) Delete(ctx context.Context, id string) error {
	return m.mutRepo.Delete(ctx, id)
}

// func (m *mtm) verifyMenu(ctx context.Context, id string, req models.MenuUploadTransaction) error {

// 	aggregatorMan, err := aggregator.NewManager(ctx, m.globalConfig, models.AggregatorName(req.Service), models.Store{}, nil, nil)
// 	if err != nil {
// 		return err
// 	}

// 	isChanged := false
// 	for i, tr := range req.ExtTransactions {
// 		tr.ID = id
// 		if tr.Status == models.PROCESSING.String() {
// 			status, err := aggregatorMan.VerifyMenu(ctx, tr)
// 			if err != nil {
// 				log.Err(err).Msgf("verify menu err in store_id %s", req.StoreID)
// 			}

// 			if status != models.PROCESSING && status != "" {
// 				isChanged = true
// 			}

// 			req.ExtTransactions[i].Status = status.String()
// 		}
// 	}

// 	req.UpdatedAt.Value = models.TimeNow()

// 	if isChanged {
// 		if err = m.mutRepo.Update(ctx, models.UpdateMenuUploadTransaction{
// 			ID:              id,
// 			RestaurantID:    &req.StoreID,
// 			ExtTransactions: req.ExtTransactions,
// 		}); err != nil {
// 			log.Err(err).Msgf("verify menu err in store_id %s", req.StoreID)
// 		}
// 	}

// 	return nil
// }
