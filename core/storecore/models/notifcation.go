package models

type Notification struct {
	Whatsapp Whatsapp `bson:"whatsapp" json:"whatsapp"`
}

type Whatsapp struct {
	Receivers []WhatsappReceiver `bson:"receivers" json:"receivers"`
}

type WhatsappReceiver struct {
	Name        string `bson:"name" json:"name"`
	PhoneNumber string `bson:"phone_number" json:"phone_number"`
	IsActive    bool   `bson:"is_active" json:"is_active"`
}
