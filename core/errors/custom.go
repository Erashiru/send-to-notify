package errors

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type Error struct {
	err *multierror.Error
}

func (ve *Error) Append(err error, errs ...error) {

	var v Error
	if errors.As(err, &v) {
		ve.err = multierror.Append(ve.err, v.err)
	} else {
		ve.err = multierror.Append(ve.err, err)
	}

	for i := range errs {
		var v Error
		if errors.As(errs[i], &v) {
			ve.err = multierror.Append(ve.err, v.err)
		} else {
			ve.err = multierror.Append(ve.err, errs[i])
		}
	}
}

func (ve Error) Error() string {
	return ve.err.Error()
}

func (ve Error) ErrorOrNil() error {

	if ve.err == nil {
		return nil
	}

	if ve.err.Len() == 0 {
		return nil
	}

	return ve
}
