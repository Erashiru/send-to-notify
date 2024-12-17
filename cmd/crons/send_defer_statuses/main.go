package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/errors"
	"log"
	"os"
)

const (
	baseUrl = "BASE_URL"
)

func main() {
	if cmd.IsLambda() {
		lambda.Start(run)
		log.Printf("that log for checking order of executing of functions after start")
	} else {
		if err := run(context.Background()); err != nil {
			log.Printf("error: %s", err)
			return
		}
	}
}

// TODO: make busy mode settable from request
func run(ctx context.Context) error {
	log.Printf("STARTING SEND-DEFER-ORDER REQUEST")

	cli := resty.New().SetBaseURL(os.Getenv(baseUrl))

	url := "/api/send-defer-orders"

	var (
		errorResp errors.ErrorResponse
	)

	resp, err := cli.R().
		SetContext(ctx).
		SetError(errorResp).
		Post(url)
	if err != nil {
		return err
	}

	if resp.IsError() {
		log.Printf("error while sending request %s", errorResp.Msg)
		return err
	}

	return nil
}
