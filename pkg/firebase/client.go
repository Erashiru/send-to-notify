package firebase_client

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/kwaaka-team/orders-core/core/models"
	"google.golang.org/api/option"
)

type MessageService struct {
	messagingClient *messaging.Client
}

func NewFirebaseMessageService(ctx context.Context, configs option.ClientOption) (*MessageService, error) {
	app, err := firebase.NewApp(ctx, nil, configs)
	if err != nil {
		return nil, err
	}

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &MessageService{
		messagingClient: messagingClient,
	}, nil
}

func (ms *MessageService) SendOrderNotification(ctx context.Context, fcmTokens []string, order models.Order) error {
	orderProducts := "Позиции в заказе: "

	for _, product := range order.Products {
		orderProducts += product.Name
		orderProducts += ", "
	}

	msg := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: order.DeliveryService + ": Новый заказ!",
			Body:  orderProducts,
		},
		Tokens: fcmTokens,
	}

	_, err := ms.messagingClient.SendMulticast(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}
