package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type UserStore interface {
	Create(ctx context.Context, userStore []models.UserStore) error
	FindUsers(ctx context.Context, query selector.User) ([]models.UserStore, error)
	DeleteByUsernameAndStoreGroupID(ctx context.Context, username string, storeGroup string) error
	UpdateUserOrderNotifications(ctx context.Context, username, fcmToken string, stores []string) error
}

type UserStoreManager struct {
	userStoreRepository drivers.UserStoreRepository
}

func NewUserStoreManager(ds drivers.Datastore) UserStore {
	return &UserStoreManager{
		userStoreRepository: ds.UserStoreRepository(),
	}
}

func (u UserStoreManager) Create(ctx context.Context, userStores []models.UserStore) error {

	return u.userStoreRepository.Insert(ctx, userStores)
}

func (u UserStoreManager) FindUsers(ctx context.Context, query selector.User) ([]models.UserStore, error) {
	res, err := u.userStoreRepository.FindUsers(ctx, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u UserStoreManager) DeleteByUsernameAndStoreGroupID(ctx context.Context, username string, storeGroup string) error {
	return u.userStoreRepository.Delete(ctx, selector.User{StoreGroupId: storeGroup, Username: username})
}

func (u UserStoreManager) UpdateUserOrderNotifications(ctx context.Context, username, fcmToken string, stores []string) error {
	return u.userStoreRepository.UpdateUserOrderNotifications(ctx, username, fcmToken, stores)
}
