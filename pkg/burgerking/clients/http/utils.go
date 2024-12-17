package http

import (
	"time"
)

const (
	baseURL = "http://service.tillypad.ru:8059"
)

const (
	retriesNumber   = 5
	retriesWaitTime = 1 * time.Second
)

const (
	acceptHeader      = "Accept"
	authHeader        = "Authorization"
	contentTypeHeader = "Content-Type"

	jsonType = "application/json"
)
