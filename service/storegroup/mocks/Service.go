// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	legalentitymodels "github.com/kwaaka-team/orders-core/service/legalentity/models"
	mock "github.com/stretchr/testify/mock"

	models "github.com/kwaaka-team/orders-core/core/storecore/models"

	selector "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

func (_m *Service) AddBrandInfo(ctx context.Context, restInfo models.BrandInfo, restGroupID string) error {
	ret := _m.Called(ctx, restInfo, restGroupID)
    if len(ret) == 0 {
        panic("no return value specified for AddBrandInfo")
    }
    var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.BrandInfo, string) error); ok {
        r0 = rf(ctx, restInfo, restGroupID)
    } else {
        r0 = ret.Error(0)
    }
    return r0
}

// CreateStoreGroup provides a mock function with given fields: ctx, group
func (_m *Service) CreateStoreGroup(ctx context.Context, group models.StoreGroup) (string, error) {
	ret := _m.Called(ctx, group)

	if len(ret) == 0 {
		panic("no return value specified for CreateStoreGroup")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.StoreGroup) (string, error)); ok {
		return rf(ctx, group)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.StoreGroup) string); ok {
		r0 = rf(ctx, group)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.StoreGroup) error); ok {
		r1 = rf(ctx, group)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllStoreGroupsIdsAndNames provides a mock function with given fields: ctx
func (_m *Service) GetAllStoreGroupsIdsAndNames(ctx context.Context) ([]models.StoreGroupIdAndName, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllStoreGroupsIdsAndNames")
	}

	var r0 []models.StoreGroupIdAndName
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.StoreGroupIdAndName, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.StoreGroupIdAndName); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.StoreGroupIdAndName)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStoreGroupByID provides a mock function with given fields: ctx, storeGroupID
func (_m *Service) GetStoreGroupByID(ctx context.Context, storeGroupID string) (models.StoreGroup, error) {
	ret := _m.Called(ctx, storeGroupID)

	if len(ret) == 0 {
		panic("no return value specified for GetStoreGroupByID")
	}

	var r0 models.StoreGroup
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.StoreGroup, error)); ok {
		return rf(ctx, storeGroupID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.StoreGroup); ok {
		r0 = rf(ctx, storeGroupID)
	} else {
		r0 = ret.Get(0).(models.StoreGroup)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, storeGroupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStoreGroupByStoreID provides a mock function with given fields: ctx, storeID
func (_m *Service) GetStoreGroupByStoreID(ctx context.Context, storeID string) (models.StoreGroup, error) {
	ret := _m.Called(ctx, storeID)

	if len(ret) == 0 {
		panic("no return value specified for GetStoreGroupByStoreID")
	}

	var r0 models.StoreGroup
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.StoreGroup, error)); ok {
		return rf(ctx, storeID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.StoreGroup); ok {
		r0 = rf(ctx, storeID)
	} else {
		r0 = ret.Get(0).(models.StoreGroup)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, storeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStoreGroupLegalEntities provides a mock function with given fields: ctx, id
func (_m *Service) GetStoreGroupLegalEntities(ctx context.Context, id string) ([]legalentitymodels.LegalEntityView, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetStoreGroupLegalEntities")
	}

	var r0 []legalentitymodels.LegalEntityView
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]legalentitymodels.LegalEntityView, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []legalentitymodels.LegalEntityView); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]legalentitymodels.LegalEntityView)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStoreGroupsWithFilter provides a mock function with given fields: ctx, query
func (_m *Service) GetStoreGroupsWithFilter(ctx context.Context, query selector.StoreGroup) ([]models.StoreGroup, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for GetStoreGroupsWithFilter")
	}

	var r0 []models.StoreGroup
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, selector.StoreGroup) ([]models.StoreGroup, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, selector.StoreGroup) []models.StoreGroup); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.StoreGroup)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, selector.StoreGroup) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsTopPartner provides a mock function with given fields: storeGroup
func (_m *Service) IsTopPartner(storeGroup models.StoreGroup) bool {
	ret := _m.Called(storeGroup)

	if len(ret) == 0 {
		panic("no return value specified for IsTopPartner")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(models.StoreGroup) bool); ok {
		r0 = rf(storeGroup)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// UpdateStoreGroup provides a mock function with given fields: ctx, group
func (_m *Service) UpdateStoreGroup(ctx context.Context, group models.UpdateStoreGroup) error {
	ret := _m.Called(ctx, group)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStoreGroup")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UpdateStoreGroup) error); ok {
		r0 = rf(ctx, group)
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


// CreateDirectPromoBanners provides a mock function with given fields: ctx, storeGroupID string, banner models.DirectPromoBanner
func (_m *Service) CreateDirectPromoBanners(ctx context.Context, storeGroupID string, banner models.DirectPromoBanner) error {
	ret := _m.Called(ctx, storeGroupID, banner)
	if len(ret) == 0 {
		panic("no return value specified for CreateDirectPromoBanners")
	}
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, models.DirectPromoBanner) error); ok {
		r0 = rf(ctx, storeGroupID, banner)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// UpdateDirectPromoBanners provides a mock function with given fields: ctx, storeGroupID string, banner models.UpdateDirectPromoBanner
func(_m *Service) UpdateDirectPromoBanners(ctx context.Context, storeGroupID string, banner models.UpdateDirectPromoBanner) error{
	ret := _m.Called(ctx, storeGroupID, banner)
	if len(ret) == 0 {
		panic("no return value specified for CreateDirectPromoBanners")
	}
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, models.UpdateDirectPromoBanner) error); ok {
		r0 = rf(ctx, storeGroupID, banner)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// GetDirectPromoBannerByID provides a mock function with given fields: ctx, storeGroupID, directPromoID string
func (_m *Service)GetDirectPromoBannerByID(ctx context.Context, storeGroupID, directPromoID string) (models.DirectPromoBanner, error){
	ret := _m.Called(ctx, storeGroupID, directPromoID)

	if len(ret) == 0 {
		panic("no return value specified for GetDirectPromoBannerByID")
	}
	var r0 models.DirectPromoBanner
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) models.DirectPromoBanner); ok {
		r0 = rf(ctx, storeGroupID,directPromoID )
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(models.DirectPromoBanner)
		}
	}
	return r0, r1
}

// GetAllDirectPromoBannersByStoreGroup provides a mock function with given fields: ctx, storeGroupID
func (_m *Service)GetAllDirectPromoBannersByStoreGroup(ctx context.Context, storeGroupID string) ([]models.DirectPromoBanner, error){
	ret := _m.Called(ctx, storeGroupID)

	if len(ret) == 0 {
		panic("no return value specified for GetAllDirectPromoBannersByStoreGroup")
	}

	var r0 []models.DirectPromoBanner
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]models.DirectPromoBanner, error)); ok {
		return rf(ctx, storeGroupID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []models.DirectPromoBanner); ok {
		r0 = rf(ctx, storeGroupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.DirectPromoBanner)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, storeGroupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteDirectPromoBannerByID provides a mock function with given fields: ctx, storeGroupID, directPromoID string
func (_m *Service)	DeleteDirectPromoBannerByID(ctx context.Context, storeGroupID, directPromoID string) error{
	ret := _m.Called(ctx, storeGroupID, directPromoID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteDirectPromoBannerByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, storeGroupID, directPromoID)
	} else {
		r0 = ret.Get(0).(error)
	}

	return r0
}