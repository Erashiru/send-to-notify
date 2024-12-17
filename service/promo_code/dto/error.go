package dto

import "github.com/pkg/errors"

var (
	ErrPromoCodeNotFound  = errors.New("promo code not found")
	ErrInvalidPromoCodeID = errors.New("invalid promocode id")
)
