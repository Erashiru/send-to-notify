package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	mongo2 "github.com/kwaaka-team/orders-core/core/database/drivers/mongo"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
)

var deliveryServices = []string{"glovo", "wolt", "express24", "qr_menu", "yandex", "emenu"}
var concurrency = 50

func main() {
	if cmd.IsLambda() {
		lambda.Start(run)
	} else {
		if err := run(); err != nil {
			log.Err(err).Msgf("failed run")
		}
	}
}

func run() error {
	ctx := context.Background()

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}

	sqsCli := que.NewSQS(sqs.NewFromConfig(cfg.AwsConfig))

	db, err := cmd.CreateMongo(ctx, cfg.DSURL, cfg.DSDB)
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

	var anotherBillRepository pos.AnotherBillRepository
	anotherBillRepository, err = pos.NewMongoAnotherBillRepository(db)
	if err != nil {
		return err
	}

	bkOfferRepository, err := mongo2.NewBKOfferRepository2(db)
	if err != nil {
		return err
	}

	posFactory, err := pos.NewFactory(
		anotherBillRepository,
		sqsCli,
		cfg.RetryConfiguration.QueueName,
		cfg.IIKOConfiguration.BaseURL, cfg.IIKOConfiguration.TransportToFrontTimeout,
		cfg.PosterConfiguration.BaseURL,
		cfg.PalomaConfiguration.BaseURL, cfg.PalomaConfiguration.Class,
		cfg.JowiConfiguration.BaseURL, cfg.JowiConfiguration.ApiKey, cfg.JowiConfiguration.ApiSecret,
		cfg.RKeeperConfiguration.RKeeperBaseURL, cfg.RKeeperConfiguration.RKeeperApiKey,
		cfg.BurgerKingConfiguration.BaseURL, bkOfferRepository,
		cfg.RKeeper7XMLConfiguration.LicenseBaseURL,
		cfg.SyrveConfiguration.BaseURL,
		cfg.YarosConfiguration.BaseURL, cfg.YarosConfiguration.InfoSystem, cfg.TillypadConfiguration.BaseUrl, cfg.Ytimes.BaseUrl, cfg.Ytimes.Token, cfg.PosistConfiguration.BaseUrl,
	)
	if err != nil {
		return err
	}

	if err = start(ctx, posFactory, storeFactory, sqsCli, models.IIKO); err != nil {
		return err
	}

	return nil
}

func start(ctx context.Context, posFactory pos.Factory, storeService store.Service, sqsCli que.SQSInterface, posType models.Pos) error {
	stores, err := storeService.FindStoresByPosType(ctx, posType.String())
	if err != nil {
		return err
	}

	stores = filterSleepStores(ctx, posFactory, storeService, stores)
	if len(stores) == 0 {
		return nil
	}

	sendNotifications(stores, sqsCli)

	awakeSleepStores(ctx, posFactory, stores)

	return nil
}

func filterSleepStores(ctx context.Context, posFactory pos.Factory, storeService store.Service, stores []coreStoreModels.Store) []coreStoreModels.Store {
	wp := NewWorkerPool(concurrency, posFactory, storeService)
	wp.Start(ctx)
	wp.Submit(stores)
	return wp.GetResult()
}

func sendNotifications(stores []coreStoreModels.Store, sqsCli que.SQSInterface) {
	bahandiStores := make([]coreStoreModels.Store, 0)
	generalStores := make([]coreStoreModels.Store, 0)

	for i := range stores {
		st := stores[i]
		if isBahandi(st) {
			bahandiStores = append(bahandiStores, st)
		} else {
			generalStores = append(generalStores, st)
		}
	}

	bahandiChatID := "-4092779578"
	generalChatID := "-4043355666"

	if err := sendNotification(sqsCli, bahandiStores, bahandiChatID); err != nil {
		log.Err(err).Msgf("send message to chat %s error", bahandiChatID)
	}
	if err := sendNotification(sqsCli, generalStores, generalChatID); err != nil {
		log.Err(err).Msgf("send message to chat %s error", generalChatID)
	}

}

func sendNotification(sqsCli que.SQSInterface, stores []coreStoreModels.Store, chatID string) error {
	if len(stores) == 0 {
		return nil
	}
	msg := constructMessage(stores)
	if err := sqsCli.SendMessage("stoplist-telegram", msg, chatID, ""); err != nil {
		return err
	}
	return nil
}

func constructMessage(stores []coreStoreModels.Store) string {
	msg := ""
	for i := range stores {
		st := stores[i]
		msg += fmt.Sprintf("Терминал у ресторана %s выключен\n\n", st.Name)
	}
	msg += "<b>DONE</b>\n\n\n"
	return msg
}

func isPosAlive(ctx context.Context, posFactory pos.Factory, storeService store.Service, store coreStoreModels.Store) (bool, error) {

	if isPosIntegrationEnabled := isPosIntegrationEnabled(storeService, store, deliveryServices); !isPosIntegrationEnabled {
		return true, nil
	}

	posService, err := posFactory.GetPosService(models.Pos(store.PosType), store)
	if err != nil {
		return false, err
	}
	isAlive, err := posService.IsAliveStatus(ctx, store)
	if err != nil {
		return false, err
	}

	return isAlive, nil
}

func isPosIntegrationEnabled(storeService store.Service, store coreStoreModels.Store, deliveryServices []string) bool {
	for i := range deliveryServices {
		isSendToPos, err := storeService.IsSendToPos(store, deliveryServices[i])
		if err != nil {
			continue
		}
		if isSendToPos {
			return true
		}
	}
	return false
}

func isBahandi(st coreStoreModels.Store) bool {
	return strings.Contains(st.Name, "bahandi") || strings.Contains(st.Name, "Bahandi")
}

func awakeSleepStores(ctx context.Context, posFactory pos.Factory, stores []coreStoreModels.Store) {
	wp := NewAwakeSleepStoresWorkerPool(concurrency, posFactory)
	wp.Start(ctx)
	wp.Submit(stores)
}

func awakeSleepStore(ctx context.Context, posFactory pos.Factory, store coreStoreModels.Store) error {
	posService, err := posFactory.GetPosService(models.Pos(store.PosType), store)
	if err != nil {
		return errors.Wrapf(err, "get posService error, storeID %s", store.ID)
	}
	err = posService.AwakeTerminal(ctx, store)
	if err != nil {
		return errors.Wrapf(err, "awake sleep store error, storeID %s", store.ID)
	}
	return nil
}
