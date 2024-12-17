package dto

import (
	"github.com/pkg/errors"
)

var (
	StatusIsNotExist     = errors.New("status is not exist")
	PosSystemIsIncorrect = errors.New("pos system is incorrect")
)

const NO = "Нет"

//go:generate stringer -type=PosStatus
type PosStatus int

const (
	ACCEPTED PosStatus = iota + 1
	NEW
	WAIT_COOKING
	READY_FOR_COOKING
	COOKING_COMPLETE
	COOKING_STARTED
	CLOSED
	READY_FOR_PICKUP
	PICKED_UP_BY_CUSTOMER
	ON_WAY
	OUT_FOR_DELIVERY
	DELIVERED
	CANCELLED_BY_POS_SYSTEM
	PAYMENT_NEW
	PAYMENT_IN_PROGRESS
	PAYMENT_SUCCESS
	PAYMENT_CANCELLED
	PAYMENT_WAITING
	PAYMENT_DELETED
	FAILED
	WAIT_SENDING
)
