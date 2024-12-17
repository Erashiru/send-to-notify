package models

type User struct {
	UID         string `bson:"uid"`
	Name        string `bson:"name"`
	PhoneNumber string `bson:"phone_number"`
	FCMToken    string `bson:"fcm_token"`
}
