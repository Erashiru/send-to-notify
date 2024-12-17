package whatsapp

import (
	"errors"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp/clients"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp/clients/http"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewWhatsappClient(conf *clients.Config) (clients.Whatsapp, error) {
	switch conf.Protocol {
	case "http":
		cli, err := http.NewClient(conf)
		if err != nil {
			return nil, err
		}
		return cli, nil
	}

	return nil, ErrInvalidProtocol
}
