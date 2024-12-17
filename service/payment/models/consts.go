package models

const (
	PAYME_CREATE_TRANSACTION_METHOD        = "PAYME_CREATE_TRANSACTION"
	PAYME_CHECK_PERFORM_TRANSACTION_METHOD = "PAYME_CHECK_PERFORM_TRANSACTION"
	PAYME_PERFORM_TRANSACTION_METHOD       = "PAYME_PERFORM_TRANSACTION"
)

const (
	CARD_APPROVED           = "CARD_APPROVED"
	CARD_DECLINED           = "CARD_DECLINED"
	PAYMENT_DECLINED        = "PAYMENT_DECLINED"
	PAYMENT_APPROVED        = "PAYMENT_APPROVED"
	PAYMENT_CAPTURED        = "PAYMENT_CAPTURED"
	PAYMENT_CANCELED        = "PAYMENT_CANCELED"
	PAYMENT_ACTION_REQUIRED = "PAYMENT_ACTION_REQUIRED"
	ORDER_EXPIRED           = "ORDER_EXPIRED"
	CAPTURE_DECLINED        = "CAPTURE_DECLINED"
	CANCEL_DECLINED         = "CANCEL_DECLINED"
	REFUND_APPROVED         = "REFUND_APPROVED"
	REFUND_DECLINED         = "REFUND_DECLINED"
)

const (
	UNPAID                = "UNPAID"
	ON_HOLD               = "ON_HOLD"
	PAID                  = "PAID"
	EXPIRED               = "EXPIRED"
	REDIRECTED_TO_PAYMENT = "REDIRECTED_TO_PAYMENT"
)

const (
	IOKA           = "ioka"
	PAYME          = "payme"
	KASPI          = "kaspi"
	WOOPPAY        = "wooppay"
	WHATSAPP       = "whatsapp"
	CallCenter     = "callcenter"
	KaspiManual    = "kaspi_manual"
	KaspiSaleScout = "kaspi_salescout"
	Cash           = "cash"
	MultiCard      = "multicard"
)
