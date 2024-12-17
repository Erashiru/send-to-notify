package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/jowi/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/stoplist"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/rs/zerolog/log"
	"sync"
)

func run() error {

	ctx := context.Background()

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}

	db, err := cmd.CreateMongo(ctx, cfg.DSURL, cfg.DSDB)
	if err != nil {
		return err
	}

	stopListService, err := stoplist.CreateStopListServiceByWebhook(db, cfg, 1)
	if err != nil {
		return err
	}

	storeRepository, err := store.NewStoreMongoRepository(db)
	if err != nil {
		return err
	}
	storeService, err := store.NewService(storeRepository)
	if err != nil {
		return err
	}

	stores, err := storeService.FindStoresByPosType(ctx, models.JOWI)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, restaurant := range stores {
		wg.Add(1)
		go func(store coreStoreModels.Store) {
			defer wg.Done()

			if err = stopListService.ActualizeStopListByStoreID(ctx, store.ID); err != nil {
				return
			}

		}(restaurant)
	}

	wg.Wait()

	return nil
}

// @title						External Client API
// @version					1.0
// @host						external-api.kwaaka.com
// @BasePath					/v1
// @schemes					https http
// @query.collection.format	multi
// @securityDefinitions.apiKey	ApiSecretAuth
// @in							header
// @name						Authorization
// @securityDefinitions.apiKey	ApiTokenAuth
// @in							header
// @name						Authorization
// @description				"Token from cognito"
func main() {
	if cmd.IsLambda() {
		lambda.Start(run)
	} else {
		if err := run(); err != nil {
			log.Err(err).Msgf("failed run")
		}
	}
}
