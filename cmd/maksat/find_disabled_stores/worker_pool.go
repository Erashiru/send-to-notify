package main

import (
	"context"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
)

type SleepStoresWorkerPool struct {
	concurrency int
	jobs        chan coreStoreModels.Store
	resChan     chan coreStoreModels.Store

	wg *sync.WaitGroup

	menuService *menu.Service
}

func NewWorkerPool(concurrency int, menuService *menu.Service) (*SleepStoresWorkerPool, error) {
	wg := &sync.WaitGroup{}
	var jobs = make(chan coreStoreModels.Store)
	var resChan = make(chan coreStoreModels.Store)
	if menuService == nil {
		return nil, errors.New("menu service is nil")
	}
	return &SleepStoresWorkerPool{
		jobs:        jobs,
		resChan:     resChan,
		concurrency: concurrency,
		wg:          wg,
		menuService: menuService,
	}, nil
}

func (wp *SleepStoresWorkerPool) Start(ctx context.Context) {
	for w := 1; w <= wp.concurrency; w++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			wp.doWork(ctx)
		}()
	}
}

func (wp *SleepStoresWorkerPool) Submit(stores []coreStoreModels.Store) {
	go func() {
		for i := range stores {
			st := stores[i]
			log.Info().Msgf("submit store number %d \n", i)
			wp.jobs <- st
		}
		close(wp.jobs)
		wp.wg.Wait()
		close(wp.resChan)
	}()
}

func (wp *SleepStoresWorkerPool) GetResult() []coreStoreModels.Store {
	result := make([]coreStoreModels.Store, 0)
	for res := range wp.resChan {
		result = append(result, res)
	}
	return result
}

func (wp *SleepStoresWorkerPool) doWork(ctx context.Context) {
	for st := range wp.jobs {
		isHasDisabledProduct, err := wp.isDisabledActivated(ctx, st)
		if err != nil {
			log.Err(err).Msgf("isDisabledActivated error")
			continue
		}
		if !isHasDisabledProduct {
			continue
		}
		wp.resChan <- st
	}
}

func (wp *SleepStoresWorkerPool) isDisabledActivated(ctx context.Context, st coreStoreModels.Store) (bool, error) {
	posMenu, err := wp.menuService.FindById(ctx, st.MenuID)
	if err != nil {
		return false, err
	}
	for i := range posMenu.Products {
		product := posMenu.Products[i]
		if product.IsDeleted {
			continue
		}
		if product.IsDisabled {
			return true, nil
		}
	}
	return false, nil
}
