package notify

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/notify/config"
	"github.com/kwaaka-team/orders-core/core/notify/managers"

	"github.com/kwaaka-team/orders-core/pkg/notify/dto"
)

type Client interface {
	SendNotification(ctx context.Context, message dto.Message) ([]dto.Result, error)
}

var _ Client = &notify{}

type notify struct {
	cfg       config.Configuration
	notifyMan managers.NotifyManager
}

func New() (Client, error) {

	opts, err := config.LoadConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return notify{
		cfg:       opts,
		notifyMan: managers.NewNotifyManager(opts),
	}, nil
}

func (n notify) SendNotification(ctx context.Context, req dto.Message) ([]dto.Result, error) {

	rsp, err := n.notifyMan.SendNotification(ctx, req.ToModel())
	if err != nil {
		return nil, err
	}
	return dto.FromResults(rsp), err
}
