package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/cmd"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCli "github.com/kwaaka-team/orders-core/pkg/menu"
	menuCoreModels "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	storeCLi "github.com/kwaaka-team/orders-core/pkg/store"
	storeCoreModel "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/rs/zerolog/log"
	"sync"
)

func run() error {

	log.Info().Msgf("starting MatchingValidate cron")

	storeCoreCli, err := storeCLi.NewClient(storeCoreModel.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	menuCoreCli, err := menuCli.New(menuCoreModels.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	stores, err := storeCoreCli.FindStores(context.Background(), storeCoreModel.StoreSelector{
		PosType: "iiko",
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	sess := session.Must(session.NewSession())
	sv3 := s3.New(sess)

	var wg sync.WaitGroup

	m := make(map[string]struct{})
	m["64c8cdccaa50cb3149a4963d"] = struct{}{}
	m["64c8ce5f427d1e7d571f4cc2"] = struct{}{}
	m["64c8df42dbf9ce3f2a92b833"] = struct{}{}
	m["64c8df96dbf9ce3f2a92b836"] = struct{}{}
	m["64c8e03ddbf9ce3f2a92b83f"] = struct{}{}
	m["64c8e03ddbf9ce3f2a92b83f"] = struct{}{}
	m["64c8e07c98f2d1abd1cd80eb"] = struct{}{}
	m["64ca189e421efbd52cf23cbd"] = struct{}{}
	m["64ca192246bf8c7680dd4997"] = struct{}{}
	m["654dbaae14f69dcf44a4e3e6"] = struct{}{}
	m["653746de2f9cda325b8c9a6d"] = struct{}{}

	for _, store := range stores {
		if store.IikoCloud.Key == "" || store.IikoCloud.Key == "new" {
			continue
		}
		if _, ok := m[store.ID]; ok {
			continue
		}

		wg.Add(1)
		go func(cur coreStoreModels.Store) {
			defer wg.Done()

			_, err := menuCoreCli.ValidateMatching(context.Background(), menuCoreModels.MenuValidateRequest{
				StoreID: cur.ID,
			}, sv3)
			if err != nil {
				fmt.Println(err)
				return
			}
		}(store)
	}

	wg.Wait()

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
