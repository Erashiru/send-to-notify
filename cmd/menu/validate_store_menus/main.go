package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCli "github.com/kwaaka-team/orders-core/pkg/menu"
	menuCliDto "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/service/aws_s3"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/stoplist"
	storeServicePkg "github.com/kwaaka-team/orders-core/service/store"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const validateStoreMenus = "validate_store_menus"

func isSmallOrBurgerKing(name string) bool {
	return strings.Contains(name, "BK") || strings.Contains(name, "SMALL")
}

type application struct {
	s3Service               aws_s3.ServiceImpl
	db                      *mongo.Database
	cfg                     general.Configuration
	menuService             *menuServicePkg.Service
	storeService            storeServicePkg.Service
	stopListService         stoplist.Service
	queueService            notifyQueue.SQSInterface
	ignorePos               map[string]struct{}
	menuCli                 menuCli.Client
	stopListServiceValidate stoplist.Service
}

func getIgnorePos() map[string]struct{} {
	return map[string]struct{}{
		"foodband":     {},
		"CTMAX":        {},
		"yaros":        {},
		"rkeeper":      {},
		"poster":       {},
		"rkeeper7_xml": {},
	}
}

func newApp(ctx context.Context) (*application, error) {

	cfg, err := general.LoadConfig(ctx)
	if err != nil {
		return nil, err
	}

	coreCfg, err := config.LoadConfig(ctx)
	if err != nil {
		return nil, err
	}

	db, err := cmd.CreateMongo(ctx, cfg.DSURL, cfg.DSDB)
	if err != nil {
		return nil, err
	}

	session := session.Must(session.NewSession(&aws.Config{}))

	s3Service := aws_s3.NewS3Service(session)

	sqsCli := notifyQueue.NewSQS(sqs.NewFromConfig(cfg.AwsConfig))

	storeRepository, err := storeServicePkg.NewStoreMongoRepository(db)
	if err != nil {
		return nil, err
	}

	stopListService, err := stoplist.CreateStopListServiceByWebhook(db, coreCfg, 1)
	if err != nil {
		return nil, err
	}

	storeFactory, err := storeServicePkg.NewService(storeRepository)
	if err != nil {
		return nil, err
	}

	menuRepo, err := menuServicePkg.NewMenuMongoRepository(db)
	if err != nil {
		return nil, err
	}

	menuService, err := menuServicePkg.NewMenuService(menuRepo, storeFactory, nil, nil)
	if err != nil {
		return nil, err
	}

	menuCli, err := menuCli.New(menuCliDto.Config{})
	if err != nil {
		return nil, err
	}

	stopListServiceValidate, err := stoplist.CreateStopListServiceByValidate(db, coreCfg, 1)
	if err != nil {
		return nil, err
	}

	return &application{
		db:                      db,
		s3Service:               *s3Service,
		menuService:             menuService,
		storeService:            storeFactory,
		queueService:            sqsCli,
		cfg:                     cfg,
		stopListService:         stopListService,
		ignorePos:               getIgnorePos(),
		menuCli:                 menuCli,
		stopListServiceValidate: stopListServiceValidate,
	}, nil
}

func run() error {
	ctx := context.Background()

	app, err := newApp(ctx)
	if err != nil {
		return err
	}

	if err := app.validateStoreMenus(ctx); err != nil {
		return err
	}

	return nil
}

type Analytics struct {
	RestaurantsCount         int `json:"restaurants_count"`
	ValidRestaurantsCount    int `json:"valid_restaurants_count"`
	NotValidRestaurantsCount int `json:"not_valid_restaurants_count"`
	ValidMenusCount          int `json:"valid_menus_count"`
	NotValidMenusCount       int `json:"not_valid_menus_count"`
}

func (app *application) validateStoreMenus(ctx context.Context) error {
	stores, err := app.storeService.FindAllStores(ctx)
	if err != nil {
		return err
	}

	currentDate := fmt.Sprintf("%v", time.Now().Format("2006-01-02"))

	currentTime := fmt.Sprintf("%v", time.Now().Format("15:04:05"))

	var (
		wg        sync.WaitGroup
		analytics Analytics
		mtx       sync.Mutex
	)

	for _, st := range stores {
		if _, ok := app.ignorePos[st.PosType]; ok {
			continue
		}

		if isSmallOrBurgerKing(st.Name) {
			continue
		}

		if st.PosType == models.PALOMA.String() {
			continue
		}

		wg.Add(1)

		go func(store coreStoreModels.Store) {
			defer wg.Done()

			if !store.IikoCloud.IsExternalMenu {
				_, err := app.menuCli.UpsertMenu(ctx, menuCliDto.MenuGroupRequest{StoreID: store.ID}, validateStoreMenus, false)
				if err != nil {
					log.Err(err).Msgf("upsert menu error, store id: %s", store.ID)
				}
				log.Info().Msgf("upsert menu success in store: %s", store.ID)
			} else {
				log.Info().Msgf("upsert skip is iiko web store: %s", store.ID)
			}

			menusReport, result, err := app.menuService.ValidateStoreMenus(ctx, store)
			if err != nil {
				log.Err(err).Msgf("validate store menus")
				return
			}

			if len(menusReport) != 0 {
				for _, r := range menusReport {
					if err := app.stopListServiceValidate.UpdateStopListForValidateStoreMenus(ctx, store.ID, r.Delivery, r.ProductDetails); err != nil {
						log.Err(err).Msgf("update stoplist by aggregator product ids for validate store menus error, store id: %s, delivery service: %s", store.ID, r.Delivery)
						continue
					}
				}

				if err = app.stopListService.ActualizeStopListByStoreID(ctx, store.ID); err != nil {
					log.Err(err).Msgf("actualize stoplist for restaurant name %s, id %s error", store.Name, store.ID)
				}
			}

			if err = app.uploadMenusReportToS3(menusReport, store, currentDate, currentTime); err != nil {
				log.Err(err).Msgf("upload to s3")
				return
			}

			if result.ValidMenuCount != 0 || result.NotValidMenuCount != 0 {
				mtx.Lock()

				analytics.RestaurantsCount += 1

				if result.ValidMenuCount != 0 {
					analytics.ValidMenusCount += result.ValidMenuCount
				}

				if result.NotValidMenuCount != 0 {
					analytics.NotValidMenusCount += result.NotValidMenuCount
					analytics.NotValidRestaurantsCount += 1
				} else {
					analytics.ValidRestaurantsCount += 1
				}

				mtx.Unlock()
			}
		}(st)
	}

	wg.Wait()

	log.Info().Msgf("Дата: %s\nКоличество активных ресторанов: %d\nКоличество валидных ресторанов: %d\nКоличество невалидных ресторанов: %d\nКоличество валидных меню: %d\nКоличество невалидных меню: %d\nВалидные меню в процентах: %d\n", currentDate, analytics.RestaurantsCount, analytics.ValidRestaurantsCount, analytics.NotValidRestaurantsCount, analytics.ValidMenusCount, analytics.NotValidMenusCount, analytics.ValidMenusCount*100/(analytics.ValidMenusCount+analytics.NotValidMenusCount))

	if err = app.queueService.SendMessage(app.cfg.Telegram, fmt.Sprintf("Дата: %s\nПроцент валидных меню: %v\n\n", currentDate, analytics.ValidMenusCount*100/(analytics.ValidMenusCount+analytics.NotValidMenusCount)), "-957120845", ""); err != nil {
		return err
	}

	return nil
}

func (app *application) uploadMenusReportToS3(details []menuServicePkg.MenuDetail, store coreStoreModels.Store, currentDate, currentTime string) error {
	for _, report := range details {
		var storeName = app.s3Service.ReduceSpaces(app.s3Service.RemoveNonAlphaNumericSymbols(store.Name))

		link := strings.TrimSpace(fmt.Sprintf("s3://%v/validation_v2/%s/%s/%s/%v", os.Getenv(models.S3_BUCKET), currentDate, currentTime, store.PosType, path.Join(storeName, report.ID)))

		if err := app.s3Service.PutObjectFromStruct(link, report, os.Getenv(models.S3_BUCKET), "application/json"); err != nil {
			log.Err(err).Msgf("put object error")
			return err
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
