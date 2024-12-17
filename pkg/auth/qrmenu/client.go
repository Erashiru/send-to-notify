package qrmenu

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/auth/config"
	"github.com/kwaaka-team/orders-core/core/auth/database"
	"github.com/kwaaka-team/orders-core/core/auth/database/datastore/drivers"
	"github.com/kwaaka-team/orders-core/core/auth/managers"
	"github.com/kwaaka-team/orders-core/core/auth/managers/validator"
	"github.com/kwaaka-team/orders-core/core/auth/models"
	"github.com/kwaaka-team/orders-core/core/auth/models/selector"
	dto2 "github.com/kwaaka-team/orders-core/pkg/auth/qrmenu/dto"
	"strconv"
)

type Client interface {
	UpdateUserInfo(ctx context.Context, user dto2.User) error
	CreateUser(ctx context.Context, user dto2.User) error
	FindUser(ctx context.Context, user dto2.User) (dto2.User, error)
	GenerateJWT(ctx context.Context, req dto2.JWTRequest) (dto2.JWTResponse, error)
	CheckJWT(ctx context.Context, req dto2.JWTRequest) (result dto2.JWTResponse, err error)
}

type authCore struct {
	authManager   managers.AuthManager
	authValidator validator.User
	conf          config.Configuration
}

func NewClient() (Client, error) {
	opts, err := config.LoadConfig(context.Background())
	if err != nil {
		return nil, err
	}

	// Connecting to DataStore
	ds, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSUSERDB,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create datastore %s: %v", opts.DSName, err)
	}

	if err = ds.Connect(); err != nil {
		return nil, fmt.Errorf("cannot connect to datastore: %s", err)
	}

	return &authCore{
		conf:        opts,
		authManager: managers.NewAuthManager(ds.AuthRepository(), validator.NewUserValidator()),
	}, nil
}

func (a *authCore) FindUser(ctx context.Context, user dto2.User) (dto2.User, error) {
	resp, err := a.authManager.FindUser(ctx, selector.NewEmptyUser().
		SetUID(user.UID).
		SetPhoneNumber(user.PhoneNumber),
	)

	if err != nil {
		return dto2.User{}, err
	}
	return dto2.ToUserDTO(resp), nil
}

func (a *authCore) CreateUser(ctx context.Context, user dto2.User) error {
	if err := a.authManager.CreateUser(ctx, user.ToModel()); err != nil {
		return err
	}
	return nil
}

func (a *authCore) GenerateJWT(ctx context.Context, jwtReq dto2.JWTRequest) (dto2.JWTResponse, error) {

	var (
		req    models.JWT
		jwtRes dto2.JWTResponse
	)
	//to/fromModel
	lifeTimeInt, err := strconv.Atoi(a.conf.JWTConfiguration.TokenLifeTime)
	if err != nil {
		return jwtRes, err
	}
	req.LifeTimeToken = lifeTimeInt
	req.SecretKey = a.conf.JWTConfiguration.SecretKey
	req.UID = jwtReq.UID

	fmt.Printf("jwt Cl request layer %v", req)

	res, err := a.authManager.GenerateJWT(ctx, req)
	if err != nil {
		return jwtRes, err
	}

	jwtRes.ExpTime = res.ExpTime
	jwtRes.Token = res.Token
	jwtRes.Phone = res.Phone

	return jwtRes, nil
}

func (a *authCore) CheckJWT(ctx context.Context, jwtReq dto2.JWTRequest) (result dto2.JWTResponse, err error) {
	req := models.JWT{
		Token:     jwtReq.Token,
		SecretKey: a.conf.JWTConfiguration.SecretKey,
	}
	res, err := a.authManager.CheckJWT(ctx, req)
	if err != nil {
		return dto2.JWTResponse{}, err
	}
	result.UID = res.UID
	result.Phone = res.Phone

	return result, nil
}

func (a *authCore) UpdateUserInfo(ctx context.Context, user dto2.User) error {
	if err := a.authManager.UpdateUserInfo(ctx, user.ToModel()); err != nil {
		return err
	}
	return nil
}
