package main

import (
	"context"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Err(err).Msg("timezone update error")
	}
}

func run() error {
	os.Setenv("SECRET_ENV", "ProdEnvs")
	os.Setenv("REGION", "eu-west-1")
	os.Setenv("SENTRY", "ProdSentry")

	ctx := context.Background()

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}

	db, err := cmd.CreateMongo(ctx, cfg.DSURL, cfg.DSDB)
	if err != nil {
		return err
	}

	repo, err := store.NewStoreMongoRepository(db)
	if err != nil {
		return err
	}
	storeService, err := store.NewService(repo)
	if err != nil {
		return err
	}

	stores, err := storeService.FindStoresByTimeZone(ctx, "Asia/Almaty")
	if err != nil {
		return err
	}

	for i := range stores {
		st := stores[i]
		if err := repo.UpdateTZ(ctx, st.ID, "Asia/Aqtobe"); err != nil {
			return err
		}
		if err := repo.UpdateOffset(ctx, st.ID, 5); err != nil {
			return err
		}
		log.Info().Msgf("done %d out of %d", i+1, len(stores))

	}

	return nil
}
