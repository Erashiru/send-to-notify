package firebase_client

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/kwaaka-team/orders-core/core/auth/config"
	"github.com/kwaaka-team/orders-core/pkg/auth/firebase_client/dto"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
)

type Client interface {
	FindUser(ctx context.Context, uid string) (dto.User, error)
}

type authCore struct {
	cli *auth.Client
}

func New() (Client, error) {
	ctx := context.Background()

	opts, err := config.LoadConfig(context.Background())
	if err != nil {
		return nil, err
	}
	// configure database URL
	conf := &firebase.Config{
		DatabaseURL: opts.FireBaseDBURL,
	}

	// fetch service account key
	resp, err := http.Get(opts.FireBaseFilePath)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	htmlData, err := io.ReadAll(resp.Body) //<--- here!

	if err != nil {
		return nil, err
	}
	opt := option.WithCredentialsJSON(htmlData)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("error in initializing firebase_client app: ", err)
	}

	cli, err := app.Auth(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase_client DB client: ", err)
	}

	return &authCore{
		cli: cli,
	}, nil
}

func (a *authCore) FindUser(ctx context.Context, uid string) (dto.User, error) {
	user, err := a.cli.GetUser(ctx, uid)
	if err != nil {
		return dto.User{}, err
	}

	return dto.User{
		UID:         uid,
		PhoneNumber: user.PhoneNumber,
	}, nil
}
