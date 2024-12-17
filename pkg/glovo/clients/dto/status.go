package dto

type Status string

const (
	SUCCESS                    Status = "SUCCESS"
	PROCESSING                 Status = "PROCESSING"
	NOT_PROCESSED              Status = "NOT_PROCESSED"
	FETCH_MENU_INVALID_PAYLOAD Status = "FETCH_MENU_INVALID_PAYLOAD"
	FETCH_MENU_SERVER_ERROR    Status = "FETCH_MENU_SERVER_ERROR"
	FETCH_MENU_UNAUTHORIZED    Status = "FETCH_MENU_UNAUTHORIZED"
	LIMIT_EXCEEDED             Status = "LIMIT_EXCEEDED"
	GLOVO_ERROR                Status = "GLOVO_ERROR"
)

func (s Status) String() string {
	return string(s)
}
