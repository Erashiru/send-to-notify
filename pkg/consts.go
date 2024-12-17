package pkg

import "time"

const (
	TokenTimeout = 59 * time.Minute
)

const (
	RetriesNumber   = 5
	RetriesWaitTime = 1 * time.Second
)

const (
	AcceptHeader      = "Accept"
	AuthHeader        = "Authorization"
	ContentTypeHeader = "Content-Type"
	JsonType          = "application/json"
	XMLType           = "application/xml"
)
