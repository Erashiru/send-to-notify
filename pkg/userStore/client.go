package userStore

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/storecore/config"
	"github.com/kwaaka-team/orders-core/core/storecore/database"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/userStore/dto"
)

type Client interface {
	Create(ctx context.Context, userStores []models.UserStore) (models.UserStore, error)
	DeleteByUsernameAndStoreGroupID(ctx context.Context, username string, storeGroup string) error
	UpdateUserOrderNotifications(ctx context.Context, username, fcmToken string, stores []string) error
}

type UserStore struct {
	userStoreManager managers.UserStore
}

func NewClient(cfg dto.Config) (Client, error) {
	opts, err := config.LoadConfig(context.Background())
	if err != nil {
		return nil, err
	}

	ds, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSDB,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create datastore %s: %v", opts.DSName, err)
	}

	if err = ds.Connect(cfg.MongoCli); err != nil {
		return nil, fmt.Errorf("cannot connect to datastore: %s", err)
	}

	return &UserStore{
		userStoreManager: managers.NewUserStoreManager(ds),
	}, nil
}

func (u *UserStore) Create(ctx context.Context, userStores []models.UserStore) (models.UserStore, error) {

	err := u.userStoreManager.Create(ctx, userStores)
	if err != nil {
		return models.UserStore{}, err
	}
	return models.UserStore{}, nil
}

func (u *UserStore) DeleteByUsernameAndStoreGroupID(ctx context.Context, username string, storeGroup string) error {
	return u.userStoreManager.DeleteByUsernameAndStoreGroupID(ctx, username, storeGroup)
}

func (u *UserStore) UpdateUserOrderNotifications(ctx context.Context, username, fcmToken string, stores []string) error {
	return u.userStoreManager.UpdateUserOrderNotifications(ctx, username, fcmToken, stores)
}
