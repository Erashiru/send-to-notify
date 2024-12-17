package storegroup

import (
	"context"
	selector2 "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/legalentity/models"
	"github.com/pkg/errors"
)

type Service interface {
	GetStoreGroupByID(ctx context.Context, storeGroupID string) (storeModels.StoreGroup, error)
	GetStoreGroupByStoreID(ctx context.Context, storeID string) (storeModels.StoreGroup, error)
	IsTopPartner(storeGroup storeModels.StoreGroup) bool
	GetStoreGroupsWithFilter(ctx context.Context, query selector2.StoreGroup) ([]storeModels.StoreGroup, error)
	CreateStoreGroup(ctx context.Context, group storeModels.StoreGroup) (string, error)
	UpdateStoreGroup(ctx context.Context, group storeModels.UpdateStoreGroup) error
	GetStoreGroupLegalEntities(ctx context.Context, id string) ([]models.LegalEntityView, error)
	GetAllStoreGroupsIdsAndNames(ctx context.Context) ([]storeModels.StoreGroupIdAndName, error)
	AddBrandInfo(ctx context.Context, brandInfo storeModels.BrandInfo, restGroupID string) error

	CreateDirectPromoBanners(ctx context.Context, storeGroupID string, banner storeModels.DirectPromoBanner) error
	UpdateDirectPromoBanners(ctx context.Context, storeGroupID string, banner storeModels.UpdateDirectPromoBanner) error
	GetDirectPromoBannerByID(ctx context.Context, storeGroupID, directPromoID string) (storeModels.DirectPromoBanner, error)
	GetAllDirectPromoBannersByStoreGroup(ctx context.Context, storeGroupID string) ([]storeModels.DirectPromoBanner, error)
	DeleteDirectPromoBannerByID(ctx context.Context, storeGroupID, directPromoID string) error
}

type serviceImpl struct {
	storeGroupRepository Repository
}

func NewService(storeGroupRepository Repository) (*serviceImpl, error) {
	if storeGroupRepository == nil {
		return nil, errors.New("store group repository is nil")
	}
	return &serviceImpl{
		storeGroupRepository: storeGroupRepository,
	}, nil
}

func (s *serviceImpl) IsTopPartner(storeGroup storeModels.StoreGroup) bool {
	return storeGroup.IsTopPartner
}

func (s *serviceImpl) GetStoreGroupByID(ctx context.Context, storeGroupID string) (storeModels.StoreGroup, error) {
	return s.storeGroupRepository.GetStoreGroupByID(ctx, storeGroupID)
}

func (s *serviceImpl) GetStoreGroupByStoreID(ctx context.Context, storeID string) (storeModels.StoreGroup, error) {
	storeGroup, err := s.storeGroupRepository.GetStoreGroupByStoreID(ctx, storeID)
	if err != nil {
		return storeModels.StoreGroup{}, err
	}

	return storeGroup, nil
}

func (s *serviceImpl) GetStoreGroupsWithFilter(ctx context.Context, query selector2.StoreGroup) ([]storeModels.StoreGroup, error) {
	storeGroups, err := s.storeGroupRepository.GetStoreGroupsWithFilter(ctx, query)
	if err != nil {
		return []storeModels.StoreGroup{}, err
	}

	return storeGroups, nil
}

func (s *serviceImpl) CreateStoreGroup(ctx context.Context, group storeModels.StoreGroup) (string, error) {
	insertedID, err := s.storeGroupRepository.CreateStoreGroup(ctx, group)
	if err != nil {
		return "", err
	}

	return insertedID, nil
}

func (s *serviceImpl) UpdateStoreGroup(ctx context.Context, group storeModels.UpdateStoreGroup) error {
	if err := s.storeGroupRepository.UpdateStoreGroup(ctx, group); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) GetStoreGroupLegalEntities(ctx context.Context, id string) ([]models.LegalEntityView, error) {
	legalEntities, err := s.storeGroupRepository.GetStoreGroupLegalEntities(ctx, id)

	if err != nil {
		return []models.LegalEntityView{}, err
	}

	return legalEntities, nil
}

func (s *serviceImpl) GetAllStoreGroupsIdsAndNames(ctx context.Context) ([]storeModels.StoreGroupIdAndName, error) {
	resp, err := s.storeGroupRepository.GetAllStoreGroupsIdsAndNames(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *serviceImpl) AddBrandInfo(ctx context.Context, brandInfo storeModels.BrandInfo, restGroupID string) error {
	return s.storeGroupRepository.AddBrandInfo(ctx, brandInfo, restGroupID)
}

func (s *serviceImpl) CreateDirectPromoBanners(ctx context.Context, storeGroupID string, banner storeModels.DirectPromoBanner) error {
	return s.storeGroupRepository.CreateDirectPromoBanners(ctx, storeGroupID, banner)
}

func (s *serviceImpl) UpdateDirectPromoBanners(ctx context.Context, storeGroupID string, banner storeModels.UpdateDirectPromoBanner) error {
	return s.storeGroupRepository.UpdateDirectPromoBanners(ctx, storeGroupID, banner)
}

func (s *serviceImpl) GetDirectPromoBannerByID(ctx context.Context, storeGroupID, directPromoID string) (storeModels.DirectPromoBanner, error) {
	storeGroup, err := s.storeGroupRepository.GetStoreGroupByID(ctx, storeGroupID)
	if err != nil {
		return storeModels.DirectPromoBanner{}, err
	}
	for _, banner := range storeGroup.DirectPromoBanners {
		if banner.ID == directPromoID {
			return banner, nil
		}
	}
	return storeModels.DirectPromoBanner{}, nil
}

func (s *serviceImpl) GetAllDirectPromoBannersByStoreGroup(ctx context.Context, storeGroupID string) ([]storeModels.DirectPromoBanner, error) {
	return s.storeGroupRepository.GetAllDirectPromoBannersByStoreGroup(ctx, storeGroupID)
}

func (s *serviceImpl) DeleteDirectPromoBannerByID(ctx context.Context, storeGroupID, directPromoID string) error {
	return s.storeGroupRepository.DeleteDirectPromoBannerByID(ctx, storeGroupID, directPromoID)
}
