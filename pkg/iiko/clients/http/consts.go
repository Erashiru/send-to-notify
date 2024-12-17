package http

import "time"

const (
	BaseURL = "https://api-ru.iiko.services"
)

const (
	retriesNumber   = 5
	retriesWaitTime = 1 * time.Second
)
