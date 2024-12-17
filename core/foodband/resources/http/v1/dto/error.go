package dto

import "github.com/pkg/errors"

var (
	ErrStoreNotFound        = errors.New("store not found")
	ErrStoreNotIntegated    = errors.New("store is not integrated with delivery service")
	ErrNotFoundActiveDsMenu = errors.New("not found active menu for delivery service")
)

type ErrorResponse struct {
	Details string `json:"details"`
}
