package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/order/dto"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

func run() error {
	ctx := context.Background()

	opts, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}

	log.Printf("options %+v \n", opts)

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

	cli, err := order.NewClient()
	if err != nil {
		log.Printf("order.NewClient error %v", err)
		return err
	}

	preorders, err := cli.GetActivePreorders(ctx, dto.ActiveOrderSelector{})
	if err != nil {
		log.Err(err).Msg("cant find active preorders")
		return err
	}

	directPreorders, err := cli.GetDirectCallCenterActivePreorders(ctx, dto.ActiveOrderSelector{})
	if err != nil {
		log.Err(err).Msg("cant find active preorders for direct, call center")
		return err
	}

	haniPreorders, err := cli.GetHaniKarimaActivePreOrders(ctx, dto.ActiveOrderSelector{})
	if err != nil {
		log.Err(err).Msg("cant find active preorders for hani")
		return err
	}

	preorders = append(preorders, directPreorders...)
	preorders = append(preorders, haniPreorders...)

	uniqueOrders := make(map[string]coreOrderModels.Order)
	for _, order := range preorders {
		uniqueOrders[order.ID] = order
	}

	var uniquePreorders []coreOrderModels.Order
	for _, order := range uniqueOrders {
		log.Info().Msgf("preorder with id %s was selected as unique. delete this log for future", order.ID)
		uniquePreorders = append(uniquePreorders, order)
	}

	if len(preorders) == 0 {
		log.Info().Msg("No active preorders to send")
	}

	var wg sync.WaitGroup
	for _, preorder := range uniquePreorders {
		wg.Add(1)
		go func(preorder coreOrderModels.Order) {
			log.Info().Msgf("sending preorder %s - %s\n", preorder.ID, preorder.RestaurantID)
			preorder, err = cli.CreateOrderInPOS(context.TODO(), preorder)
			if err != nil {
				log.Err(err).Msg("error while creating order in POS")
			} else {
				log.Info().Msgf("Preorder %s successfully sent to POS, POS id %s", preorder.ID, preorder.PosOrderID)
			}
			wg.Done()
		}(preorder)
	}
	wg.Wait()

	return nil
}

func main() {
	if cmd.IsLambda() {
		log.Info().Msg("Starting lambda")
		lambda.Start(run)
	} else {
		log.Info().Msg("Starting locally")

		log.Info().Msg("Setting environment variables")
		os.Setenv("SECRET_ENV", "StageEnvs")
		os.Setenv("REGION", "eu-west-1")
		os.Setenv("SENTRY", "StageSentry")

		if err := run(); err != nil {
			log.Err(err).Msgf("failed run")
		}
	}

}
