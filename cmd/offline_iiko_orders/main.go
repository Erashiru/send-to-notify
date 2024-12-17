package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/offline_orders"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func run(event events.SQSEvent) error {
	ctx := context.Background()

	log.Info().Msgf("Messages: %v", event.Records)

	opts, err := config.LoadConfig(ctx)

	if err != nil {
		log.Err(err).Msgf("error loadConfig message %s", err.Error())
		return err
	}

	db, err := sql.Open("postgres", opts.PostgreSqlConfiguration.ConnectionString)
	if err != nil {
		log.Err(err).Msgf("error db.Open message %s", err.Error())
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Err(err).Msgf("error db.Ping message %s", err.Error())
		return err
	}

	log.Info().Msgf("Connected to DB")

	var offlineOrders []models.ExtendedOrderEvent
	for _, message := range event.Records {
		var offlineOrder models.ExtendedOrderEvent
		err = json.Unmarshal([]byte(message.Body), &offlineOrder)
		if err != nil {
			log.Err(err).Msgf("error unmarshalling message %s", err.Error())
			return err
		}
		offlineOrders = append(offlineOrders, offlineOrder)
	}

	srv := offline_orders.NewOfflineOrdersService(db)

	err = srv.SaveOrders(ctx, offlineOrders)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(run)
}
