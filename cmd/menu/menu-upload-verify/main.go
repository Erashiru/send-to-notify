package main

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/utils"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/rs/zerolog/log"

	menuCli "github.com/kwaaka-team/orders-core/pkg/menu"
	menuCoreModels "github.com/kwaaka-team/orders-core/pkg/menu/dto"
)

func run() error {
	log.Info().Msgf("starting cron")

	menuCoreCli, err := menuCli.New(menuCoreModels.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Only chocofood transactions, if needed add another delivery service
	transactions, err := menuCoreCli.GetProcessingMenuUploadTransactions(context.TODO(), menuCoreModels.GetMenuUploadTransactions{
		DeliveryService: models.CHOCOFOOD.String(),
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	transactionsGlovo, err := menuCoreCli.GetProcessingMenuUploadTransactions(context.TODO(), menuCoreModels.GetMenuUploadTransactions{
		DeliveryService: models.GLOVO.String(),
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	transactions = append(transactions, transactionsGlovo...)

	var wg sync.WaitGroup

	for _, transaction := range transactions {

		wg.Add(1)
		go func(req menuCoreModels.MenuUploadTransaction) {
			defer wg.Done()
			fmt.Printf("Starting store %s - %s\n", req.ID, req.Service)
			trx, err := menuCoreCli.VerifyUploadMenu(context.TODO(), menuCoreModels.MenuUploadVerifyRequest{
				TransactionId: req.ID,
			})

			if err != nil {
				fmt.Println(err)
			}
			utils.Beautify("Finished transaction", trx)
		}(transaction)
	}
	wg.Wait()

	return nil
}

func main() {
	if cmd.IsLambda() {
		log.Info().Msg("Starting lambda")
		lambda.Start(run)
	} else {
		log.Info().Msg("Starting locally")

		log.Info().Msg("Setting environment variables")
		os.Setenv("SECRET_ENV", "StageEnvs")
		os.Setenv("REGION", "eu-west-1")

		if err := run(); err != nil {
			log.Err(err).Msgf("failed run")
		}
	}

}
