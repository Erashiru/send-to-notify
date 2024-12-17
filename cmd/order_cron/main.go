package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/order/dto"
	"github.com/rs/zerolog/log"
	"time"
)

func Run(ctx context.Context) error {

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

	cli, err := order.NewClient()
	if err != nil {
		log.Printf("order.NewClient error %v", err)
		return err
	}

	orders, err := cli.GetActiveOrders(ctx, dto.ActiveOrderSelector{
		PosType: "paloma",
	})
	if err != nil {
		log.Err(err).Msg("get active paloma orders error")
		return err
	}

	for _, order := range orders {
		log.Info().Msgf("manual order update for id=%s, order_id=%s, order_code=%s, delivery=%s", order.ID, order.OrderID, order.OrderCode, order.DeliveryService)

		if err := cli.ManualUpdateStatus(ctx, order); err != nil {
			log.Err(err).Msgf("manual update status error")
			utils.Beautify("manual update status error for: ", order)
			continue
		}
	}

	log.Info().Msgf("manual order update finished")

	posterOrders, _, err := cli.GetOrdersWithFilters(ctx, dto.OrderSelector{
		PosType:       "poster",
		OnlyActive:    true,
		OrderTimeFrom: time.Now().UTC().Add(-2 * time.Hour),
	})

	if err != nil {
		log.Err(err).Msg("get active poster orders error")
		return err
	}

	for _, order := range posterOrders {
		if order.CookingCompleteTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
			continue
		}
		if order.CookingCompleteTime.After(time.Now().UTC()) {
			continue
		}

		if hasStatusInStatusesHistory("COOKING_COMPLETE", order.StatusesHistory) {
			continue
		}

		log.Info().Msgf("update order status in aggregator for id=%s, order_id=%s, order_code=%s, delivery=%s", order.ID, order.OrderID, order.OrderCode, order.DeliveryService)

		err := cli.UpdateOrderStatus(ctx, order.PosOrderID, order.PosType, "ready", "")
		if err != nil {
			log.Err(err).Msgf("update order status error for id=%s, order_id=%s, order_code=%s, delivery=%s", order.ID, order.OrderID, order.OrderCode, order.DeliveryService)
		}
	}

	log.Info().Msgf("update poster order status in aggregator finished")

	return nil
}

func hasStatusInStatusesHistory(status string, history []coreOrderModels.OrderStatusUpdate) bool {
	for _, s := range history {
		if s.Name == status {
			return true
		}
	}
	return false
}

func main() {
	lambda.Start(Run)
}
