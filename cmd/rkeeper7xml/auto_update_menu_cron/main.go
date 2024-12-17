package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/config/general"
	menuClient "github.com/kwaaka-team/orders-core/pkg/menu"
	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	storeServicePkg "github.com/kwaaka-team/orders-core/service/store"
	"github.com/rs/zerolog/log"
)

func run() error {
	ctx := context.TODO()

	config, err := general.LoadConfig(ctx)
	if err != nil {
		return err
	}

	db, err := cmd.CreateMongo(ctx, config.DSURL, config.DSDB)
	if err != nil {
		return err
	}

	storeRepository, err := storeServicePkg.NewStoreMongoRepository(db)
	if err != nil {
		return err
	}

	storeFactory, err := storeServicePkg.NewService(storeRepository)
	if err != nil {
		return err
	}

	var menuRepo menuServicePkg.Repository

	menuRepo, err = menuServicePkg.NewMenuMongoRepository(db)
	if err != nil {
		return err
	}

	menuService, err := menuServicePkg.NewMenuService(menuRepo, storeFactory, nil, nil)
	if err != nil {
		return err
	}

	stores, err := storeFactory.FindStoresByPosType(ctx, "rkeeper7_xml")
	if err != nil {
		return err
	}

	menuCli, err := menuClient.New(dto.Config{})
	if err != nil {
		return err
	}

	awsS3 := s3.New(config.AwsSession)

	for _, store := range stores {
		if store.RKeeper7XML.TradeGroupId == "" || store.RKeeper7XML.Domain == "" {
			continue
		}

		if _, err = menuCli.UpsertMenu(ctx, dto.MenuGroupRequest{
			StoreID: store.ID,
		}, "cron auto-update-aggregator-menu-rkeeper7xml", true); err != nil {
			log.Err(err).Msgf("upsert menu for %s with id %s finished with error", store.Name, store.ID)
			continue
		}

		for _, menuDS := range store.Menus {
			if !menuDS.IsActive {
				continue
			}

			err = menuService.AutoUpdateAggregatorMenu(ctx, store, menuDS.ID)
			if err != nil {
				log.Err(err).Msgf("auto update aggregator menu with id %s finished with error", menuDS.ID)
			}

			if _, err = menuCli.UploadMenu(ctx, dto.MenuUploadRequest{
				StoreId:      store.ID,
				MenuId:       menuDS.ID,
				DeliveryName: menuDS.Delivery,
				Sv3:          awsS3,
			}); err != nil {
				log.Err(err).Msgf("upload menu to aggregator with id %s finished with error", menuDS.ID)
			}
		}
	}

	return nil
}

func main() {
	//os.Setenv("REGION", "eu-west-1")
	//os.Setenv("SENTRY", "ProdSentry")
	//os.Setenv("SECRET_ENV", "ProdEnvs")
	//os.Setenv("S3_BUCKET", "kwaaka-menu-files")

	if cmd.IsLambda() {
		lambda.Start(run)
	} else {
		if err := run(); err != nil {
			log.Err(err).Msgf("failed run")
		}
	}
}
