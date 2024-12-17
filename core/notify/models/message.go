package models

type Message struct {
	Title       string
	Description string
	Services    []Service
	TaskList    TaskList
}

type Result struct {
	Service Service
	Status  Status
	Message string
}
