package main

import (
	"context"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/stoplist"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
)

type WorkerPool struct {
	jobs            chan coreStoreModels.Store
	concurrency     int
	stopListService stoplist.Service
	menuService     *menu.Service
	menuRepo        menu.Repository

	wg *sync.WaitGroup
}

func NewWorkerPool(concurrency int, stopListService stoplist.Service,
	menuService *menu.Service, menuRepo menu.Repository) (*WorkerPool, error) {

	if concurrency <= 0 {
		return nil, errors.New("concurrency value is invalid")
	}
	if stopListService == nil || menuService == nil || menuRepo == nil {
		return nil, errors.New("constructor error")
	}

	wg := &sync.WaitGroup{}
	var jobs = make(chan coreStoreModels.Store)
	return &WorkerPool{
		jobs:            jobs,
		concurrency:     concurrency,
		stopListService: stopListService,
		menuService:     menuService,
		menuRepo:        menuRepo,
		wg:              wg,
	}, nil
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for w := 1; w <= wp.concurrency; w++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			wp.doWork(ctx)
		}()
	}
}

func (wp *WorkerPool) doWork(ctx context.Context) {
	for st := range wp.jobs {
		if err := process(ctx, st, wp.stopListService, wp.menuService, wp.menuRepo); err != nil {
			log.Err(err).Msgf("")
		}
	}
}

func (wp *WorkerPool) Submit(stores []coreStoreModels.Store) {
	go func() {
		for i := range stores {
			st := stores[i]
			wp.jobs <- st
		}
		close(wp.jobs)
	}()
	wp.wg.Wait()
}
