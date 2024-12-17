package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/auth/database/datastore/drivers"
	"github.com/kwaaka-team/orders-core/core/auth/managers/validator"
	models2 "github.com/kwaaka-team/orders-core/core/auth/models"
	"github.com/kwaaka-team/orders-core/core/auth/models/selector"
)

type AuthManager interface {
	UpdateUserInfo(ctx context.Context, user models2.User) error
	CreateUser(ctx context.Context, user models2.User) error
	FindUser(ctx context.Context, query selector.User) (models2.User, error)
	GenerateJWT(ctx context.Context, req models2.JWT) (models2.JWT, error)
	CheckJWT(ctx context.Context, req models2.JWT) (models2.JWT, error)
}

type auth struct {
	userRepository drivers.UserRepository
	userValidator  validator.User
}

func NewAuthManager(userRepository drivers.UserRepository, userValidator validator.User) AuthManager {
	return &auth{
		userRepository: userRepository,
		userValidator:  userValidator,
	}
}

func (a auth) CreateUser(ctx context.Context, user models2.User) error {
	if err := a.userValidator.ValidateUID(ctx, selector.NewEmptyUser().SetUID(user.UID)); err != nil {
		return err
	}
	return a.userRepository.CreateUser(ctx, user)
}

func (a auth) FindUser(ctx context.Context, query selector.User) (models2.User, error) {
	if err := a.userValidator.ValidateUID(ctx, query); err != nil {
		return models2.User{}, err
	}
	return a.userRepository.GetUser(ctx, query)
}

func (a auth) UpdateUserInfo(ctx context.Context, user models2.User) error {
	if err := a.userValidator.ValidateUID(ctx, selector.NewEmptyUser().SetUID(user.UID)); err != nil {
		return err
	}
	return a.userRepository.UpdateUserInfo(ctx, user)
}
