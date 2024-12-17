// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	models "github.com/kwaaka-team/orders-core/core/models"
	mock "github.com/stretchr/testify/mock"

	storecoremodels "github.com/kwaaka-team/orders-core/core/storecore/models"

	telegram "github.com/kwaaka-team/orders-core/core/managers/telegram"
)

// TelegramService is an autogenerated mock type for the TelegramService type
type TelegramService struct {
	mock.Mock
}

// SendMessageToQueue provides a mock function with given fields: notificationType, _a1, store, err, msg
func (_m *TelegramService) SendMessageToQueue(notificationType telegram.NotificationType, _a1 models.Order, store storecoremodels.Store, err, msg, extraMsg string, product models2.Product) error {
	ret := _m.Called(notificationType, _a1, store, err, msg)

	if len(ret) == 0 {
		panic("no return value specified for SendMessageToQueue")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(telegram.NotificationType, models.Order, storecoremodels.Store, string, string) error); ok {
		r0 = rf(notificationType, _a1, store, err, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendMessageToRestaurant provides a mock function with given fields: notificationType, _a1, store, err
func (_m *TelegramService) SendMessageToRestaurant(notificationType telegram.NotificationType, _a1 models.Order, store storecoremodels.Store, err string) error {
	ret := _m.Called(notificationType, _a1, store, err)

	if len(ret) == 0 {
		panic("no return value specified for SendMessageToRestaurant")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(telegram.NotificationType, models.Order, storecoremodels.Store, string) error); ok {
		r0 = rf(notificationType, _a1, store, err)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTelegramService creates a new instance of TelegramService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTelegramService(t interface {
	mock.TestingT
	Cleanup(func())
}) *TelegramService {
	mock := &TelegramService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
