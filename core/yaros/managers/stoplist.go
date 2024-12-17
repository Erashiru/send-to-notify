package managers

import (
	menuCli "github.com/kwaaka-team/orders-core/pkg/menu"
	orderCli "github.com/kwaaka-team/orders-core/pkg/order"
	storeCli "github.com/kwaaka-team/orders-core/pkg/store"
)

type YarosManager interface {
}

type Manager struct {
	menuCli  menuCli.Client
	orderCli orderCli.Client
	storeCli storeCli.Client
}

func NewYarosManager(storeCli storeCli.Client, menuCli menuCli.Client, orderCli orderCli.Client) YarosManager {
	return &Manager{
		menuCli:  menuCli,
		storeCli: storeCli,
		orderCli: orderCli,
	}
}
