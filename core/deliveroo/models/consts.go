package models

type DeliveryService string

const (
	DELIVEROO DeliveryService = "deliveroo"

	OrderCreate = "order.new"
	OrderUpdate = "order.status_update"
	MenuUpload  = "menu.upload_result"
)

func (d DeliveryService) String() string {
	return string(d)
}
