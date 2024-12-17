package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	netHttp "github.com/kwaaka-team/orders-core/pkg/net-http-client/http"
	"log"
	"os"
)

const (
	baseUrlEnv = "BASE_URL"
)

func run(ctx context.Context, event Event) error {
	cli := netHttp.NewHTTPClient(os.Getenv(baseUrlEnv))

	req := dto.UpdateOrderStatusCronRequest{
		PosTypes: event.PosTypes,
	}

	body, err := json.Marshal(req)
	if err != nil {
		log.Printf("marshal body error: %v", err)
		return err
	}

	status, response, err := cli.Post("/api/update-order-status", body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		log.Printf("cron order status update error: %v", err)
		return err
	}

	if status >= 400 {
		log.Printf("cron order status update finished with http status %d, response %v", status, string(response))
		return fmt.Errorf("status code: %d, response: %v", status, string(response))
	}

	return nil
}

func main() {
	lambda.Start(run)
}

type Event struct {
	PosTypes []string `json:"pos_types"`
}
