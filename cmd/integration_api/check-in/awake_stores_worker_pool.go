package main

import (
	"context"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/pos"
	"github.com/rs/zerolog/log"
	"sync"
)

type AwakeSleepStoresWorkerPool struct {
	jobs        chan coreStoreModels.Store
	concurrency int
	posFactory  pos.Factory

	wg *sync.WaitGroup
}

func NewAwakeSleepStoresWorkerPool(concurrency int, posFactory pos.Factory) *AwakeSleepStoresWorkerPool {
	wg := &sync.WaitGroup{}
	var jobs = make(chan coreStoreModels.Store)
	return &AwakeSleepStoresWorkerPool{
		jobs:        jobs,
		concurrency: concurrency,
		posFactory:  posFactory,
		wg:          wg,
	}
}

func (wp *AwakeSleepStoresWorkerPool) Start(ctx context.Context) {
	for w := 1; w <= wp.concurrency; w++ {
		wp.wg.Add(1)
		go func() {
			wp.doWork(ctx)
		}()
	}
}

func (wp *AwakeSleepStoresWorkerPool) doWork(ctx context.Context) {
	defer wp.wg.Done()

	for store := range wp.jobs {
		if err := awakeSleepStore(ctx, wp.posFactory, store); err != nil {
			log.Err(err).Msgf("")
		}
	}
}

func (wp *AwakeSleepStoresWorkerPool) Submit(stores []coreStoreModels.Store) {
	go func() {
		for i := range stores {
			st := stores[i]
			wp.jobs <- st
		}
		close(wp.jobs)
	}()
	wp.wg.Wait()
}
