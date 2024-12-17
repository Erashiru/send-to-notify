package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	baseUrlEnv = "BASE_URL"
)

func main() {
	lambda.Start(run)
}

func run(ctx context.Context) error {
	client := resty.New().
		SetBaseURL(os.Getenv(baseUrlEnv))

	resp, err := client.R().
		Post("/api/kaspi-salescout/create-order")

	if err != nil {
		log.Printf("cron order status get error: %v", err)
		return err
	}

	if resp.StatusCode() >= 400 {
		log.Printf("cron kaspi salescout finished with http status %d, response %v", resp.StatusCode(), resp.Status())
		return fmt.Errorf("status code: %d, response: %v", resp.StatusCode(), resp.Status())
	}

	return nil
}
