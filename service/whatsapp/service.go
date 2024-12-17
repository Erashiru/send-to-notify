package whatsapp

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/order/dto"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp/clients"
	paymentModels "github.com/kwaaka-team/orders-core/service/payment/models"
	paymentDto "github.com/kwaaka-team/orders-core/service/payment/whatsapp/dto"
	"github.com/kwaaka-team/orders-core/service/store"
	storeGroupServicePkg "github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/kwaaka-team/orders-core/service/whatsapp/repository"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	Pending           = "Pending"
	Completed         = "Completed"
	autoReplyCooldown = 3 * time.Hour
)

var ErrWpp = errors.New("whatsapp success status")

type Service interface {
	SendNewsletter(ctx context.Context, restGroupId string, text string, name string) error
	SendMessage(ctx context.Context, to, message, storeId string) error
	GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (paymentModels.WebhookEvent, error)
	SendFilePdf(ctx context.Context, to, fileName, message string, pdfFile []byte) error
	SendMessageFromBaseEnvs(ctx context.Context, to, message string) error
}

type ServiceImpl struct {
	NewsletterRepo    *repository.Repository
	WhatsappClient    clients.Whatsapp
	StoreService      store.Service
	storeGroupService storeGroupServicePkg.Service
	OrderService      order.Client
	redisCli          *redis.Client
	Instance          string
	AuthToken         string
	BaseUrl           string
}

func NewWhatsappService(wsClient clients.Whatsapp, Instance, AuthToken, BaseUrl string, newsletterRepo *repository.Repository, StoreService store.Service, OrderService order.Client, storeGroupService storeGroupServicePkg.Service, redisClient *redis.Client) (Service, error) {
	return &ServiceImpl{
		NewsletterRepo:    newsletterRepo,
		WhatsappClient:    wsClient,
		StoreService:      StoreService,
		OrderService:      OrderService,
		Instance:          Instance,
		AuthToken:         AuthToken,
		BaseUrl:           BaseUrl,
		storeGroupService: storeGroupService,
		redisCli:          redisClient,
	}, nil
}

func (s *ServiceImpl) SendNewsletter(ctx context.Context, restGroupId string, text string, name string) error {
	restGroup, err := s.storeGroupService.GetStoreGroupByID(ctx, restGroupId)
	if err != nil {
		return err
	}

	instanceId, authToken, exist, err := s.getWppSettingsIfExist(ctx, restGroup)
	if err != nil {
		return err
	}
	if exist {
		s.WhatsappClient, err = whatsapp.NewWhatsappClient(&clients.Config{
			Insecure:  true,
			Protocol:  "http",
			Instance:  instanceId,
			AuthToken: authToken,
		})
		if err != nil {
			return err
		}
	}

	recipients, err := s.getRecipientsOfRestaurantGroup(ctx, restGroupId)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(recipients))
	wg.Add(len(recipients))

	for _, recipient := range recipients {
		go func(to, msg string) {
			defer wg.Done()
			if err := s.WhatsappClient.SendMessage(ctx, to, text); err != nil {
				errors <- err
			} else {
				errors <- nil
			}
		}(recipient, text)
	}

	wg.Wait()
	close(errors)

	var sendError error
	for range recipients {
		if err := <-errors; err != nil && sendError == nil {
			sendError = err
		}
	}
	if sendError != nil {
		return sendError
	}

	newsletter := models.Newsletter{
		Name:              name,
		Text:              text,
		RestaurantGroupId: restGroupId,
		Recipients:        recipients,
		Status:            Completed,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	if err := s.NewsletterRepo.CreateNewsletter(ctx, newsletter); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) getWppSettingsIfExist(ctx context.Context, restGroup storeModels.StoreGroup) (string, string, bool, error) {
	for _, storeId := range restGroup.StoreIds {
		store, err := s.StoreService.GetByID(ctx, storeId)
		if err != nil {
			return "", "", false, err
		}
		if store.WhatsappConfig.InstanceId != "" && store.WhatsappConfig.AuthToken != "" {
			return store.WhatsappConfig.InstanceId, store.WhatsappConfig.AuthToken, true, nil
		}
	}

	return "", "", false, nil
}

func (s *ServiceImpl) getRecipientsOfRestaurantGroup(ctx context.Context, restGroupId string) ([]string, error) {
	stores, err := s.StoreService.GetStoresByStoreGroupID(ctx, restGroupId)
	if err != nil {
		return []string{}, err
	}

	var allOrders []models.Order
	for _, store := range stores {
		orders, _, err := s.OrderService.GetOrdersWithFilters(ctx, dto.OrderSelector{
			DeliveryService: "qr_menu",
			StoreID:         store.ID,
		})
		if err != nil {
			return []string{}, err
		}
		allOrders = append(allOrders, orders...)
	}

	uniquePhones := make(map[string]models.Customer)
	for _, order := range allOrders {
		phone := order.Customer.PhoneNumber
		if _, exists := uniquePhones[phone]; !exists {
			uniquePhones[phone] = order.Customer
		}
	}

	var phoneNumbers []string
	for phone := range uniquePhones {
		phoneNumbers = append(phoneNumbers, phone)
	}

	return phoneNumbers, nil
}

func (s *ServiceImpl) SendMessage(ctx context.Context, to, message, storeId string) error {
	wppClient, err := s.initWppClient(ctx, storeId)
	if err != nil {
		return err
	}

	return wppClient.SendMessage(ctx, to, message)
}

func (s *ServiceImpl) initWppClient(ctx context.Context, storeId string) (clients.Whatsapp, error) {

	if storeId != "" {
		st, err := s.StoreService.GetByID(ctx, storeId)
		if err != nil {
			log.Info().Msgf("Error getting store by ID to initialize WhatsApp client: %s, storeId: %s", err, storeId)
			return s.WhatsappClient, nil
		}

		if st.WhatsappConfig.InstanceId != "" && st.WhatsappConfig.AuthToken != "" {
			wppClient, err := whatsapp.NewWhatsappClient(&clients.Config{
				Instance:  st.WhatsappConfig.InstanceId,
				AuthToken: st.WhatsappConfig.AuthToken,
				Protocol:  "http",
				Insecure:  true,
			})
			if err != nil {
				log.Info().Msgf("Error initializing WhatsApp client: %s, storeId: %s", err, storeId)
				return s.WhatsappClient, nil
			}

			return wppClient, nil
		}

		return s.WhatsappClient, nil
	}

	return s.WhatsappClient, nil
}

func (s *ServiceImpl) GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (paymentModels.WebhookEvent, error) {
	ctx := context.Background()

	wppWebhook, ok := r.(paymentDto.WebhookEvent)
	if !ok {
		return paymentModels.WebhookEvent{}, errors.New("casting error")
	}

	var (
		res paymentModels.WebhookEvent
		err error
	)
	switch wppWebhook.EventType {
	case paymentDto.MessageReceived:
		if strings.Contains(wppWebhook.Data.QuoteMsg.Body, "Ответьте на это сообщение “Да”, чтобы подтвердить оплату клиентом.") && wppWebhook.Data.Body == "Да" {
			res, err = s.constructSuccessOrderWebhook(wppWebhook)
			if err != nil {
				return paymentModels.WebhookEvent{}, err
			}
		} else {
			st, err := s.StoreService.GetStoresByWppPhoneNum(ctx, wppWebhook.Data.To)
			if err != nil {
				return paymentModels.WebhookEvent{}, err
			}

			key := fmt.Sprintf("last_message_timestamp:%s:%s", wppWebhook.Data.From, wppWebhook.Data.To)
			lastMessageTimestamp, err := s.redisCli.Get(ctx, key).Result()
			if err != nil && err != redis.Nil {
				return paymentModels.WebhookEvent{}, err
			}

			if lastMessageTimestamp != "" {
				lastMessageTime, err := time.Parse(time.RFC3339, lastMessageTimestamp)
				if err != nil {
					return paymentModels.WebhookEvent{}, err
				}

				if time.Since(lastMessageTime) < autoReplyCooldown {
					return res, ErrWpp
				}
			}

			if err := s.redisCli.Set(ctx, key, time.Now().Format(time.RFC3339), autoReplyCooldown).Err(); err != nil {
				return paymentModels.WebhookEvent{}, nil
			}

			// if we could not find any store, just cancel action, otherwise take st[0], because all the stores
			// that have same phone number must have same setting, but if the st[0]
			// would not have proper setting, but st[1], will have it, it is unlucko(
			if len(st) == 0 {
				return paymentModels.WebhookEvent{}, ErrWpp
			}

			idleMsg, err := s.constructIdleMessage(ctx, wppWebhook.Data.From, st[0])
			if err != nil {
				return paymentModels.WebhookEvent{}, err
			}

			if err := s.SendMessage(ctx, wppWebhook.Data.From, idleMsg, st[0].ID); err != nil {
				return paymentModels.WebhookEvent{}, err
			}

			return paymentModels.WebhookEvent{}, ErrWpp
		}
	}
	return res, nil
}

func (s *ServiceImpl) constructIdleMessage(ctx context.Context, to string, st storeModels.Store) (string, error) {
	parts := strings.Split(to, "@")
	if len(parts) == 0 {
		return "", errors.New("split error")
	}

	if parts[1] != "c.us" {
		return "", errors.New("cannot send message to group")
	}

	msg := "Здравствуйте! 👋🏻\n\nЯ бот, который помогает отслеживать статусы ваших заказов. К сожалению, я не умею отвечать на вопросы :( \n\n"

	stGroup, err := s.storeGroupService.GetStoreGroupByStoreID(ctx, st.ID)
	if err != nil {
		return "", err
	}

	if stGroup.DomainName != "" {
		msg += fmt.Sprintf("Если вы хотите сделать заказ - пожалуйста, воспользуйтесь этой ссылкой: %s", stGroup.DomainName)
	}

	return url.QueryEscape(msg), nil
}

func (s *ServiceImpl) constructSuccessOrderWebhook(r paymentDto.WebhookEvent) (paymentModels.WebhookEvent, error) {
	paymentOrderId, err := s.getPaymentOrderIdFromMsg(r.Data.QuoteMsg.Body)
	if err != nil {
		return paymentModels.WebhookEvent{}, err
	}

	res := paymentModels.WebhookEvent{
		Event: paymentModels.PAYMENT_APPROVED,
		PaymentEvent: paymentModels.PaymentEvent{
			OrderID: paymentOrderId,
		},
		OrderEvent: paymentModels.OrderEvent{
			Status: "PAID",
		},
	}

	return res, nil
}

func (s *ServiceImpl) getPaymentOrderIdFromMsg(msg string) (string, error) {
	re := regexp.MustCompile(`ID Оплаты: (\S+)`)

	match := re.FindStringSubmatch(msg)

	if len(match) > 1 {
		return match[1], nil
	} else {
		return "", errors.New("payment id not found in quoted message")
	}
}

func (s *ServiceImpl) SendFilePdf(ctx context.Context, to, fileName, message string, pdfFile []byte) error {
	defWhatsappCli, err := whatsapp.NewWhatsappClient(&clients.Config{
		Instance:  s.Instance,
		AuthToken: s.AuthToken,
		BaseURL:   s.BaseUrl,
		Insecure:  true,
		Protocol:  "http",
	})
	if err != nil {
		return err
	}

	pdfBase64 := base64.StdEncoding.EncodeToString(pdfFile)

	if err := defWhatsappCli.SendFilePdf(ctx, to, fileName, message, pdfBase64); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) SendMessageFromBaseEnvs(ctx context.Context, to, message string) error {
	defWhatsappCli, err := whatsapp.NewWhatsappClient(&clients.Config{
		Instance:  s.Instance,
		AuthToken: s.AuthToken,
		BaseURL:   s.BaseUrl,
		Insecure:  true,
		Protocol:  "http",
	})
	if err != nil {
		return err
	}
	if err := defWhatsappCli.SendMessage(ctx, to, url.QueryEscape(message)); err != nil {
		return err
	}

	return nil
}
