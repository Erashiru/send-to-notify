package detector

import (
	errs "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/pkg/errors"

	"net/http"
)

func ErrorHandler(err error) (code int, obj any) {

	switch {
	case errors.Is(err, errs.ErrEmpty) || errors.Is(err, errs.ErrNotFound):
		return http.StatusBadRequest, errs.ErrorResponse{
			Msg:  err.Error(),
			Code: http.StatusBadRequest,
		}
	default:
		return http.StatusInternalServerError, errs.ErrorResponse{
			Msg:  err.Error(),
			Code: http.StatusInternalServerError,
		}
	}
}
