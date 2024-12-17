package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/order/dto"
	"github.com/rs/zerolog/log"
)

func run() error {
	ctx := context.TODO()

	cli, err := order.NewClient()
	if err != nil {
		return err
	}

	orders, err := cli.GetActiveOrders(ctx, dto.ActiveOrderSelector{
		PosType: models.RKEEPER7XML.String(),
	})
	if err != nil {
		return err
	}

	for _, order := range orders {
		log.Info().Msgf("manual update status for order with id=%s, current status=%s", order.ID, order.Status)

		if err := cli.ManualUpdateStatus(ctx, order); err != nil {
			log.Info().Msgf("manual update status rkeeper7xml error: %s", err.Error())
			continue
		}
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
