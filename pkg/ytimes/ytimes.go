package ytimes

import (
	"github.com/kwaaka-team/orders-core/pkg/ytimes/clients"
	"github.com/kwaaka-team/orders-core/pkg/ytimes/clients/http"
)

func New(cfg clients.Config) clients.Client {
	return http.NewClient(cfg)
}
