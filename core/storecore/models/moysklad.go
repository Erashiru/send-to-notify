package models

type StoreMoySkladConfig struct {
	UserName       string                     `bson:"username" json:"username"`
	Password       string                     `bson:"password" json:"password"`
	OrderID        string                     `bson:"order_id" json:"order_id"`
	OrganizationID string                     `bson:"organization_id" json:"organization_id"`
	Status         MoySkladStatus             `bson:"status" json:"status"`
	SendToPos      bool                       `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketPlace  bool                       `bson:"is_marketplace" json:"is_marketplace"`
	PaymentTypes   DeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	Code           string                     `bson:"code" json:"code"`
}

type MoySkladStatus struct {
	ID         string `bson:"id" json:"id"`
	Name       string `bson:"name" json:"name"`
	StatusType string `bson:"status_type" json:"status_type"`
}
