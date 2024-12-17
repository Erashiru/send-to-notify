package utils

import (
	"github.com/golang-jwt/jwt/v4"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type JWTService interface {
	GenerateJWTToken(clientID, clientSecret, service string) string
	ValidateToken(encodedToken string) (*jwt.Token, *JWTClaims, error)
}

type jwtServices struct {
	secretKey string
}

type JWTClaims struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Service      string `json:"service"`
	jwt.RegisteredClaims
}

func JWTAuthService(secretKey string) JWTService {
	return &jwtServices{
		secretKey: secretKey,
	}
}

func (svc *jwtServices) GenerateJWTToken(clientID, clientSecret, service string) string {
	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &JWTClaims{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Service:      service,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(svc.secretKey))

	if err != nil {
		log.Err(err).Msg("Unexpected error")
		// If there is an error in creating the JWT return an core server error
		return ""
	}

	return tokenString
}

func (svc *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, *JWTClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errs.ErrTokenIsNotValid
		}
		return []byte(svc.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(encodedToken, &JWTClaims{}, keyFunc)

	if err != nil {
		log.Err(err).Msg("Invalid signature")
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errs.ErrExpiredToken) {
			return nil, nil, errs.ErrExpiredToken
		}
		return nil, nil, errs.ErrTokenIsNotValid
	}

	claims, ok := jwtToken.Claims.(*JWTClaims)

	if !ok {
		return nil, nil, errs.ErrTokenIsNotValid
	}

	return jwtToken, claims, nil
}
