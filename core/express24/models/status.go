package models

type OrderStatus string

const (
	STATUS_NEW OrderStatus = "NEW"
)

func (s OrderStatus) String() string {
	return string(s)
}
