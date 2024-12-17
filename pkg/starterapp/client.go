package starterapp

import (
	"github.com/kwaaka-team/orders-core/pkg/starterapp/clients"
	"github.com/kwaaka-team/orders-core/pkg/starterapp/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewStarterAppClient(cfg *clients.Config) (clients.StarterApp, error) {
	switch cfg.Protocol {
	case "http":
		cli, err := http.NewClient(cfg)
		if err != nil {
			return nil, err
		}
		return cli, nil
	}
	return nil, ErrInvalidProtocol
}
