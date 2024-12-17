package http

import "time"

const (
	BaseURL      = "https://api.paloma365.com"
	tokenTimeout = 59 * time.Minute
)

const (
	retriesNumber   = 5
	retriesWaitTime = 1 * time.Second
)

const (
	acceptHeader      = "Accept"
	authHeader        = "Authorization"
	contentTypeHeader = "Content-Type"
	jsonType          = "application/json"
)
