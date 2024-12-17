package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/general"
	menuClient "github.com/kwaaka-team/orders-core/pkg/menu"
	"github.com/rs/zerolog/log"

	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
	"os"
)

func main() {
	ctx := context.TODO()

	os.Setenv("REGION", "eu-west-1")
	os.Setenv("SENTRY", "ProdSentry")
	os.Setenv("SECRET_ENV", "ProdEnvs")
	os.Setenv("S3_BUCKET", "kwaaka-menu-files")

	config, err := general.LoadConfig(ctx)
	if err != nil {
		log.Err(err).Msg("Error loading config")
		return
	}

	menuCli, err := menuClient.New(dto.Config{})
	if err != nil {
		log.Err(err).Msg("Error creating menu client")
		return
	}

	awsS3 := s3.New(config.AwsSession)

	trID, err := menuCli.UploadMenu(ctx, dto.MenuUploadRequest{
		StoreId:      "66a782475faef412fa6eaf04",
		MenuId:       "66a889de9463700202a4c4de",
		DeliveryName: "starter_app",
		Sv3:          awsS3,
	})

	fmt.Println("trID: ", trID)

	if err != nil {
		log.Err(err).Msg("Error uploading menu")
		return
	}
}
