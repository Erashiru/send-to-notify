package models

type OrderStatus string

const (
	STATUS_NEW OrderStatus = "NEW"
)

func (s OrderStatus) String() string {
	return string(s)
}

type WoltStatus string

const (
	Created  WoltStatus = "CREATED"
	Canceled WoltStatus = "CANCELED"
	Rejected WoltStatus = "rejected"
)

func (status WoltStatus) String() string {
	return string(status)
}
