package dto

import "time"

type JWTRequest struct {
	SecretKey     string `json:"secret_key"`
	UID           string `json:"uid"`
	Token         string `json:"token"`
	LifeTimeToken int
}
type JWTResponse struct {
	ExpTime time.Time `json:"exp_time"`
	Token   string    `json:"token"`
	UID     string    `json:"uid"`
	Phone   string    `json:"phone"`
}
