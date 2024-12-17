package grafana

import (
	"github.com/kwaaka-team/orders-core/pkg/grafana/clients"
	"github.com/kwaaka-team/orders-core/pkg/grafana/config"
)

func NewGrafanaClient(config config.Config) clients.Grafana {
	return clients.NewGrafanaClient(config)
}
