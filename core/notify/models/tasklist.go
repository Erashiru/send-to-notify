package models

import "github.com/pkg/errors"

var ErrNoTaskList = errors.New("no such task list")

type TaskList int

const (
	ORDER_ERROR TaskList = iota + 1
	MENU_ERROR
)
