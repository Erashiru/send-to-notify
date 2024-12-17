package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/rs/zerolog/log"
)

func main() {
	if cmd.IsLambda() {
		log.Info().Msgf("[INFO] Starting Lambda function")
		lambda.Start(run)
		log.Info().Msgf("[INFO] Lambda function finished work")
	} else {
		log.Info().Msgf("[INFO] Running locally")
		if err := run(); err != nil {
			log.Err(err).Msgf("[ERROR] %s", err.Error())
			return
		}
		log.Info().Msgf("[INFO] Finished running locally")
	}
}

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
		log.Err(err).Msgf("cannot create datastore %s: %v", opts.DSName, err)
		return fmt.Errorf("cannot create datastore %s: %v", opts.DSName, err)
	}

	if err = ds.Connect(ds.Client()); err != nil {
		log.Err(err).Msgf("cannot connect to datastore: %s", err)
		return fmt.Errorf("cannot connect to datastore: %s", err)
	}
	db := ds.Client().Database(opts.DSDB)

	//defer ds.Close(ctx)

	log.Info().Msgf("[INFO] connected to %s", ds.Name())

	storeCli, err := store.NewClient(storeModels.Config{
		MongoCli: ds.Client(),
	})
	if err != nil {
		log.Err(err).Msgf("new storeClient error %v", err)
		return err
	}
	orderCli, err := order.NewClient()
	if err != nil {
		log.Err(err).Msgf("creating new order-core client error: %v", err)
		return err
	}

	_, _, posFactory, _, _, err := order.CreateServices(db, opts)
	if err != nil {
		log.Err(err).Msgf("creating posFactory error: %v", err)
		return err
	}

	orderAutoCloseStores, err := storeCli.FindStores(ctx, storeModels.StoreSelector{OrderAutoClose: &[]bool{true}[0]})
	if err != nil {
		log.Err(err).Msg("find orderAutoCloseStores error")
		return err
	}

	for _, storeObj := range orderAutoCloseStores {

		orders, err := orderCli.GetOrdersForAutoCloseCron(ctx, storeObj.OrderAutoCloseSettings.OrderAutoCloseTime, storeObj.ID)
		if err != nil {
			log.Err(err).Msg("orderCli.GetOrdersForAutoCloseCron")
			return err
		}
		for _, orderObj := range orders {
			posSystem, err := posFactory.GetPosService(models.Pos(storeObj.PosType), storeObj)
			if err != nil {
				log.Err(err).Msgf("posFactory.GetPosService for pos type: %s", storeObj.PosType)
				return err
			}
			if err := posSystem.CloseOrder(ctx, orderObj.PosOrderID); err != nil {
				log.Err(err).Msgf("posSystem.CloseOrder for pos order id: %s", orderObj.PosOrderID)
				return err
			}
			if err := orderCli.UpdateOrderStatusByID(ctx, orderObj.OrderID, orderObj.PosType, models.CLOSED.String()); err != nil {
				log.Err(err).Msgf("update order status for order id: %s", orderObj.OrderID)
				return err
			}
		}
	}

	log.Info().Msgf("cron finished successfully")

	return nil
}

//Env:DB_SECRET=ProdKwaakaDB;SECRET_ENV=ProdEnvs
