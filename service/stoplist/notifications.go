package stoplist

import (
	"context"
	"fmt"
	"strings"
	"time"

	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/rs/zerolog/log"
)

func (s *ServiceImpl) sendToNotify(ctx context.Context, store storeModels.Store, trd menuModels.TransactionData) {
	const stopListQueue = "stoplist-telegram"

	message := getMessage(ctx, store, trd)
	fmt.Println("send to notify 01")
	if err := s.notifyCli.SendMessage(stopListQueue, message, "-1001830413167", store.Telegram.TelegramBotToken); err != nil {
		log.Err(err).Msgf("could not send message %s to sqs chat_id: %s", message, store.Telegram.StopListChatID)
	}
}

func getMessage(ctx context.Context, store storeModels.Store, trd menuModels.TransactionData) string {
	var msg strings.Builder
	msg.WriteString("<b>[❌] Обновление стоп-листов</b>\n")

	msg.WriteString("<b>Ресторан: ")
	msg.WriteString(store.Name)
	msg.WriteString("</b>\n")

	msg.WriteString("<b>Город: ")
	msg.WriteString(store.Address.City)
	msg.WriteString("</b>\n")

	msg.WriteString("<b>Дата создания: ")
	msg.WriteString(time.Now().Format("2006.01.02 15:04:05"))
	msg.WriteString("</b>\n\n")

	msg.WriteString("<b>Аггрегатор: ")
	msg.WriteString(trd.Delivery)
	msg.WriteString("</b>\n")

	msg.WriteString("<b>Аггрегатор ID: ")
	msg.WriteString(trd.StoreID)
	msg.WriteString("</b>\n")

	msg.WriteString("<b>Статус: ")
	msg.WriteString(trd.Status.String())
	msg.WriteString("</b>\n")
	msg.WriteString("<b>Ошибка: ")
	msg.WriteString(trd.Message)
	msg.WriteString("</b>\n")

	return msg.String()
}
