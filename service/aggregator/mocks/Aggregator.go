// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	menumodels "github.com/kwaaka-team/orders-core/core/menu/models"
	coremodels "github.com/kwaaka-team/orders-core/core/models"

	mock "github.com/stretchr/testify/mock"

	models "github.com/kwaaka-team/orders-core/core/storecore/models"
)

// Aggregator is an autogenerated mock type for the Aggregator type
type Aggregator struct {
	mock.Mock
}

// GetSystemCreateOrderRequestByAggregatorRequest provides a mock function with given fields: req, store
func (_m *Aggregator) GetSystemCreateOrderRequestByAggregatorRequest(req interface{}, store models.Store) (coremodels.Order, error) {
	ret := _m.Called(req, store)

	if len(ret) == 0 {
		panic("no return value specified for GetSystemCreateOrderRequestByAggregatorRequest")
	}

	var r0 coremodels.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, models.Store) (coremodels.Order, error)); ok {
		return rf(req, store)
	}
	if rf, ok := ret.Get(0).(func(interface{}, models.Store) coremodels.Order); ok {
		r0 = rf(req, store)
	} else {
		r0 = ret.Get(0).(coremodels.Order)
	}

	if rf, ok := ret.Get(1).(func(interface{}, models.Store) error); ok {
		r1 = rf(req, store)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MapSystemStatusToAggregatorStatus provides a mock function with given fields: order, posStatus, store
func (_m *Aggregator) MapSystemStatusToAggregatorStatus(order coremodels.Order, posStatus coremodels.PosStatus, store models.Store) string {
	ret := _m.Called(order, posStatus, store)

	if len(ret) == 0 {
		panic("no return value specified for MapSystemStatusToAggregatorStatus")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(coremodels.Order, coremodels.PosStatus, models.Store) string); ok {
		r0 = rf(order, posStatus, store)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// UpdateOrderInAggregator provides a mock function with given fields: ctx, order, store, aggregatorStatus
func (_m *Aggregator) UpdateOrderInAggregator(ctx context.Context, order coremodels.Order, store models.Store, aggregatorStatus string) error {
	ret := _m.Called(ctx, order, store, aggregatorStatus)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrderInAggregator")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, coremodels.Order, models.Store, string) error); ok {
		r0 = rf(ctx, order, store, aggregatorStatus)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateStopListByAttributesBulk provides a mock function with given fields: ctx, aggregatorStoreID, attributes
func (_m *Aggregator) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menumodels.Attribute) (string, error) {
	ret := _m.Called(ctx, aggregatorStoreID, attributes)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStopListByAttributesBulk")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []menumodels.Attribute) (string, error)); ok {
		return rf(ctx, aggregatorStoreID, attributes)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []menumodels.Attribute) string); ok {
		r0 = rf(ctx, aggregatorStoreID, attributes)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []menumodels.Attribute) error); ok {
		r1 = rf(ctx, aggregatorStoreID, attributes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStopListByProducts provides a mock function with given fields: ctx, aggregatorStoreID, products, isAvailable
func (_m *Aggregator) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menumodels.Product, isAvailable bool) (string, error) {
	ret := _m.Called(ctx, aggregatorStoreID, products, isAvailable)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStopListByProducts")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []menumodels.Product, bool) (string, error)); ok {
		return rf(ctx, aggregatorStoreID, products, isAvailable)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []menumodels.Product, bool) string); ok {
		r0 = rf(ctx, aggregatorStoreID, products, isAvailable)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []menumodels.Product, bool) error); ok {
		r1 = rf(ctx, aggregatorStoreID, products, isAvailable)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStopListByProductsBulk provides a mock function with given fields: ctx, aggregatorStoreID, products
func (_m *Aggregator) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menumodels.Product) (string, error) {
	ret := _m.Called(ctx, aggregatorStoreID, products)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStopListByProductsBulk")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []menumodels.Product) (string, error)); ok {
		return rf(ctx, aggregatorStoreID, products)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []menumodels.Product) string); ok {
		r0 = rf(ctx, aggregatorStoreID, products)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []menumodels.Product) error); ok {
		r1 = rf(ctx, aggregatorStoreID, products)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAggregator creates a new instance of Aggregator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAggregator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Aggregator {
	mock := &Aggregator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
