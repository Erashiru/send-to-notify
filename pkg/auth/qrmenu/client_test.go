package qrmenu

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/auth/managers"
	"github.com/kwaaka-team/orders-core/core/auth/managers/validator"
	"github.com/kwaaka-team/orders-core/pkg/auth/qrmenu/dto"
	"testing"
)

func Test_authCore_GenerateJWT(t *testing.T) {
	type fields struct {
		authManager   managers.AuthManager
		authValidator validator.User
	}
	type args struct {
		ctx    context.Context
		jwtReq dto.JWTRequest
	}
	am := managers.NewAuthManager(nil, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.JWTResponse
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				authManager:   am,
				authValidator: nil,
			},
			args: args{
				ctx: context.TODO(),
				jwtReq: dto.JWTRequest{
					SecretKey:     "secret",
					UID:           "devstackq",
					LifeTimeToken: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authCore{
				authManager:   tt.fields.authManager,
				authValidator: tt.fields.authValidator,
			}

			got, err := a.GenerateJWT(tt.args.ctx, tt.args.jwtReq)

			fmt.Printf("res test GenerateJWT %v %v %v", got, err, a)
		})
	}
}
