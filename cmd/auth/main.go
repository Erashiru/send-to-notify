package main

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/auth/qrmenu"
	"github.com/kwaaka-team/orders-core/pkg/auth/qrmenu/dto"
	"log"
	"os"
)

func main() {

	os.Setenv("REGION", "eu-west-1")
	os.Setenv("SECRET_ENV", "LocalEnvs")

	cli, err := qrmenu.NewClient()
	if err != nil {
		log.Fatalf("Create new client error: %v", err)
	}

	res, err := cli.GenerateJWT(context.TODO(), dto.JWTRequest{
		SecretKey:     "secret",
		UID:           "devstackq",
		LifeTimeToken: 1,
	})
	if err != nil {
		log.Fatalf("GenerateJWT error: %v", err)
	}

	jwtRes, err := cli.CheckJWT(context.TODO(), dto.JWTRequest{
		Token: res.Token},
	)

	fmt.Println("res checkJwt", err, jwtRes)

}
