package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	models3 "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	orderServicePkg "github.com/kwaaka-team/orders-core/service/order"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

const (
	baseUrlEnv    = "BASE_URL"
	intervalDay   = "day"
	intervalQuery = "interval"
)

func main() {
	if cmd.IsLambda() {
		lambda.Start(run)
		log.Printf("that log for checking order of executing of functions after start")
	} else {
		if err := run(context.Background()); err != nil {
			panic(err)
		}
	}
}

func run(ctx context.Context) error {
	opts, err := general.LoadConfig(ctx)
	if err != nil {
		return err
	}

	var stat models.OrderStat

	client := resty.New().
		SetBaseURL(os.Getenv(baseUrlEnv))

	resp, err := client.R().
		SetResult(&stat).
		SetQueryParam(intervalQuery, intervalDay).
		Get("/api/get-stat")
	if err != nil {
		log.Printf("cron order status get error: %v", err)
		return err
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		log.Printf("cron order status update finished with http status %d, response %v", resp.StatusCode(), resp.Status())
		return fmt.Errorf("status code: %d, response: %v", resp.StatusCode(), resp.Status())
	}

	sqsCli := notifyQueue.NewSQS(sqs.NewFromConfig(opts.AwsConfig))
	telegramService, err := orderServicePkg.NewTelegramService(sqsCli, opts.QueConfiguration.Telegram, opts.NotificationConfiguration, &telegram.Repository{})
	if err != nil {
		return err
	}

	err = telegramService.SendMessageToQueue(telegram.OrderStat, models.Order{}, storeModels.Store{}, "", telegram.ConstructOrderReport(telegram.Telegram, stat), "", models3.Product{})
	if err != nil {
		return err
	}

	return nil
}
