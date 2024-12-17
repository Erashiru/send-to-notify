package http

import (
	"time"
)

const (
	retriesNumber   = 3
	retriesWaitTime = 1 * time.Second
)

const (
	acceptHeader      = "Accept"
	authHeader        = "Authorization"
	contentTypeHeader = "Content-Type"
	jsonType          = "application/json"
)
