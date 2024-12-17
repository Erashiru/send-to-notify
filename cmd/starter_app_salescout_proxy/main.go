package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/core/config"
	v1 "github.com/kwaaka-team/orders-core/core/salescout_proxy/resources/v1"
	"github.com/kwaaka-team/orders-core/core/salescout_proxy/service"
	"log"
	"net/http"
	"os"
)

func isLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func run() error {
	log.Println("RUNNING STARTERAPP SALESCOUT PROXY")

	ctx := context.Background()
	opts, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}
	salescoutProxyService, err := salescout_proxy.NewKaspiSaleScoutService(opts.KaspiSaleScoutConfiguration.BaseUrl, opts.KaspiSaleScoutConfiguration.Token, opts.KaspiSaleScoutConfiguration.MerchantID)
	if err != nil {
		log.Fatal(err)
	}

	server := v1.NewServer(salescoutProxyService, opts.StarterAppConfiguration.StarterAppSaleScoutProxyToken)
	if isLambda() {
		lambda.Start(server.GinProxy)
		log.Println("saleScout proxy gin lambda started")
	} else {
		if err = http.ListenAndServe(":8081", server.Router); err != nil {
			log.Fatal("listen and serve error: ", err)
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
