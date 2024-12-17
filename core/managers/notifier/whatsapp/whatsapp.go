package whatsapp

import (
	"encoding/json"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models"
	storecoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"net/url"
	"strings"
)

type WhatsappNotifier struct {
	sqsCli   que.SQSInterface
	queueUrl string
}

func NewWhatsappNotifier(sqsCli que.SQSInterface, queueUrl string) WhatsappNotifier {
	return WhatsappNotifier{
		sqsCli:   sqsCli,
		queueUrl: queueUrl,
	}
}

func (wa WhatsappNotifier) Notify(ctx context.Context, status string, order models.Order, storeGroup storecoreModels.StoreGroup, store storecoreModels.Store) error {
	storeDomain := "https://qr.kwaaka.app/order/"
	if order.DeliveryService == models.KWAAKA_ADMIN.String() {
		storeDomain = "https://status.kwaaka.direct/order/"
	} else if storeGroup.DomainName != "" {
		storeDomain = storeGroup.DomainName + "order/"
	}

	msg := constructStatusChangeMessage(status, order, store, storeDomain)
	if msg == "" {
		return nil
	}
	messageToSQS := Message{
		CustomerPhone: order.Customer.PhoneNumber,
		Message:       msg,
		InstanceId:    store.WhatsappConfig.InstanceId,
		AuthToken:     store.WhatsappConfig.AuthToken,
	}
	message, err := json.Marshal(messageToSQS)
	if err != nil {
		return err
	}

	log.Info().Msgf("queue message body: %s", string(message))

	return wa.sqsCli.SendSQSMessageToFIFO(ctx, wa.queueUrl, string(message), order.ID)
}

func constructStatusChangeMessage(status string, order models.Order, store storecoreModels.Store, domain string) string {
	orderID := strings.Trim(order.ID, string('"'))

	var message string

	switch status {
	case models.ACCEPTED.String():
		if order.DeliveryService == models.KWAAKA_ADMIN.String() {
			var products string
			for _, product := range order.Products {
				products += fmt.Sprintf("%dx %s\n\n", product.Quantity, product.Name)
			}

			message = fmt.Sprintf("✅ Ваш заказ успешно создан и скоро начнет готовиться!\n\nВы заказали:\n\n%sОтследить статус заказа можно по ссылке: %s%s\n\n", products, domain, orderID)
		}
	case models.COOKING_STARTED.String():
		message = fmt.Sprintf("🍽️ Начали готовить ваш заказ.\n\nОтследить статус заказа можно по ссылке: %s%s\n\n", domain, orderID)
	default:
		return ""
	}

	message += fmt.Sprintf("Номер поддержки: %[1]s \nЧат поддержки: https://wa.me/%[1]s", store.StorePhoneNumber)

	return url.QueryEscape(message)
}
