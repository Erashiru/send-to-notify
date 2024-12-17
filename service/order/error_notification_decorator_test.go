package order

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	mocks2 "github.com/kwaaka-team/orders-core/service/order/mocks"
	"github.com/kwaaka-team/orders-core/service/store/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestErrorNotificationDecorator_CreateOrder_NotError(t *testing.T) {
	storeService := &mocks.Service{}
	telegramService := &mocks2.TelegramService{}
	whatsappService := &mocks2.Whatsapp{}
	service := &mocks2.CreationService{}
	s, err := NewErrorNotificationDecorator(service, storeService, telegramService, whatsappService)
	if err != nil {
		t.Fatal(err)
	}
	expected := models.Order{
		ID: "id",
	}
	service.On("CreateOrder", context.Background(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expected, nil)
	actual, err := s.CreateOrder(context.Background(), "", "", nil, "")
	if err != nil {
		t.Fatal(err)
	}

	storeService.AssertNotCalled(t, "GetByExternalIdAndDeliveryService", mock.Anything, mock.Anything, mock.Anything)
	telegramService.AssertNotCalled(t, "SendMessageToQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatal("expected and actual are not equal")
	}

}

func TestErrorNotificationDecorator_CreateOrder_ErrorNotNotification(t *testing.T) {
	storeService := &mocks.Service{}
	telegramService := &mocks2.TelegramService{}
	whatsappService := &mocks2.Whatsapp{}
	service := &mocks2.CreationService{}
	s, err := NewErrorNotificationDecorator(service, storeService, telegramService, whatsappService)
	if err != nil {
		t.Fatal(err)
	}
	expected := models.Order{
		ID: "id",
	}
	service.On("CreateOrder", context.Background(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expected, errors.New("custom error"))
	actual, err := s.CreateOrder(context.Background(), "", "", nil, "")
	if err == nil {
		t.Fatal(err)
	}
	if errors.Is(err, errWithNotification) {
		t.Fatal("err is not errWithNotification")
	}

	storeService.AssertNotCalled(t, "GetByExternalIdAndDeliveryService", mock.Anything, mock.Anything, mock.Anything)
	telegramService.AssertNotCalled(t, "SendMessageToQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatal("expected and actual are not equal")
	}

}

func TestErrorNotificationDecorator_CreateOrder_ErrorNotification_StoreError(t *testing.T) {
	storeService := &mocks.Service{}
	telegramService := &mocks2.TelegramService{}
	service := &mocks2.CreationService{}
	whatsappService := &mocks2.Whatsapp{}
	s, err := NewErrorNotificationDecorator(service, storeService, telegramService, whatsappService)
	if err != nil {
		t.Fatal(err)
	}
	expected := models.Order{
		ID: "id",
	}

	storeNotFoundErr := errors.New("store not found error")
	service.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expected, storeNotFoundErr)

	storeService.On("GetByExternalIdAndDeliveryService", mock.Anything, mock.Anything, mock.Anything).Return(storeModels.Store{}, errors.New("store not found error"))

	actual, err := s.CreateOrder(context.Background(), "", "", nil, "")
	if err == nil {
		t.Fatal(err)
	}
	if errors.Is(err, errWithNotification) {
		t.Fatal("err is not errWithNotification")
	}

	storeService.AssertCalled(t, "GetByExternalIdAndDeliveryService", mock.Anything, mock.Anything, mock.Anything)

	telegramService.AssertNotCalled(t, "SendMessageToQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatal("expected and actual are not equal")
	}

	if !errors.Is(err, storeNotFoundErr) {
		t.Fatal("store error")
	}

}

func TestErrorNotificationDecorator_CreateOrder_ErrorNotification(t *testing.T) {
	storeService := &mocks.Service{}
	telegramService := &mocks2.TelegramService{}
	whatsappService := &mocks2.Whatsapp{}
	service := &mocks2.CreationService{}
	s, err := NewErrorNotificationDecorator(service, storeService, telegramService, whatsappService)
	if err != nil {
		t.Fatal(err)
	}
	expected := models.Order{
		ID: "id",
	}

	expectedErr := errors.New("expected error")
	service.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expected, expectedErr)

	expectedStore := storeModels.Store{}
	storeService.On("GetByExternalIdAndDeliveryService", mock.Anything, mock.Anything, mock.Anything).Return(expectedStore, nil)

	telegramService.On("SendMessageToQueue", telegram.CreateOrder, expected, expectedStore, expectedErr.Error()).Return(nil)

	actual, err := s.CreateOrder(context.Background(), "", "", nil, "")
	if err == nil {
		t.Fatal(err)
	}
	if errors.Is(err, errWithNotification) {
		t.Fatal("err is not errWithNotification")
	}

	storeService.AssertCalled(t, "GetByExternalIdAndDeliveryService", mock.Anything, mock.Anything, mock.Anything)
	telegramService.AssertCalled(t, "SendMessageToQueue", telegram.CreateOrder, expected, expectedStore, expectedErr.Error())

	if !reflect.DeepEqual(expected, actual) {
		t.Fatal("expected and actual are not equal")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatal("store error")
	}

}
