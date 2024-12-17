package http

import "time"

const (
	tokenTimeout = 59 * time.Minute
)

const (
	acceptHeader      = "Accept"
	authHeader        = "Authorization-token"
	contentTypeHeader = "Content-Type"
	jsonType          = "application/json"
)
