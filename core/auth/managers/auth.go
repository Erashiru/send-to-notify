package managers

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kwaaka-team/orders-core/core/auth/models"
	"github.com/pkg/errors"
	"time"
)

type Claims struct {
	UID   string `json:"uid"`
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}
type JWTData struct {
	jwt.Claims
	CustomClaims map[string]string `json:"custom_claims"`
}

func (a auth) GenerateJWT(ctx context.Context, req models.JWT) (models.JWT, error) {
	//set expiration time
	var res models.JWT

	expirationTime := time.Now().Add(time.Duration(req.LifeTimeToken) * time.Minute)
	claims := &Claims{
		UID:   req.UID,
		Phone: req.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(req.SecretKey))
	if err != nil {
		return res, err
	}

	res.Token = tokenString
	res.ExpTime = expirationTime

	return res, nil
}

func (a auth) CheckJWT(ctx context.Context, req models.JWT) (models.JWT, error) {
	var (
		claims = &Claims{}
		res    models.JWT
	)

	tkn, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(req.SecretKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return res, err
		}
		return res, err
	}
	if !tkn.Valid {
		return res, errors.New("invalid token")
	}
	res.UID = claims.UID

	return res, nil
}
