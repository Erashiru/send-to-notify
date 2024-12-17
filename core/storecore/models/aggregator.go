package models

const (
	WOLT         AggregatorName = "wolt"
	GLOVO        AggregatorName = "glovo"
	YANDEX       AggregatorName = "yandex"
	QRMENU       AggregatorName = "qr_menu"
	SELFDELIVERY AggregatorName = "self-delivery"
)

type AggregatorName string

func (a AggregatorName) String() string {
	return string(a)
}
