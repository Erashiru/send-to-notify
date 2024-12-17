package storeStatus

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp/clients"
	"github.com/rs/zerolog/log"
	"net/url"
	"time"
)

type WhatsAppObserver struct {
	WhatsAppClient clients.Whatsapp
}

func (s WhatsAppObserver) Notify(ctx context.Context, restaurant models.Store, externalStoreID, deliveryService string, storeIsOpened bool) error {
	if len(restaurant.Notification.Whatsapp.Receivers) == 0 {
		log.Printf("whats app receiver did not filled in restaurant by id: %s", restaurant.ID)
		return nil
	}
	message := constructStoreClosedToNotify(restaurant, deliveryService, storeIsOpened)
	for _, v := range restaurant.Notification.Whatsapp.Receivers {
		if err := s.WhatsAppClient.SendMessage(ctx, v.PhoneNumber, message); err != nil {
			log.Err(err).Msgf("error occured during sending message to the whats app service ")
		}
	}
	return nil
}

func (s WhatsAppObserver) NotifyStatusReport(ctx context.Context, restaurant models.Store, durations []models.OpenTimeDuration) error {
	return nil
}

func (s WhatsAppObserver) NotifyStatusChange(ctx context.Context, status string, phone string) error {
	message := constructStatusChangeMessage(status)
	if err := s.WhatsAppClient.SendMessage(ctx, phone, message); err != nil {
		return err
	}
	return nil
}

func constructStoreClosedToNotify(store models.Store, deliveryService string, storeIsOpened bool) string {
	location, err := time.LoadLocation(store.Settings.TimeZone.TZ)
	if err != nil {
		log.Err(err).Msgf("location parsing error occured so it will use asia/almaty timezone to send message to whatsapp")
		location, _ = time.LoadLocation("Asia/Almaty")
	}
	closedTime := time.Now().In(location).Format("2006-01-02 15:04")

	var message string
	switch {
	case storeIsOpened:
		message = fmt.Sprintf("Ресторан закрыт и был открыт по автооткрытию\nТочка: %s\nАгрегатор: %s\nГород: %s\nВремя закрытия: %s\n", store.Name, deliveryService, store.City, closedTime)
	default:
		message = fmt.Sprintf("Ресторан закрыт:\nТочка: %s\nАгрегатор: %s\nГород: %s\nВремя закрытия: %s\n", store.Name, deliveryService, store.City, closedTime)
	}

	return url.QueryEscape(message)
}

func constructStatusChangeMessage(status string) string {
	switch status {
	case "COOKING_STARTED":
		return "Ваш заказ готовится, примерное время ожидания - 20 минут"
	case "COOKING_COMPLETE":
		return "Ваш заказ готов!"
	default:
		return "Неизвестный статус заказа, обратитесь в службу поддержки ресторана"
	}
}
