package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/service/order"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	lambda.Start(run)
}
func run() error {
	ctx := context.Background()

	opts, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}

	db, err := createDB(ctx, opts)
	if err != nil {
		return err
	}

	log.Printf("[INFO] connected to %s", db.Name())

	orderService, err := createOrderService(db)
	if err != nil {
		return err
	}

	err = orderService.UpdateOrdersWithTwoMoreHours(ctx)
	if err != nil {
		return err
	}

	log.Info().Msgf("update order status after 2 more hours finished")

	return nil
}

func createOrderService(db *mongo.Database) (*order.ServiceImpl, error) {

	if db == nil {
		return nil, errors.New("db is nil")
	}

	orderRepo, err := order.NewMongoRepository(db)
	if err != nil {
		return nil, err
	}

	orderService, err := order.NewServiceImpl(orderRepo)
	if err != nil {
		return nil, err
	}

	return orderService, nil
}

func createDB(ctx context.Context, cfg config.Configuration) (*mongo.Database, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DSURL))
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.DSDB), nil
}
