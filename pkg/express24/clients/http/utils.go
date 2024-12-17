package http

import "time"

const (
	tokenTimeout = 59 * time.Minute
)

const (
	retriesNumber   = 3
	retriesWaitTime = 1 * time.Second
)

const (
	acceptHeader      = "Accept"
	authHeader        = "Authorization-token"
	contentTypeHeader = "Content-Type"
	jsonType          = "application/json"
)
