package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/pkg/storeStatus"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/order"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/kwaaka-team/orders-core/service/store/repository/storeclosedtime"
	storeGroupServicePkg "github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/rs/zerolog/log"
)

func run() error {
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

	menuRepository, err := menu.NewMenuMongoRepository(db)
	if err != nil {
		return err
	}

	menuService, err := menu.NewMenuService(menuRepository, nil, nil, nil)
	if err != nil {
		return err
	}

	telegramService, err := order.NewTelegramService(sqsCli, opts.QueConfiguration.Telegram, opts.NotificationConfiguration, &telegram.Repository{})
	if err != nil {
		return err
	}

	aggFactory, err := aggregator.NewFactory(
		opts.WoltConfiguration.BaseURL,
		opts.GlovoConfiguration.BaseURL, opts.GlovoConfiguration.Token,
		opts.TalabatConfiguration.MiddlewareBaseURL, opts.TalabatConfiguration.MenuBaseUrl,
		opts.Express24Configuration.BaseURL, opts.StarterAppConfiguration.BaseUrl,
		menuService, nil, config.Configuration{},
	)
	if err != nil {
		return err
	}

	storeRepository, err := store.NewStoreMongoRepository(db)
	if err != nil {
		return err
	}

	storeGroupRepository, err := storeGroupServicePkg.NewMongoRepository(db)
	if err != nil {
		return err
	}

	storeGroupService, err := storeGroupServicePkg.NewService(storeGroupRepository)
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

	status, err := storeStatus.NewStoreStatus(aggFactory, storeFactory, subject, storeActiveTimeRepository, storeGroupService)
	if err != nil {
		return fmt.Errorf("main - fn run - fn NewStoreStatus - %w", err)
	}

	subject.AddObserver(storeStatus.TelegramObserver{TelegramClient: telegramService})

	log.Print("Store Status Report started")

	if err = status.ReportStatuses(ctx); err != nil {
		return err
	}

	log.Print("Store Status Report finished")

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
