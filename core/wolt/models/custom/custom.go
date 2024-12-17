package custom

import (
	"github.com/hashicorp/go-multierror"
)

type Error struct {
	err *multierror.Error
}

func (ve Error) Error() string {
	return ve.err.Error()
}
