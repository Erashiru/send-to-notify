package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	glovoCli "github.com/kwaaka-team/orders-core/pkg/glovo"
	glovoConfig "github.com/kwaaka-team/orders-core/pkg/glovo/clients"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/kwaaka-team/orders-core/service/order/scheduledupdatestatus"
	"github.com/rs/zerolog/log"
)

func run() error {
	ctx := context.Background()
	opts, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}
	ds, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSDB,
	})
	if err != nil {
		return fmt.Errorf("cannot create datastore %s: %v", opts.DSName, err)
	}

	if err = ds.Connect(ds.Client()); err != nil {
		return fmt.Errorf("cannot connect to datastore: %s", err)
	}
	defer ds.Close(ctx)

	log.Printf("[INFO] connected to %s", ds.Name())

	orderCli, err := order.NewClient()
	if err != nil {
		log.Printf("order.NewClient error %v", err)
		return err
	}
	storeCli, err := store.NewClient(storeModels.Config{
		MongoCli: ds.Client(),
	})
	if err != nil {
		log.Printf("new storeClient error %v", err)
		return err
	}

	glovoCLi, err := glovoCli.NewGlovoClient(&glovoConfig.Config{
		Protocol: "http",
		BaseURL:  opts.GlovoConfiguration.BaseURL,
		ApiKey:   opts.GlovoConfiguration.Token,
	})
	if err != nil {
		return err
	}
	deliveryService, err := scheduledupdatestatus.NewDeliveryService(glovoCLi, orderCli)
	if err != nil {
		return err
	}

	scheduledUpdate := scheduledupdatestatus.New(storeCli, deliveryService)

	if err = scheduledUpdate.UpdateToReady(ctx); err != nil {
		return err
	}
	return nil
}

func main() {
	if cmd.IsLambda() {
		lambda.Start(run)
	} else {
		if err := run(); err != nil {
			log.Err(err).Msgf("failed run")
		}
	}
}
