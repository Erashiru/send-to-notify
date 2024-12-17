package auth_client

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/externalapi/config"
	"github.com/kwaaka-team/orders-core/core/externalapi/database"
	"github.com/kwaaka-team/orders-core/core/externalapi/database/drivers"
	"github.com/kwaaka-team/orders-core/core/externalapi/managers"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/pkg/externalapi/auth_client/dto"
	"github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
)

type AuthClient struct {
	AuthClient managers.AuthClient
}

type Client interface {
	SetCreds(ctx context.Context, req dto.SetCredentialsRequest) (string, error)
	GetListID(ctx context.Context) ([]dto.ClientIDResponse, error)
	SetCredentianal(ctx context.Context, req dto.CredentianRequest) error
}

func NewAuthClient() (Client, error) {
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

	if err = ds.Connect(); err != nil {
		return nil, fmt.Errorf("cannot connect to datastore: %s", err)
	}
	storeCli, err := store.NewClient(storeModels.Config{})
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		AuthClient: managers.NewAuthClientManager(ds, storeCli, opts.AppSecret, opts.EmenuGlobalConfiguration),
	}, nil
}

func (cli *AuthClient) SetCreds(ctx context.Context, req dto.SetCredentialsRequest) (string, error) {
	res, err := cli.AuthClient.GenerateToken(ctx, models.Credentials{
		RestID:  req.RestID,
		Service: req.Service,
		AuthenticateData: models.AuthenticateData{
			ClientID:     req.ClientID,
			ClientSecret: req.ClientSecret,
			GrantType:    req.GrantType,
			Scope:        req.Scope,
		},
	})
	if err != nil {
		return "", err
	}
	return res, nil
}

func (cli *AuthClient) GetListID(ctx context.Context) ([]dto.ClientIDResponse, error) {
	res, err := cli.AuthClient.GetListID(ctx)
	if err != nil {
		return nil, err
	}
	return dto.ToDTO(res), nil
}

func (cli *AuthClient) SetCredentianal(ctx context.Context, req dto.CredentianRequest) error {
	return cli.AuthClient.SetCredential(ctx, req.ToModel())
}
