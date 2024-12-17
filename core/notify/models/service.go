package models

import "github.com/pkg/errors"

var ErrNoService = errors.New("no such service")

type Service int

const (
	BITRIX Service = iota + 1
	SQS
	TELEGRAM
	WHATSAPP
	CLICKUP
)
