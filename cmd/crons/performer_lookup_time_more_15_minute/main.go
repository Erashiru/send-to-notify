package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	baseUrlEnv = "BASE_URL"
)

func run() error {
	url := os.Getenv(baseUrlEnv)
	if url == "" {
		err := fmt.Errorf("no base url\n")
		log.Err(err).Msgf("[ERROR], get base url")
		return err
	}

	path := "/api/delivery/performer-lookup-time"
	client := resty.New().SetBaseURL(url)

	var errResponse errors.ErrorResponse
	resp, err := client.R().
		SetError(&errResponse).
		Post(path)
	if err != nil {
		log.Err(err).Msgf("[ERROR] no response")
		return err
	}

	if resp.IsError() {
		log.Err(fmt.Errorf("%v", resp.Error())).Msgf("[ERROR] get error response: %+v", errResponse)
		return fmt.Errorf("get error response, %v", errResponse)
	}

	return nil
}

func main() {
	if cmd.IsLambda() {
		log.Info().Msgf("[INFO] Starting Lambda function")
		lambda.Start(run)
		log.Info().Msgf("[INFO] Lambda function finished work")
	} else {
		log.Info().Msgf("[INFO] Running locally")
		if err := run(); err != nil {
			log.Err(err).Msgf("[ERROR] %s", err.Error())
			return
		}
		log.Info().Msgf("[INFO] Finished running locally")
	}
}
