// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	coremodels "github.com/kwaaka-team/orders-core/core/models"

	mock "github.com/stretchr/testify/mock"

	models "github.com/kwaaka-team/orders-core/service/kwaaka_3pl/models"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Cancel3plOrder provides a mock function with given fields: ctx, orderID
func (_m *Service) Cancel3plOrder(ctx context.Context, orderID string) error {
	ret := _m.Called(ctx, orderID)

	if len(ret) == 0 {
		panic("no return value specified for Cancel3plOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, orderID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create3plOrder provides a mock function with given fields: ctx, req
func (_m *Service) Create3plOrder(ctx context.Context, req models.CreateDeliveryRequest) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for Create3plOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateDeliveryRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDeliveryInfoByOrderId provides a mock function with given fields: ctx, orderId
func (_m *Service) GetDeliveryInfoByOrderId(ctx context.Context, orderId string) (coremodels.GetDeliveryInfoResp, error) {
	ret := _m.Called(ctx, orderId)

	if len(ret) == 0 {
		panic("no return value specified for GetDeliveryInfoByOrderId")
	}

	var r0 coremodels.GetDeliveryInfoResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (coremodels.GetDeliveryInfoResp, error)); ok {
		return rf(ctx, orderId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) coremodels.GetDeliveryInfoResp); ok {
		r0 = rf(ctx, orderId)
	} else {
		r0 = ret.Get(0).(coremodels.GetDeliveryInfoResp)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, orderId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetKwaaka3plDispatcher provides a mock function with given fields: ctx, req
func (_m *Service) SetKwaaka3plDispatcher(ctx context.Context, req coremodels.SetKwaaka3plDispatcherRequest) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for SetKwaaka3plDispatcher")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, coremodels.SetKwaaka3plDispatcherRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
