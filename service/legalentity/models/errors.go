package models

import "errors"

var (
	ErrInvalidID     = errors.New("invalid id")
	ErrDuplicateData = errors.New("some field duplicates existing data")
	ErrNotFound      = errors.New("document not found")
	ErrInvalidInput  = errors.New("some fields in input are wrong")
)
