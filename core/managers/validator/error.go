package validator

import "github.com/pkg/errors"

var (
	ErrFailed             = errors.New("Order failed")
	ErrPassed             = errors.New("Order passed")
	ErrCastingPos         = errors.New("cant cast POS order")
	ErrSendingToQue       = errors.New("Sending message to telegram by using que error ")
	ErrNotifyingInBrowser = errors.New("Sending message to browser with firebase error ")
	ErrIgnoringPos        = errors.New("need heal order, bad order restrictions")
	ErrPosInitialize      = errors.New("pos initialize error")
)
