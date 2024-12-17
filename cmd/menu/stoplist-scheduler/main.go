package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/database"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/managers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
	storeCore "github.com/kwaaka-team/orders-core/pkg/store"
	storeDto "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/kwaaka-team/orders-core/service/entity_changes_history"
	entityChangesHistoryModels "github.com/kwaaka-team/orders-core/service/entity_changes_history/models"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	"sync"
)

// receive sqs events -> update stoplist []rst:[]products;
func Run(ctx context.Context, sqsEvent events.SQSEvent) error {

	opts, err := menu.LoadConfig(context.Background(), "", "")
	if err != nil {
		return err
	}

	var (
		restaurantsToUpdateStopList []models.StopListScheduler
		wg                          sync.WaitGroup
		cfg                         dto.Config
	)

	ds, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSDB,
	})
	if err != nil {
		return fmt.Errorf("cannot create datastore %s: %v", opts.DSName, err)
	}

	if err = ds.Connect(cfg.MongoCli); err != nil {
		return fmt.Errorf("cannot connect to datastore: %s", err)
	}

	storeCli, err := storeCore.NewClient(storeDto.Config{})
	if err != nil {
		return err
	}
	stopListMan := managers.NewStopListTransactionManager(ds.StopListTransactionRepository())

	entityChangesHistoryRepo, err := entity_changes_history.NewEntityChangesHistoryMongoRepository(ds.DataBase())
	if err != nil {
		return err
	}

	menuMan := managers.NewMenuManager(opts, ds, ds.MenuRepository(entityChangesHistoryRepo), ds.StoreRepository(), nil, nil, stopListMan, nil, nil, nil, nil, storeCli, nil, menuServicePkg.MongoRepository{}, ds.RestGroupMenuRepository())

	//init manager, write DRY logic code; update aggregatorMenu, POSMenu & aggregatorWEB
	for idx, message := range sqsEvent.Records {

		if err = json.Unmarshal([]byte(message.Body), &restaurantsToUpdateStopList); err != nil {
			return err
		}

		fmt.Printf(" rsts: %+v index %v len %v", restaurantsToUpdateStopList, idx, len(restaurantsToUpdateStopList))

		for _, stopListSetting := range restaurantsToUpdateStopList {
			wg.Add(1)
			go func(wg *sync.WaitGroup, rst models.StopListScheduler) { //errGroup?
				defer wg.Done()

				fmt.Println("run in gorutine rst: ", rst.ID, rst.RstID)
				if err = menuMan.StopListSchedule(ctx, rst, entityChangesHistoryModels.EntityChangesHistoryRequest{
					Author:   "cron",
					TaskType: "cmd/menu/stoplist-scheduler/main.go - Run",
				}); err != nil {
					fmt.Printf("Failed to update schedule: %v", err)
					return
				}

			}(&wg, stopListSetting)
		}
	}

	wg.Wait()

	fmt.Println("done")

	return nil
}

// response for client? []setttings -> by rst_group - []by_time ?
func main() {
	lambda.Start(Run)
}
