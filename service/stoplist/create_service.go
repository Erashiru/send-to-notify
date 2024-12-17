package stoplist

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/kwaaka-team/orders-core/core/config"
	mongo2 "github.com/kwaaka-team/orders-core/core/database/drivers/mongo"
	"github.com/kwaaka-team/orders-core/pkg/que"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateStopListServiceByWebhook(db *mongo.Database, cfg config.Configuration, concurrencyLevel int) (Service, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	storeService, storeGroupService, menuService, aggFactory, posFactory, stopListRepository, concurrencyLevel, sqsCli, err := getServices(db, cfg, concurrencyLevel)
	if err != nil {
		return nil, err
	}

	s, err := NewStopListServicePosWebhook(storeService, storeGroupService, menuService, aggFactory,
		posFactory, stopListRepository, cfg.WoltConfiguration, concurrencyLevel, sqsCli)

	if err != nil {
		return nil, err
	}
	return s, nil
}

func getServices(db *mongo.Database, cfg config.Configuration, concurrencyLevel int) (store.Service, storegroup.Service, *menu.Service, aggregator.Factory, pos.Factory, Repository, int, que.SQSInterface, error) {
	menuRepo, err := menu.NewMenuMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}
	menuService, err := menu.NewMenuService(menuRepo, nil, nil, nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}

	storeRepository, err := store.NewStoreMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}
	storeService, err := store.NewService(storeRepository)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}

	storeGroupRepository, err := storegroup.NewMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}
	storeGroupService, err := storegroup.NewService(storeGroupRepository)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}

	aggFactory, err := aggregator.NewFactory(
		cfg.WoltConfiguration.BaseURL,
		cfg.GlovoConfiguration.BaseURL, cfg.GlovoConfiguration.Token,
		cfg.TalabatConfiguration.MiddlewareBaseURL, cfg.TalabatConfiguration.MenuBaseUrl,
		cfg.Express24Configuration.BaseURL, cfg.StarterAppConfiguration.BaseUrl, menuService, nil, cfg,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}

	var anotherBillRepository pos.AnotherBillRepository
	anotherBillRepository, err = pos.NewMongoAnotherBillRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}

	bkOfferRepository, err := mongo2.NewBKOfferRepository2(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}

	sqsCli := notifyQueue.NewSQS(sqs.NewFromConfig(cfg.AwsConfig))

	posFactory, err := pos.NewFactory(
		anotherBillRepository, sqsCli, cfg.RetryConfiguration.QueueName,
		cfg.IIKOConfiguration.BaseURL, cfg.IIKOConfiguration.TransportToFrontTimeout, cfg.PosterConfiguration.BaseURL,
		cfg.PalomaConfiguration.BaseURL, cfg.PalomaConfiguration.Class, cfg.JowiConfiguration.BaseURL,
		cfg.JowiConfiguration.ApiKey, cfg.JowiConfiguration.ApiSecret, cfg.RKeeperBaseURL, cfg.RKeeperApiKey,
		cfg.BurgerKingConfiguration.BaseURL, bkOfferRepository, cfg.RKeeper7XMLConfiguration.LicenseBaseURL,
		cfg.SyrveConfiguration.BaseURL, cfg.YarosConfiguration.BaseURL, cfg.YarosConfiguration.InfoSystem,
		cfg.TillypadConfiguration.BaseUrl, cfg.Ytimes.BaseUrl, cfg.Ytimes.Token, cfg.PosistConfiguration.BaseUrl)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}

	stopListRepository, err := NewStoplistTransactionMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, 0, nil, err
	}

	return storeService, storeGroupService, menuService, aggFactory, posFactory, stopListRepository, concurrencyLevel, sqsCli, nil
}

func CreateStopListServiceByValidate(db *mongo.Database, cfg config.Configuration, concurrencyLevel int, sqsCli que.SQSInterface) (Service, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	storeService, storeGroupService, menuService, aggFactory, posFactory, stopListRepository, concurrencyLevel, _, err := getServices(db, cfg, concurrencyLevel)
	if err != nil {
		return nil, err
	}

	s, err := NewStopListServiceValidate(storeService, storeGroupService, menuService, aggFactory, posFactory, stopListRepository, concurrencyLevel, cfg.WoltConfiguration, sqsCli)
	if err != nil {
		return nil, err
	}

	return s, nil
}
