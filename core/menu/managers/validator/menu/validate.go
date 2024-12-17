package menu

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/pkg/errors"
)

type ValidatorMenu interface {
	Validate(ctx context.Context, menu models.Menu) error
}

func NewValidatorMenu(delivery string) (ValidatorMenu, error) {

	switch delivery {
	case models.GLOVO.String():
		return newGlovo(), nil
	case models.WOLT.String():
		return newWolt(), nil
	case models.TALABAT.String():
		return newTalabat(), nil
	case models.YANDEX.String():
		return newYandex(), nil
	case models.EXPRESS24.String():
		return newExpress24(), nil
	case "starter_app":
		return newstarterApp(), nil
	}
	return nil, errors.New("MenuValidate: Invalid delivery service")
}
