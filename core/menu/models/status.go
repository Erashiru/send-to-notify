package models

type Status string

const (
	SUCCESS                    Status = "SUCCESS"             // Загружен в агрегатор
	PARTIALLY_PROCESSED        Status = "PARTIALLY_PROCESSED" // Частично загружен в агрегатор
	PROCESSING                 Status = "PROCESSING"          // Публикуется
	NOT_PROCESSED              Status = "NOT_PROCESSED"       // Не загружен в аггрегатор
	READY                      Status = "READY"               // Готов к публикации
	NOT_READY                  Status = "NOT_READY"           // Не готов к публикации
	HAS_PROMO                  Status = "HAS_PROMO"           // Идет промоакция
	ERROR                      Status = "ERROR"
	FAILED                     Status = "FAILED"
	FETCH_MENU_INVALID_PAYLOAD Status = "FETCH_MENU_INVALID_PAYLOAD"
	FETCH_MENU_SERVER_ERROR    Status = "FETCH_MENU_SERVER_ERROR"
	FETCH_MENU_UNAUTHORIZED    Status = "FETCH_MENU_UNAUTHORIZED"
	LIMIT_EXCEEDED             Status = "LIMIT_EXCEEDED"
	DELIVERY_SERVICE_ERROR     Status = "DELIVERY_SERVICE_ERROR"
	KWAAKA_ERROR               Status = "KWAAKA_ERROR"
)

var TransactionStatuses = []Status{
	KWAAKA_ERROR,
	DELIVERY_SERVICE_ERROR,
	LIMIT_EXCEEDED,
	NOT_PROCESSED,
	FETCH_MENU_UNAUTHORIZED,
	FETCH_MENU_SERVER_ERROR,
	FETCH_MENU_INVALID_PAYLOAD,
	FAILED,
	ERROR,
	NOT_READY,
	READY,
	PROCESSING,
	PARTIALLY_PROCESSED,
	SUCCESS,
}

func (s Status) String() string {
	return string(s)
}

func (s Status) ValidStatus(validStatuses []Status) bool {
	for _, valid := range validStatuses {
		if valid == s {
			return true
		}
	}
	return false
}
