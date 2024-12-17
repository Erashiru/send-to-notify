package main

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/jowi/config"
	"github.com/kwaaka-team/orders-core/core/jowi/managers"
	httpprotocol "github.com/kwaaka-team/orders-core/core/jowi/resource/http"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"os"
)

func setLambdaLevel() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func run() error {
	ctx := context.Background()
	opts, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}

	orderCli, err := order.NewClient()
	if err != nil {
		log.Fatalf("creating new order-core client error: %v", err)
	}

	storeCli, err := store.NewClient(storeModels.Config{
		MongoCli: orderCli.Client(),
	})
	if err != nil {
		log.Fatalf("creating new store-core client error: %v", err)
	}

	jowiManager := managers.NewJowiManager(storeCli, orderCli)

	server := httpprotocol.NewServer(jowiManager, opts)

	switch setLambdaLevel() {
	case false:
		return http.ListenAndServe(":8080", server.Router)
	default:
		lambda.Start(server.GinProxy)
		return nil
	}
}

//	@title						External Client API
//	@version					1.0
//	@host						external-api.kwaaka.com
//	@BasePath					/v1
//	@schemes					https http
//	@query.collection.format	multi
//	@securityDefinitions.apiKey	ApiSecretAuth
//	@in							header
//	@name						Authorization
//	@securityDefinitions.apiKey	ApiTokenAuth
//	@in							header
//	@name						Authorization
//	@description				"Token from cognito"
func main() {
	if err := run(); err != nil {
		log.Fatalf("app run error: %v", err)
	}
}
