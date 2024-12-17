package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/errors"
	"log"
	"os"
)

const baseUrlEnv = "BASE_URL"

func run(ctx context.Context) error {
	path := "/api/report/delivery-dispatcher-price"

	baseURL := os.Getenv(baseUrlEnv)
	if baseURL == "" {
		log.Printf("[ERROR] Environment variable %s is not set", baseUrlEnv)
		return fmt.Errorf("environment variable %s is not set", baseUrlEnv)
	}

	client := resty.New().SetBaseURL(baseURL)

	var errResponse errors.ErrorResponse

	resp, err := client.R().SetError(&errResponse).Post(path)
	if err != nil {
		log.Printf("[ERROR] Failed to send a request to %s, %s", baseURL, err)
		return err
	}

	if resp.IsError() {
		log.Printf("[ERROR] API response error: %v", resp.Error())
		return fmt.Errorf("%s", resp.Error())
	}

	return nil
}

func main() {
	if cmd.IsLambda() {
		log.Printf("[INFO] Starting Lambda function")
		lambda.Start(run)
		log.Printf("[INFO] Lambda function finished work")
	} else {
		log.Printf("[INFO] Running locally")
		if err := run(context.Background()); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
		log.Printf("[INFO] Finished running locally")
	}
}
