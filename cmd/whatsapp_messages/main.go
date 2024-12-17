package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/managers/notifier/whatsapp"
	whatsappConfig "github.com/kwaaka-team/orders-core/pkg/whatsapp/clients"
	whatsappClient "github.com/kwaaka-team/orders-core/pkg/whatsapp/clients/http"
	"github.com/rs/zerolog/log"
)

func main() {
	lambda.Start(ReceiveMessage)
}

func ReceiveMessage(event events.SQSEvent) error {
	log.Info().Msgf("Messages: %v", event.Records)

	ctx := context.Background()
	opts, err := config.LoadConfig(ctx)
	if err != nil {
		log.Err(err).Msgf("Error loading config")
		return err
	}

	for _, record := range event.Records {
		var message whatsapp.Message
		if err = json.Unmarshal([]byte(record.Body), &message); err != nil {
			log.Err(err).Msg("error unmarshalling message")
			return err
		}

		// defining client according to store settings
		var client whatsappConfig.Whatsapp
		if message.InstanceId != "" && message.AuthToken != "" {
			client, err = whatsappClient.NewClient(&whatsappConfig.Config{
				Protocol:  "http",
				AuthToken: message.AuthToken,
				Instance:  message.InstanceId,
				BaseURL:   opts.WhatsAppConfiguration.BaseUrl,
			})
		} else {
			client, err = whatsappClient.NewClient(&whatsappConfig.Config{
				Protocol:  "http",
				AuthToken: opts.WhatsAppConfiguration.AuthToken,
				Instance:  opts.WhatsAppConfiguration.Instance,
				BaseURL:   opts.WhatsAppConfiguration.BaseUrl,
			})
		}
		if err != nil {
			return err
		}

		if err := client.SendMessage(ctx, message.CustomerPhone, message.Message); err != nil {
			log.Err(err).Msg("error sending message via whatsapp")
			return err
		}
	}

	return nil
}
