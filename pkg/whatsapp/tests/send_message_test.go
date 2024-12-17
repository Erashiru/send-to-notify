package tests

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp/clients"
	"testing"
)

func TestSendMessage(t *testing.T) {
	cli, err := whatsapp.NewWhatsappClient(&clients.Config{
		Protocol:  "http",
		AuthToken: "512fbnbbkpna1890",
		Instance:  "instance63594",
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err = cli.SendMessage(context.TODO(), "+77066776235", "Hello Buddy"); err != nil {
		t.Error(err)
		return
	}

}
