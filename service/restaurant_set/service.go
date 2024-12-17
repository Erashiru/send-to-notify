package restaurant_set

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	"github.com/kwaaka-team/orders-core/core/models"
	storeGroupServicePkg "github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/rs/zerolog/log"
	"sync"
)

type Service interface {
	GetRestaurantSetById(ctx context.Context, id string) (models.RestaurantSet, error)
	CreateRestaurantSet(ctx context.Context, set models.RestaurantSet) (string, error)
	GetRestaurantSetInfoWithRestGroup(ctx context.Context, id string) (dto.RestaurantGroupSetResponse, error)
	GetRestaurantSetByDomainName(ctx context.Context, domainName string) (dto.RestaurantGroupSetResponse, error)
}

type ServiceImpl struct {
	storeGroupService storeGroupServicePkg.Service
	repo              Repository
}

func NewService(r Repository, storeGroupService storeGroupServicePkg.Service) (*ServiceImpl, error) {
	return &ServiceImpl{storeGroupService: storeGroupService, repo: r}, nil
}

func (s *ServiceImpl) GetRestaurantSetById(ctx context.Context, id string) (models.RestaurantSet, error) {
	return s.repo.GetRestaurantSetById(ctx, id)
}

func (s *ServiceImpl) CreateRestaurantSet(ctx context.Context, set models.RestaurantSet) (string, error) {
	return s.repo.CreateRestaurantSet(ctx, set)
}

func (s *ServiceImpl) GetRestaurantSetInfoWithRestGroup(ctx context.Context, id string) (dto.RestaurantGroupSetResponse, error) {
	restGroupSet, err := s.repo.GetRestaurantSetById(ctx, id)
	if err != nil {
		return dto.RestaurantGroupSetResponse{}, err
	}

	res := dto.RestaurantGroupSetResponse{
		ID:               restGroupSet.ID,
		Name:             restGroupSet.Name,
		Logo:             restGroupSet.Logo,
		DomainName:       restGroupSet.DomainName,
		HeaderImage:      restGroupSet.HeaderImage,
		RestaurantGroups: make([]dto.RestGroupResponse, 0),
	}

	resultCh := make(chan dto.RestGroupResponse, len(restGroupSet.RestaurantGroupIds))
	errCh := make(chan error, len(restGroupSet.RestaurantGroupIds))
	var wg sync.WaitGroup

	for _, restGroupId := range restGroupSet.RestaurantGroupIds {
		wg.Add(1)

		go func(restGroupId string) {
			defer wg.Done()
			restGroup, err := s.storeGroupService.GetStoreGroupByID(ctx, restGroupId)
			if err != nil {
				errCh <- err
				return
			}

			resultCh <- dto.RestGroupResponse{
				Id:                  restGroup.ID,
				Name:                restGroup.Name,
				ColumnView:          restGroup.ColumnView,
				Description:         restGroup.Description,
				HeaderImage:         restGroup.HeaderImage,
				ExtraLogo:           restGroup.ExtraLogo,
				DomainName:          restGroup.DomainName,
				DefaultRestaurantId: restGroup.DefaultRestaurantId,
				DefaultCity:         restGroup.DefaultCity,
				Tags:                restGroup.Tags,
				Logo:                restGroup.Logo,
			}
		}(restGroupId)
	}

	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	for range restGroupSet.RestaurantGroupIds {
		select {
		case err := <-errCh:
			log.Err(err).Msgf("error while getting restaurant group info: %s", err.Error())
		case restGroupResp := <-resultCh:
			res.RestaurantGroups = append(res.RestaurantGroups, restGroupResp)
		}
	}

	return res, nil
}

func (s *ServiceImpl) GetRestaurantSetByDomainName(ctx context.Context, domainName string) (dto.RestaurantGroupSetResponse, error) {
	restaurantSet, err := s.repo.GetRestaurantSetByDomainName(ctx, domainName)
	if err != nil {
		return dto.RestaurantGroupSetResponse{}, err
	}

	res := dto.RestaurantGroupSetResponse{
		ID:               restaurantSet.ID,
		Name:             restaurantSet.Name,
		Logo:             restaurantSet.Logo,
		DomainName:       restaurantSet.DomainName,
		HeaderImage:      restaurantSet.HeaderImage,
		RestaurantGroups: make([]dto.RestGroupResponse, 0),
	}

	resultCh := make(chan dto.RestGroupResponse, len(restaurantSet.RestaurantGroupIds))
	errCh := make(chan error, len(restaurantSet.RestaurantGroupIds))
	var wg sync.WaitGroup

	for _, restGroupId := range restaurantSet.RestaurantGroupIds {
		wg.Add(1)

		go func(restGroupId string) {
			defer wg.Done()
			restGroup, err := s.storeGroupService.GetStoreGroupByID(ctx, restGroupId)
			if err != nil {
				errCh <- err
				return
			}

			resultCh <- dto.RestGroupResponse{
				Id:                  restGroup.ID,
				Name:                restGroup.Name,
				ColumnView:          restGroup.ColumnView,
				Description:         restGroup.Description,
				HeaderImage:         restGroup.HeaderImage,
				Logo:                restGroup.Logo,
				ExtraLogo:           restGroup.ExtraLogo,
				DomainName:          restGroup.DomainName,
				DefaultRestaurantId: restGroup.DefaultRestaurantId,
				DefaultCity:         restGroup.DefaultCity,
				Tags:                restGroup.Tags,
			}
		}(restGroupId)
	}

	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	for range restaurantSet.RestaurantGroupIds {
		select {
		case err := <-errCh:
			log.Err(err).Msgf("error while getting restaurant group info: %s", err.Error())
		case restGroupResp := <-resultCh:
			res.RestaurantGroups = append(res.RestaurantGroups, restGroupResp)
		}
	}

	return res, nil
}
