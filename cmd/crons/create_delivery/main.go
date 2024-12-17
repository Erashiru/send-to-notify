package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	baseUrlEnv  = "BASE_URL"
	indriveTime = "INDRIVE_TIME"
)

func main() {
	if cmd.IsLambda() {
		lambda.Start(run)
		log.Printf("that log for checking order of executing of functions after start")
	} else {
		if err := run(context.Background()); err != nil {
			panic(err)
		}
	}
}

func run(ctx context.Context) error {
	client := resty.New().
		SetBaseURL(os.Getenv(baseUrlEnv))

	query := os.Getenv(indriveTime)

	resp, err := client.R().
		SetQueryParam("call_time", query).
		Post("/api/3pl/create-delivery")
	if err != nil {
		log.Printf("cron create indrive error: %v", err)
		return err
	}

	if resp.IsError() {
		log.Printf("cron create indrive successfully ended, check logs of integration api")
	}

	return nil
}
