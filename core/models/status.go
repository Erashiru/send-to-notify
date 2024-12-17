package models

import (
	"github.com/pkg/errors"
)

var (
	StatusIsNotExist      = errors.New("status is not exist")
	PosSystemIsIncorrect  = errors.New("pos system is incorrect")
	InvalidStatusPriority = errors.New("invalid status priority")
)

const NO = "Нет"

//go:generate stringer -type=PosStatus
type PosStatus int

const (
	DELIVERY_SERVICE_MODE = 3
)

const (
	Default3plCustomerName  = "Client"
	Default3plCustomerEmail = "dispatch@kwaaka.com"
	Default3plCustomerPhone = "+77777777777"
)

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

const (
	OUT_FOR_DELIVERY_str      = "OUT_FOR_DELIVERY"
	PICKED_UP_BY_CUSTOMER_str = "PICKED_UP_BY_CUSTOMER"
)

//go:generate stringer -type=AggregatorStatus
type AggregatorStatus int

const (
	Accept AggregatorStatus = iota + 1
	Reject
	Ready
	Confirm
	Delivered
)

type ChocoStatus string

const (
	Accepted           ChocoStatus = "accepted"
	Rejected           ChocoStatus = "reject"
	ReadyForPickup     ChocoStatus = "ready_for_pickup"
	OutForDelivery     ChocoStatus = "out_for_delivery"
	PickedUpByCustomer ChocoStatus = "picked_up_by_customer"
)

func (a ChocoStatus) String() string {
	return string(a)
}

type TalabatStatus string

const (
	OrderAccepted TalabatStatus = "order_accepted"
	OrderRejected TalabatStatus = "order_rejected"
	OrderPickedUp TalabatStatus = "order_picked_up"
	OrderPrepared TalabatStatus = "order_prepared"
)

func (a TalabatStatus) String() string {
	return string(a)
}

type DeliverooStatus string

type DeliverooReason string

const (
	AcceptedDeliveroo  DeliverooStatus = "accepted"
	RejectedDeliveroo  DeliverooStatus = "rejected"
	GonfirmedDeliveroo DeliverooStatus = "confirmed"
	FailedDeliveroo    DeliverooStatus = "failed"

	LocationNotSupported DeliverooReason = "location_not_supported"
	OtherDeliveroo       DeliverooReason = "other"
)

func (d DeliverooReason) String() string {
	return string(d)
}

func (d DeliverooStatus) String() string {
	return string(d)
}

type StarterAppStatus string

const (
	Created      StarterAppStatus = "created"
	Canceled     StarterAppStatus = "canceled"
	Draft        StarterAppStatus = "draft"
	NotConfirmed StarterAppStatus = "notConfirmed"
	Checked      StarterAppStatus = "checked"
	InProgress   StarterAppStatus = "inProgress"
	Cooked       StarterAppStatus = "cooked"
	OnTheWay     StarterAppStatus = "onTheWay"
	Done         StarterAppStatus = "done"
)

func (s StarterAppStatus) String() string {
	return string(s)
}
