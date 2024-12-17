package http

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/iiko/clients"
	"testing"

	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func TestClient_GetStopList(t *testing.T) {
	ctx := context.Background()
	cli, err := New(&clients.Config{
		ApiLogin: "1442801e",
		BaseURL:  "https://api-ru.iiko.services",
	})

	if err != nil {
		t.Error(err)
	}
	defer cli.Close(ctx)

	orgs, err := cli.GetOrganizations(ctx)
	if err != nil {
		t.Error(err)
	}

	org := orgs[0]

	stoplistItems, err := cli.GetStopList(ctx, models.StopListRequest{
		Organizations: []string{org.ID},
	})
	if err != nil {
		t.Error(err)
	}

	for _, item := range stoplistItems.TerminalGroups {
		t.Log(item)
	}
}
