package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	menuCli "github.com/kwaaka-team/orders-core/pkg/menu"
	menuCoreModels "github.com/kwaaka-team/orders-core/pkg/menu/dto"
)

func run() error {
	log.Info().Msgf("start create row")

	menuCoreCli, err := menuCli.New(menuCoreModels.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := menuCoreCli.AddRowToAttributeGroup(context.TODO(), "64d4ce8515e1dbec83007fbb"); err != nil {
		return err
	}

	return nil
}

func main() {
	log.Info().Msg("Starting locally")

	log.Info().Msg("Setting environment variables")
	os.Setenv("SECRET_ENV", "StageEnvs")
	os.Setenv("REGION", "eu-west-1")

	if err := run(); err != nil {
		log.Err(err).Msgf("failed run")
	}

}
