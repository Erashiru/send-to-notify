package models

type Aggregator string

const (
	WOLT         Aggregator = "wolt"
	GLOVO        Aggregator = "glovo"
	YANDEX       Aggregator = "yandex"
	QRMENU       Aggregator = "qr_menu"
	EMENU        Aggregator = "emenu"
	MOYSKLAD     Aggregator = "moysklad"
	CHOCOFOOD    Aggregator = "chocofood"
	TALABAT      Aggregator = "talabat"
	DELIVEROO    Aggregator = "deliveroo"
	EXPRESS24    Aggregator = "express24"
	KWAAKA_ADMIN Aggregator = "kwaaka_admin"
	STARTERAPP   Aggregator = "starter_app"
)

func (a Aggregator) String() string {
	switch a {
	case WOLT:
		return "wolt"
	case GLOVO:
		return "glovo"
	case YANDEX:
		return "yandex"
	case QRMENU:
		return "qr_menu"
	case EMENU:
		return "emenu"
	case MOYSKLAD:
		return "moysklad"
	case CHOCOFOOD:
		return "chocofood"
	case DELIVEROO:
		return "deliveroo"
	case TALABAT:
		return "talabat"
	case EXPRESS24:
		return "express24"
	case KWAAKA_ADMIN:
		return "kwaaka_admin"
	case STARTERAPP:
		return "starter_app"
	}

	return ""
}
