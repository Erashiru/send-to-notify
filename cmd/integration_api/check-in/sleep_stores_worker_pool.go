package main

import (
	"context"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/store"
	"sync"
)

type SleepStoresWorkerPool struct {
	jobs         chan coreStoreModels.Store
	resChan      chan coreStoreModels.Store
	concurrency  int
	posFactory   pos.Factory
	storeService store.Service

	wg *sync.WaitGroup
}

func NewWorkerPool(concurrency int, posFactory pos.Factory, storeService store.Service) *SleepStoresWorkerPool {
	wg := &sync.WaitGroup{}
	var jobs = make(chan coreStoreModels.Store)
	var resChan = make(chan coreStoreModels.Store)
	return &SleepStoresWorkerPool{
		jobs:         jobs,
		resChan:      resChan,
		concurrency:  concurrency,
		posFactory:   posFactory,
		storeService: storeService,
		wg:           wg,
	}
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
		isAliveVal, err := wp.isAlive(ctx, wp.posFactory, wp.storeService, st)
		if err != nil {
			continue
		}
		if isAliveVal {
			continue
		}
		wp.resChan <- st
	}
}

func (wp *SleepStoresWorkerPool) isAlive(ctx context.Context, posFactory pos.Factory, storeService store.Service, st coreStoreModels.Store) (bool, error) {
	isAliveVal, err := isPosAlive(ctx, posFactory, storeService, st)
	if err != nil {
		return false, err
	}
	return isAliveVal, nil

}
