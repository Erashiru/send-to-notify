package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/service/legal_entity_payment"
	wppService "github.com/kwaaka-team/orders-core/service/whatsapp"
	"github.com/rs/zerolog/log"
)

func run() error {
	ctx := context.Background()
	opts, err := general.LoadConfig(ctx)
	if err != nil {
		log.Err(err).Msgf("[ERROR] get config error")
		return err
	}

	wppService, err := wppService.NewWhatsappService(nil, opts.WhatsAppConfiguration.Instance, opts.WhatsAppConfiguration.AuthToken, opts.WhatsAppConfiguration.BaseUrl, nil, nil, nil, nil, nil)
	if err != nil {
		return err
	}

	legalEntityPaymentService, err := legal_entity_payment.NewService(opts, nil, nil, nil, wppService)
	if err != nil {
		log.Err(err).Msgf("[ERROR] get legal entity payment service error")
		return err
	}

	if err := legalEntityPaymentService.SendPayment(ctx); err != nil {
		log.Err(err).Msgf("[ERROR] send payment to whatsapp error")
		return err
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
