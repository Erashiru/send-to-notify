package stoplist

import (
	"context"
	"errors"
	"github.com/kwaaka-team/orders-core/core/config"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	aggregatorMock "github.com/kwaaka-team/orders-core/service/aggregator/mocks"
	"github.com/kwaaka-team/orders-core/service/menu"
	menuMock "github.com/kwaaka-team/orders-core/service/menu/mocks"
	"github.com/kwaaka-team/orders-core/service/pos"
	posMock "github.com/kwaaka-team/orders-core/service/pos/mocks"
	stopListMock "github.com/kwaaka-team/orders-core/service/stoplist/mocks"
	storeMock "github.com/kwaaka-team/orders-core/service/store/mocks"
	storeGroupMock "github.com/kwaaka-team/orders-core/service/storegroup/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestServiceImpl_UpdateStopListByPosProductID(t *testing.T) {
	var storeService = &storeMock.Service{}
	var storeGroupService = &storeGroupMock.Service{}

	var menuRepository = &menuMock.Repository{}
	menuService, err := menu.NewMenuService(menuRepository, nil, nil, nil)
	if err != nil {
		t.Error(err)
		return
	}

	var aggFactory = &aggregatorMock.Factory{}
	var posFactory pos.Factory = &posMock.Factory{}
	var repo = &stopListMock.Repository{}
	var woltCfg = config.WoltConfiguration{}

	stopListService, err := NewStopListServiceCron(storeService, storeGroupService, menuService, aggFactory, posFactory, repo, woltCfg, 1)
	if err != nil {
		t.Error(err)
		return
	}

	storeDSMenus := make([]storeModels.StoreDSMenu, 0)

	aggMenuID1 := "aggMenuID1"
	storeDSMenu1 := storeModels.StoreDSMenu{
		Delivery: "delivery_1",
		ID:       aggMenuID1,
		IsActive: true,
	}
	storeDSMenus = append(storeDSMenus, storeDSMenu1)

	aggMenuID2 := "aggMenuID2"
	storeDSMenu2 := storeModels.StoreDSMenu{
		Delivery: "delivery_2",
		ID:       aggMenuID2,
		IsActive: true,
	}
	storeDSMenus = append(storeDSMenus, storeDSMenu2)

	storeID := "storeID"
	posMenuID := "posMenuID"
	store := storeModels.Store{
		ID:      storeID,
		PosType: "iiko",
		MenuID:  posMenuID,
		Menus:   storeDSMenus,
	}
	expectedPosMenu := menuModels.Menu{
		ID:       posMenuID,
		IsActive: true,
		IsSync:   true,
		Products: generatePosMenuProducts(),
	}
	storeService.On("GetByID", mock.Anything, storeID).
		Return(store, nil)
	storeService.On("GetByID", mock.Anything, mock.Anything).
		Return(storeModels.Store{}, errors.New("store not found"))

	menuRepository.On("FindById", mock.Anything, posMenuID).
		Return(&expectedPosMenu, nil)

	if err = stopListService.UpdateStopListByPosProductID(context.Background(), true, storeID, "pos_5_1_deleted"); err != nil {
		t.Error(err)
		return
	}
	menuRepository.AssertNotCalled(t, "BulkUpdateProductsAvailability", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	if err = stopListService.UpdateStopListByPosProductID(context.Background(), true, storeID, "pos_5_2_deleted"); err != nil {
		t.Error(err)
		return
	}
	menuRepository.AssertNotCalled(t, "BulkUpdateProductsAvailability", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	ctx := context.Background()

	aggMenu1 := menuModels.Menu{
		ID:       aggMenuID1,
		Products: generateAggMenuProducts1(),
	}
	menuRepository.On("FindById", mock.Anything, aggMenuID1).Return(&aggMenu1, nil)

	aggMenu2 := menuModels.Menu{
		ID: aggMenuID2,
	}
	menuRepository.On("FindById", mock.Anything, aggMenuID2).Return(&aggMenu2, nil)

	expectedErr := errors.New("BulkUpdateProductsAvailability error")
	menuRepository.On("BulkUpdateProductsAvailability", mock.Anything, mock.Anything, []string{"system_1"}, true).Return(expectedErr)
	if err = stopListService.UpdateStopListByPosProductID(ctx, true, storeID, "pos_1"); err == nil {
		t.Error("expected error, got nil")
		return
	} else {
		if !errors.Is(expectedErr, err) {
			t.Error(err)
			return
		}
	}

	menuRepository.On("BulkUpdateProductsAvailability", mock.Anything, mock.Anything, []string{"system_2_1", "system_2_2"}, true).Return(nil)

	storeService.On("GetStoreExternalIds", store, "delivery_1").Return([]string{"delivery_1_store_1", "delivery_1_store_2"}, nil)
	aggregator1 := aggregatorMock.Aggregator{}
	//aggFactory.On("GetAggregator", "delivery_1", store).Return(aggregator1, nil)
	aggFactory.On("GetAggregator", mock.Anything, mock.Anything).Return(&aggregator1, nil)
	aggregator1.On("UpdateStopListByProductsBulk", mock.Anything, mock.Anything, mock.Anything).Return("transactionID1", nil)

	repo.On("InsertStopListTransaction", mock.Anything, mock.Anything).Return(nil)

	menuRepository.On("BulkUpdateProductsAvailability", mock.Anything, mock.Anything, []string{"delivery_2_1", "delivery_2_2"}, true).Return(nil)

	menuRepository.On("BulkUpdateProductsDisabledStatus", mock.Anything, "posMenuID", []string{"system_2_1", "system_2_2"}, false).Return(nil)
	menuRepository.On("BulkUpdateProductsDisabledStatus", mock.Anything, "aggMenuID1", []string{"delivery_2_1", "delivery_2_2"}, false).Return(nil)

	if err = stopListService.UpdateStopListByPosProductID(ctx, true, storeID, "pos_2"); err != nil {
		t.Error(err)
		return
	}

	menuRepository.AssertCalled(t, "BulkUpdateProductsAvailability", mock.Anything, posMenuID, []string{"system_1"}, true)
	menuRepository.AssertCalled(t, "BulkUpdateProductsAvailability", mock.Anything, aggMenu1.ID, []string{"delivery_2_1", "delivery_2_2"}, true)
	aggregator1.AssertCalled(t, "UpdateStopListByProductsBulk", ctx, "delivery_1_store_1", mock.Anything)
}

func generatePosMenuProducts() []menuModels.Product {
	products := make([]menuModels.Product, 0)

	p1 := menuModels.Product{
		ProductID: "pos_1",
		ExtID:     "system_1",
	}
	products = append(products, p1)

	p2 := menuModels.Product{
		ProductID: "pos_2",
		ExtID:     "system_2_1",
	}
	products = append(products, p2)

	p3 := menuModels.Product{
		ProductID: "pos_2",
		ExtID:     "system_2_2",
	}
	products = append(products, p3)

	p4 := menuModels.Product{
		ProductID: "pos_4",
		ExtID:     "",
	}
	products = append(products, p4)

	p51Deleted := menuModels.Product{
		ProductID: "pos_5_1_deleted",
		ExtID:     "system_5_1_deleted",
		IsDeleted: true,
	}
	products = append(products, p51Deleted)

	p52Deleted := menuModels.Product{
		ProductID: "pos_5_2_deleted",
		ExtID:     "",
		IsDeleted: true,
	}
	products = append(products, p52Deleted)

	p6Extra := menuModels.Product{
		ProductID: "pos_6_extra",
		ExtID:     "system_6_extra",
	}
	products = append(products, p6Extra)

	p7Extra := menuModels.Product{
		ProductID: "pos_7_extra",
		ExtID:     "system_7_extra",
	}
	products = append(products, p7Extra)

	return products
}

func generateAggMenuProducts1() []menuModels.Product {
	products := make([]menuModels.Product, 0)

	p1 := menuModels.Product{
		ExtID: "delivery_2_1",
		PosID: "system_2_1",
	}
	products = append(products, p1)

	p2 := menuModels.Product{
		ExtID: "delivery_2_2",
		PosID: "system_2_2",
	}
	products = append(products, p2)

	p3 := menuModels.Product{
		ExtID: "delivery_2_3",
		PosID: "system_2_3",
	}
	products = append(products, p3)

	p41Deleted := menuModels.Product{
		ExtID:     "delivery_2_2_deleted_1",
		PosID:     "system_2_2",
		IsDeleted: true,
	}
	products = append(products, p41Deleted)

	p42Deleted := menuModels.Product{
		ExtID:     "delivery_2_2_deleted_2",
		PosID:     "system_2_2",
		IsDeleted: true,
	}
	products = append(products, p42Deleted)

	return products
}
