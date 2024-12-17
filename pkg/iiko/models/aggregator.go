package models

type Aggregator string

const (
	WOLT         Aggregator = "wolt"
	GLOVO        Aggregator = "glovo"
	QRMENU       Aggregator = "qr_menu"
	YANDEX       Aggregator = "yandex"
	MOYSKLAD     Aggregator = "moysklad"
	EMENU        Aggregator = "emenu"
	EXPRESS24    Aggregator = "express24"
	KWAAKA_ADMIN Aggregator = "kwaaka_admin"
	STARTERAPP   Aggregator = "starter_app"
)

func (a Aggregator) String() string {
	return string(a)
}
