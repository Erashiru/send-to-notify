// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/menu/client.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	s3 "github.com/aws/aws-sdk-go/service/s3"
	gomock "github.com/golang/mock/gomock"
	models "github.com/kwaaka-team/orders-core/core/integration_api/models"
	models0 "github.com/kwaaka-team/orders-core/core/menu/models"
	models1 "github.com/kwaaka-team/orders-core/core/storecore/models"
	dto "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	dto0 "github.com/kwaaka-team/orders-core/pkg/store/dto"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// AddRowToAttributeGroup mocks base method.
func (m *MockClient) AddRowToAttributeGroup(ctx context.Context, menuId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRowToAttributeGroup", ctx, menuId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRowToAttributeGroup indicates an expected call of AddRowToAttributeGroup.
func (mr *MockClientMockRecorder) AddRowToAttributeGroup(ctx, menuId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRowToAttributeGroup", reflect.TypeOf((*MockClient)(nil).AddRowToAttributeGroup), ctx, menuId)
}

// AttributesStopList mocks base method.
func (m *MockClient) AttributesStopList(ctx context.Context, storeId string, attributes []dto.StopListItem, author string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttributesStopList", ctx, storeId, attributes, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// AttributesStopList indicates an expected call of AttributesStopList.
func (mr *MockClientMockRecorder) AttributesStopList(ctx, storeId, attributes, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttributesStopList", reflect.TypeOf((*MockClient)(nil).AttributesStopList), ctx, storeId, attributes, author)
}

// AutoUpdateMenuDescriptions mocks base method.
func (m *MockClient) AutoUpdateMenuDescriptions(ctx context.Context, storeId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AutoUpdateMenuDescriptions", ctx, storeId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AutoUpdateMenuDescriptions indicates an expected call of AutoUpdateMenuDescriptions.
func (mr *MockClientMockRecorder) AutoUpdateMenuDescriptions(ctx, storeId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AutoUpdateMenuDescriptions", reflect.TypeOf((*MockClient)(nil).AutoUpdateMenuDescriptions), ctx, storeId)
}

// AutoUpdateMenuPrices mocks base method.
func (m *MockClient) AutoUpdateMenuPrices(ctx context.Context, storeId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AutoUpdateMenuPrices", ctx, storeId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AutoUpdateMenuPrices indicates an expected call of AutoUpdateMenuPrices.
func (mr *MockClientMockRecorder) AutoUpdateMenuPrices(ctx, storeId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AutoUpdateMenuPrices", reflect.TypeOf((*MockClient)(nil).AutoUpdateMenuPrices), ctx, storeId)
}

// AutoUploadMenuByPOS mocks base method.
func (m *MockClient) AutoUploadMenuByPOS(ctx context.Context, req models0.Menu) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AutoUploadMenuByPOS", ctx, req)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AutoUploadMenuByPOS indicates an expected call of AutoUploadMenuByPOS.
func (mr *MockClientMockRecorder) AutoUploadMenuByPOS(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AutoUploadMenuByPOS", reflect.TypeOf((*MockClient)(nil).AutoUploadMenuByPOS), ctx, req)
}

// CreateAttributeGroup mocks base method.
func (m *MockClient) CreateAttributeGroup(ctx context.Context, menuID, attrGroupName string, min, max int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAttributeGroup", ctx, menuID, attrGroupName, min, max)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAttributeGroup indicates an expected call of CreateAttributeGroup.
func (mr *MockClientMockRecorder) CreateAttributeGroup(ctx, menuID, attrGroupName, min, max interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAttributeGroup", reflect.TypeOf((*MockClient)(nil).CreateAttributeGroup), ctx, menuID, attrGroupName, min, max)
}

// CreateGlovoSuperCollection mocks base method.
func (m *MockClient) CreateGlovoSuperCollection(ctx context.Context, menuId string, superCollections dto.MenuSuperCollections) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGlovoSuperCollection", ctx, menuId, superCollections)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateGlovoSuperCollection indicates an expected call of CreateGlovoSuperCollection.
func (mr *MockClientMockRecorder) CreateGlovoSuperCollection(ctx, menuId, superCollections interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGlovoSuperCollection", reflect.TypeOf((*MockClient)(nil).CreateGlovoSuperCollection), ctx, menuId, superCollections)
}

// CreateMenuByAggregatorApi mocks base method.
func (m *MockClient) CreateMenuByAggregatorApi(ctx context.Context, storeId, aggregator string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMenuByAggregatorApi", ctx, storeId, aggregator)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMenuByAggregatorApi indicates an expected call of CreateMenuByAggregatorApi.
func (mr *MockClientMockRecorder) CreateMenuByAggregatorApi(ctx, storeId, aggregator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMenuByAggregatorApi", reflect.TypeOf((*MockClient)(nil).CreateMenuByAggregatorApi), ctx, storeId, aggregator)
}

// CreateMenuUploadTransaction mocks base method.
func (m *MockClient) CreateMenuUploadTransaction(ctx context.Context, req dto.MenuUploadTransaction) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMenuUploadTransaction", ctx, req)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMenuUploadTransaction indicates an expected call of CreateMenuUploadTransaction.
func (mr *MockClientMockRecorder) CreateMenuUploadTransaction(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMenuUploadTransaction", reflect.TypeOf((*MockClient)(nil).CreateMenuUploadTransaction), ctx, req)
}

// DeleteAttributeGroupFromDB mocks base method.
func (m *MockClient) DeleteAttributeGroupFromDB(ctx context.Context, req dto.DeleteAttributeGroup) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAttributeGroupFromDB", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAttributeGroupFromDB indicates an expected call of DeleteAttributeGroupFromDB.
func (mr *MockClientMockRecorder) DeleteAttributeGroupFromDB(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAttributeGroupFromDB", reflect.TypeOf((*MockClient)(nil).DeleteAttributeGroupFromDB), ctx, req)
}

// DeleteProducts mocks base method.
func (m *MockClient) DeleteProducts(ctx context.Context, req dto.DeleteProducts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProducts", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProducts indicates an expected call of DeleteProducts.
func (mr *MockClientMockRecorder) DeleteProducts(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProducts", reflect.TypeOf((*MockClient)(nil).DeleteProducts), ctx, req)
}

// DeleteProductsFromDB mocks base method.
func (m *MockClient) DeleteProductsFromDB(ctx context.Context, req dto.DeleteProducts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProductsFromDB", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProductsFromDB indicates an expected call of DeleteProductsFromDB.
func (mr *MockClientMockRecorder) DeleteProductsFromDB(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProductsFromDB", reflect.TypeOf((*MockClient)(nil).DeleteProductsFromDB), ctx, req)
}

// GetAttributeForUpdate mocks base method.
func (m *MockClient) GetAttributeForUpdate(ctx context.Context, query dto.AttributeSelector) (dto.AttributesUpdate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttributeForUpdate", ctx, query)
	ret0, _ := ret[0].(dto.AttributesUpdate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttributeForUpdate indicates an expected call of GetAttributeForUpdate.
func (mr *MockClientMockRecorder) GetAttributeForUpdate(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttributeForUpdate", reflect.TypeOf((*MockClient)(nil).GetAttributeForUpdate), ctx, query)
}

// GetEmptyProducts mocks base method.
func (m *MockClient) GetEmptyProducts(ctx context.Context, menuID string, page, limit int64) ([]models0.Product, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmptyProducts", ctx, menuID, page, limit)
	ret0, _ := ret[0].([]models0.Product)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetEmptyProducts indicates an expected call of GetEmptyProducts.
func (mr *MockClientMockRecorder) GetEmptyProducts(ctx, menuID, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmptyProducts", reflect.TypeOf((*MockClient)(nil).GetEmptyProducts), ctx, menuID, page, limit)
}

// GetMenu mocks base method.
func (m *MockClient) GetMenu(ctx context.Context, externalStoreID string, deliveryService dto.DeliveryService) (models0.Menu, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMenu", ctx, externalStoreID, deliveryService)
	ret0, _ := ret[0].(models0.Menu)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMenu indicates an expected call of GetMenu.
func (mr *MockClientMockRecorder) GetMenu(ctx, externalStoreID, deliveryService interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMenu", reflect.TypeOf((*MockClient)(nil).GetMenu), ctx, externalStoreID, deliveryService)
}

// GetMenuByID mocks base method.
func (m *MockClient) GetMenuByID(ctx context.Context, menuID string) (models0.Menu, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMenuByID", ctx, menuID)
	ret0, _ := ret[0].(models0.Menu)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMenuByID indicates an expected call of GetMenuByID.
func (mr *MockClientMockRecorder) GetMenuByID(ctx, menuID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMenuByID", reflect.TypeOf((*MockClient)(nil).GetMenuByID), ctx, menuID)
}

// GetMenuGroups mocks base method.
func (m *MockClient) GetMenuGroups(ctx context.Context, menuID string) ([]dto.MenuGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMenuGroups", ctx, menuID)
	ret0, _ := ret[0].([]dto.MenuGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMenuGroups indicates an expected call of GetMenuGroups.
func (mr *MockClientMockRecorder) GetMenuGroups(ctx, menuID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMenuGroups", reflect.TypeOf((*MockClient)(nil).GetMenuGroups), ctx, menuID)
}

// GetMenuStatus mocks base method.
func (m *MockClient) GetMenuStatus(ctx context.Context, storeId string, isDeleted bool) ([]dto0.StoreDsMenuDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMenuStatus", ctx, storeId, isDeleted)
	ret0, _ := ret[0].([]dto0.StoreDsMenuDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMenuStatus indicates an expected call of GetMenuStatus.
func (mr *MockClientMockRecorder) GetMenuStatus(ctx, storeId, isDeleted interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMenuStatus", reflect.TypeOf((*MockClient)(nil).GetMenuStatus), ctx, storeId, isDeleted)
}

// GetMenuUploadTransaction mocks base method.
func (m *MockClient) GetMenuUploadTransaction(ctx context.Context, req dto.MenuUploadTransaction) (dto.MenuUploadTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMenuUploadTransaction", ctx, req)
	ret0, _ := ret[0].(dto.MenuUploadTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMenuUploadTransaction indicates an expected call of GetMenuUploadTransaction.
func (mr *MockClientMockRecorder) GetMenuUploadTransaction(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMenuUploadTransaction", reflect.TypeOf((*MockClient)(nil).GetMenuUploadTransaction), ctx, req)
}

// GetMenuUploadTransactions mocks base method.
func (m *MockClient) GetMenuUploadTransactions(ctx context.Context, req dto.GetMenuUploadTransactions) ([]models0.MenuUploadTransaction, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMenuUploadTransactions", ctx, req)
	ret0, _ := ret[0].([]models0.MenuUploadTransaction)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMenuUploadTransactions indicates an expected call of GetMenuUploadTransactions.
func (mr *MockClientMockRecorder) GetMenuUploadTransactions(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMenuUploadTransactions", reflect.TypeOf((*MockClient)(nil).GetMenuUploadTransactions), ctx, req)
}

// GetPosDiscounts mocks base method.
func (m *MockClient) GetPosDiscounts(ctx context.Context, storeID string, deliveryService dto.DeliveryService) (dto.PosDiscount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPosDiscounts", ctx, storeID, deliveryService)
	ret0, _ := ret[0].(dto.PosDiscount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPosDiscounts indicates an expected call of GetPosDiscounts.
func (mr *MockClientMockRecorder) GetPosDiscounts(ctx, storeID, deliveryService interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPosDiscounts", reflect.TypeOf((*MockClient)(nil).GetPosDiscounts), ctx, storeID, deliveryService)
}

// GetProcessingMenuUploadTransactions mocks base method.
func (m *MockClient) GetProcessingMenuUploadTransactions(ctx context.Context, req dto.GetMenuUploadTransactions) ([]dto.MenuUploadTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProcessingMenuUploadTransactions", ctx, req)
	ret0, _ := ret[0].([]dto.MenuUploadTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProcessingMenuUploadTransactions indicates an expected call of GetProcessingMenuUploadTransactions.
func (mr *MockClientMockRecorder) GetProcessingMenuUploadTransactions(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProcessingMenuUploadTransactions", reflect.TypeOf((*MockClient)(nil).GetProcessingMenuUploadTransactions), ctx, req)
}

// GetPromos mocks base method.
func (m *MockClient) GetPromos(ctx context.Context, externalStoreID string, deliveryService dto.DeliveryService) (models0.Promo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPromos", ctx, externalStoreID, deliveryService)
	ret0, _ := ret[0].(models0.Promo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPromos indicates an expected call of GetPromos.
func (mr *MockClientMockRecorder) GetPromos(ctx, externalStoreID, deliveryService interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPromos", reflect.TypeOf((*MockClient)(nil).GetPromos), ctx, externalStoreID, deliveryService)
}

// GetStorePromos mocks base method.
func (m *MockClient) GetStorePromos(ctx context.Context, query dto.GetPromosSelector) ([]dto.PromoDiscount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorePromos", ctx, query)
	ret0, _ := ret[0].([]dto.PromoDiscount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStorePromos indicates an expected call of GetStorePromos.
func (mr *MockClientMockRecorder) GetStorePromos(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorePromos", reflect.TypeOf((*MockClient)(nil).GetStorePromos), ctx, query)
}

// GetStores mocks base method.
func (m *MockClient) GetStores(ctx context.Context, deliveryService dto.DeliveryService) ([]models1.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStores", ctx, deliveryService)
	ret0, _ := ret[0].([]models1.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStores indicates an expected call of GetStores.
func (mr *MockClientMockRecorder) GetStores(ctx, deliveryService interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStores", reflect.TypeOf((*MockClient)(nil).GetStores), ctx, deliveryService)
}

// InsertMenu mocks base method.
func (m *MockClient) InsertMenu(ctx context.Context, menu models0.Menu) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertMenu", ctx, menu)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertMenu indicates an expected call of InsertMenu.
func (mr *MockClientMockRecorder) InsertMenu(ctx, menu interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMenu", reflect.TypeOf((*MockClient)(nil).InsertMenu), ctx, menu)
}

// ListStoresByProduct mocks base method.
func (m *MockClient) ListStoresByProduct(ctx context.Context, req dto.GetStoreByProductRequest) ([]models1.Store, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListStoresByProduct", ctx, req)
	ret0, _ := ret[0].([]models1.Store)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListStoresByProduct indicates an expected call of ListStoresByProduct.
func (mr *MockClientMockRecorder) ListStoresByProduct(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStoresByProduct", reflect.TypeOf((*MockClient)(nil).ListStoresByProduct), ctx, req)
}

// MergeMenus mocks base method.
func (m *MockClient) MergeMenus(ctx context.Context, restaurantID string, restaurantIDs []string, author string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MergeMenus", ctx, restaurantID, restaurantIDs, author)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MergeMenus indicates an expected call of MergeMenus.
func (mr *MockClientMockRecorder) MergeMenus(ctx, restaurantID, restaurantIDs, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MergeMenus", reflect.TypeOf((*MockClient)(nil).MergeMenus), ctx, restaurantID, restaurantIDs, author)
}

// PosIntegrationUpdateStopList mocks base method.
func (m *MockClient) PosIntegrationUpdateStopList(ctx context.Context, storeId string, request models.StopListRequest, author string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PosIntegrationUpdateStopList", ctx, storeId, request, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// PosIntegrationUpdateStopList indicates an expected call of PosIntegrationUpdateStopList.
func (mr *MockClientMockRecorder) PosIntegrationUpdateStopList(ctx, storeId, request, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PosIntegrationUpdateStopList", reflect.TypeOf((*MockClient)(nil).PosIntegrationUpdateStopList), ctx, storeId, request, author)
}

// RecoveryMenu mocks base method.
func (m *MockClient) RecoveryMenu(ctx context.Context, menuId, author string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecoveryMenu", ctx, menuId, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecoveryMenu indicates an expected call of RecoveryMenu.
func (mr *MockClientMockRecorder) RecoveryMenu(ctx, menuId, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecoveryMenu", reflect.TypeOf((*MockClient)(nil).RecoveryMenu), ctx, menuId, author)
}

// RenewPositionsInVirtualStore mocks base method.
func (m *MockClient) RenewPositionsInVirtualStore(ctx context.Context, restaurantID, originalRestaurantID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenewPositionsInVirtualStore", ctx, restaurantID, originalRestaurantID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RenewPositionsInVirtualStore indicates an expected call of RenewPositionsInVirtualStore.
func (mr *MockClientMockRecorder) RenewPositionsInVirtualStore(ctx, restaurantID, originalRestaurantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenewPositionsInVirtualStore", reflect.TypeOf((*MockClient)(nil).RenewPositionsInVirtualStore), ctx, restaurantID, originalRestaurantID)
}

// StopPositionsInVirtualStore mocks base method.
func (m *MockClient) StopPositionsInVirtualStore(ctx context.Context, restaurantID, originalRestaurantID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopPositionsInVirtualStore", ctx, restaurantID, originalRestaurantID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopPositionsInVirtualStore indicates an expected call of StopPositionsInVirtualStore.
func (mr *MockClientMockRecorder) StopPositionsInVirtualStore(ctx, restaurantID, originalRestaurantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopPositionsInVirtualStore", reflect.TypeOf((*MockClient)(nil).StopPositionsInVirtualStore), ctx, restaurantID, originalRestaurantID)
}

// UpdateMatchingProduct mocks base method.
func (m *MockClient) UpdateMatchingProduct(ctx context.Context, req models0.MatchingProducts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMatchingProduct", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMatchingProduct indicates an expected call of UpdateMatchingProduct.
func (mr *MockClientMockRecorder) UpdateMatchingProduct(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMatchingProduct", reflect.TypeOf((*MockClient)(nil).UpdateMatchingProduct), ctx, req)
}

// UpdateMenu mocks base method.
func (m *MockClient) UpdateMenu(ctx context.Context, req models0.Menu) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMenu", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMenu indicates an expected call of UpdateMenu.
func (mr *MockClientMockRecorder) UpdateMenu(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMenu", reflect.TypeOf((*MockClient)(nil).UpdateMenu), ctx, req)
}

// UpdateMenuName mocks base method.
func (m *MockClient) UpdateMenuName(ctx context.Context, query dto.UpdateMenuName) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMenuName", ctx, query)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMenuName indicates an expected call of UpdateMenuName.
func (mr *MockClientMockRecorder) UpdateMenuName(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMenuName", reflect.TypeOf((*MockClient)(nil).UpdateMenuName), ctx, query)
}

// UpdateMenuUploadTransaction mocks base method.
func (m *MockClient) UpdateMenuUploadTransaction(ctx context.Context, req dto.MenuUploadTransaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMenuUploadTransaction", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMenuUploadTransaction indicates an expected call of UpdateMenuUploadTransaction.
func (mr *MockClientMockRecorder) UpdateMenuUploadTransaction(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMenuUploadTransaction", reflect.TypeOf((*MockClient)(nil).UpdateMenuUploadTransaction), ctx, req)
}

// UpdateProductAvailableStatus mocks base method.
func (m *MockClient) UpdateProductAvailableStatus(ctx context.Context, menuID, productID string, status bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProductAvailableStatus", ctx, menuID, productID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProductAvailableStatus indicates an expected call of UpdateProductAvailableStatus.
func (mr *MockClientMockRecorder) UpdateProductAvailableStatus(ctx, menuID, productID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProductAvailableStatus", reflect.TypeOf((*MockClient)(nil).UpdateProductAvailableStatus), ctx, menuID, productID, status)
}

// UpdateProductByFields mocks base method.
func (m *MockClient) UpdateProductByFields(ctx context.Context, menuId, productID string, req models0.ProductUpdateRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProductByFields", ctx, menuId, productID, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProductByFields indicates an expected call of UpdateProductByFields.
func (mr *MockClientMockRecorder) UpdateProductByFields(ctx, menuId, productID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProductByFields", reflect.TypeOf((*MockClient)(nil).UpdateProductByFields), ctx, menuId, productID, req)
}

// UpdateStopListStores mocks base method.
func (m *MockClient) UpdateStopListStores(ctx context.Context, req dto.ProductStopList, author string) (dto.StoreProductStopLists, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStopListStores", ctx, req, author)
	ret0, _ := ret[0].(dto.StoreProductStopLists)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStopListStores indicates an expected call of UpdateStopListStores.
func (mr *MockClientMockRecorder) UpdateStopListStores(ctx, req, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStopListStores", reflect.TypeOf((*MockClient)(nil).UpdateStopListStores), ctx, req, author)
}

// UploadMenu mocks base method.
func (m *MockClient) UploadMenu(ctx context.Context, req dto.MenuUploadRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadMenu", ctx, req)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadMenu indicates an expected call of UploadMenu.
func (mr *MockClientMockRecorder) UploadMenu(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadMenu", reflect.TypeOf((*MockClient)(nil).UploadMenu), ctx, req)
}

// UpsertDeliveryMenu mocks base method.
func (m *MockClient) UpsertDeliveryMenu(ctx context.Context, req dto.MenuGroupRequest, author string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertDeliveryMenu", ctx, req, author)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertDeliveryMenu indicates an expected call of UpsertDeliveryMenu.
func (mr *MockClientMockRecorder) UpsertDeliveryMenu(ctx, req, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertDeliveryMenu", reflect.TypeOf((*MockClient)(nil).UpsertDeliveryMenu), ctx, req, author)
}

// UpsertMenu mocks base method.
func (m *MockClient) UpsertMenu(ctx context.Context, req dto.MenuGroupRequest, author string, upsertToAggrMenu bool) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertMenu", ctx, req, author, upsertToAggrMenu)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertMenu indicates an expected call of UpsertMenu.
func (mr *MockClientMockRecorder) UpsertMenu(ctx, req, author, upsertToAggrMenu interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertMenu", reflect.TypeOf((*MockClient)(nil).UpsertMenu), ctx, req, author, upsertToAggrMenu)
}

// UpsertMenuByFields mocks base method.
func (m *MockClient) UpsertMenuByFields(ctx context.Context, fields models0.UpdateFields, agg models0.UpdateFieldsAggregators, req dto.MenuGroupRequest, author string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertMenuByFields", ctx, fields, agg, req, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertMenuByFields indicates an expected call of UpsertMenuByFields.
func (mr *MockClientMockRecorder) UpsertMenuByFields(ctx, fields, agg, req, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertMenuByFields", reflect.TypeOf((*MockClient)(nil).UpsertMenuByFields), ctx, fields, agg, req, author)
}

// ValidateAggProductErr mocks base method.
func (m *MockClient) ValidateAggProductErr(ctx context.Context, menuID, storeID string, limit int) ([]models0.Product, []models0.Product, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAggProductErr", ctx, menuID, storeID, limit)
	ret0, _ := ret[0].([]models0.Product)
	ret1, _ := ret[1].([]models0.Product)
	ret2, _ := ret[2].(int)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ValidateAggProductErr indicates an expected call of ValidateAggProductErr.
func (mr *MockClientMockRecorder) ValidateAggProductErr(ctx, menuID, storeID, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAggProductErr", reflect.TypeOf((*MockClient)(nil).ValidateAggProductErr), ctx, menuID, storeID, limit)
}

// ValidateAttributeGroupName mocks base method.
func (m *MockClient) ValidateAttributeGroupName(ctx context.Context, menuId, name string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAttributeGroupName", ctx, menuId, name)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateAttributeGroupName indicates an expected call of ValidateAttributeGroupName.
func (mr *MockClientMockRecorder) ValidateAttributeGroupName(ctx, menuId, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAttributeGroupName", reflect.TypeOf((*MockClient)(nil).ValidateAttributeGroupName), ctx, menuId, name)
}

// ValidateMatching mocks base method.
func (m *MockClient) ValidateMatching(ctx context.Context, req dto.MenuValidateRequest, sv3 *s3.S3) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateMatching", ctx, req, sv3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateMatching indicates an expected call of ValidateMatching.
func (mr *MockClientMockRecorder) ValidateMatching(ctx, req, sv3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateMatching", reflect.TypeOf((*MockClient)(nil).ValidateMatching), ctx, req, sv3)
}

// ValidateVirtualStoreMatching mocks base method.
func (m *MockClient) ValidateVirtualStoreMatching(ctx context.Context, req dto.MenuValidateRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateVirtualStoreMatching", ctx, req)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateVirtualStoreMatching indicates an expected call of ValidateVirtualStoreMatching.
func (mr *MockClientMockRecorder) ValidateVirtualStoreMatching(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateVirtualStoreMatching", reflect.TypeOf((*MockClient)(nil).ValidateVirtualStoreMatching), ctx, req)
}

// VerifyUploadMenu mocks base method.
func (m *MockClient) VerifyUploadMenu(ctx context.Context, req dto.MenuUploadVerifyRequest) (dto.MenuUploadTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyUploadMenu", ctx, req)
	ret0, _ := ret[0].(dto.MenuUploadTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyUploadMenu indicates an expected call of VerifyUploadMenu.
func (mr *MockClientMockRecorder) VerifyUploadMenu(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyUploadMenu", reflect.TypeOf((*MockClient)(nil).VerifyUploadMenu), ctx, req)
}
