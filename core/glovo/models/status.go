package models

type OrderStatus string

const (
	STATUS_NEW                           OrderStatus = "NEW"
	STATUS_PENDING                       OrderStatus = "PENDING"
	STATUS_PROCESSING                    OrderStatus = "PROCESSING"
	STATUS_ACCEPTED                      OrderStatus = "ACCEPTED"
	STATUS_SKIPPED                       OrderStatus = "SKIPPED"
	STATUS_COOKING_STARTED               OrderStatus = "COOKING_STARTED"
	STATUS_COOKING_COMPLETE              OrderStatus = "COOKING_COMPLETE"
	STATUS_READY_FOR_PICKUP              OrderStatus = "READY_FOR_PICKUP"
	STATUS_OUT_FOR_DELIVERY              OrderStatus = "OUT_FOR_DELIVERY"
	STATUS_PICKED_UP_BY_CUSTOMER         OrderStatus = "PICKED_UP_BY_CUSTOMER"
	STATUS_DELIVERED                     OrderStatus = "DELIVERED"
	STATUS_CLOSED                        OrderStatus = "CLOSED"
	STATUS_CANCELLED_BY_DELIVERY_SERVICE OrderStatus = "CANCELLED_BY_DELIVERY_SERVICE"
	STATUS_CANCELLED_BY_POS_SYSTEM       OrderStatus = "CANCELLED_BY_POS_SYSTEM"
	STATUS_FAILED                        OrderStatus = "FAILED"
)

func (s OrderStatus) String() string {
	return string(s)
}
