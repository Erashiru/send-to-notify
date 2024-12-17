package dto

type DeliveryService string

const (
	GLOVO  DeliveryService = "glovo"
	WOLT   DeliveryService = "wolt"
	Yandex DeliveryService = "yandex"
)

func (ds DeliveryService) String() string {
	return string(ds)
}
