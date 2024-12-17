package dto

import "time"

type Response struct {
	InvalidParams []InvalidParam `json:"invalid-params"`
	Title         string         `json:"title"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type AcceptOrderRequest struct {
	CommittedPreparationTime time.Time `json:"committedPreparationTime"`
}
