package main

import (
	"context"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/service/stoplist"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

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
	_ = stopListService

	if err = stopListService.UpdateStopListByPosProductID(ctx, true, "", ""); err != nil {
		return err
	}

	return nil
}
