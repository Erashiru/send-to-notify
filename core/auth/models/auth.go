package models

import "time"

type JWT struct {
	UID           string `json:"uid"`
	Phone         string `json:"phone"`
	SecretKey     string `json:"secret_key"`
	LifeTimeToken int    `json:"expiration"`

	ExpTime time.Time `json:"exp_time"`
	Token   string    `json:"token"`
}
