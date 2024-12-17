package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hashicorp/go-multierror"
	"github.com/kwaaka-team/orders-core/cmd"
	netHttp "github.com/kwaaka-team/orders-core/pkg/net-http-client/http"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

const (
	minutesBeforeCheck = 3
	notificationCount  = 1
	baseUrlEnv         = "BASE_URL"
)

type Errors struct {
	Errs *multierror.Error `json:"err"`
}

func main() {
	if cmd.IsLambda() {
		lambda.Start(run)
	} else {
		if err := run(); err != nil {
			log.Err(err).Msgf("run failed: %v", err)
		}
	}
}

func run() error {
	cli := netHttp.NewHTTPClient(os.Getenv(baseUrlEnv))

	path := fmt.Sprintf("/v1/qr-menu/notify-unpaid-customers?minutes_before_check=%d&notification_count=%d", minutesBeforeCheck, notificationCount)
	status, body, err := cli.Get(path, nil, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		log.Err(err).Msgf("error making request to path %s: %v", path, err)
		return err
	}
	if status != http.StatusOK {
		var errorResponse Errors
		if err = json.Unmarshal(body, &errorResponse); err != nil {
			log.Err(err).Msg("error unmarshalling response")
			return err
		}
		if errorResponse.Errs.ErrorOrNil() != nil {
			for _, errMsg := range errorResponse.Errs.Errors {
				log.Err(errMsg).Msg("error from notify all unpaid customers request")
			}
		}
	}

	log.Info().Msg("Send notifications to unpaid customers cron successfully end")

	return nil
}
