package yandexDeliveryProtocolTillypad

import (
	"github.com/kwaaka-team/orders-core/pkg/tillypad/yandexDeliveryProtocolTillypad/clients"
	"github.com/kwaaka-team/orders-core/pkg/tillypad/yandexDeliveryProtocolTillypad/clients/http"
)

func NewTillypadClient(cfg clients.Config) (clients.Client, error) {
	return http.NewClient(cfg)
}
