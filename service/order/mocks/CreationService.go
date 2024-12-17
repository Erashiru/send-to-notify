// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/kwaaka-team/orders-core/core/models"
	mock "github.com/stretchr/testify/mock"
)

// CreationService is an autogenerated mock type for the CreationService type
type CreationService struct {
	mock.Mock
}

// CreateOrder provides a mock function with given fields: ctx, externalStoreID, deliveryService, aggReq, storeSecret
func (_m *CreationService) CreateOrder(ctx context.Context, externalStoreID string, deliveryService string, aggReq interface{}, storeSecret string) (models.Order, error) {
	ret := _m.Called(ctx, externalStoreID, deliveryService, aggReq, storeSecret)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrder")
	}

	var r0 models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, interface{}, string) (models.Order, error)); ok {
		return rf(ctx, externalStoreID, deliveryService, aggReq, storeSecret)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, interface{}, string) models.Order); ok {
		r0 = rf(ctx, externalStoreID, deliveryService, aggReq, storeSecret)
	} else {
		r0 = ret.Get(0).(models.Order)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, interface{}, string) error); ok {
		r1 = rf(ctx, externalStoreID, deliveryService, aggReq, storeSecret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCreationService creates a new instance of CreationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCreationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *CreationService {
	mock := &CreationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
