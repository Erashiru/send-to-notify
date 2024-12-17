package managers

import (
	"context"
	"encoding/json"
	"fmt"
	generalConfig "github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	models3 "github.com/kwaaka-team/orders-core/core/menu/models"
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/service/iiko/clients/order"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	pkg "github.com/kwaaka-team/orders-core/pkg/iiko"
	"github.com/kwaaka-team/orders-core/pkg/iiko/clients"
	orderCoreCli "github.com/kwaaka-team/orders-core/pkg/order"
	orderModels "github.com/kwaaka-team/orders-core/pkg/order/dto"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	models2 "github.com/kwaaka-team/orders-core/service/error_solutions/models"
	telegramSvc "github.com/kwaaka-team/orders-core/service/order"
	pos2 "github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/stoplist"
	"github.com/kwaaka-team/orders-core/service/store"
	"strings"
	"time"

	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
	menuCoreCli "github.com/kwaaka-team/orders-core/pkg/menu"
	"github.com/kwaaka-team/orders-core/pkg/notify"
	notifyModels "github.com/kwaaka-team/orders-core/pkg/notify/dto"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"

	storeCore "github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/rs/zerolog/log"
)

type Event interface {
	WebhookEvent(ctx context.Context, token string, webhooks models.WebhookEvents, pos models.Pos) ([]string, error)
	GetCustomerDiscounts(ctx context.Context, storeId, phoneNumber string) (coreOrderModels.Balance, error)
	GetDiscountHistory(ctx context.Context, storeId, phoneNumber string) (coreOrderModels.Transactions, error)
	GetDiscountsForStore(ctx context.Context, storeID string) (models.StoreDiscountsResponse, error)
}

type eventImpl struct {
	storeCoreCli             storeCore.Client
	menuCoreCli              menuCoreCli.Client
	orderCoreCli             orderCoreCli.Client
	orderCli                 order.Order
	notifyCli                notify.Client
	notifyQueue              notifyQueue.SQSInterface
	configIIKOBaseURL        string
	configRetryQueueName     string
	configRetryCount         generalConfig.RetryCount
	configSQSOfflineOrderURL string
	stopListService          stoplist.Service
	storeService             store.Service
	telegramService          telegramSvc.TelegramService
	errorSolution            error_solutions.Service
}

func NewEvent(
	storeCoreCli storeCore.Client,
	menuCoreCli menuCoreCli.Client,
	notifyCli notify.Client,
	orderCoreCli orderCoreCli.Client,
	notifyQueue notifyQueue.SQSInterface,
	configIIKOBaseURL string,
	configRetryQueueName string,
	configRetryCount generalConfig.RetryCount,
	configIntegrationBaseURL string,
	configSQSOfflineOrderURL string,
	stopListService stoplist.Service,
	storeService store.Service,
	telegramService telegramSvc.TelegramService,
	errorSolution error_solutions.Service,
) Event {
	return &eventImpl{
		storeCoreCli:             storeCoreCli,
		menuCoreCli:              menuCoreCli,
		orderCoreCli:             orderCoreCli,
		notifyCli:                notifyCli,
		orderCli:                 order.NewOrderClient(configIntegrationBaseURL),
		notifyQueue:              notifyQueue,
		configIIKOBaseURL:        configIIKOBaseURL,
		configRetryQueueName:     configRetryQueueName,
		configRetryCount:         configRetryCount,
		configSQSOfflineOrderURL: configSQSOfflineOrderURL,
		stopListService:          stopListService,
		storeService:             storeService,
		telegramService:          telegramService,
		errorSolution:            errorSolution,
	}
}

func (man *eventImpl) WebhookEvent(ctx context.Context, token string, webhooks models.WebhookEvents, pos models.Pos) ([]string, error) {

	if len(webhooks) == 0 {
		return nil, fmt.Errorf("get webhooks err: %w", errors.ErrEmpty)
	}

	var (
		success = make([]string, 0, len(webhooks))
		errs    = make([]string, 0, len(webhooks))
		results = make([]string, 0, len(webhooks))
	)

	for _, webhook := range webhooks {

		switch webhook.EventType {
		case models.StopListUpdate:
			if err := man.updateStoplist(ctx, token, webhook); err != nil {
				errs = append(errs, webhook.EventType.String())
				continue
			}
			success = append(success, webhook.EventType.String())
		case models.DeliveryOrderUpdate, models.DeliveryOrderError:
			if err := man.orderStatus(ctx, webhook, pos); err != nil {
				errs = append(errs, webhook.EventType.String())
				continue
			}
			success = append(success, webhook.EventType.String())
		case models.TableOrderUpdate:
			store, err := man.storeCoreCli.FindStore(ctx, storeModels.StoreSelector{
				PosOrganizationID: webhook.OrganizationID,
			})
			if err != nil {
				log.Err(err).Msgf("Couldn't Find Store by organization id: %v", webhook.OrganizationID)
				return nil, err
			}
			if err := man.sendTableOrder(ctx, webhook.EventInfo.Order, webhook.EventInfo.ID, store.ID, store.Name, store.IikoCloud.Key, store.IikoCloud.OrganizationID); err != nil {
				errs = append(errs, webhook.EventType.String())
				log.Info().Msgf("Sending offline order to SQS error: %v ", err)
			}
			if err := man.orderTableStatus(ctx, webhook); err != nil {
				errs = append(errs, webhook.EventType.String())
				continue
			}
			success = append(success, webhook.EventType.String())
		default:
			errs = append(errs, errors.ErrEventType.Error())
		}
	}

	if len(success) == 0 && len(errs) == 0 {
		return nil, errors.ErrEmpty
	}

	results = append(results,
		fmt.Sprintf("Success webhook types %s", success),
		fmt.Sprintf("Errors webhook types %s", errs))

	return results, nil
}

func (man *eventImpl) GetCustomerDiscounts(ctx context.Context, storeId, phoneNumber string) (coreOrderModels.Balance, error) {
	store, err := man.storeCoreCli.FindStore(ctx, storeModels.StoreSelector{
		ID: storeId,
	})
	if err != nil {
		return coreOrderModels.Balance{}, err
	}

	cli, err := pkg.NewClient(&clients.Config{
		Protocol: "http",
		BaseURL:  man.configIIKOBaseURL,
		ApiLogin: store.IikoCloud.Key,
	})
	if err != nil {
		return coreOrderModels.Balance{}, err
	}

	customer, err := cli.GetCustomerInfo(ctx, models.GetCustomerInfoRequest{
		Phone:          phoneNumber,
		Type:           "phone",
		OrganizationId: store.IikoCloud.OrganizationID,
	})
	if err != nil {
		return coreOrderModels.Balance{}, err
	}

	var res coreOrderModels.Balance
	for _, balance := range customer.WalletBalances {
		if balance.Id == store.IikoCloud.DiscountBalanceId {
			res.Balance = balance.Balance
			return res, nil
		}
	}

	return coreOrderModels.Balance{}, errors.ErrNotFound
}

func (man *eventImpl) GetDiscountHistory(ctx context.Context, storeId, phoneNumber string) (coreOrderModels.Transactions, error) {
	store, err := man.storeCoreCli.FindStore(ctx, storeModels.StoreSelector{
		ID: storeId,
	})
	if err != nil {
		return coreOrderModels.Transactions{}, err
	}

	cli, err := pkg.NewClient(&clients.Config{
		Protocol: "http",
		BaseURL:  man.configIIKOBaseURL,
		ApiLogin: store.IikoCloud.Key,
	})
	if err != nil {
		return coreOrderModels.Transactions{}, err
	}

	customer, err := cli.GetCustomerInfo(ctx, models.GetCustomerInfoRequest{
		Phone:          phoneNumber,
		Type:           "phone",
		OrganizationId: store.IikoCloud.OrganizationID,
	})
	if err != nil {
		return coreOrderModels.Transactions{}, err
	}

	transactions, err := cli.GetCustomerTransactions(ctx, models.GetTransactionInfoReq{
		CustomerId:     customer.Id,
		PageSize:       10000,
		OrganizationId: store.IikoCloud.OrganizationID,
	})
	if err != nil {
		return coreOrderModels.Transactions{}, err
	}

	var resp coreOrderModels.Transactions
	for _, transaction := range transactions.Transactions {
		switch transaction.TypeName {
		case "PayFromWallet", "RefillWalletFromOrder":
			resp.Transactions = append(resp.Transactions, coreOrderModels.Transaction{
				TypeName:      transaction.TypeName,
				BalanceBefore: transaction.BalanceBefore,
				BalanceAfter:  transaction.BalanceAfter,
			})
		}
	}

	return resp, nil
}

func (man *eventImpl) GetDiscountsForStore(ctx context.Context, storeID string) (models.StoreDiscountsResponse, error) {
	store, err := man.storeCoreCli.FindStore(ctx, storeModels.StoreSelector{
		ID: storeID,
	})
	if err != nil {
		log.Err(err).Msgf("can not find store with restaurant_id: %s", storeID)
		return models.StoreDiscountsResponse{}, fmt.Errorf("event.go - fn GetDiscountsForStore - fn FindStore: %w", err)
	}

	cli, err := pkg.NewClient(&clients.Config{
		Protocol: "http",
		BaseURL:  man.configIIKOBaseURL,
		ApiLogin: store.IikoCloud.Key,
	})
	if err != nil {
		log.Err(err).Msgf("can not initialize iiko client, base_url: %s, api_login: %s", man.configIIKOBaseURL, store.IikoCloud.Key)
		return models.StoreDiscountsResponse{}, fmt.Errorf("event.go - fn GetDiscountsForStore - fn pkg.NewClient: %w", err)
	}

	discounts, err := cli.GetDiscounts(ctx, store.IikoCloud.OrganizationID)
	if err != nil {
		log.Err(err).Msgf("coundn't extract discounts from iiko using organization_id: %s", store.IikoCloud.OrganizationID)
		return models.StoreDiscountsResponse{}, fmt.Errorf("event.go - fn GetDiscountsForStore - fn cli.GetDiscounts: %w", err)
	}

	return discounts, nil
}

func (man *eventImpl) updateStoplist(ctx context.Context, token string, webhook models.WebhookEvent) error {

	log.Info().Msgf("webhook update stop list in restaurant token %s", token)

	if webhook.OrganizationID != "" {
		stores, err := man.storeService.GetStoresByIIKOOrganizationId(ctx, webhook.OrganizationID)
		if err == nil {
			for _, store := range stores {
				if store.Token != token {
					if err = man.stopListService.ActualizeStopListByToken(ctx, store.Token); err != nil {
						log.Err(err).Msgf("error update stop list by organization id by store token %s", store.Token)
						continue
					}

					log.Info().Msgf("success update stoplist for different token, but similar organization id: name=%s, token=%s", store.Name, store.Token)
				}
			}
		}
	}

	if err := man.stopListService.ActualizeStopListByToken(ctx, token); err != nil {
		log.Err(err).Msgf("error update stop list in token %s", token)
		return err
	}

	return nil
}

func (man *eventImpl) orderStatus(ctx context.Context, webhook models.WebhookEvent, pos models.Pos) error {

	log.Info().Msgf("webhook event type %s event info %+v", webhook.EventType, webhook.EventInfo)

	getOrder, err := man.orderCoreCli.GetOrder(ctx, orderModels.OrderSelector{
		PosOrderID: webhook.EventInfo.ID,
	})
	if err != nil {
		log.Trace().Err(err).Msgf("pos order id not found in database: %v", webhook.EventInfo.ID)
		return err
	}

	if webhook.EventInfo.CreationsStatus == CreationStatusError {

		storeCli, err := storeCore.NewClient(storeModels.Config{})
		if err != nil {
			return err
		}

		st, err := storeCli.FindStore(ctx, storeModels.StoreSelector{
			ID: getOrder.RestaurantID,
		})
		if err != nil {
			log.Err(err).Msgf("can not find store with restaurant_id: %s", getOrder.RestaurantID)
			return err
		}

		cli, err := pkg.NewClient(&clients.Config{
			Protocol: "http",
			BaseURL:  man.configIIKOBaseURL,
			ApiLogin: st.IikoCloud.Key,
		})
		if err != nil {
			log.Err(err).Msgf("can not initialize iiko client, base_url: %s, api_login: %s", man.configIIKOBaseURL, st.IikoCloud.Key)
			return err
		}

		retrievedOrder, err := cli.RetrieveDeliveryOrder(ctx, st.IikoCloud.OrganizationID, getOrder.PosOrderID)
		if err != nil {
			log.Err(err).Msgf("can not retrieve iiko order with organization_id: %s, pos_order_id: %s", st.IikoCloud.OrganizationID, getOrder.PosOrderID)
			return err
		}

		log.Info().Msgf("retrieved order: %+v", retrievedOrder)

		switch retrievedOrder.CreationStatus {
		case "Error":
			log.Info().Msgf("wh: error msg: %v count order %v", webhook.EventInfo.Error.Message, getOrder.RetryCount)
			//write error message
			getOrder.CreationResult.ErrorDescription = webhook.EventInfo.Error.Message

			getOrder.Errors = append(getOrder.Errors, coreOrderModels.Error{
				CreatedAt: time.Now().UTC(),
				Code:      0,
				Message:   webhook.EventInfo.Error.Message,
			})

			restGroup, err := storeCli.FindStoreGroup(ctx, selector.StoreGroup{StoreIDs: []string{st.ID}})
			if err != nil {
				return err
			}
			count := restGroup.RetryCount

			//for retry case; else status = failed - send msg to tg
			if webhook.EventInfo.Error.Message == TimeoutErr && getOrder.RetryCount <= count {
				log.Info().Msgf("timeout case, count <=%d", count)
				getOrder.Status = "NEW"
				webhook.EventInfo.Order.Status = "NEW"
			} else {
				getOrder.Status = "FAILED"
			}

			getOrder, err := man.SetErrorSolutionForError(ctx, getOrder, retrievedOrder, st)
			if err != nil {
				return err
			}

			log.Info().Msgf("getOrder.FailReason.Code: %s, getOrder.FailReason.Message: %s", getOrder.FailReason.Code, getOrder.FailReason.Message)
			//update order in DB
			if err = man.orderCoreCli.UpdateOrder(ctx, getOrder); err != nil {
				log.Trace().Err(err).Msgf("UpdateOrder error ")
				return err
			}

			if webhook.EventInfo.Error.Message == TimeoutErr && getOrder.RetryCount <= count {
				log.Info().Msgf("timeout case: run RETRY")
				if err = man.OrderRetry(ctx, webhook); err != nil {
					log.Trace().Err(err).Msgf("OrderRetry error")
					return err
				}
			} else {
				log.Info().Msg("webhook event info is error, send to bitrix")
				go man.sendNotification(context.Background(), webhook)
				webhook.EventInfo.Order.Status = "Error"

				if err = man.telegramService.SendMessageToRestaurant(telegram.CreateOrder, getOrder, st, retrievedOrder.ErrorInfo.Message); err != nil {
					log.Warn().Msgf("(WebhookEvent -> orderStatus) error sending message to restaurant: %v", err)
				}
			}
			log.Info().Msgf("order INFO for UpdateOrder in DB %+v orderDB.status %v event.Status %v ",
				getOrder.CreationResult.ErrorDescription, getOrder.Status, webhook.EventInfo.Order.Status)
		default:
			log.Info().Msgf("fake creation timeout")
			webhook.EventInfo.CreationsStatus = "Success"
			webhook.EventInfo.Order.Status = "WaitCooking"
			webhook.EventInfo.Error.Message = ""
		}

	}

	switch getOrder.DeliveryService {
	case models.WOLT.String(), models.GLOVO.String(), models.QRMENU.String(), models.YANDEX.String(), models.MOYSKLAD.String(), models.EMENU.String(), models.EXPRESS24.String(), models.KWAAKA_ADMIN.String(), models.STARTERAPP.String():
		if err = man.orderCoreCli.UpdateOrderStatus(ctx, webhook.EventInfo.ID, pos.String(), webhook.EventInfo.Order.Status, webhook.EventInfo.Error.Message); err != nil {
			log.Trace().Err(err).Msgf("updating %v status unsuccessful, id=%v, status=%v", getOrder.DeliveryService, webhook.EventInfo.ID, webhook.EventInfo.Order.Status)
			return err
		}
	default:
		if err = man.orderCli.OrderStatus(ctx, *webhook.EventInfo); err != nil {
			return err
		}
	}

	log.Info().Msgf("done UpdateOrderStatus, status for update %v", webhook.EventInfo.Order.Status)

	return nil
}

func (man *eventImpl) OrderRetry(ctx context.Context, webhook models.WebhookEvent) error {
	log.Info().Msgf("webhook error is creation timeout, queue name: %v", man.configRetryQueueName)
	if err := man.notifyQueue.SendSQSMessage(ctx, man.configRetryQueueName, webhook.EventInfo.ID); err != nil {
		log.Trace().Err(err).Msgf("SendSQSMessage creation-timeout error: %v", webhook.EventInfo.ID)
		return err
	}
	return nil
}

func (man *eventImpl) orderTableStatus(ctx context.Context, webhook models.WebhookEvent) error {

	log.Info().Msgf("webhook event type %s event info %v", webhook.EventType, webhook.EventInfo)

	if webhook.EventInfo.CreationsStatus == CreationStatusError {
		log.Info().Msg("webhook event info is error")
		go man.sendNotification(context.Background(), webhook)
	}

	if err := man.orderCli.OrderTableStatus(ctx, *webhook.EventInfo); err != nil {
		return err
	}

	return nil
}

func (man *eventImpl) sendNotification(ctx context.Context, webhook models.WebhookEvent) {

	if err := skipErrors(webhook.EventInfo.Error.Message); err != nil {
		log.Info().Msg("webhook event error is creation timeout")
		return
	}

	log.Info().Msgf("send notification to bitrix %s", webhook.EventType)

	restaurant, err := man.storeCoreCli.FindStore(ctx, storeModels.StoreSelector{
		PosOrganizationID: webhook.EventInfo.OrganizationID,
	})
	if err != nil {
		log.Err(err).Msgf("get restaurant by pos org id %s", webhook.EventInfo.OrganizationID)
		return
	}

	code := "N/A"

	title := fmt.Sprintf("Ошибка в ресторане %s, Адрес: г.%s, %s\n", restaurant.Name, restaurant.Address.City, restaurant.Address.Street)

	description := fmt.Sprintf("ID ресторана: %s\n Номер Заказа: [%s]\n"+
		"ID организации: %s\n Ошибка[%s]:\n"+
		"[MESSAGE] - %s\n", restaurant.ID, code, webhook.EventInfo.OrganizationID, webhook.EventInfo.Error.Code, webhook.EventInfo.Error.Message)

	_, err = man.notifyCli.SendNotification(ctx, notifyModels.Message{
		Title:       title,
		Description: description,
		Services: []notifyModels.Service{
			notifyModels.BITRIX,
			notifyModels.CLICKUP,
		},
		TaskList: notifyModels.ORDER_ERROR,
	})
	if err != nil {
		log.Err(err).Msg("add task in bitrix/clickup error")
		return
	}
}

func skipErrors(msg string) error {

	var res map[string]interface{}
	if err := json.NewDecoder(strings.NewReader(msg)).Decode(&res); err != nil {
		if msg == TimeoutErr {
			log.Info().Msgf("skip creation timeout error")
			return errors.ErrTimeout
		}
		log.Err(err).Msg("unmarshall error")
		return nil
	}

	switch res["message"].(type) {
	case string:
		if res["message"] == TimeoutErr {
			log.Info().Msgf("skip creation timeout error")
			return errors.ErrTimeout
		}
	}

	return nil
}

func (man *eventImpl) sendTableOrder(ctx context.Context, order models.OrderEvent, eventId, storeId, storeName, storeIikoCloudKey, organizationId string) error {
	extendedOrder := models.ExtendedOrderEvent{
		Order:     order,
		StoreId:   storeId,
		StoreName: storeName,
		EventId:   eventId,
	}

	if order.Customer.Type == "regular" {
		cli, err := pkg.NewClient(&clients.Config{
			Protocol: "http",
			BaseURL:  man.configIIKOBaseURL,
			ApiLogin: storeIikoCloudKey,
		})
		if err != nil {
			return err
		}
		customerInfo, err := cli.GetCustomerInfo(ctx, models.GetCustomerInfoRequest{
			Phone:          order.Phone,
			Type:           "phone",
			OrganizationId: organizationId,
		})
		if err != nil {
			return err
		}
		extendedOrder.RegularCustomerInfo = customerInfo
	}

	orderJson, err := json.Marshal(extendedOrder)
	if err != nil {
		return err
	}
	log.Trace().Msgf("%s", string(orderJson))
	err = man.notifyQueue.SendSQSMessage(ctx, man.configSQSOfflineOrderURL, string(orderJson))
	if err != nil {
		return err
	}
	return nil
}

func (man *eventImpl) fillFailReason(ctx context.Context, order coreOrderModels.Order, store coreStoreModels.Store) (models2.ErrorSolution, coreOrderModels.Order, bool, error) {

	errSolution, sendStopListStatus, err := man.errorSolution.GetErrorSolutionByCode(ctx, store, order.FailReason.Code)
	if err != nil {
		return models2.ErrorSolution{}, coreOrderModels.Order{}, false, err
	}

	order.FailReason.BusinessName, order.FailReason.Reason, order.FailReason.Solution = errSolution.BusinessName, errSolution.Reason, errSolution.Solution

	return errSolution, order, sendStopListStatus, nil
}

func (man *eventImpl) SetErrorSolutionForError(ctx context.Context, order coreOrderModels.Order, retrievedOrder models.RetrieveOrder, st coreStoreModels.Store) (coreOrderModels.Order, error) {

	var product models3.Product

	errorSolutions, err := man.errorSolution.GetAllErrorSolutions(ctx)
	if err != nil {
		log.Err(err).Msgf("eventImpl manager error: GetAllErrorSolutions")
		return coreOrderModels.Order{}, err
	}

	order.FailReason.Code = pos2.MatchingCodes(retrievedOrder.ErrorInfo.Message+retrievedOrder.ErrorInfo.Description, errorSolutions)
	order.FailReason.Message = retrievedOrder.ErrorInfo.Message
	errSolution, order, sendStopListStatus, err := man.fillFailReason(ctx, order, st)
	if err != nil {
		log.Err(err).Msg("fillFailReason error")
		return coreOrderModels.Order{}, err
	}

	//если true, тогда получаем продукт и ставим на стоп во всех доступных агрегаторах
	var productID string
	if sendStopListStatus {

		productID = pos2.GetProductIDFromRegexp(retrievedOrder.ErrorInfo.Message, errSolution)
		productErrorCodes := map[string]bool{
			"21": true, "4": true, "1": true,
		}
		attributeErrorCodes := map[string]bool{
			"25": true, "5": true, "7": true, "21": true, "27": true, "28": true,
		}
		if len(productID) > 0 {
			var err error
			switch {
			case productErrorCodes[errSolution.Code]:
				err = man.stopListService.UpdateStopListByPosProductID(ctx, false, st.ID, productID)
				if err == nil {
					log.Info().Msgf("successfully put product with id: %s to stop with error solution code: %s for store_id : %s", productID, errSolution.Code, st.ID)
				}
				for _, orderProduct := range order.Products {
					if orderProduct.ID == productID {
						product.ExtID = orderProduct.ID
						product.Name = append(product.Name, models3.LanguageDescription{
							Value: orderProduct.Name,
						})
					}
				}

			case attributeErrorCodes[errSolution.Code]:
				err = man.stopListService.UpdateStopListByAttributeID(ctx, false, st.ID, productID)
				if err == nil {
					log.Info().Msgf("successfully put attribute with id: %s to stop with error solution code: %s for store_id : %s", productID, errSolution.Code, st.ID)
				}
				for _, orderProduct := range order.Products {
					if len(orderProduct.Attributes) > 0 {
						for _, orderAttribute := range orderProduct.Attributes {
							if orderAttribute.ID == productID {
								product.ExtID = orderAttribute.ID
								product.Name = append(product.Name, models3.LanguageDescription{
									Value: orderAttribute.Name,
								})
							}
						}
					}
				}

			default:
				return coreOrderModels.Order{}, fmt.Errorf("unsupported error code to update stoplist bu pos product/attribute id : %s", errSolution.Code)
			}
			if err != nil {
				return coreOrderModels.Order{}, err
			}
			log.Info().Msgf("send stoplist status true, for product/attribute id: %s, store_id:%s, error solution code:%s", productID, st.ID, errSolution.Code)

			if errSolution.SendToTelegram {
				if err := man.telegramService.SendMessageToQueue(telegram.PutProductToStopListWithErrSolution, order, st, order.FailReason.BusinessName, "", "", product); err != nil {
					return coreOrderModels.Order{}, err
				}
			}
		}
	}

	return order, nil
}
