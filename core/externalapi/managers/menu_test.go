package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/externalapi/database/drivers"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/pkg/menu"
	"github.com/kwaaka-team/orders-core/pkg/store"
	"reflect"
	"testing"
)

func TestMenuClientManager_GetMenu(t *testing.T) {
	type fields struct {
		ds       drivers.DataStore
		menuCli  menu.Client
		storeCli store.Client
	}
	type args struct {
		ctx          context.Context
		storeID      string
		service      string
		clientSecret string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Menu
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &MenuClientManager{
				ds:       tt.fields.ds,
				menuCli:  tt.fields.menuCli,
				storeCli: tt.fields.storeCli,
			}
			got, err := manager.GetMenu(tt.args.ctx, tt.args.storeID, tt.args.service, tt.args.clientSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMenu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMenu() got = %v, want %v", got, tt.want)
			}
		})
	}
}
