package dto

import (
	"github.com/kwaaka-team/orders-core/core/auth/models"
)

type User struct {
	UID         string `json:"uid"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	FCMToken    string `json:"fcm_token"`
}

func (u User) ToModel() models.User {
	return models.User{
		UID:         u.UID,
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		FCMToken:    u.FCMToken,
	}
}

func ToUserDTO(req models.User) User {
	return User{
		UID:         req.UID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		FCMToken:    req.FCMToken,
	}
}
