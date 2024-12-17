// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/clients/iiko.go

// Package mock_clients is a generated GoMock package.
package mock_clients

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

// MockIIKO is a mock of IIKO interface.
type MockIIKO struct {
	ctrl     *gomock.Controller
	recorder *MockIIKOMockRecorder
}

// MockIIKOMockRecorder is the mock recorder for MockIIKO.
type MockIIKOMockRecorder struct {
	mock *MockIIKO
}

// NewMockIIKO creates a new mock instance.
func NewMockIIKO(ctrl *gomock.Controller) *MockIIKO {
	mock := &MockIIKO{ctrl: ctrl}
	mock.recorder = &MockIIKOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIIKO) EXPECT() *MockIIKOMockRecorder {
	return m.recorder
}

// Auth mocks base method.
func (m *MockIIKO) Auth(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Auth indicates an expected call of Auth.
func (mr *MockIIKOMockRecorder) Auth(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockIIKO)(nil).Auth), ctx)
}

// Close mocks base method.
func (m *MockIIKO) Close(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close", ctx)
}

// Close indicates an expected call of Close.
func (mr *MockIIKOMockRecorder) Close(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIIKO)(nil).Close), ctx)
}

// CreateDeliveryOrder mocks base method.
func (m *MockIIKO) CreateDeliveryOrder(ctx context.Context, req models.CreateDeliveryRequest) (models.CreateDeliveryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDeliveryOrder", ctx, req)
	ret0, _ := ret[0].(models.CreateDeliveryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDeliveryOrder indicates an expected call of CreateDeliveryOrder.
func (mr *MockIIKOMockRecorder) CreateDeliveryOrder(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDeliveryOrder", reflect.TypeOf((*MockIIKO)(nil).CreateDeliveryOrder), ctx, req)
}

// GetMenu mocks base method.
func (m *MockIIKO) GetMenu(ctx context.Context, organizationID string) (models.GetMenuResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMenu", ctx, organizationID)
	ret0, _ := ret[0].(models.GetMenuResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMenu indicates an expected call of GetMenu.
func (mr *MockIIKOMockRecorder) GetMenu(ctx, organizationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMenu", reflect.TypeOf((*MockIIKO)(nil).GetMenu), ctx, organizationID)
}

// GetOrganizations mocks base method.
func (m *MockIIKO) GetOrganizations(ctx context.Context) ([]models.Info, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizations", ctx)
	ret0, _ := ret[0].([]models.Info)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizations indicates an expected call of GetOrganizations.
func (mr *MockIIKOMockRecorder) GetOrganizations(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizations", reflect.TypeOf((*MockIIKO)(nil).GetOrganizations), ctx)
}

// GetStopList mocks base method.
func (m *MockIIKO) GetStopList(ctx context.Context, req models.StopListRequest) (models.StopListResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStopList", ctx, req)
	ret0, _ := ret[0].(models.StopListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStopList indicates an expected call of GetStopList.
func (mr *MockIIKOMockRecorder) GetStopList(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStopList", reflect.TypeOf((*MockIIKO)(nil).GetStopList), ctx, req)
}

// GetWebhookSetting mocks base method.
func (m *MockIIKO) GetWebhookSetting(ctx context.Context, organizationID string) (models.GetWebhookSettingResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWebhookSetting", ctx, organizationID)
	ret0, _ := ret[0].(models.GetWebhookSettingResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWebhookSetting indicates an expected call of GetWebhookSetting.
func (mr *MockIIKOMockRecorder) GetWebhookSetting(ctx, organizationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWebhookSetting", reflect.TypeOf((*MockIIKO)(nil).GetWebhookSetting), ctx, organizationID)
}

// RetrieveDeliveryOrder mocks base method.
func (m *MockIIKO) RetrieveDeliveryOrder(ctx context.Context, organizationID, orderID string) (models.RetrieveOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveDeliveryOrder", ctx, organizationID, orderID)
	ret0, _ := ret[0].(models.RetrieveOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveDeliveryOrder indicates an expected call of RetrieveDeliveryOrder.
func (mr *MockIIKOMockRecorder) RetrieveDeliveryOrder(ctx, organizationID, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveDeliveryOrder", reflect.TypeOf((*MockIIKO)(nil).RetrieveDeliveryOrder), ctx, organizationID, orderID)
}

// UpdateWebhookSetting mocks base method.
func (m *MockIIKO) UpdateWebhookSetting(ctx context.Context, request models.UpdateWebhookRequest) (models.CorID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWebhookSetting", ctx, request)
	ret0, _ := ret[0].(models.CorID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateWebhookSetting indicates an expected call of UpdateWebhookSetting.
func (mr *MockIIKOMockRecorder) UpdateWebhookSetting(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWebhookSetting", reflect.TypeOf((*MockIIKO)(nil).UpdateWebhookSetting), ctx, request)
}
