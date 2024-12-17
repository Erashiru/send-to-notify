package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/wolt_discount_run/resources/v1"
	woltDiscountRun "github.com/kwaaka-team/orders-core/core/wolt_discount_run/service"
	"log"
)

func run() error {
	log.Println("RUNNING WOLT DISCOUNT RUN")

	ctx := context.Background()
	opts, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}

	service := woltDiscountRun.NewService(opts.AdminBaseURL, opts.WoltDiscountRunToken)

	server := v1.NewServer(service)

	if cmd.IsLambda() {
		lambda.Start(server.SqsProxy)
		log.Println("wolt discount run lambda started")
	} else {
		if err := server.SqsProxy(ctx, events.SQSEvent{}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
