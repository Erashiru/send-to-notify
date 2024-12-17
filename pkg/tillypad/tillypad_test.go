package tillypad

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/tillypad/yandexDeliveryProtocolTillypad"
	"github.com/kwaaka-team/orders-core/pkg/tillypad/yandexDeliveryProtocolTillypad/clients"
	"testing"
)

func TestTillypad(t *testing.T) {

	clientId := "tillypad"
	clientSecret := "S&drqg_9T6FNMuyJ"

	ctx := context.Background()

	cli, err := yandexDeliveryProtocolTillypad.NewTillypadClient(clients.Config{
		BaseURL:      "http://service.tillypad.ru:8059",
		Protocol:     "http",
		ClientId:     clientId,
		ClientSecret: clientSecret,
		PathPrefix:   "/yandex-eda/yami-alm-wolt",
	})
	if err != nil {
		t.Error(err)
		return
	}

	menu, err := cli.GetMenu(ctx, "955240CB-0815-EC4A-8444-F83FE036E6FC")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(menu)
}
