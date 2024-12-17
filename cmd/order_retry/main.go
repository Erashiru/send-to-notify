package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/order/dto"
	"github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/kwaaka-team/orders-core/service/pos"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"strings"
)

func Run(ctx context.Context, sqsEvent events.SQSEvent) error {
	opts, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}

	log.Info().Msgf("order_retry, configuration: %+v", opts)

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

	log.Info().Msgf("connected to %s", ds.Name())

	cli, err := order.NewClient()
	if err != nil {
		log.Err(err).Msgf("order.NewClient error")
		return err
	}

	storeCli, err := store.NewClient(storeModels.Config{MongoCli: cli.Client()})
	if err != nil {
		log.Err(err).Msgf("new storeClient error")
		return err
	}

	for idx, message := range sqsEvent.Records {
		log.Info().Msgf("The message '%s' for event source: %s => msg_body: %s", message.MessageId, message.EventSource, message.Body)

		orderSelector := dto.OrderSelector{PosOrderID: message.Body}
		split := strings.Split(message.Body, "_")
		if len(split) == 2 && split[0] == "yaros" {
			orderSelector.OrderID = split[1]
		}

		getOrder, err := cli.GetOrder(ctx, orderSelector)
		if err != nil {
			log.Err(err).Msgf("GetOrder error")
			continue
		}
		log.Info().Msgf("orderInfo: %s, %s, %s, %s, %s, %d", getOrder.RestaurantID, getOrder.StoreID, message.Body, getOrder.DeliveryService,
			getOrder.PosType, idx)

		restGroup, err := storeCli.FindStoreGroup(ctx, selector.StoreGroup{StoreIDs: []string{getOrder.RestaurantID}})
		if err != nil {
			log.Err(err).Msgf("FindStoreGroup error")
			continue
		}

		maxRetryCount := restGroup.RetryCount

		if getOrder.RetryCount >= maxRetryCount {
			if err = cli.UpdateOrderStatus(ctx, message.Body, models.IIKO.String(), "Error", getOrder.CreationResult.Message+fmt.Sprintf("\n\nКоличество попыток для создания заказа было равным %d\n", maxRetryCount)); err != nil {
				log.Err(err).Msgf("getOrder.RetryCount >= maxRetryCount -> update order status failed")
				continue
			}
			log.Info().Msgf("getOrder.RetryCount >= maxRetryCount -> success update order")
		} else if getOrder.RetryCount < maxRetryCount {
			if err = retryOrder(ctx, getOrder, cli); err != nil {
				continue
			}
		}
	}

	return nil
}

func retryOrder(ctx context.Context, getOrder models.Order, cli *order.OrderCoreClient) error {
	getOrder.IsRetry = true

	createdOrder, err := cli.CreateOrderInPOS(ctx, getOrder)
	if err != nil {
		log.Err(err).Msgf("CreateOrderInPOS error")
		if errors.Is(err, pos.ErrRetry) {
			log.Warn().Msgf("OrderRetry CreateOrderInPOS failed, retry again, orderId: %s", createdOrder.OrderID)
		}
		return err
	}

	log.Info().Msgf("order again sent to pos; pos_id: %s", createdOrder.PosOrderID)

	createdOrder.RetryCount++
	createdOrder.Status = "NEW"
	createdOrder.IsRetry = true

	if err = cli.UpdateOrder(ctx, createdOrder); err != nil {
		log.Err(err).Msgf("UpdateOrder error")
		return err
	}

	log.Info().Msgf("UpdateOrder prev_pos_id: %v; new_pos_id %v; retryCount is: %v; prevOrder.status %v",
		getOrder.PosOrderID, createdOrder.PosOrderID, createdOrder.RetryCount, getOrder.Status)

	if createdOrder.DeliveryService == models.WOLT.String() {
		if err = cli.UpdateOrderStatusInDS(ctx, createdOrder.ID, dto.ACCEPTED); err != nil {
			log.Err(err).Msgf("UpdateOrder status error")
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(Run)
}
