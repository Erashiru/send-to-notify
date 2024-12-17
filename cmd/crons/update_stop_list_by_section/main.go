package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	netHttp "github.com/kwaaka-team/orders-core/pkg/net-http-client/http"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

const (
	baseUrlEnv = "BASE_URL"
)

func run(ctx context.Context, req dto.UpdateStopListBySectionCronRequest) error {
	cli := netHttp.NewHTTPClient(os.Getenv(baseUrlEnv))

	body, err := json.Marshal(req)
	if err != nil {
		log.Printf("marshal body error: %v", err)
		return err
	}

	status, response, err := cli.Post("/api/update-stoplist-by-section", body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		log.Printf("post error: %v", err)
		return err
	}

	if status != http.StatusNoContent {
		log.Printf("cron stoplist by section update status code: %d\nresponse: %v", status, response)
		return fmt.Errorf("status code: %d\nresponse: %v", status, response)
	}

	return nil
}

func main() {
	lambda.Start(run)
}
