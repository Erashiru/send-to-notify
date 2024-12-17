package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/pkg/storeStatus"
	whatsappConfig "github.com/kwaaka-team/orders-core/pkg/whatsapp/clients"
	whatsappClient "github.com/kwaaka-team/orders-core/pkg/whatsapp/clients/http"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	orderServicePkg "github.com/kwaaka-team/orders-core/service/order"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/kwaaka-team/orders-core/service/store/repository/storeclosedtime"
	"github.com/rs/zerolog/log"
)

type app struct {
	targets []string
}

func newApp() *app {
	targets := []string{"glovo", "wolt"}

	return &app{
		targets: targets,
	}
}

func run() error {
	app := newApp()

	ctx := context.Background()

	opts, err := general.LoadConfig(ctx)
	if err != nil {
		return err
	}

	db, err := cmd.CreateMongo(ctx, opts.DSURL, opts.DSDB)
	if err != nil {
		return err
	}

	sqsCli := notifyQueue.NewSQS(sqs.NewFromConfig(opts.AwsConfig))

	var aggFactory aggregator.Factory

	aggFactory, err = aggregator.NewFactory(
		opts.WoltConfiguration.BaseURL,
		opts.GlovoConfiguration.BaseURL, opts.GlovoConfiguration.Token,
		opts.TalabatConfiguration.MiddlewareBaseURL, opts.TalabatConfiguration.MenuBaseUrl,
		opts.Express24Configuration.BaseURL, opts.StarterAppConfiguration.BaseUrl, nil, nil, config.Configuration{},
	)
	if err != nil {
		return err
	}

	telegramService, err := orderServicePkg.NewTelegramService(sqsCli, opts.QueConfiguration.Telegram, opts.NotificationConfiguration, &telegram.Repository{})
	if err != nil {
		return err
	}

	whatsAppClient, err := whatsappClient.NewClient(&whatsappConfig.Config{
		Protocol:  "http",
		AuthToken: opts.WhatsAppConfiguration.AuthToken,
		Instance:  opts.WhatsAppConfiguration.Instance,
		BaseURL:   opts.WhatsAppConfiguration.BaseUrl,
	})
	if err != nil {
		return err
	}

	storeRepository, err := store.NewStoreMongoRepository(db)
	if err != nil {
		return err
	}

	storeFactory, err := store.NewService(storeRepository)
	if err != nil {
		return err
	}

	storeActiveTimeRepository, err := storeclosedtime.NewMongoRepository(db)
	if err != nil {
		return err
	}

	subject := &storeStatus.Subject{}

	status, err := storeStatus.NewStoreStatus(aggFactory, storeFactory, subject, storeActiveTimeRepository, nil)
	if err != nil {
		return err
	}

	whatsappObserver := &storeStatus.WhatsAppObserver{WhatsAppClient: whatsAppClient}
	telegramObserver := &storeStatus.TelegramObserver{TelegramClient: telegramService}
	datastoreObserver := &storeStatus.DatastoreObserver{DatastoreClient: storeActiveTimeRepository}

	subject.AddObserver(whatsappObserver)
	subject.AddObserver(telegramObserver)
	subject.AddObserver(datastoreObserver)

	log.Print("UpdateStoreSchedule started")

	for _, deliveryService := range app.targets {
		if err = status.UpdateStoresSchedule(ctx, deliveryService); err != nil {
			return err
		}
	}

	log.Print("UpdateStoreSchedule finished")

	return nil
}

func main() {
	if cmd.IsLambda() {
		lambda.Start(run)
	} else {
		if err := run(); err != nil {
			log.Err(err).Msgf("run failed")
		}
	}
}
