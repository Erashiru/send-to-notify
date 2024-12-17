package dto

type TaskType string

const (
	GetOrder    TaskType = "GetOrder"
	CreateOrder TaskType = "CreateOrder"
	Task        TaskType = "GetTaskResponse"
	CancelOrder TaskType = "CancelOrder"
	PayOrder    TaskType = "PayOrder"
)

func (t TaskType) String() string {
	return string(t)
}
